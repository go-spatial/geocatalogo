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

    "flag"

    "github.com/tomkralidis/geocatalogo"
    "github.com/tomkralidis/geocatalogo/metadata/parsers"
)

func main() {

    if len(os.Args) == 1 {
        fmt.Printf("Usage: %s <command> [<args>]\n", os.Args[0])
        fmt.Println("Commands: ")
        fmt.Println(" index: add a metadata record to the index")
        fmt.Println(" search: search the index")
        return
    }

    mycatalogo := geocatalogo.New()
    fmt.Println("geocatalogo version: " + mycatalogo.Version)

    insertCommand := flag.NewFlagSet("insert", flag.ExitOnError)
    fileFlag := insertCommand.String("file", "", "Path to metadata file")

    searchCommand := flag.NewFlagSet("search", flag.ExitOnError)
    termFlag := searchCommand.String("term", "", "Search term(s)")

    switch os.Args[1] {
        case "insert":
            insertCommand.Parse(os.Args[2:])
        case "search":
            searchCommand.Parse(os.Args[2:])
        default:
            fmt.Printf("%q is not valid command.\n", os.Args[1])
            os.Exit(2)
    }

    if insertCommand.Parsed() {
        if *fileFlag == "" {
            fmt.Println("Please supply path to metadata file")
            return
        }
        fmt.Printf("Indexing: %q\n", *fileFlag)
        source, err := ioutil.ReadFile(*fileFlag)
        if err != nil {
            panic(err)
        }
        metadataRecord, err := parsers.ParseCSWRecord(source)
        result := mycatalogo.Index(metadataRecord)
        if result == true {
            fmt.Println(result)
        }
    } else if searchCommand.Parsed() {
        if *termFlag == "" {
            fmt.Println("Please supply path to metadata file")
            return
        }
        results := mycatalogo.Search(*termFlag)
        fmt.Println(results)
    }
    return
}
