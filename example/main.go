package main

import (
	"fmt"
	"github.com/miracl/conflate"
	"path"
	"runtime"
)

func init() {
	// define the unmarshallers for the given file extensions, blank extension is the global unmarshaller
	conflate.Unmarshallers = conflate.UnmarshallerMap{
		".json": {customJSONUnmarshal},
		".jsn":  {conflate.JSONUnmarshal},
		".yaml": {conflate.YAMLUnmarshal},
		".yml":  {conflate.YAMLUnmarshal},
		".toml": {conflate.TOMLUnmarshal},
		".tml":  {conflate.TOMLUnmarshal},
		"":      {conflate.JSONUnmarshal, conflate.YAMLUnmarshal, conflate.TOMLUnmarshal},
	}
}

// example of a custom unmarshaller for JSON
func customJSONUnmarshal(data []byte, out interface{}) error {
	fmt.Println("Using custom JSON Unmarshaller")
	return conflate.JSONUnmarshal(data, out)
}

func main() {
	_, thisFile, _, _ := runtime.Caller(0)
	thisDir := path.Dir(thisFile)

	// merge multiple config files
	c, err := conflate.FromFiles(path.Join(thisDir, "../testdata/valid_parent.json"))
	if err != nil {
		fmt.Println(err)
		return
	}
	// load a json schema
	err = c.SetSchemaFile(path.Join(thisDir, "../testdata/test.schema.json"))
	if err != nil {
		fmt.Println(err)
		return
	}
	// apply defaults defined in schema to merged data
	err = c.ApplyDefaults()
	if err != nil {
		fmt.Println(err)
		return
	}
	// validate merged data against schema
	err = c.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}
	// unmarshal merged data to a struct/interface
	var data interface{}
	err = c.Unmarshal(&data)
	if err != nil {
		fmt.Println(err)
		return
	}
	// output merged data as json
	json, err := c.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	// output merged data as yaml
	yaml, err := c.MarshalYAML()
	if err != nil {
		fmt.Println(err)
		return
	}
	// output merged data as toml
	toml, err := c.MarshalTOML()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(json))
	fmt.Println("")
	fmt.Println(string(yaml))
	fmt.Println("")
	fmt.Println(string(toml))
}
