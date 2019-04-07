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
const VERSION string = "0.6.2"

type Properties struct {
	start    *time.Time `json:"start,omitempty"`
	end      *time.Time `json:"end,omitempty"`
	provider string     `json:"provider,omitempty"`
	license  string     `json:"license,omitempty"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
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
	Assets     []assets            `json:"assets,omitempty"`
}

type STACItemCollection struct {
	Type          string     `json:"type"`
	NextPageToken int        `json:"nextPageToken,omitempty"`
	Items         []STACItem `json:"items"`
}

type STACCatalogDefinition struct {
	Version     string `json:"stac_version"`
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Links       []Link `json:"links,omitempty"`
}

// STACAPIDescription provides the API description
func STACAPIDescription(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
	var jsonBytes []byte
	var scd STACCatalogDefinition

	scd.Version = VERSION
	scd.Title = cat.Config.Metadata.Identification.Title

	jsonBytes = Struct2JSON(&scd, false)

	EmitResponse(w, 200, cat.Config.Server.MimeType, jsonBytes)
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
	var stacItemCollection STACItemCollection
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
			jsonBytes = Struct2JSON(exception, cat.Config.Server.PrettyPrint)
			EmitResponse(w, 400, cat.Config.Server.MimeType, jsonBytes)
			return
		}
		for _, bt := range bboxTokens {
			bt, _ := strconv.ParseFloat(bt, 64)
			bbox = append(bbox, bt)
		}
	}
	value, _ = kvp["time"]
	if len(value) > 0 {
		for _, t := range strings.Split(value[0], ",") {
			timestep, err := time.Parse(time.RFC3339, t)
			if err != nil {
				exception := search.Exception{
					Code:        20002,
					Description: "time format error (should be ISO 8601/RFC3339)"}
				jsonBytes = Struct2JSON(exception, cat.Config.Server.PrettyPrint)
				EmitResponse(w, 400, cat.Config.Server.MimeType, jsonBytes)
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

	stacItemCollection = STACItemCollection{}

	Results2STACItemCollection(&results, &stacItemCollection)

	jsonBytes = Struct2JSON(stacItemCollection, cat.Config.Server.PrettyPrint)

	EmitResponse(w, 200, cat.Config.Server.MimeType, jsonBytes)
	return
}

// STACRouter provides STAC API Routing
func STACRouter(cat *geocatalogo.GeoCatalogue) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		STACAPIDescription(w, r, cat)
	}).Methods("GET")

	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		source, _ := ioutil.ReadFile(cat.Config.Server.OpenAPIDef)
		w.Header().Set("Content-Type", cat.Config.Server.MimeType)
		fmt.Fprintf(w, "%s", source)
	}).Methods("GET")

	router.HandleFunc("/stac/search", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET", "POST")

	router.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET")

	return router
}

func Struct2JSON(iface interface{}, prettyPrint bool) []byte {
	var jsonBytes []byte

	if prettyPrint == true {
		jsonBytes, _ = json.MarshalIndent(iface, "", "    ")
	} else {
		jsonBytes, _ = json.Marshal(iface)
	}
	return jsonBytes
}

// EmitResponseOK provides HTTP response for successful requests
func EmitResponse(w http.ResponseWriter, code int, mime string, response []byte) {
	w.Header().Set("Content-Type", mime)
	if code != 200 {
		w.WriteHeader(code)
	}
	fmt.Fprintf(w, "%s", response)
	return
}

func Results2STACItemCollection(r *search.Results, s *STACItemCollection) {
	s.Type = "ItemCollection"
	s.NextPageToken = r.NextRecord
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
		s.Items = append(s.Items, si)
	}
	return
}
