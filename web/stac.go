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
	"fmt"
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

type assets struct {
	name string `json:"name,omitempty"`
	href string `json:"href,omitempty"`
}

type STACItem struct {
	Type       string              `json:"type,omitempty"`
	Id         string              `json:"id,omitempty"`
	BBox       [4]float64          `json:"bbox,omitempty"`
	Geometry   metadata.Geometry   `json:"geometry,omitempty"`
	Properties metadata.Properties `json:"properties,omitempty"`
	Links      []Link              `json:"links,omitempty"`
	Assets     []Link              `json:"assets,omitempty"`
}

type STACFeatureCollection struct {
	Type     string     `json:"type"`
	Features []STACItem `json:"features"`
	Links    []Link     `json:"links"`
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

	geocatalogo.EmitResponse(w, 200, cat.Config.Server.MimeType, jsonBytes)
	return
}

// STACOpenAPI generates an OpenAPI document or Swagger representation
func STACOpenAPI(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	f := r.URL.Query().Get("f")
	if f != "" && f == "json" {
		bytes, _ := GenerateOpenAPIDocument(cat.Config)
		geocatalogo.EmitResponse(w, 200, cat.Config.Server.MimeType, bytes)
	} else {

		data := map[string]interface{}{"config": cat.Config}
		content, _ := geocatalogo.RenderTemplate(SwaggerHTML, data)

		geocatalogo.EmitResponse(w, 200, "text/html", content)
	}
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
	var next int
	var identifiers []string
	var results search.Results
	var stacFeatureCollection STACFeatureCollection
	var tmp string

	kvp := make(map[string][]string)

	for k, v := range r.URL.Query() {
		kvp[strings.ToLower(k)] = v
	}

	vars := mux.Vars(r)
	tmp, ok := vars["id"]

	if ok == true {
		identifiers = append(identifiers, tmp)
	}

	value, _ = kvp["bbox"]
	if len(value) > 0 {
		bboxTokens := strings.Split(value[0], ",")
		if len(bboxTokens) != 4 {
			exception := search.Exception{
				Code:        20002,
				Description: "bbox format error (should be minx,miny,maxx,maxy)"}
			jsonBytes = geocatalogo.Struct2JSON(exception, cat.Config.Server.PrettyPrint)
			geocatalogo.EmitResponse(w, 400, cat.Config.Server.MimeType, jsonBytes)
			return
		}
		for _, bt := range bboxTokens {
			bt, _ := strconv.ParseFloat(bt, 64)
			bbox = append(bbox, bt)
		}
	}
	value, _ = kvp["datetime"]
	if len(value) > 0 {
		for _, t := range strings.Split(value[0], ",") {
			timestep, err := time.Parse(time.RFC3339, t)
			if err != nil {
				exception := search.Exception{
					Code:        20002,
					Description: "time format error (should be ISO 8601/RFC3339)"}
				jsonBytes = geocatalogo.Struct2JSON(exception, cat.Config.Server.PrettyPrint)
				geocatalogo.EmitResponse(w, 400, cat.Config.Server.MimeType, jsonBytes)
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

	value, _ = kvp["next"]
	if len(value) > 0 {
		next, _ = strconv.Atoi(value[0])
	}

	if len(identifiers) > 0 {
		results = cat.Get(identifiers)
	} else {
		results = cat.Search(filter, bbox, timeVal, next, limit)
	}

	stacFeatureCollection = STACFeatureCollection{}

	Results2STACFeatureCollection(&results, &stacFeatureCollection)

	jsonBytes = geocatalogo.Struct2JSON(stacFeatureCollection, cat.Config.Server.PrettyPrint)

	geocatalogo.EmitResponse(w, 200, cat.Config.Server.MimeType, jsonBytes)
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

	router.HandleFunc("/stac/search", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET", "POST")

	router.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET")

	return router
}

func Results2STACFeatureCollection(r *search.Results, s *STACFeatureCollection) {
	s.Type = "FeatureCollection"
	for _, rec := range r.Records {
		si := STACItem{}
		si.Type = "Feature"
		si.Id = rec.Identifier
		si.BBox = rec.Geometry.Bounds()
		si.Geometry = rec.Geometry
		si.Properties = rec.Properties
		for _, link := range rec.Links {
			sil := Link{Rel: "self", Href: link.URL}
			si.Links = append(si.Links, sil)
		}
		for _, asset := range rec.Assets {
			sila := Link{Href: asset.URL}
			si.Assets = append(si.Assets, sila)
		}
		s.Features = append(s.Features, si)
	}
	nextLink := Link{Rel: "next"}
	nextLink.Href = fmt.Sprintf("%s/stac/search?next=%d", "http://URL", r.NextRecord)
	s.Links = append(s.Links, nextLink)
	return
}
