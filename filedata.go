package conflate

import (
	pkgurl "net/url"
	"os"
	"path/filepath"
	"strings"
)

type filedata struct {
	url      pkgurl.URL
	data     []byte
	obj      map[string]interface{}
	includes []string
}

var emptyFiledata = filedata{}

type filedatas []filedata

// UnmarshallerFunc defines the type of function used for unmarshalling data
type UnmarshallerFunc func([]byte, interface{}) error

// UnmarshallerFuncs defines the type for a slice of UnmarshallerFunc
type UnmarshallerFuncs []UnmarshallerFunc

// UnmarshallerMap defines the type of a map of string to UnmarshallerFuncs
type UnmarshallerMap map[string]UnmarshallerFuncs

// Unmarshallers is a list of unmarshalling functions to be used for given file extensions. The unmarshaller slice for the blank file extension is used when no match is found.
var Unmarshallers = UnmarshallerMap{
	".json": {JSONUnmarshal},
	".jsn":  {JSONUnmarshal},
	".yaml": {YAMLUnmarshal},
	".yml":  {YAMLUnmarshal},
	".toml": {TOMLUnmarshal},
	".tml":  {TOMLUnmarshal},
	"":      {JSONUnmarshal, YAMLUnmarshal, TOMLUnmarshal},
}

func newFiledata(data []byte, url pkgurl.URL) (filedata, error) {
	fd := filedata{data: data, url: url}
	err := fd.unmarshal()
	if err != nil {
		return emptyFiledata, err
	}
	err = fd.validate()
	if err != nil {
		return emptyFiledata, err
	}
	err = fd.extractIncludes()
	if err != nil {
		return emptyFiledata, err
	}
	return fd, nil
}

func newExpandedFiledata(data []byte, url pkgurl.URL) (filedata, error) {
	return newFiledata(recursiveExpand(data), url)
}

func (fd *filedata) wrapError(err error) error {
	if fd.url == emptyURL {
		return err
	}
	return wrapError(err, "Error processing %v", fd.url.String())
}

func (fd *filedata) validate() error {
	return fd.wrapError(validate(fd.obj, getSchema()))
}

func (fd *filedata) unmarshal() error {
	ext := strings.ToLower(filepath.Ext(fd.url.Path))
	unmarshallers, ok := Unmarshallers[ext]
	if !ok {
		unmarshallers = Unmarshallers[""]
	}
	err := makeError("Could not unmarshal data")
	for _, unmarshal := range unmarshallers {
		uerr := unmarshal(fd.data, &fd.obj)
		if uerr == nil {
			return nil
		}
		err = wrapError(uerr, err.Error())
	}
	return err
}

func (fd *filedata) extractIncludes() error {
	if Includes == "" {
		return nil
	}
	err := jsonMarshalUnmarshal(fd.obj[Includes], &fd.includes)
	if err != nil {
		return wrapError(err, "Could not extract includes")
	}
	delete(fd.obj, Includes)
	return nil
}

func (fds filedatas) objs() []interface{} {
	var objs []interface{}
	for _, fd := range fds {
		objs = append(objs, fd.obj)
	}
	return objs
}

func (fd *filedata) isEmpty() bool {
	return fd == nil || fd.obj == nil
}

func recursiveExpand(b []byte) []byte {
	const maxExpansions = 10
	var c int
	for i := 0; i < maxExpansions; i++ {
		b, c = expand(b)
		if c == 0 {
			return b
		}
	}
	return b
}

func expand(b []byte) ([]byte, int) {
	var c int
	return []byte(os.Expand(string(b),
		func(name string) string {
			val, ok := os.LookupEnv(name)
			if ok {
				c++
				return val
			}
			return "$" + name
		})), c
}

var getSchema = getDefaultSchema

func getDefaultSchema() map[string]interface{} {
	return map[string]interface{}{
		"anyOf": []interface{}{
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					Includes: map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
			map[string]interface{}{
				"type": "null",
			},
		},
	}
}
