package conflate

import (
	pkgurl "net/url"
	"path/filepath"
	"strings"
)

type filedata struct {
	url      pkgurl.URL
	bytes    []byte
	obj      map[string]interface{}
	includes []string
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

func newFiledata(bytes []byte, url pkgurl.URL) (filedata, error) {
	fd := filedata{bytes: bytes, url: url}
	err := fd.unmarshal()
	if err != nil {
		return fd, err
	}
	err = fd.extractIncludes()
	return fd, err
}

func wrapFiledata(bytes []byte) (filedata, error) {
	return newFiledata(bytes, emptyURL)
}

func wrapFiledatas(bytes ...[]byte) (filedatas, error) {
	var fds []filedata
	for _, b := range bytes {
		fd, err := wrapFiledata(b)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fd)
	}
	return fds, nil
}

func (fd *filedata) unmarshal() error {
	ext := strings.ToLower(filepath.Ext(fd.url.Path))
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

func (fd *filedata) extractIncludes() error {
	err := jsonMarshalUnmarshal(fd.obj["includes"], &fd.includes)
	if err != nil {
		return wrapError(err, "Could not extract includes")
	}
	delete(fd.obj, "includes")
	return nil
}

func (fds filedatas) objs() []interface{} {
	var objs []interface{}
	for _, fd := range fds {
		objs = append(objs, fd.obj)
	}
	return objs
}
