<p align="center"><img src="gophers.png" alt="gophers" style="width: 50%; height: 50%"></p>

# CONFLATE

_Library providing routines to merge and validate JSON, YAML, TOML files and/or structs ([godoc](https://godoc.org/github.com/miracl/conflate))_

_Typical use case: Make your application configuration files **multi-format**, **modular**, **templated**, **sparse**, **location-independent** and **validated**_

[![Build Status](https://secure.travis-ci.org/miracl/conflate.png?branch=master)](https://travis-ci.org/miracl/conflate?branch=master)
[![Coverage Status](https://coveralls.io/repos/miracl/conflate/badge.svg?branch=master&service=github)](https://coveralls.io/github/miracl/conflate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/miracl/conflate)](https://goreportcard.com/report/github.com/miracl/conflate)

## Description

Conflate is a library and cli-tool, that provides the following features :

* merge data from multiple formats (JSON/YAML/TOML/go structs) and multiple locations (filesystem paths and urls)
* validate the merged data against a JSON schema
* apply any default values defined in a JSON schema to the merged data
* expand environment variables inside the data
* marshal merged data to multiple formats (JSON/YAML/TOML/go structs)

Improvements, ideas and bug fixes are welcomed.

## Getting started

Run the following command, which will build and install the latest binary in $GOPATH/bin

```
go get github.com/miracl/conflate/...
```
Alternatively, you can install one of the pre-built release binaries from https://github.com/miracl/conflate/releases

## Usage of Library

Please refer to the [godoc](https://godoc.org/github.com/miracl/conflate) and the [example code](./example/main.go)

## Usage of CLI Tool

Help can be obtained in the usual way :

```bash
$conflate --help
Usage of conflate:
  -data value
    	The path/url of JSON/YAML/TOML data, or 'stdin' to read from standard input
  -defaults
    	Apply defaults from schema to data
  -expand
    	Expand environment variables in files
  -format string
    	Output format of the data JSON/YAML/TOML
  -includes string
    	Name of includes array. Blank string suppresses expansion of includes arrays (default "includes")
  -noincludes
    	Switches off conflation of includes. Overrides any --includes setting.
  -schema string
    	The path/url of a JSON v4 schema file
  -validate
    	Validate the data against the schema
  -version
    	Display the version number
```

To conflate the following file ... :

```bash
$cat ./testdata/valid_parent.json
{
  "includes": [
    "valid_child.json",
    "valid_sibling.json"
  ],
  "parent_only" : "parent", 
  "parent_child" : "parent", 
  "parent_sibling" : "parent", 
  "all": "parent"
}
```

...run the following command, which will merge [valid_parent.json](https://raw.githubusercontent.com/miracl/conflate/master/testdata/valid_parent.json), 
[valid_child.json](https://raw.githubusercontent.com/miracl/conflate/master/testdata/valid_child.json), [valid_sibling.json](https://raw.githubusercontent.com/miracl/conflate/master/testdata/valid_sibling.json),  :

```bash
$conflate -data ./testdata/valid_parent.json -format JSON
{
  "all": "parent",
  "child_only": "child",
  "parent_child": "parent",
  "parent_only": "parent",
  "parent_sibling": "parent",
  "sibling_child": "sibling",
  "sibling_only": "sibling"
}
```
Note how the `includes` are loaded remotely as relative paths.

Also, note values in a file override values in any included files, and that an included file overrides values in any included file above it in the `includes` list.

If you instead host a file somewhere else, then just use a URL :

```bash
$conflate -data https://raw.githubusercontent.com/miracl/conflate/master/testdata/valid_parent.json -format JSON
{
  "all": "parent",
  "child_only": "child",
  "parent_child": "parent",
  "parent_only": "parent",
  "parent_sibling": "parent",
  "sibling_child": "sibling",
  "sibling_only": "sibling"
}

```

The `includes` here are also loaded as relative urls and follow exactly the same merging rules.

To output in a different format use the `-format` option, e.g. TOML :

```bash
$conflate -data ./testdata/valid_parent.json -format TOML
all = "parent"
child_only = "child"
parent_child = "parent"
parent_only = "parent"
parent_sibling = "parent"
sibling_child = "sibling"
sibling_only = "sibling"
```

To additionally use defaults from a JSON [schema](https://raw.githubusercontent.com/miracl/conflate/master/testdata/test.schema.json) and validate the conflated data against the schema, use `-defaults` and `-validate` respectively :

```bash
$cat ./testdata/blank.yaml

$conflate -data ./testdata/blank.yaml -schema ./testdata/test.schema.json -validate -format YAML
Schema validation failed : The document is not valid against the schema : Invalid type. Expected: object, given: null (#)

$conflate -data ./testdata/blank.yaml -schema ./testdata/test.schema.json -defaults -validate -format YAML
all: parent
child_only: child
parent_child: parent
parent_only: parent
parent_sibling: parent
sibling_child: sibling
sibling_only: sibling
```

Note any defaults are applied before validation is performed, as you would expect.

If you don't want to intrusively embed an `"includes"` array inside your JSON, you can instead provide multiple data files which are processed from left-to-right :

```bash
$conflate -data ./testdata/valid_child.json -data ./testdata/valid_sibling.json -format JSON
{
  "all": "sibling",
  "child_only": "child",
  "parent_child": "child",
  "parent_sibling": "sibling",
  "sibling_child": "sibling",
  "sibling_only": "sibling"
}
```

Or alternatively, you can create a top-level JSON file containing only the `includes` array. For fun, lets choose to use YAML for the top-level file, and output TOML :

```bash
$cat toplevel.yaml 
includes:
  - testdata/valid_child.json
  - testdata/valid_sibling.json

$conflate -data toplevel.yaml -format TOML
all = "sibling"
child_only = "child"
parent_child = "child"
parent_sibling = "sibling"
sibling_child = "sibling"
sibling_only = "sibling"
```

If you want to read a file from stdin you can do the following. Here we pipe in some TOML to override a value to demonstrate :

```bash
$echo 'all="MY OVERRIDDEN VALUE"' |  conflate -data ./testdata/valid_parent.json -data stdin  -format JSON
{
  "all": "MY OVERRIDDEN VALUE",
  "child_only": "child",
  "parent_child": "parent",
  "parent_only": "parent",
  "parent_sibling": "parent",
  "sibling_child": "sibling",
  "sibling_only": "sibling"
}
```

Note that in all cases `-data` sources are processed from left-to-right, with values in right files overriding values in left files, so the following doesnt work :

```bash
$echo 'all="MY OVERRIDDEN VALUE"' |  conflate -data stdin -data ./testdata/valid_parent.json  -format JSON
{
  "all": "parent",
  "child_only": "child",
  "parent_child": "parent",
  "parent_only": "parent",
  "parent_sibling": "parent",
  "sibling_child": "sibling",
  "sibling_only": "sibling"
}
```

You can optionally expand environment variables in the files like this :

```bash
$export echo MYVALUE="some value"
$export echo MYJSONMAP='{ "item1" : "value1" }'
$echo '{ "my_value": "$MYVALUE", "my_map": $MYJSONMAP }' | conflate -data stdin -expand -format JSON
{
  "my_map": {
    "item1": "value1"
  },
  "my_value": "some value"
}
```

# Acknowledgements

Images derived from originals by Renee French https://golang.org/doc/gopher/
