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
	"net/http"

	"github.com/go-spatial/geocatalogo/search"
)

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
