///////////////////////////////////////////////////////////////////////////////
//
// The MIT License (MIT)
// Copyright (c) 2017 Tom Kralidis
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
// USE OR OTHER DEALINGS IN THE SOFTWARE.
//
///////////////////////////////////////////////////////////////////////////////

// Package geocatalogo provides the main interactions
// with the geospatial catalogue
package geocatalogo

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/tomkralidis/geocatalogo/config"
	"github.com/tomkralidis/geocatalogo/metadata"
	"github.com/tomkralidis/geocatalogo/repository"
	"github.com/tomkralidis/geocatalogo/search"
)

// VERSION provides the geocatalogo version installed.
const VERSION string = "0.1.0"

var log = logrus.New()

// GeoCatalogue provides the core structure
type GeoCatalogue struct {
	Config     config.Config
	Repository repository.Repository
}

// New provides the initializing functionality
func New() GeoCatalogue {

	c := GeoCatalogue{}
	c.Config = config.LoadFromEnv()

	// setup logging
	InitLog(&c.Config, log)

	log.Info("geocatalogo version " + VERSION)
	log.Info("Configuration: " + os.Getenv("GEOCATALOGO_CONFIG"))

	log.Info("Loading repository")
	c.Repository = repository.Open(c.Config, log)

	return c
}

// Index adds a metadata record to the Index
func (c *GeoCatalogue) Index(record metadata.Record) bool {
	log.Info("Indexing " + record.Identifier)
	err := c.Repository.Insert(record)
	if err != nil {
		log.Errorf("Indexing failed: %v", err)
		return false
	}
	return true
}

// UnIndex removes a metadata record to the Index
func (c *GeoCatalogue) UnIndex() bool {
	return c.Repository.Delete()
}

// Search performs a search/query against the Index
func (c *GeoCatalogue) Search(term string, from int, size int) search.SearchResults {
	sr := search.SearchResults{}
	log.Info("Searching index")
	err := c.Repository.Query(term, &sr, from, size)
	if err != nil {
		return sr
	}
	return sr
}

// Get retrieves a single metadata record from the Index
func (c *GeoCatalogue) Get() bool {
	return c.Repository.Get()
}
