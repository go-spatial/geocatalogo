# geocatalogo

[![Report Card](https://goreportcard.com/badge/github.com/tomkralidis/geocatalogo)](https://goreportcard.com/report/github.com/tomkralidis/geocatalogo)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/tomkralidis/geocatalogo)



# pyshadoz

[![Build Status](https://travis-ci.org/WMO-ET-WDC/pyshadoz.png)](https://travis-ci.org/WMO-ET-WDC/pyshadoz)
[![Coverage Status](https://coveralls.io/repos/github/WMO-ET-WDC/pyshadoz/badge.svg?branch=master)](https://coveralls.io/github/WMO-ET-WDC/pyshadoz?branch=master)

## Overview

Geospatial Catalogue in Go

## Installation

```bash
# install geocatalogo
go get github.com/tomkralidis/geocatalogo/...
# set configuration
# (sample in $GOPATH/src/github.com/tomkralidis/geocatalogo/geocatalogo-config.yml)
export GEOCATALOGO_CONFIG=/path/to/geocatalogo-config.yml
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

