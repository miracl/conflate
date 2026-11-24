package conflate

import (
	"fmt"
	pkgurl "net/url"
	"os"
	"path/filepath"
	"strings"
)

type filedata struct {
	url      *pkgurl.URL
	data     []byte
	obj      map[string]any
	includes []string
}

var emptyFiledata = filedata{}

type filedatas []filedata

// UnmarshallerFunc defines the type of function used for unmarshalling data.
type UnmarshallerFunc func([]byte, any) error

// UnmarshallerFuncs defines the type for a slice of UnmarshallerFunc.
type UnmarshallerFuncs []UnmarshallerFunc

// UnmarshallerMap defines the type of a map of string to UnmarshallerFuncs.
type UnmarshallerMap map[string]UnmarshallerFuncs

// Unmarshallers is a list of unmarshalling functions to be used for given file extensions.
// The unmarshaller slice for the blank file extension is used when no match is found.
var Unmarshallers = UnmarshallerMap{
	".json": {JSONUnmarshal},
	".jsn":  {JSONUnmarshal},
	".yaml": {YAMLUnmarshal},
	".yml":  {YAMLUnmarshal},
	".toml": {TOMLUnmarshal},
	".tml":  {TOMLUnmarshal},
	"":      {JSONUnmarshal, YAMLUnmarshal, TOMLUnmarshal},
}

func newFiledata(data []byte, url *pkgurl.URL) (filedata, error) {
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

func newExpandedFiledata(data []byte, url *pkgurl.URL) (filedata, error) {
	return newFiledata(recursiveExpand(data), url)
}

func (fd *filedata) wrapError(err error) error {
	if fd == nil || fd.url == nil || *fd.url == emptyURL || err == nil {
		return err
	}

	return fmt.Errorf("error processing %v: %w", fd.url.String(), err)
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

	var err error

	for _, unmarshal := range unmarshallers {
		uerr := unmarshal(fd.data, &fd.obj)
		if uerr == nil {
			return nil
		}

		err = fmt.Errorf("could not unmarshal data: %w", uerr)
	}

	return err
}

func (fd *filedata) extractIncludes() error {
	if Includes == "" {
		return nil
	}

	err := jsonMarshalUnmarshal(fd.obj[Includes], &fd.includes)
	if err != nil {
		return fmt.Errorf("could not extract includes: %w", err)
	}

	delete(fd.obj, Includes)

	return nil
}

func (fds filedatas) objs() (objs []any) {
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

	for range maxExpansions {
		b, c = expand(b)
		if c == 0 {
			return b
		}
	}

	return b
}

func expand(b []byte) (result []byte, count int) {
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

//nolint:goconst // explicit definition
func getDefaultSchema() map[string]any {
	return map[string]any{
		"anyOf": []any{
			map[string]any{
				"type": "object",
				"properties": map[string]any{
					Includes: map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
			},
			map[string]any{
				"type": "null",
			},
		},
	}
}
