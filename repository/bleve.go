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
	"errors"
	"fmt"
	"strconv"
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

// Bleve provides an object model for repository.
type Bleve struct {
	Type     string
	URL      string
	Mappings map[string]string
	Index    bleve.Index
}

// New creates a repository
func NewBleve(cfg config.Config, log *logrus.Logger) error {

	kvconfig := map[string]interface{}{
		"mossLowerLevelStoreName": "mossStore",
	}

	log.Debug("Creating Repository" + cfg.Repository.URL)
	log.Debug("Type: " + cfg.Repository.Type)
	log.Debug("URL: " + cfg.Repository.URL)

	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewUsing(cfg.Repository.URL, mapping, upsidedown.Name, moss.Name, kvconfig)

	if err != nil {
		errorText := fmt.Sprintf("Cannot create repository: %v\n", err)
		log.Errorf(errorText)
		return errors.New(errorText)
	}
	log.Debug("Persisting moss kv index")
	time.Sleep(30 * time.Second)
	index.Close()

	return nil
}

// Open loads a repository
func OpenBleve(cfg config.Config, log *logrus.Logger) (Bleve, error) {
	log.Debug("Loading Repository " + cfg.Repository.URL)
	log.Debug("Type: " + cfg.Repository.Type)
	log.Debug("URL: " + cfg.Repository.URL)
	s := Bleve{
		Type:     cfg.Repository.Type,
		URL:      cfg.Repository.URL,
		Mappings: cfg.Repository.Mappings,
	}

	index, err := bleve.Open(cfg.Repository.URL)

	if err != nil {
		return s, err
	}
	s.Index = index

	return s, nil
}

// Insert inserts a record into the repository
func (r *Bleve) Insert(record metadata.Record) error {
	record.Properties.Geocatalogo.Inserted = time.Now()
	return r.Index.Index(record.Properties.Identifier, record)
}

// Update updates a record into the repository
func (r *Bleve) Update() bool {
	return true
}

// Delete deletes a record into the repository
func (r *Bleve) Delete() bool {
	return true
}

// Query performs a search against the repository
func (r *Bleve) Query(term string, sr *search.Results, from int, size int) error {
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
		sr.Records = append(sr.Records, transformResultToRecord(rec))
	}
	return nil
}

// Get gets specified metadata records from the repository
func (r *Bleve) Get(identifiers []string, sr *search.Results) error {
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
		sr.Records = append(sr.Records, transformResultToRecord(rec))
	}
	return nil
}

// transformResultToRecord provides a helper function to transform a
// bleveSearch.DocumentMatch result to a metadata.Record
func transformResultToRecord(rec *bleveSearch.DocumentMatch) metadata.Record {
	var ll [2]float64
	var ul [2]float64
	var ur [2]float64
	var lr [2]float64

	mr := metadata.Record{}
	mr.Type = fmt.Sprintf("%v", rec.Fields["type"])

	if rec.Fields["geometry.coordinates"] != nil {
		mr.Geometry.Type = fmt.Sprintf("%v", rec.Fields["geometry.type"])
		coords := rec.Fields["geometry.coordinates"].([]interface{})
		ll[0] = coords[0].(float64)
		ll[1] = coords[1].(float64)
		ul[0] = coords[2].(float64)
		ul[1] = coords[3].(float64)
		ur[0] = coords[4].(float64)
		ur[1] = coords[5].(float64)
		lr[0] = coords[6].(float64)
		lr[1] = coords[7].(float64)

		mr.Geometry.Coordinates = [][2]float64{ll, ul, ur, lr, ll}
	}

	mr.Properties.Identifier = fmt.Sprintf("%v", rec.Fields["properties.identifier"])
	mr.Properties.Type = fmt.Sprintf("%v", rec.Fields["properties.type"])
	mr.Properties.Title = fmt.Sprintf("%v", rec.Fields["properties.title"])
	mr.Properties.Abstract = fmt.Sprintf("%v", rec.Fields["properties.abstract"])
	mr.Properties.Language = fmt.Sprintf("%v", rec.Fields["properties.language"])

	product_id := fmt.Sprintf("%v", rec.Fields["properties.product_info.product_id"])
	scene_id := fmt.Sprintf("%v", rec.Fields["properties.product_info.scene_id"])
	acquisition_date, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v", rec.Fields["properties.product_info.acquisition_date"]))
	cloud_cover, _ := strconv.ParseFloat(fmt.Sprintf("%v", rec.Fields["properties.product_info.path"]), 64)
	processing_level := fmt.Sprintf("%v", rec.Fields["properties.product_info.processing_level"])
	path, _ := strconv.ParseUint(fmt.Sprintf("%v", rec.Fields["properties.product_info.path"]), 10, 64)
	row, _ := strconv.ParseUint(fmt.Sprintf("%v", rec.Fields["properties.product_info.row"]), 10, 64)

	pi := &metadata.ProductInfo{
		ProductIdentifier: product_id,
		SceneIdentifier:   scene_id,
		AcquisitionDate:   &acquisition_date,
		CloudCover:        cloud_cover,
		ProcessingLevel:   processing_level,
		Path:              path,
		Row:               row,
	}

	mr.Properties.ProductInfo = pi

	mr.Properties.Geocatalogo.Source = fmt.Sprintf("%v", rec.Fields["properties._geocatalogo.source"])
	mr.Properties.Geocatalogo.Schema = fmt.Sprintf("%v", rec.Fields["properties._geocatalogo.schema"])
	mr.Properties.Geocatalogo.Inserted, _ = time.Parse(time.RFC3339, fmt.Sprintf("%v", rec.Fields["properties._geocatalogo.inserted"]))

	return mr
}
