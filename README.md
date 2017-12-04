# geocatalogo

[![Report Card](https://goreportcard.com/badge/github.com/tomkralidis/geocatalogo)](https://goreportcard.com/report/github.com/tomkralidis/geocatalogo)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/tomkralidis/geocatalogo)

## Overview

Geospatial Catalogue in Go

## Installation

```bash
# install geocatalogo
go get github.com/tomkralidis/geocatalogo/...
# set configuration
# (sample in $GOPATH/src/github.com/tomkralidis/geocatalogo/default.env)
. default.env
```

## Running

```bash
# list commands
geocatalogo

# index a metadata record
geocatalogo index --file=/path/to/record.xml

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

