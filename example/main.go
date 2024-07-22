// Package main demonstrates a sample conflate usage.
package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/miracl/conflate"
)

// example of a custom unmarshaller for JSON.
func customJSONUnmarshal(data []byte, out interface{}) error {
	fmt.Println("Using custom JSON Unmarshaller")

	return conflate.JSONUnmarshal(data, out)
}

func main() {
	// define the unmarshallers for the given file extensions, blank extension is the global unmarshaller
	conflate.Unmarshallers = conflate.UnmarshallerMap{
		".json": conflate.UnmarshallerFuncs{customJSONUnmarshal},
		".jsn":  conflate.UnmarshallerFuncs{conflate.JSONUnmarshal},
		".yaml": conflate.UnmarshallerFuncs{conflate.YAMLUnmarshal},
		".yml":  conflate.UnmarshallerFuncs{conflate.YAMLUnmarshal},
		".toml": conflate.UnmarshallerFuncs{conflate.TOMLUnmarshal},
		".tml":  conflate.UnmarshallerFuncs{conflate.TOMLUnmarshal},
		"":      conflate.UnmarshallerFuncs{conflate.JSONUnmarshal, conflate.YAMLUnmarshal, conflate.TOMLUnmarshal},
	}

	_, thisFile, _, _ := runtime.Caller(0) //nolint:dogsled // ok for an example
	thisDir := path.Dir(thisFile)

	// merge multiple config files
	c, err := conflate.FromFiles(path.Join(thisDir, "../testdata/valid_parent.json"))
	if err != nil {
		fmt.Println(err)

		return
	}

	// load a json schema
	schema, err := conflate.NewSchemaFile(path.Join(thisDir, "../testdata/test.schema.json"))
	if err != nil {
		fmt.Println(err)

		return
	}

	// apply defaults defined in schema to merged data
	err = c.ApplyDefaults(schema)
	if err != nil {
		fmt.Println(err)

		return
	}

	// validate merged data against schema
	err = c.Validate(schema)
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
