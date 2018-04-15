package conflate

import (
	. "net/url"
)

type include struct {
	Path string
	URL  URL
}

var emptyInclude = include{}

type includes []include

func newIncludeFromURL(url URL) include {
	return include{
		URL: url,
	}
}

func newIncludesFromURLs(urls ...URL) []include {
	var incs []include
	for _, url := range urls {
		incs = append(incs, newIncludeFromURL(url))
	}
	return incs
}

func newIncludeFromPath(root URL, path string) (include, error) {
	url, err := toURL(root, path)
	if err != nil {
		return include{}, err
	}
	return newIncludeFromURL(url), nil
}

func newIncludesFromPaths(root URL, paths ...string) ([]include, error) {
	urls, err := toURLs(root, paths...)
	if err != nil {
		return nil, err
	}
	return newIncludesFromURLs(urls...), nil
}

func (inc include) isEmpty() bool {
	return inc == include{}
}

func (inc include) isEqual(other include) bool {
	return inc.URL == other.URL
}

func (inc include) load() ([]byte, error) {
	return loadURL(inc.URL)
}

func (incs includes) contains(search include) bool {
	if search.isEmpty() {
		return false
	}
	for _, inc := range incs {
		if inc.isEqual(search) {
			return true
		}
	}
	return false
}
