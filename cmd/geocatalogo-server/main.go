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

	"github.com/tomkralidis/geocatalogo"
)

var mycatalogo = geocatalogo.New()

func handler(w http.ResponseWriter, r *http.Request) {
	q, ok := r.URL.Query()["q"]
	//startPosition, ok := r.URL.Query()["startPosition"]
	//maxRecords, ok := r.URL.Query()["maxRecords"]"

	if !ok || len(q) < 1 {
		fmt.Fprintf(w, "Url Param 'q' is missing")
		return
	}
	results := mycatalogo.Search(q[0], 10, 10)
	b, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		fmt.Fprintf(w, "ERROR: %s", err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
}
