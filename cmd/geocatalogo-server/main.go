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

// Package main - simple HTTP Wrapper
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"

	"github.com/tomkralidis/geocatalogo"
	"github.com/tomkralidis/geocatalogo/search"
)

var mycatalogo, _ = geocatalogo.NewFromEnv()

func handler(w http.ResponseWriter, r *http.Request) {
	var q string
	var recordids []string
	var startPosition int
	var maxRecords int = 10
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR: one of q or recordids are required")
		return
	}

	if q != "" && len(recordids) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR: q and recordids are mutually exclusive")
		return
	}

	if q != "" {
		results = mycatalogo.Search(q, startPosition, maxRecords)
	}

	if len(recordids) > 0 {
		results = mycatalogo.Get(recordids)
	}

	b, err := json.MarshalIndent(results, "", "    ")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR: %s", err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

func main() {
	var port int = 8000
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
