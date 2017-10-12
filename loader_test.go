package conflate

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
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

func TestToURL_RelativePathCwd(t *testing.T) {
	root, err := url.Parse("/home/username/service/")
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
	url, err := toURL(root, "/full/path/file")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Path, "/full/path/file")
}

func TestToURL_FullyQualifiedFileUrl(t *testing.T) {
	root, err := url.Parse("/no/matter/what/path")
	url, err := toURL(root, "file:/full/path/file")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, url.Path, "/full/path/file")
	assert.Equal(t, url.Scheme, "file")
}

func TestToURL_FullyQualifiedHttpPath(t *testing.T) {
	root, err := url.Parse("/no/matter/what/path/")
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

func TestLoad_Error(t *testing.T) {
	data, err := loadURL(url.URL{})
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

func TestLoad(t *testing.T) {
	url, err := url.Parse("http://www.miracl.com")
	data, err := loadURL(*url)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	root, err := workingDir()
	test, err := toURL(root, "./testdata/valid_parent.json")
	assert.Nil(t, err)
	assert.NotNil(t, test)
	data, err = loadURL(*url)
	assert.Nil(t, err)
	assert.NotNil(t, data)
}

// --------

func TestLoadAll_Error(t *testing.T) {
	data, err := loadAll(url.URL{})
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

func TestLoadAll(t *testing.T) {
	url1, err := url.Parse("http://www.miracl.com")
	root, err := workingDir()
	url2, err := toURL(root, "./testdata/valid_parent.json")
	data, err := loadAll(*url1, url2)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, len(data), 2)
}

// --------

func TestNewClient(t *testing.T) {
	c := newClient()
	assert.NotNil(t, c)
	assert.NotNil(t, c.Transport)
}

// --------

func TestExtractIncludes_Error(t *testing.T) {
	data := []byte{1, 2, 3}
	paths, err := extractIncludes(data)
	assert.NotNil(t, err)
	assert.Nil(t, paths)
}

func TestExtractIncludes(t *testing.T) {
	data := []byte(`{ "includes": [ "inc1", "inc2", "inc3"] }`)
	paths, err := extractIncludes(data)
	assert.Nil(t, err)
	assert.Equal(t, []string{"inc1", "inc2", "inc3"}, paths)
}

func TestExtractIncludes_NilData(t *testing.T) {
	data := []byte{}
	paths, err := extractIncludes(data)
	assert.Nil(t, err)
	assert.Nil(t, paths)
}

// --------

func TestLoadURLs_LoadError(t *testing.T) {
	data, err := loadURLs(url.URL{})
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

func TestLoadURLs_IncludesError(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "loader.go")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := loadURLs(url)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Could not extract includes")
	assert.Nil(t, data)
}

func TestLoadURLs_BadUrlInInclude(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/bad_url_in_include.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := loadURLs(url)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Could not parse path")
	assert.Nil(t, data)
}

func TestLoadURLs_MissingFileInInclude(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/missing_file_in_include.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := loadURLs(url)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to load url")
	assert.Nil(t, data)
}

func TestLoadURLs_RecursiveInclude(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/recursive_include_parent.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := loadURLs(url)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The url recursively includes itself")
	assert.Nil(t, data)
}

func TestLoadURLs(t *testing.T) {
	root, err := workingDir()
	assert.Nil(t, err)
	assert.NotNil(t, root)
	url, err := toURL(root, "testdata/valid_parent.json")
	assert.Nil(t, err)
	assert.NotNil(t, url)
	data, err := loadURLs(url)
	assert.Nil(t, err)
	assert.NotNil(t, data)
}
