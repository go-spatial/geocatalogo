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
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/sirupsen/logrus"

	"github.com/tomkralidis/geocatalogo/config"
	"github.com/tomkralidis/geocatalogo/metadata"
)

// Repository provides an object model for repository.
type Repository struct {
	Type     string
	URL      string
	Mappings map[string]string
	Index    bleve.Index
}

// Open loads a repository
func Open(cfg config.Config, log *logrus.Logger) Repository {
	log.Info("TOMK")
	log.Debug("Loading Repository" + cfg.Repository.URL)
	log.Debug("Type: " + cfg.Repository.Type)
	log.Debug("URL: " + cfg.Repository.URL)
	s := Repository{
		Type:     cfg.Repository.Type,
		URL:      cfg.Repository.URL,
		Mappings: cfg.Repository.Mappings,
	}

	index, err := bleve.Open(cfg.Repository.URL)

	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		index, err := bleve.New(cfg.Repository.URL, mapping)
		if err != nil {
			panic(err)
		}
		s.Index = index
	} else {
		s.Index = index
	}

	return s
}

func (r *Repository) Insert(record metadata.Record) error {
	err := r.Index.Index(record.Identifier, record)
	return err
}

func (r *Repository) Update() bool {
	return true
}

func (r *Repository) Delete() bool {
	return true
}

func (r *Repository) Query(term string) (sr *bleve.SearchResult, err error) {
	query := bleve.NewMatchQuery(term)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchResult, err := r.Index.Search(searchRequest)
	if err != nil {
		return searchResult, err
	}
	for _, rec := range searchResult.Hits {
		bytes, _ := json.Marshal(metadata.Record{
			Identifier: fmt.Sprintf("%v", rec.Fields["Identifier"]),
			Type:       fmt.Sprintf("%v", rec.Fields["Type"]),
			Title:      fmt.Sprintf("%v", rec.Fields["Title"]),
			Abstract:   fmt.Sprintf("%v", rec.Fields["Abstract"]),
			Language:   fmt.Sprintf("%v", rec.Fields["Language"]),
		})
		fmt.Printf(string(bytes))
	}
	return searchResult, nil
}

func (r *Repository) Get() bool {
	return true
}
