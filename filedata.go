package conflate

import (
	. "net/url"
	"os"
	"path/filepath"
	"strings"
)

type filedata struct {
	include  include
	bytes    []byte
	obj      map[string]interface{}
	includes []include
}

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

func newFiledata(bytes []byte, inc include) (filedata, error) {
	fd := filedata{bytes: bytes, include: inc}
	err := fd.unmarshal()
	if err != nil {
		return filedata{}, fd.wrapError(err)
	}
	err = fd.extractIncludes()
	if err != nil {
		return filedata{}, fd.wrapError(wrapError(err, "Could not extract includes"))
	}
	return fd, nil
}

func newExpandedFiledata(bytes []byte, inc include) (filedata, error) {
	return newFiledata(recursiveExpand(bytes), inc)
}

func (fd *filedata) wrapError(err error) error {
	if fd.include.isEmpty() {
		return err
	}
	return wrapError(err, "Error processing %v", fd.include.URL.String())
}

func (fd *filedata) unmarshal() error {
	ext := strings.ToLower(filepath.Ext(fd.include.URL.Path))
	unmarshallers, ok := Unmarshallers[ext]
	if !ok {
		unmarshallers = Unmarshallers[""]
	}
	err := makeError("Could not unmarshal data")
	for _, unmarshal := range unmarshallers {
		uerr := unmarshal(fd.bytes, &fd.obj)
		if uerr == nil {
			return nil
		}
		err = wrapError(uerr, err.Error())
	}
	return err
}

func unmarshalInclude(root URL, data []byte) (include, error) {
	var path string
	err := JSONUnmarshal(data, &path)
	if err == nil {
		return newIncludeFromPath(root, path)
	}
	var inc include
	err = JSONUnmarshal(data, &inc)

	return inc, err
}

func (fd *filedata) extractIncludes() error {
	if Includes == "" {
		return nil
	}
	if fd.obj[Includes] == nil {
		delete(fd.obj, Includes)
		return nil
	}
	objs, ok := fd.obj[Includes].([]interface{})
	if !ok {
		return makeError("Includes must be an array")
	}
	var incs includes
	for _, obj := range objs {
		data, err := jsonMarshal(obj)
		if err != nil {
			return err
		}
		inc, err := unmarshalInclude(fd.include.URL, data)
		if err != nil {
			return err
		}
		incs = append(incs, inc)
	}
	fd.includes = incs
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
