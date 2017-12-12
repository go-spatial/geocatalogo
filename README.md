# geocatalogo

[![Build Status](https://travis-ci.org/tomkralidis/geocatalogo.png)](https://travis-ci.org/tomkralidis/geocatalogo)
[![Report Card](https://goreportcard.com/badge/github.com/tomkralidis/geocatalogo)](https://goreportcard.com/report/github.com/tomkralidis/geocatalogo)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/tomkralidis/geocatalogo)

## Overview

Geospatial Catalogue in Go

## Installation

```bash
# install dependencies
go get golang.org/x/text/encoding
# install geocatalogo
go get github.com/tomkralidis/geocatalogo/...
# install utilities/helpers
go get github.com/blevesearch/bleve/...
go install github.com/golang/lint/...
# set configuration
# (sample in $GOPATH/src/github.com/tomkralidis/geocatalogo/geocatalogo-config.env)
. geocatalogo-config.env
```

## Running

```bash
# list commands
geocatalogo

# index a metadata record
geocatalogo index --file=/path/to/record.xml

# index a directory of metadata records
geocatalogo index --dir=/path/to/dir

# search index
geocatalogo search --term=landsat

# get version
geocatalogo version
```

### Using the API

TODO

## Development

### Running Tests

## Releasing

### Bugs and Issues

All bugs, enhancements and issues are managed on [GitHub](https://github.com/tomkralidis/geocatalogo).

## Contact

* [Tom Kralidis](https://github.com/tomkralidis)

