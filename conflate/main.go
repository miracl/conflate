package main

import (
	"flag"
	"fmt"
	"github.com/miracl/conflate"
	"io/ioutil"
	"os"
	"strings"
)

var version = "devel"

func failIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

	var data dataFlag
	flag.Var(&data, "data", "The path/url of JSON/YAML/TOML data, or 'stdin' to read from standard input")
	schema := flag.String("schema", "", "The path/url of a JSON v4 schema file")
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
			b, err := ioutil.ReadAll(os.Stdin)
			failIfError(err)
			err = c.AddData(b)
			failIfError(err)
		} else {
			err := c.AddFiles(d)
			failIfError(err)
		}
	}

	if *schema != "" {
		err := c.SetSchemaFile(*schema)
		failIfError(err)
	}
	if *defaults {
		err := c.ApplyDefaults()
		failIfError(err)
	}
	if *validate {
		err := c.Validate()
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
