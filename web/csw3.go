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
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-spatial/geocatalogo"
	"github.com/go-spatial/geocatalogo/search"
	"github.com/gorilla/mux"
)

// CSW3OpenSearchHandler provides a default HTTP API
func CSW3OpenSearchHandler(w http.ResponseWriter, r *http.Request, cat *geocatalogo.GeoCatalogue) {
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

// CSW3OpenSearchRouter provides CSW 3 OpenSearch Routing
func CSW3OpenSearchRouter(cat *geocatalogo.GeoCatalogue) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		CSW3OpenSearchHandler(w, r, cat)
	}).Methods("GET")
	return router
}
