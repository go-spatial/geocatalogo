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

// Package web - simple HTTP Wrapper
package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-spatial/geocatalogo"
	"github.com/go-spatial/geocatalogo/metadata"
	"github.com/go-spatial/geocatalogo/search"
	"github.com/gorilla/mux"
)

// VERSION provides the supported version of the STAC API specification.
const VERSION string = "0.8.0"

type STACSearch struct {
	Limit       int        `json:"limit,omitempty"`
	Datetime    string     `json:"datetime,omitempty"`
	Collections []string   `json:"collections,omitempty"`
	Bbox        [4]float64 `json:"bbox,omitempty"`
}

type Properties struct {
	start    *time.Time `json:"start,omitempty"`
	end      *time.Time `json:"end,omitempty"`
	provider string     `json:"provider,omitempty"`
	license  string     `json:"license,omitempty"`
}

type Link struct {
	Rel   string `json:"rel"`
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
	Href  string `json:"href"`
}

type SearchMetadata struct {
	Next     string `json:"next"`
	Returned int    `json:"returned"`
	Limit    int    `json:"limit,omitempty"`
	Matched  int    `json:"matched,omitempty"`
}

//type assets struct {
//	name string `json:"name,omitempty"`
//	href string `json:"href,omitempty"`
//}

type STACItem struct {
	Type        string              `json:"type,omitempty"`
	Id          string              `json:"id,omitempty"`
	StacVersion string              `json:"stac_version"`
	BBox        [4]float64          `json:"bbox,omitempty"`
	Geometry    metadata.Geometry   `json:"geometry,omitempty"`
	Properties  metadata.Properties `json:"properties,omitempty"`
	Links       []Link              `json:"links,omitempty"`
	Assets      map[string]Link     `json:"assets,omitempty"`
}

type STACFeatureCollection struct {
	Type           string         `json:"type"`
	Features       []STACItem     `json:"features"`
	Links          []Link         `json:"links"`
	NumberMatched  int            `json:"numberMatched"`
	NumberReturned int            `json:"numberReturned"`
	SearchMetadata SearchMetadata `json:"search:metadata"`
}

type STACCatalogDefinition struct {
	Version     string `json:"stac_version"`
	Id          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description"`
	Links       []Link `json:"links"`
}

// STACAPIDescription provides the API description
func STACAPIDescription(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	var jsonBytes []byte
	var scd STACCatalogDefinition

	scd.Id = "geocatalogo"
	scd.Version = VERSION
	scd.Title = cat.Config.Metadata.Identification.Title
	scd.Description = cat.Config.Metadata.Identification.Abstract

	var searchLink = Link{}
	searchLink.Rel = "search"
	searchLink.Type = "application/json"
	searchLink.Title = "search"
	searchLink.Href = fmt.Sprintf("%s/stac/search", cat.Config.Server.URL)

	scd.Links = append(scd.Links, searchLink)

	jsonBytes = geocatalogo.Struct2JSON(&scd, false)

	geocatalogo.EmitResponse(cat, w, 200, jsonBytes)
	return
}

// STACOpenAPI generates an OpenAPI document or Swagger representation
func STACOpenAPI(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	f := r.URL.Query().Get("f")
	if f != "" && f == "json" {
		source, _ := ioutil.ReadFile(cat.Config.Server.OpenAPI)
		data := map[string]interface{}{"config": cat.Config}
		content, _ := geocatalogo.RenderTemplate(string(source), data)
		w.Header().Set("Content-Type", cat.Config.Server.MimeType)
		fmt.Fprintf(w, "%s", content)
	} else {
		data := map[string]interface{}{"config": cat.Config}
		content, _ := geocatalogo.RenderTemplate(SwaggerHTML, data)

		cat.Config.Server.MimeType = "text/html"
		geocatalogo.EmitResponse(cat, w, 200, content)
	}
	return
}

// STACCollections provides STAC compliant collection descriptions
func STACCollections(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	return
}

// STACItems provides STAC compliant Items matching filters
func STACItems(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	var jsonBytes []byte
	var value []string
	var filter string
	var bbox []float64
	var timeVal []time.Time
	var limit = 10
	var page int = 1
	var from int
	var ids []string
	var collections []string
	var results search.Results
	var stacFeatureCollection STACFeatureCollection

	kvp := make(map[string][]string)

	if r.Method == "GET" {
		for k, v := range r.URL.Query() {
			kvp[strings.ToLower(k)] = v
		}
	} else if r.Method == "POST" {
		stacSearch := STACSearch{}

		err := json.NewDecoder(r.Body).Decode(&stacSearch)
		if err != nil {
			fmt.Println(err)
			exception := search.Exception{
				Code:        20002,
				Description: "JSON parsing error"}
			jsonBytes = geocatalogo.Struct2JSON(exception, cat.Config.Server.PrettyPrint)
			geocatalogo.EmitResponse(cat, w, 400, jsonBytes)
			return
		}
		if stacSearch.Limit > 0 {
			kvp["limit"] = []string{strconv.Itoa(stacSearch.Limit)}
		}
		if len(stacSearch.Datetime) > 0 {
			kvp["datetime"] = []string{stacSearch.Datetime}
		}
		if stacSearch.Collections != nil {
			kvp["collections"] = stacSearch.Collections
		}
		if stacSearch.Bbox != [4]float64{0, 0, 0, 0} {
			tmp := fmt.Sprintf("%f,%f,%f,%f", stacSearch.Bbox[0], stacSearch.Bbox[1], stacSearch.Bbox[2], stacSearch.Bbox[3])
			kvp["bbox"] = []string{tmp}
		}
	}

	value, _ = kvp["bbox"]
	if len(value) > 0 {
		bboxTokens := strings.Split(value[0], ",")
		fmt.Println(bboxTokens)
		if len(bboxTokens) != 4 {
			exception := search.Exception{
				Code:        20002,
				Description: "bbox format error (should be minx,miny,maxx,maxy)"}
			jsonBytes = geocatalogo.Struct2JSON(exception, cat.Config.Server.PrettyPrint)
			geocatalogo.EmitResponse(cat, w, 400, jsonBytes)
			return
		}
		for _, bt := range bboxTokens {
			bt, _ := strconv.ParseFloat(bt, 64)
			bbox = append(bbox, bt)
		}
	}
	value, _ = kvp["datetime"]
	if len(value) > 0 {
		for _, t := range strings.Split(value[0], "/") {
			timestep, err := time.Parse(time.RFC3339, t)
			if err != nil {
				exception := search.Exception{
					Code:        20002,
					Description: "time format error (should be ISO 8601/RFC3339)"}
				jsonBytes = geocatalogo.Struct2JSON(exception, cat.Config.Server.PrettyPrint)
				geocatalogo.EmitResponse(cat, w, 400, jsonBytes)
				return
			}
			timeVal = append(timeVal, timestep)
		}
	}

	value, _ = kvp["filter"]
	if len(value) > 0 {
		filter = value[0]
	}

	value, _ = kvp["limit"]
	if len(value) > 0 {
		limit, _ = strconv.Atoi(value[0])
	}

	value, _ = kvp["page"]
	if len(value) > 0 {
		page, _ = strconv.Atoi(value[0])
	}

	if page == 1 {
		from = 0
	} else {
		from = page*limit - (limit - 1)
	}

	value, _ = kvp["ids"]
	if len(value) > 0 {
		ids = strings.Split(value[0], ",")
	}

	value, _ = kvp["collections"]
	if len(value) > 0 {
		collections = strings.Split(value[0], ",")
	}
	fmt.Println(collections == nil)

	if len(ids) > 0 {
		results = cat.Get(ids)
	} else {
		results = cat.Search(collections, filter, bbox, timeVal, from, limit)
	}

	stacFeatureCollection = STACFeatureCollection{}

	Results2STACFeatureCollection(cat.Config.Server.Limit, cat.Config.Server.URL, &results, &stacFeatureCollection)

	jsonBytes = geocatalogo.Struct2JSON(stacFeatureCollection, cat.Config.Server.PrettyPrint)

	geocatalogo.EmitResponse(cat, w, 200, jsonBytes)
	return
}

// STACRouter provides STAC API Routing
func STACRouter(cat *geocatalogo.GeoCatalogue) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/stac", func(w http.ResponseWriter, r *http.Request) {
		STACAPIDescription(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		STACOpenAPI(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		STACCollections(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/collections{collectionId}", func(w http.ResponseWriter, r *http.Request) {
		STACCollections(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/stac/search", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET", "POST", "OPTIONS")

	router.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET", "POST")

	router.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET")

	return router
}

func Results2STACFeatureCollection(limit int, url string, r *search.Results, s *STACFeatureCollection) {
	s.Type = "FeatureCollection"
	for _, rec := range r.Records {
		si := STACItem{}
		si.Type = "Feature"
		si.Id = rec.Identifier
		si.StacVersion = "0.8.0"
		si.BBox = rec.Geometry.Bounds()
		si.Geometry = rec.Geometry
		//si.Datetime = rec.Properties.ProductInfo.AcquisitionDate
		//si.Collection = rec.Properties.ProductInfo.Collection
		si.Properties = rec.Properties
		for _, link := range rec.Links {
			sil := Link{Rel: "self", Href: link.URL}
			si.Links = append(si.Links, sil)
		}
		si.Assets = make(map[string]Link)
		for _, asset := range rec.Assets {
			si.Assets[asset.Name] = Link{Type: asset.Type, Href: asset.URL}
		}
		s.Features = append(s.Features, si)
	}
	nextLink := Link{Rel: "next"}
	nextLink.Href = fmt.Sprintf("%s/stac/search?next=%d", url, r.NextRecord)
	s.Links = append(s.Links, nextLink)
	s.NumberMatched = r.Matches
	s.NumberReturned = r.Returned
	s.SearchMetadata.Next = strconv.Itoa(r.NextRecord)
	s.SearchMetadata.Limit = limit
	s.SearchMetadata.Matched = r.Matches
	s.SearchMetadata.Returned = r.Returned
	return
}
