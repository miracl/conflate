package conflate

import (
	"io/ioutil"
	"net"
	"net/http"
	pkgurl "net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	goos        = runtime.GOOS
	emptyURL    = pkgurl.URL{}
	getwd       = os.Getwd
	driveLetter = regexp.MustCompile(`^[A-Za-z]:.*$`)
)

type loader struct {
	newFiledata func([]byte, pkgurl.URL) (filedata, error)
}

func (l *loader) loadURLsRecursive(parentUrls []pkgurl.URL, urls ...pkgurl.URL) (filedatas, error) {
	var allData filedatas
	for _, url := range urls {
		data, err := l.loadURLRecursive(parentUrls, url)
		if err != nil {
			return nil, err
		}
		allData = append(allData, data...)
	}
	return allData, nil
}

func (l *loader) loadURLRecursive(parentUrls []pkgurl.URL, url pkgurl.URL) (filedatas, error) {
	data, err := loadURL(url)
	if err != nil {
		return nil, err
	}
	fdata, err := l.newFiledata(data, url)
	if err != nil {
		return nil, err
	}
	return l.loadDatumRecursive(parentUrls, &url, fdata)
}

func (l *loader) loadDataRecursive(parentUrls []pkgurl.URL, data ...filedata) (filedatas, error) {
	var allData filedatas
	for _, datum := range data {
		childData, err := l.loadDatumRecursive(parentUrls, nil, datum)
		if err != nil {
			return nil, err
		}
		allData = append(allData, childData...)
	}
	return allData, nil
}

func (l *loader) loadDatumRecursive(parentUrls []pkgurl.URL, url *pkgurl.URL, data filedata) (filedatas, error) {
	if data.isEmpty() {
		return nil, nil
	}
	if containsURL(url, parentUrls) {
		return nil, makeError("The url recursively includes itself (%v)", url)
	}
	childUrls, err := toURLs(url, data.includes...)
	if err != nil {
		return nil, err
	}
	var newParentUrls []pkgurl.URL
	newParentUrls = append(newParentUrls, parentUrls...)
	if url != nil {
		newParentUrls = append(newParentUrls, *url)
	}
	childData, err := l.loadURLsRecursive(newParentUrls, childUrls...)
	if err != nil {
		return nil, err
	}
	var allData filedatas
	allData = append(allData, childData...)
	allData = append(allData, data)
	return allData, nil
}

func (l *loader) wrapFiledata(bytes []byte) (filedata, error) {
	return l.newFiledata(bytes, emptyURL)
}

func (l *loader) wrapFiledatas(bytes ...[]byte) (filedatas, error) {
	var fds []filedata
	for _, b := range bytes {
		fd, err := l.wrapFiledata(b)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fd)
	}
	return fds, nil
}

func loadURL(url pkgurl.URL) ([]byte, error) {
	if url.Scheme == "file" {
		// attempt to load locally handling case where we are loading from fifo etc
		b, err := ioutil.ReadFile(getPath(url.Path))
		if err == nil {
			return b, nil
		}
	}
	client := http.Client{Transport: newTransport()}
	resp, err := client.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, makeError("Failed to load url : %v : %v", resp.StatusCode, url.String())
	}
	return data, err
}

func newTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	return transport
}

func toURLs(rootURL *pkgurl.URL, paths ...string) ([]pkgurl.URL, error) {
	var urls []pkgurl.URL
	for _, path := range paths {
		url, err := toURL(rootURL, path)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func toURL(rootURL *pkgurl.URL, path string) (pkgurl.URL, error) {
	if path == "" {
		return emptyURL, makeError("The file path is blank")
	}
	var err error
	if rootURL == nil {
		rootURL, err = workingDir()
		if err != nil {
			return emptyURL, err
		}
	}
	url, err := pkgurl.Parse(setPath(path))
	if err != nil {
		return emptyURL, wrapError(err, "Could not parse path")
	}
	if !url.IsAbs() {
		url = rootURL.ResolveReference(url)
		url.RawQuery = rootURL.RawQuery
	}
	return *url, nil
}

func containsURL(searchURL *pkgurl.URL, urls []pkgurl.URL) bool {
	if searchURL == nil {
		return false
	}
	for _, url := range urls {
		if url == *searchURL {
			return true
		}
	}
	return false
}

func workingDir() (*pkgurl.URL, error) {
	rootPath, err := getwd()
	if err != nil {
		return nil, err
	}
	rootURL, err := pkgurl.Parse("file://" + setPath(rootPath) + "/")
	if err != nil {
		return nil, err
	}
	return rootURL, nil
}

func setPath(path string) string {
	if goos == "windows" {
		// https://blogs.msdn.microsoft.com/ie/2006/12/06/file-uris-in-windows/
		path = strings.Replace(path, `\`, `/`, -1)
		path = strings.TrimLeft(path, `/`)
		if driveLetter.MatchString(path) {
			path = `/` + path
		}
	}
	return path
}

func getPath(path string) string {
	if goos == "windows" {
		// https://blogs.msdn.microsoft.com/ie/2006/12/06/file-uris-in-windows/
		path = strings.TrimLeft(path, `/`)
		if !driveLetter.MatchString(path) {
			path = `//` + path
		}
		path = strings.Replace(path, `/`, `\`, -1)
	}
	return path
}
