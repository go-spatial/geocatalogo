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

package repository

import (
	"fmt"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/store/moss"
	"github.com/blevesearch/bleve/index/upsidedown"
	bleveSearch "github.com/blevesearch/bleve/search"
	"github.com/sirupsen/logrus"

	"github.com/tomkralidis/geocatalogo/config"
	"github.com/tomkralidis/geocatalogo/metadata"
	"github.com/tomkralidis/geocatalogo/search"
)

// BleveRepository provides an object model for repository.
type BleveRepository struct {
	Type     string
	URL      string
	Mappings map[string]string
	Index    bleve.Index
}

// New creates a repository
func New(cfg config.Config, log *logrus.Logger) bool {

	kvconfig := map[string]interface{}{
		"mossLowerLevelStoreName": "mossStore",
	}

	log.Debug("Creating Repository" + cfg.Repository.URL)
	log.Debug("Type: " + cfg.Repository.Type)
	log.Debug("URL: " + cfg.Repository.URL)

	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewUsing(cfg.Repository.URL, mapping, upsidedown.Name, moss.Name, kvconfig)

	if err != nil {
		log.Errorf("Cannot create repository: %v\n", err)
		return false
	}
	log.Debug("Persisting moss kv index")
	time.Sleep(30 * time.Second)
	index.Close()

	return true
}

// Open loads a repository
func Open(cfg config.Config, log *logrus.Logger) BleveRepository {
	log.Debug("Loading Repository" + cfg.Repository.URL)
	log.Debug("Type: " + cfg.Repository.Type)
	log.Debug("URL: " + cfg.Repository.URL)
	s := BleveRepository{
		Type:     cfg.Repository.Type,
		URL:      cfg.Repository.URL,
		Mappings: cfg.Repository.Mappings,
	}

	index, err := bleve.Open(cfg.Repository.URL)

	if err != nil {
		panic(err)
	}
	s.Index = index

	return s
}

// Insert inserts a record into the repository
func (r *BleveRepository) Insert(record metadata.Record) error {
	err := r.Index.Index(record.Identifier, record)
	return err
}

// Update updates a record into the repository
func (r *BleveRepository) Update() bool {
	return true
}

// Delete deletes a record into the repository
func (r *BleveRepository) Delete() bool {
	return true
}

// Query performs a search against the repository
func (r *BleveRepository) Query(term string, sr *search.SearchResults, from int, size int) error {
	query := bleve.NewQueryStringQuery(term)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchRequest.From = from
	searchRequest.Size = size

	searchResult, err := r.Index.Search(searchRequest)

	if err != nil {
		return err

	}

	sr.ElapsedTime = int(searchResult.Took / time.Millisecond)
	sr.Matches = int(searchResult.Total)
	sr.Returned = size
	sr.NextRecord = size + 1
	sr.Records = make([]metadata.Record, 0)

	if sr.Matches < size {
		sr.Returned = sr.Matches
		sr.NextRecord = 0
	}

	for _, rec := range searchResult.Hits {
		sr.Records = append(sr.Records, TransformSearchResultToRecord(rec))
	}
	return nil
}

// Get gets specified metadata records from the repository
func (r *BleveRepository) Get(identifiers []string, sr *search.SearchResults) error {
	query := bleve.NewDocIDQuery(identifiers)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchResult, err := r.Index.Search(searchRequest)
	if err != nil {
		return err
	}

	sr.Matches = int(searchResult.Total)
	sr.Returned = sr.Matches
	sr.NextRecord = 0
	sr.Records = make([]metadata.Record, 0)

	for _, rec := range searchResult.Hits {
		sr.Records = append(sr.Records, TransformSearchResultToRecord(rec))
	}
	return nil
}

// TransformSearchResultToRecord provides a helper function to transform a
// bleveSearch.DocumentMatch result to a metadata.Record
func TransformSearchResultToRecord(rec *bleveSearch.DocumentMatch) metadata.Record {
	mr := metadata.Record{
		Identifier: fmt.Sprintf("%v", rec.Fields["Identifier"]),
		Type:       fmt.Sprintf("%v", rec.Fields["Type"]),
		Title:      fmt.Sprintf("%v", rec.Fields["Title"]),
		Abstract:   fmt.Sprintf("%v", rec.Fields["Abstract"]),
		Language:   fmt.Sprintf("%v", rec.Fields["Language"]),
	}
	return mr
}
