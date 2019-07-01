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

package geocatalogo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Struct2JSON generates a byte representation from a struct
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
