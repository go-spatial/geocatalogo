///////////////////////////////////////////////////////////////////////////////
//
// The MIT License (MIT)
// Copyright (c) 2018 Tom Kralidis
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
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"

	"github.com/tomkralidis/geocatalogo/config"
	"github.com/tomkralidis/geocatalogo/metadata"
	"github.com/tomkralidis/geocatalogo/search"
)

// Elasticsearch provides an object model for repository.
type Elasticsearch struct {
	Type      string
	URL       string
	Mappings  map[string]string
	Index     elastic.Client
	IndexName string
	TypeName  string
}

func createClient(repoURL string) (*elastic.Client, error) {
	var esURL string

	u, err := url.Parse(repoURL)
	esURL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	client, err := elastic.NewClient(elastic.SetURL(esURL))
	if err != nil {
		return nil, err
	}
	return client, nil

}

// New creates a repository
func New(cfg config.Config, log *logrus.Logger) error {

	indexMapping := `{
			"mappings": {
				"scene": {
					"properties": {
						"geometry": {
							"type": "geo_shape"
						}
					}
				}
			}
		}`

	ctx := context.Background()

	client, err := createClient(cfg.Repository.URL)
	if err != nil {
		panic(err)
	}

	tokens := strings.Split(cfg.Repository.URL, "/")
	indexName := tokens[len(tokens)-2]

	createIndex, err := client.CreateIndex(indexName).Body(indexMapping).Do(ctx)
	if err != nil {
		errorText := fmt.Sprintf("Cannot create repository: %v\n", err)
		log.Errorf(errorText)
		return errors.New(errorText)
	}
	if !createIndex.Acknowledged {
		return errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}

	log.Debug("Creating Repository" + cfg.Repository.URL)
	log.Debug("Type: " + cfg.Repository.Type)
	log.Debug("URL: " + cfg.Repository.URL)

	return nil
}

// Open loads a repository
func Open(cfg config.Config, log *logrus.Logger) (Elasticsearch, error) {
	log.Debug("Loading Repository " + cfg.Repository.URL)
	log.Debug("Type: " + cfg.Repository.Type)
	log.Debug("URL: " + cfg.Repository.URL)
	s := Elasticsearch{
		Type:      cfg.Repository.Type,
		URL:       cfg.Repository.URL,
		Mappings:  cfg.Repository.Mappings,
		IndexName: getIndexName(cfg.Repository.URL),
		TypeName:  getTypeName(cfg.Repository.URL),
	}
	log.Debug("IndexName: " + s.IndexName)
	log.Debug("TypeName: " + s.TypeName)

	client, err := createClient(cfg.Repository.URL)
	if err != nil {
		return s, err
	}

	s.Index = *client

	return s, nil
}

// Insert inserts a record into the repository
func (r *Elasticsearch) Insert(record metadata.Record) error {
	ctx := context.Background()
	record.Properties.Geocatalogo.Inserted = time.Now()
	_, err := r.Index.Index().
		Index(r.IndexName).
		Type(r.TypeName).
		Id(record.Properties.Identifier).
		BodyJson(record).
		Do(ctx)

	if err != nil {
		return err
	}
	return nil
}

// Update updates a record into the repository
func (r *Elasticsearch) Update() bool {
	return true
}

// Delete deletes a record into the repository
func (r *Elasticsearch) Delete() bool {
	return true
}

// Query performs a search against the repository
func (r *Elasticsearch) Query(term string, bbox []float64, timeVal []time.Time, from int, size int, sr *search.Results) error {
	var mr metadata.Record
	//	var query elastic.Query
	ctx := context.Background()

	query := elastic.NewBoolQuery()

	if term == "" {
		query = query.Must(elastic.NewMatchAllQuery())
	} else {
		query = query.Must(elastic.NewQueryStringQuery(term))
	}
	if len(timeVal) > 0 {
		if len(timeVal) == 1 { // exact match
			query = query.Must(elastic.NewTermQuery("properties.product_info.acquisition_date", timeVal[0]))
		} else if len(timeVal) == 2 { // range
			rangeQuery := elastic.NewRangeQuery("properties.product_info.acquisition_date").
				From(timeVal[0]).
				To(timeVal[1])
			query = query.Must(rangeQuery)
		}
	}
	if len(bbox) == 4 {
		// workaround for issuing a RawStringQuery until
		// GeoShape queries are supported (https://github.com/olivere/elastic/pull/276)
		var tpl bytes.Buffer
		vars := map[string]interface{}{
			"bbox":  bbox,
			"field": "geometry",
		}
		rawStringQueryTemplate, _ := template.New("J").Parse(`{   
          "geo_shape": {
            "{{ .field }}": {
              "shape": {
                "type": "envelope",
                "coordinates": [
                  [   
                    {{ index .bbox 0 }}, 
                    {{ index .bbox 1 }}
                  ],  
                  [   
                    {{ index .bbox 2 }}, 
                    {{ index .bbox 3 }}
                  ]   
                ]
              },
              "relation": "within"
            }   
          }   
        }`)
		rawStringQueryTemplate.Execute(&tpl, vars)

		query = query.Must(elastic.NewRawStringQuery(tpl.String()))
	}

	//src, err := query.Source()
	//data, err := json.Marshal(src)
	//fmt.Println(string(data))

	searchResult, err := r.Index.Search().
		Index(r.IndexName).
		Type(r.TypeName).
		From(from).
		Size(size).
		Query(query).Do(ctx)

	if err != nil {
		return err
	}

	sr.ElapsedTime = int(searchResult.TookInMillis)
	sr.Matches = int(searchResult.TotalHits())
	sr.Returned = size
	sr.NextRecord = size + 1

	if sr.Matches < size {
		sr.Returned = sr.Matches
		sr.NextRecord = 0
	}

	for _, item := range searchResult.Each(reflect.TypeOf(mr)) {
		if t, ok := item.(metadata.Record); ok {
			sr.Records = append(sr.Records, t)
		}
	}

	return nil
}

// Get gets specified metadata records from the repository
func (r *Elasticsearch) Get(identifiers []string, sr *search.Results) error {
	var mr metadata.Record

	idsQuery := elastic.NewIdsQuery(r.TypeName).Ids(identifiers...)
	ctx := context.Background()
	searchResult, err := r.Index.Search().
		Index(r.IndexName).
		Type(r.TypeName).
		Query(idsQuery).Do(ctx)

	if err != nil {
		return err
	}

	sr.Matches = int(searchResult.TotalHits())
	sr.Returned = sr.Matches
	sr.NextRecord = 0

	for _, item := range searchResult.Each(reflect.TypeOf(mr)) {
		if t, ok := item.(metadata.Record); ok {
			sr.Records = append(sr.Records, t)
		}
	}

	return nil
}

// getTypeName returns the name of the ES Index
func getIndexName(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-2]
}

// getTypeName returns the name of the ES Type
func getTypeName(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}
