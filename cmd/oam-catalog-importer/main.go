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

// Package main - OpenAerialMap Catalog importer
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-spatial/geocatalogo"
	"github.com/go-spatial/geocatalogo/metadata/parsers"
)

func main() {
	//var acquisitionDateLayout = "2006-01-02 15:04:05"
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s -file </path/to/oam.json>\n", os.Args[0])
		return
	}

	fileFlag := flag.String("file", "", "Path to oam.json")
	flag.Parse()

	if *fileFlag == "" {
		fmt.Println("Missing file flag")
		os.Exit(1)
	}

	raw, err := ioutil.ReadFile(*fileFlag)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cat, err := geocatalogo.NewFromEnv()

	var results parsers.OAMCatalogResults

	json.Unmarshal(raw, &results)

	for _, res := range results.Result {
		rec, err := parsers.ParseOAMCatalogResult(res)
		if err != nil {
			fmt.Println(err)
		}
		result := cat.Index(rec)
		if !result {
			fmt.Println("ERROR Indexing " + rec.Identifier)
		}
	}
	return
}
