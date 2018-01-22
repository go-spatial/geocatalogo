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

// Package main - simple Wrapper
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"flag"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/go-spatial/geocatalogo"
	"github.com/go-spatial/geocatalogo/config"
	"github.com/go-spatial/geocatalogo/metadata/parsers"
	"github.com/go-spatial/geocatalogo/repository"
)

func main() {
	var router *mux.Router
	var plural = ""
	var bbox []float64
	var timeVal []time.Time
	var fileCount = 0
	var fileCounter = 1
	fileList := []string{}

	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <command> [<args>]\n", os.Args[0])
		fmt.Println("Commands: ")
		fmt.Println(" createindex: add a metadata record to the index")
		fmt.Println(" index: add a metadata record to the index")
		fmt.Println(" search: search the index")
		fmt.Println(" get: get metadata record by id")
		fmt.Println(" serve: run web server")
		fmt.Println(" version: geocatalogo version")
		return
	}

	createIndexCommand := flag.NewFlagSet("createindex", flag.ExitOnError)

	indexCommand := flag.NewFlagSet("index", flag.ExitOnError)
	fileFlag := indexCommand.String("file", "", "Path to metadata file")
	dirFlag := indexCommand.String("dir", "", "Path to directory of metadata files")

	searchCommand := flag.NewFlagSet("search", flag.ExitOnError)
	termFlag := searchCommand.String("term", "", "Search term(s)")
	bboxFlag := searchCommand.String("bbox", "", "Bounding box (minx,miny,maxx,maxy)")
	timeFlag := searchCommand.String("time", "", "Time (t1[,t2]), RFC3339 format")
	fromFlag := searchCommand.Int("from", 0, "Start position / offset (default=0)")
	sizeFlag := searchCommand.Int("size", 10, "Number of results to return (default=10)")

	getCommand := flag.NewFlagSet("get", flag.ExitOnError)
	idFlag := getCommand.String("id", "", "list of identifiers (comma-separated)")

	serveCommand := flag.NewFlagSet("serve", flag.ExitOnError)
	portFlag := serveCommand.Int("port", 8000, "port")
	apiFlag := serveCommand.String("api", "default", "API to serve (default, stac)")

	versionCommand := flag.NewFlagSet("version", flag.ExitOnError)

	switch os.Args[1] {
	case "createindex":
		createIndexCommand.Parse(os.Args[2:])
	case "index":
		indexCommand.Parse(os.Args[2:])
	case "search":
		searchCommand.Parse(os.Args[2:])
	case "get":
		getCommand.Parse(os.Args[2:])
	case "serve":
		serveCommand.Parse(os.Args[2:])
	case "version":
		versionCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(10001)
	}

	if versionCommand.Parsed() {
		osinfo := runtime.GOOS + "/" + runtime.GOARCH

		fmt.Println("geocatalogo version " + geocatalogo.VERSION + " " + osinfo)
		return
	}

	if createIndexCommand.Parsed() {
		testLog := logrus.New()

		testConfig := config.LoadFromEnv()

		err := repository.New(testConfig, testLog)

		if err != nil {
			fmt.Println("Repository not created")
		} else {
			fmt.Println("Repository created")
		}
		return
	}

	cat, err := geocatalogo.NewFromEnv()

	if err != nil {
		fmt.Println(err)
		os.Exit(10002)
	}

	if indexCommand.Parsed() {
		if *fileFlag == "" && *dirFlag == "" {
			fmt.Println("Please supply path to metadata file(s) via -file or -dir")
			os.Exit(10003)
		}
		if *fileFlag != "" && *dirFlag != "" {
			fmt.Println("Only one of -file or -dir is allowed")
			os.Exit(10004)
		}
		if *fileFlag != "" {
			fileList = append(fileList, *fileFlag)
		} else if *dirFlag != "" {
			filepath.Walk(*dirFlag, func(path string, f os.FileInfo, err error) error {
				if !f.IsDir() {
					fileList = append(fileList, path)
				}
				return nil
			})
		}

		fileCount = len(fileList)

		if fileCount != 1 {
			plural = "s"
		}

		fmt.Printf("Indexing %d file%s\n", len(fileList), plural)

		for _, file := range fileList {
			start := time.Now()
			parseStart := time.Now()
			fmt.Printf("Indexing file %d of %d: %q\n", fileCounter, fileCount, file)
			source, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Could not read file: %s\n", err)
				continue
			}
			metadataRecord, err := parsers.ParseCSWRecord(source)
			if err != nil {
				fmt.Printf("Could not parse metadata: %s\n", err)
				continue
			}
			parseElapsed := time.Since(parseStart)
			indexStart := time.Now()
			result := cat.Index(metadataRecord)
			if !result {
				fmt.Println("Error Indexing")
			}
			indexElapsed := time.Since(indexStart)
			elapsed := time.Since(start)
			fmt.Printf("Function took %s (parse: %s, index: %s)\n", elapsed, parseElapsed, indexElapsed)
			fileCounter++
		}
	} else if searchCommand.Parsed() {
		if *bboxFlag != "" {
			bboxTokens := strings.Split(*bboxFlag, ",")
			if len(bboxTokens) != 4 {
				fmt.Println("bbox format error (should be minx,miny,maxx,maxy)")
				os.Exit(10006)
			}
			for _, b := range bboxTokens {
				b_, _ := strconv.ParseFloat(b, 64)
				bbox = append(bbox, b_)
			}
		}
		if *timeFlag != "" {
			for _, t := range strings.Split(*timeFlag, ",") {
				timestep, err := time.Parse(time.RFC3339, t)
				if err != nil {
					fmt.Println("time format error (should be ISO 8601/RFC3339)")
					os.Exit(10007)
				}
				timeVal = append(timeVal, timestep)
			}
		}
		results := cat.Search(*termFlag, bbox, timeVal, *fromFlag, *sizeFlag)
		fmt.Printf("Found %d records\n", results.Matches)
		for _, result := range results.Records {
			fmt.Printf("    %s - %s\n", result.Properties.Identifier, result.Properties.Title)
		}
	} else if serveCommand.Parsed() {
		fmt.Printf("Serving on port %d\n", *portFlag)
		if *apiFlag == "stac" {
			router = geocatalogo.STACRouter(cat)
		} else { // csw3-opensearch is the default
			router = geocatalogo.CSW3OpenSearchRouter(cat)
		}
		if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), router); err != nil {
			fmt.Println(err)
			os.Exit(10008)
		}
	} else if getCommand.Parsed() {
		if *idFlag == "" {
			fmt.Println("Please provide identifier")
			os.Exit(10009)
		}
		recordids := strings.Split(*idFlag, ",")
		results := cat.Get(recordids)
		for _, result := range results.Records {
			b, _ := json.MarshalIndent(result, "", "    ")
			fmt.Printf("%s\n", b)
		}
	}
	return
}
