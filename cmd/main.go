package main

import (
	"flag"
	"fmt"
	"github.com/miracl/conflate"
	"os"
	"strings"
)

func main() {

	var data dataFlag
	flag.Var(&data, "data", "The path/url of JSON/YAML/TOML data")
	schema := flag.String("schema", "", "The path/url of a JSON v4 schema file")
	defaults := flag.Bool("defaults", false, "Apply defaults from schema to data")
	validate := flag.Bool("validate", false, "Validate the data against the schema")
	format := flag.String("format", "", "Output format of the data JSON/YAML/TOML")

	flag.Parse()

	c, err := conflate.FromFiles(data...)
	if err != nil {
		fmt.Println(err)
		return
	}
	if *schema != "" {
		err = c.SetSchemaFile(*schema)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if *defaults {
		err = c.ApplyDefaults()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if *validate {
		err = c.Validate()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if *format != "" {
		var data interface{}
		err = c.Unmarshal(&data)
		if err != nil {
			fmt.Println(err)
			return
		}

		var out []byte
		switch *format {
		case "JSON":
			out, err = c.MarshalJSON()
		case "YAML":
			out, err = c.MarshalYAML()
		case "TOML":
			out, err = c.MarshalTOML()
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		os.Stdout.Write(out)
	}
}

type dataFlag []string

func (f *dataFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *dataFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}
