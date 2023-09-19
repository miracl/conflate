// Package main constructs the conflate library main functionality.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/miracl/conflate"
)

var version = "devel"

func failIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//nolint:funlen // that's ok
func main() {
	var data dataFlag

	flag.Var(&data, "data", "The path/url of JSON/YAML/TOML data, or 'stdin' to read from standard input")
	schemaFile := flag.String("schema", "", "The path/url of a JSON v4 schema file")
	defaults := flag.Bool("defaults", false, "Apply defaults from schema to data")
	validate := flag.Bool("validate", false, "Validate the data against the schema")
	format := flag.String("format", "", "Output format of the data JSON/YAML/TOML")
	includes := flag.String("includes", "includes", "Name of includes array. Blank string suppresses expansion of includes arrays")
	noincludes := flag.Bool("noincludes", false, "Switches off conflation of includes. Overrides any --includes setting.")
	expand := flag.Bool("expand", false, "Expand environment variables in files")
	showVersion := flag.Bool("version", false, "Display the version number")

	flag.Parse()

	if *showVersion {
		fmt.Println(version)

		return
	}

	conflate.Includes = *includes
	if *noincludes {
		conflate.Includes = ""
	}

	c := conflate.New()
	c.Expand(*expand)

	if len(data) == 0 {
		data = append(data, "stdin")
	}

	for _, d := range data {
		if d == "stdin" {
			b, err := io.ReadAll(os.Stdin)
			failIfError(err)

			err = c.AddData(b)
			failIfError(err)
		} else {
			err := c.AddFiles(d)
			failIfError(err)
		}
	}

	var schema *conflate.Schema

	if *schemaFile != "" {
		s, err := conflate.NewSchemaFile(*schemaFile)
		failIfError(err)

		schema = s
	}

	if *defaults {
		err := c.ApplyDefaults(schema)
		failIfError(err)
	}

	if *validate {
		err := c.Validate(schema)
		failIfError(err)
	}

	if *format != "" {
		var data interface{}
		err := c.Unmarshal(&data)
		failIfError(err)

		var out []byte

		switch strings.ToUpper(*format) {
		case "JSON":
			out, err = c.MarshalJSON()
		case "YAML":
			out, err = c.MarshalYAML()
		case "TOML":
			out, err = c.MarshalTOML()
		}

		failIfError(err)

		_, err = os.Stdout.Write(out)
		if err != nil {
			fmt.Printf("err when formatting: %v", err.Error())
		}
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
