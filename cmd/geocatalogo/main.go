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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"flag"
	"github.com/sirupsen/logrus"
	"github.com/tomkralidis/geocatalogo"
	"github.com/tomkralidis/geocatalogo/config"
	"github.com/tomkralidis/geocatalogo/metadata/parsers"
	"github.com/tomkralidis/geocatalogo/repository"
)

func main() {

	fileList := []string{}
	var plural = ""

	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <command> [<args>]\n", os.Args[0])
		fmt.Println("Commands: ")
		fmt.Println(" createindex: add a metadata record to the index")
		fmt.Println(" index: add a metadata record to the index")
		fmt.Println(" search: search the index")
		fmt.Println(" version: geocatalogo version")
		return
	}

	createIndexCommand := flag.NewFlagSet("createindex", flag.ExitOnError)

	indexCommand := flag.NewFlagSet("index", flag.ExitOnError)
	fileFlag := indexCommand.String("file", "", "Path to metadata file")
	dirFlag := indexCommand.String("dir", "", "Path to directory of metadata files")

	searchCommand := flag.NewFlagSet("search", flag.ExitOnError)
	termFlag := searchCommand.String("term", "", "Search term(s)")

	versionCommand := flag.NewFlagSet("version", flag.ExitOnError)

	switch os.Args[1] {
	case "createindex":
		createIndexCommand.Parse(os.Args[2:])
	case "index":
		indexCommand.Parse(os.Args[2:])
	case "search":
		searchCommand.Parse(os.Args[2:])
	case "version":
		versionCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if versionCommand.Parsed() {
		osinfo := runtime.GOOS + "/" + runtime.GOARCH

		fmt.Println("geocatalogo version " + geocatalogo.VERSION + " " + osinfo)
		return
	}

	if createIndexCommand.Parsed() {
		testLog := logrus.New()

		testConfig := config.LoadFromEnv()

		status := repository.New(testConfig, testLog)

		if !status {
			fmt.Println("Repository not created")
		} else {
			fmt.Println("Repository created")
		}
		return
	}

	mycatalogo := geocatalogo.New()

	if indexCommand.Parsed() {
		if *fileFlag == "" && *dirFlag == "" {
			fmt.Println("Please supply path to metadata file(s) via -file or -dir")
			os.Exit(3)
		}
		if *fileFlag != "" && *dirFlag != "" {
			fmt.Println("Only one of -file or -dir is allowed")
			os.Exit(4)
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
		if len(fileList) != 1 {
			plural = "s"
		}
		fmt.Printf("Indexing %d file%s\n", len(fileList), plural)

		for _, file := range fileList {
			fmt.Printf("Indexing: %q\n", file)
			source, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}
			metadataRecord, err := parsers.ParseCSWRecord(source)
			result := mycatalogo.Index(metadataRecord)
			if result {
				fmt.Println(result)
			}
		}
	} else if searchCommand.Parsed() {
		if *termFlag == "" {
			fmt.Println("Please provide search term")
			os.Exit(5)
		}
		results := mycatalogo.Search(*termFlag)
		fmt.Println(results)
	}
	return
}
