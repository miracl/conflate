# CONFLATE

_Library providing routines to merge and validate JSON, YAML, TOML files and/or structs ([godoc](https://godoc.org/github.com/miracl/conflate))_

[![Build Status](https://secure.travis-ci.org/miracl/conflate.png?branch=master)](https://travis-ci.org/miracl/conflate?branch=master)
[![Coverage Status](https://coveralls.io/repos/miracl/conflate/badge.svg?branch=master&service=github)](https://coveralls.io/github/miracl/conflate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/miracl/conflate)](https://goreportcard.com/report/github.com/miracl/conflate)

## Description

Conflate is a library that provides the following features :

* merge data from multiple formats (JSON/YAML/TOML/go structs) and multiple locations (filesystem paths and urls)
* validate the merged data against a JSON schema
* apply any default values defined in a JSON schema to the merged data
* expand environment variables inside the data
* marshal merged data to multiple formats (JSON/YAML/TOML/go structs)

Data files can include other files using the `includes` array, meaning that they will be merged (i.e. in JSON this is simply a string array at the top-level). The `includes` array can support
multiple path types (see below), including relative paths to local or remote files. Values present in the containing file override those from any file in the `includes` array, and values from any file in 
the `includes` array override those from files before it in the `includes` array :

```json
{
  "includes": [
    "./myfile.json",
    "/etc/myfolder/myfile.yaml",
    "file://etc/myfolder/myfile.toml",
    "http://mydomain/myfile"
  ],
  "mydata": 123
}
```

## Usage

For basic usage refer to the [example](./example/main.go), or use the [cli tool](./conflate)
