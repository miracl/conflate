package conflate

import (
	"io/ioutil"
	"net/http"
	. "net/url"
)

var emptyURL = URL{}

func loadURL(url URL) ([]byte, error) {
	if url.Scheme == "file" {
		// attempt to load locally handling case where we are loading from fifo etc
		b, err := ioutil.ReadFile(url.Path)
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

func toURLs(root URL, paths ...string) ([]URL, error) {
	var urls []URL
	for _, path := range paths {
		url, err := toURL(root, path)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func toURL(root URL, path string) (URL, error) {
	if path == "" {
		return emptyURL, makeError("The file path is blank")
	}
	var err error
	if root == emptyURL {
		root, err = workingURL()
		if err != nil {
			return emptyURL, err
		}
	}
	url, err := Parse(path)
	if err != nil {
		return emptyURL, wrapError(err, "Could not parse path")
	}
	if !url.IsAbs() {
		url = root.ResolveReference(url)
		url.RawQuery = root.RawQuery
	}
	return *url, nil
}

func workingURL() (URL, error) {
	wd, err := getwd()
	if err != nil {
		return emptyURL, err
	}
	wdURL, err := Parse("file://" + wd + "/")
	if err != nil {
		return emptyURL, err
	}
	return *wdURL, nil
}
