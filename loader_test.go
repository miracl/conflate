package conflate

import (
	gocontext "context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

// --------

func TestWorkingDir_NoRootPath(t *testing.T) {
	oldGetwd := getwd
	getwd = func() (dir string, err error) {
		return "", makeError("No root error")
	}
	defer func() { getwd = oldGetwd }()
	url, err := workingDir()
	assert.NotNil(t, err)
	assert.Nil(t, url)
}

func TestWorkingDir_ParseError(t *testing.T) {
	oldGetwd := getwd
	getwd = func() (dir string, err error) {
		return "#^/\\^&%*&%", nil
	}
	defer func() { getwd = oldGetwd }()
	url, err := workingDir()
	assert.NotNil(t, err)
	assert.Nil(t, url)
}

func TestWorkingDir(t *testing.T) {
	url, err := workingDir()
	assert.NotNil(t, url)
	assert.Nil(t, err)
}

// --------

func TestToURL_Error(t *testing.T) {
	url, err := toURL(&emptyURL, "\\^&%")
	assert.NotNil(t, err)
	assert.Equal(t, url, emptyURL)
}

func TestToURL_Blank(t *testing.T) {
	_, err := toURL(nil, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The file path is blank")
}

func TestToURL_RelativePathCwd(t *testing.T) {
	root, err := url.Parse("/home/username/service/")
	assert.Nil(t, err)
	url, err := toURL(root, "../../fileName")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Path, "/home/fileName")

	url, err = toURL(root, "fileName")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Path, "/home/username/service/fileName")

	url, err = toURL(root, "./fileName")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Path, "/home/username/service/fileName")
}

func TestToURL_FullyQualifiedPath(t *testing.T) {
	root, err := url.Parse("/no/matter/what/path")
	assert.Nil(t, err)
	url, err := toURL(root, "/full/path/file")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Path, "/full/path/file")
}

func TestToURL_FullyQualifiedFileUrl(t *testing.T) {
	root, err := url.Parse("/no/matter/what/path")
	assert.Nil(t, err)
	url, err := toURL(root, "file:/full/path/file")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Path, "/full/path/file")
	assert.Equal(t, url.Scheme, "file")
}

func TestToURL_FullyQualifiedHttpPath(t *testing.T) {
	root, err := url.Parse("/no/matter/what/path/")
	assert.Nil(t, err)
	url, err := toURL(root, "http://www.some.url.com")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Scheme, "http")
	assert.Equal(t, url.Host, "www.some.url.com")

	url, err = toURL(root, "http://www.some.url.com/file")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Scheme, "http")
	assert.Equal(t, url.Host, "www.some.url.com")
	assert.Equal(t, url.Path, "/file")
}

func TestToURL_RelativeHttpUrl(t *testing.T) {
	root, err := url.Parse("https://www.some.url.com/path/inside/")
	assert.Nil(t, err)
	url, err := toURL(root, "./file")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Scheme, "https")
	assert.Equal(t, url.Host, "www.some.url.com")
	assert.Equal(t, url.Path, "/path/inside/file")

	url, err = toURL(root, "../file")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Scheme, "https")
	assert.Equal(t, url.Host, "www.some.url.com")
	assert.Equal(t, url.Path, "/path/file")
}

// --------

func TestToURLs_Error(t *testing.T) {
	urls, err := toURLs(&emptyURL, "", "\\^&%")
	assert.NotNil(t, err)
	assert.Nil(t, urls)
}

func TestToURLs(t *testing.T) {
	root, err := url.Parse("https://www.some.url.com/path/inside/")
	assert.Nil(t, err)
	urls, err := toURLs(root, "./one", "../two", "three")
	assert.Nil(t, err)
	assert.NotNil(t, urls)
	assert.Equal(t, len(urls), 3)
	assert.Equal(t, urls[0].Host, "www.some.url.com")
	assert.Equal(t, urls[0].Scheme, "https")
	assert.Equal(t, urls[0].Path, "/path/inside/one")
	assert.Equal(t, urls[1].Host, "www.some.url.com")
	assert.Equal(t, urls[1].Scheme, "https")
	assert.Equal(t, urls[1].Path, "/path/two")
	assert.Equal(t, urls[2].Host, "www.some.url.com")
	assert.Equal(t, urls[2].Scheme, "https")
	assert.Equal(t, urls[2].Path, "/path/inside/three")
}

// --------

func TestLoadURLError(t *testing.T) {
	data, err := loadURL(url.URL{})
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

func testServer() func() {
	var wg sync.WaitGroup
	wg.Add(1)
	server := &http.Server{
		Addr:    "0.0.0.0:9999",
		Handler: http.FileServer(http.Dir("./testdata")),
	}
	go func() {
		defer wg.Done()
		server.ListenAndServe()
	}()
	return func() {
		server.Shutdown(gocontext.Background())
	}
}
func testWaitForURL(t *testing.T, url string) {
	// wait for a couple of seconds for server to come up
	for i := 0; i < 4; i++ {
		_, err := http.Get(url)
		if err == nil {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	assert.FailNow(t, "could not connect to url : "+url)
}

func TestLoadURL(t *testing.T) {
	shutdown := testServer()
	defer shutdown()
	testWaitForURL(t, "http://0.0.0.0:9999")
	url, err := url.Parse("http://0.0.0.0:9999/valid_parent.json")
	assert.Nil(t, err)
	data, err := loadURL(*url)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Contains(t, string(data), "parent")
}

func TestLoadURL_Relative(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	url, err := toURL(root, "./testdata/valid_parent.json")
	assert.Nil(t, err)
	data, err := loadURL(url)
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Contains(t, string(data), "parent")
}

// --------

var testLoader = loader{newFiledata: newFiledata}

func TestLoadURLsRecursive_LoadError(t *testing.T) {
	data, err := testLoader.loadURLsRecursive(nil, url.URL{})
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

func TestLoadURLsRecursive_IncludesError(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "loader.go")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := testLoader.loadURLsRecursive(nil, url)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Could not unmarshal")
	assert.Nil(t, data)
}

func TestLoadURLsRecursive_BadUrlInInclude(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/bad_url_in_include.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := testLoader.loadURLsRecursive(nil, url)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Could not parse path")
	assert.Nil(t, data)
}

func TestLoadURLsRecursive_MissingFileInInclude(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/missing_file_in_include.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := testLoader.loadURLsRecursive(nil, url)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to load url")
	assert.Nil(t, data)
}

func TestLoadURLsRecursive_RecursiveInclude(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/recursive_include_parent.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := testLoader.loadURLsRecursive(nil, url)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The url recursively includes itself")
	assert.Nil(t, data)
}

func TestLoadURLsRecursive(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/valid_parent.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := testLoader.loadURLsRecursive(nil, url)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 3, len(data))
	assert.Contains(t, data[0].url.String(), "valid_child.json")
	assert.Contains(t, data[1].url.String(), "valid_sibling.json")
	assert.Contains(t, data[2].url.String(), "valid_parent.json")
}

func TestLoadURLsRecursive_BlankChildYaml(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/parent_blank.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := testLoader.loadURLsRecursive(nil, url)
	assert.Nil(t, err)
	assert.NotNil(t, data)
}

func TestLoadURLsRecursive_BlankChildJson(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/parent_blank.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := testLoader.loadURLsRecursive(nil, url)
	assert.Nil(t, err)
	assert.NotNil(t, data)
}

func TestLoadURLsRecursive_BlankChildToml(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/parent_blank.toml")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := testLoader.loadURLsRecursive(nil, url)
	assert.Nil(t, err)
	assert.NotNil(t, data)
}

func testPath(t *testing.T, urlPath string, filePath string) {
	assert.Equal(t, urlPath, setPath(filePath))
	assert.Equal(t, filePath, getPath(urlPath))
}

func TestPath_Windows(t *testing.T) {
	old := goos
	goos = "windows"
	defer func() { goos = old }()
	testPath(t, `/C:/`, `C:\`)
	testPath(t, `/C:/a`, `C:\a`)
	testPath(t, `/C:/a/`, `C:\a\`)
	testPath(t, `/C:/a`, `C:\a`)
	testPath(t, `/C:/a/`, `C:\a\`)
	testPath(t, `/c:/`, `c:\`)
	testPath(t, `/c:/a`, `c:\a`)
	testPath(t, `/c:/a/`, `c:\a\`)
	testPath(t, `/c:/a`, `c:\a`)
	testPath(t, `/c:/a/`, `c:\a\`)
	testPath(t, `unc`, `\\unc`)
	testPath(t, `unc/a`, `\\unc\a`)
	testPath(t, `unc/a/`, `\\unc\a\`)
}
