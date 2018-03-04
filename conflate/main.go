package main

import (
	"flag"
	"fmt"
	"github.com/miracl/conflate"
	"os"
	"strings"
)

func failIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

	var data dataFlag
	flag.Var(&data, "data", "The path/url of JSON/YAML/TOML data")
	schema := flag.String("schema", "", "The path/url of a JSON v4 schema file")
	defaults := flag.Bool("defaults", false, "Apply defaults from schema to data")
	validate := flag.Bool("validate", false, "Validate the data against the schema")
	format := flag.String("format", "", "Output format of the data JSON/YAML/TOML")
	expand := flag.Bool("expand", false, "Expand environment variables in files")

	flag.Parse()

	c := conflate.New()
	c.Expand(*expand)

	err := c.AddFiles(data...)
	failIfError(err)

	if *schema != "" {
		err = c.SetSchemaFile(*schema)
		failIfError(err)
	}
	if *defaults {
		err = c.ApplyDefaults()
		failIfError(err)
	}
	if *validate {
		err = c.Validate()
		failIfError(err)
	}
	if *format != "" {
		var data interface{}
		err = c.Unmarshal(&data)
		failIfError(err)

		var out []byte
		switch *format {
		case "JSON":
			out, err = c.MarshalJSON()
		case "YAML":
			out, err = c.MarshalYAML()
		case "TOML":
			out, err = c.MarshalTOML()
		}
		failIfError(err)
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
