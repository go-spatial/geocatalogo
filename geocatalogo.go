///////////////////////////////////////////////////////////////////////////////
//
// The MIT License (MIT)
// Copyright (c) 2019 Tom Kralidis
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
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-spatial/geocatalogo/config"
	"github.com/go-spatial/geocatalogo/metadata"
	"github.com/go-spatial/geocatalogo/repository"
	"github.com/go-spatial/geocatalogo/search"
)

// VERSION provides the geocatalogo version installed.
const VERSION string = "0.2-dev"

var log = logrus.New()

// GeoCatalogue provides the core structure
type GeoCatalogue struct {
	Config     config.Config
	Repository *repository.Elasticsearch
}

// New provides the initializing functionality
func New(cfg *config.Config) (*GeoCatalogue, error) {

	c := GeoCatalogue{}
	c.Config = *cfg

	// setup logging
	InitLog(&c.Config, log)

	log.Info("geocatalogo version " + VERSION)
	log.Info("Configuration: " + os.Getenv("GEOCATALOGO_CONFIG"))

	log.Info("Loading repository")
	repo, err := repository.Open(c.Config, log)
	if err != nil {
		return &c, err
	}
	c.Repository = &repo

	return &c, nil
}

// NewFromEnv provides the initializing functionality
// using configuration from the environment
func NewFromEnv() (*GeoCatalogue, error) {
	cfg := config.LoadFromEnv()
	return New(&cfg)
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
func (c *GeoCatalogue) Search(collections []string, term string, bbox []float64, timeVal []time.Time, from int, size int) search.Results {
	sr := search.Results{}
	log.Info("Searching index")
	err := c.Repository.Query(collections, term, bbox, timeVal, from, size, &sr)
	if err != nil {
		log.Warn(err)
		return sr
	}
	return sr
}

// Get retrieves a single metadata record from the Index
func (c *GeoCatalogue) Get(identifiers []string) search.Results {
	sr := search.Results{}
	log.Info("Searching index")
	err := c.Repository.Get(identifiers, &sr)
	if err != nil {
		log.Warn(err)
		return sr
	}
	return sr
}
