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

// Package geocatalogo - simple HTTP Wrapper
package geocatalogo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-spatial/geocatalogo/search"
	"github.com/gorilla/mux"
)

// CSW3OpenSearchHandler provides a default HTTP API
func CSW3OpenSearchHandler(w http.ResponseWriter, r *http.Request, cat *GeoCatalogue) {
	var q string
	var recordids []string
	var bbox []float64
	var timeVal []time.Time
	var startPosition int
	var maxRecords = 10
	var value []string
	var results search.Results

	kvp := make(map[string][]string)

	for k, v := range r.URL.Query() {
		kvp[strings.ToLower(k)] = v
	}

	value, _ = kvp["startposition"]
	if len(value) > 0 {
		startPosition, _ = strconv.Atoi(value[0])
	}

	value, _ = kvp["maxrecords"]
	if len(value) > 0 {
		maxRecords, _ = strconv.Atoi(value[0])
	}

	value, _ = kvp["q"]
	if len(value) > 0 {
		q = value[0]
	}

	value, _ = kvp["recordids"]
	if len(value) > 0 {
		recordids = strings.Split(value[0], ",")
	}

	if q == "" && len(recordids) < 1 {
		exception := search.Exception{
			Code:        20001,
			Description: "ERROR: one of q or recordids are required"}
		EmitResponseNotOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &exception)
		return
	}

	if q != "" && len(recordids) > 0 {
		exception := search.Exception{
			Code:        20002,
			Description: "ERROR: q and recordids are mutually exclusive"}
		EmitResponseNotOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &exception)
		return
	}

	if q != "" {
		results = cat.Search(q, bbox, timeVal, startPosition, maxRecords)
	}

	if len(recordids) > 0 {
		results = cat.Get(recordids)
	}

	EmitResponseOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &results)

	return
}

// STACAPIDescription provides the API description
func STACAPIDescription(w http.ResponseWriter, r *http.Request, cat *GeoCatalogue) {
	fmt.Fprintf(w, "hi there")
}

// STACItems provides STAC compliant Items matching filters
func STACItems(w http.ResponseWriter, r *http.Request, cat *GeoCatalogue) {
	var value []string
	var filter string
	var bbox []float64
	var timeVal []time.Time
	var limit = 10
	var next int
	var identifiers []string
	var results search.Results
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
			EmitResponseNotOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &exception)
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
				EmitResponseNotOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &exception)
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

	EmitResponseOK(w, cat.Config.Server.MimeType, cat.Config.Server.PrettyPrint, &results)
	return
}

// CSW3OpenSearchRouter provides CSW 3 OpenSearch Routing
func CSW3OpenSearchRouter(cat *GeoCatalogue) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		CSW3OpenSearchHandler(w, r, cat)
	}).Methods("GET")
	return router
}

// STACRouter provides STAC API Routing
func STACRouter(cat *GeoCatalogue) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/stac/api", 301)
	}).Methods("GET")

	router.HandleFunc("/stac", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/stac/api", 301)
	}).Methods("GET")

	router.HandleFunc("/stac/api", func(w http.ResponseWriter, r *http.Request) {
		source, _ := ioutil.ReadFile(filepath.Join(cat.Config.Server.ApiDataBasedir, "stac-api.json"))
		w.Header().Set("Content-Type", cat.Config.Server.MimeType)
		fmt.Fprintf(w, "%s", source)
	}).Methods("GET")

	router.HandleFunc("/stac/items", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET", "POST")

	router.HandleFunc("/stac/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		STACItems(w, r, cat)
	}).Methods("GET")

	return router
}

// EmitResponseOK provides HTTP response for successful requests
func EmitResponseOK(w http.ResponseWriter, contentType string, prettyPrint bool, results *search.Results) {
	var jsonBytes []byte

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(200)
	if prettyPrint == true {
		jsonBytes, _ = json.MarshalIndent(results, "", "    ")
	} else {
		jsonBytes, _ = json.Marshal(results)
	}
	fmt.Fprintf(w, "%s", jsonBytes)
	return
}

// EmitResponseNotOK provides HTTP response for unsuccessful requests
func EmitResponseNotOK(w http.ResponseWriter, contentType string, prettyPrint bool, exception *search.Exception) {
	var jsonBytes []byte

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(400)

	if prettyPrint == true {
		jsonBytes, _ = json.MarshalIndent(exception, "", "    ")
	} else {
		jsonBytes, _ = json.Marshal(exception)
	}
	fmt.Fprintf(w, "%s", jsonBytes)
	return
}
