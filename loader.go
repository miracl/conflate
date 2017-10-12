package conflate

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var emptyURL = url.URL{}
var getwd = os.Getwd

func loadURLs(urls ...url.URL) ([][]byte, error) {
	return loadURLsRecursive(urls, nil)
}

func loadURLsRecursive(urls []url.URL, parentURLs []url.URL) ([][]byte, error) {
	var data [][]byte
	for _, url := range urls {
		for _, parentURL := range parentURLs {
			if url == parentURL {
				return nil, makeError("The url recursively includes itself (%s)", url.String())
			}
		}
		parentData, err := loadURL(url)
		if err != nil {
			return nil, err
		}
		childPaths, err := extractIncludes(parentData)
		if err != nil {
			return nil, err
		}
		childUrls, err := toURLs(&url, childPaths...)
		if err != nil {
			return nil, err
		}
		childData, err := loadURLsRecursive(childUrls, append(parentURLs, url))
		if err != nil {
			return nil, err
		}
		data = append(data, childData...)
		data = append(data, parentData)
	}
	return data, nil
}

func loadURL(url url.URL) ([]byte, error) {
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

func newClient() http.Client {
	return http.Client{Transport: newTransport()}
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

func loadAll(urls ...url.URL) ([][]byte, error) {
	var data [][]byte
	for _, url := range urls {
		datum, err := loadURL(url)
		if err != nil {
			return nil, err
		}
		data = append(data, datum)
	}
	return data, nil
}

func toURLs(rootURL *url.URL, paths ...string) ([]url.URL, error) {
	var urls []url.URL
	for _, path := range paths {
		url, err := toURL(rootURL, path)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func toURL(rootURL *url.URL, path string) (url.URL, error) {
	var err error
	if rootURL == nil {
		rootURL, err = workingDir()
		if err != nil {
			return emptyURL, err
		}
	}
	url, err := url.Parse(path)
	if err != nil {
		return emptyURL, wrapError(err, "Could not parse path")
	}
	if !url.IsAbs() {
		url = rootURL.ResolveReference(url)
		url.RawQuery = rootURL.RawQuery
	}
	return *url, nil
}

func extractIncludes(data []byte) ([]string, error) {
	out := struct {
		Includes []string
	}{}
	err := unmarshal(data, &out)
	if err != nil {
		return nil, wrapError(err, "Could not extract includes")
	}
	return out.Includes, nil
}

func workingDir() (*url.URL, error) {
	rootPath, err := getwd()
	if err != nil {
		return nil, err
	}
	rootURL, err := url.Parse("file://" + rootPath + "/")
	if err != nil {
		return nil, err
	}
	return rootURL, nil
}
