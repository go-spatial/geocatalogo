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

// Package main - Landsat 8 on Amazone AWS importer
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tomkralidis/geocatalogo"
	"github.com/tomkralidis/geocatalogo/metadata"
)

func main() {
	var acquisitionDateLayout = "2006-01-02 15:04:05"
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s -file </path/to/scene-list>\n", os.Args[0])
		return
	}

	sceneListFlag := flag.String("file", "", "Path to scene_list csv")
	flag.Parse()

	if *sceneListFlag == "" {
		fmt.Println("Missing file flag")
		os.Exit(1)
	}

	f, err := os.Open(*sceneListFlag)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cat, err := geocatalogo.NewFromEnv()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	for lineno, line := range lines {
		if lineno == 0 { // skip header
			continue
		}
		acquisitionDate, _ := time.Parse(acquisitionDateLayout, line[2])
		cloudCover, _ := strconv.ParseFloat(line[3], 64)
		path, _ := strconv.ParseUint(line[5], 10, 64)
		row, _ := strconv.ParseUint(line[6], 10, 64)
		minLat, _ := strconv.ParseFloat(line[7], 64)
		minLon, _ := strconv.ParseFloat(line[8], 64)
		maxLat, _ := strconv.ParseFloat(line[9], 64)
		maxLon, _ := strconv.ParseFloat(line[10], 64)
		downloadURL := line[11]
		metadataURL := strings.Replace(line[11], "/index.html", "/"+line[0]+"_MTL.json", 1)

		metadataRecord := metadata.Record{}

		metadataRecord.Type = "Feature"

		metadataRecord.Properties.Identifier = line[0]
		metadataRecord.Properties.Title = line[1]
		metadataRecord.Properties.Abstract = "Landsat 8 scene " + line[1]
		metadataRecord.Properties.Links = append(metadataRecord.Properties.Links, metadata.Link{URL: downloadURL})
		metadataRecord.Properties.Links = append(metadataRecord.Properties.Links, metadata.Link{URL: metadataURL})

		pi := &metadata.ProductInfo{
			ProductIdentifier: line[0],
			SceneIdentifier:   line[1],
			AcquisitionDate:   &acquisitionDate,
			CloudCover:        cloudCover,
			ProcessingLevel:   line[4],
			Path:              path,
			Row:               row,
		}
		metadataRecord.Properties.ProductInfo = pi

		metadataRecord.Properties.Geocatalogo.Inserted = time.Now()

		metadataRecord.Geometry.Type = "Polygon"

		var coordinates = [][][2]float64{{
			{minLon, minLat},
			{minLon, maxLat},
			{maxLon, maxLat},
			{maxLon, minLat},
			{minLon, minLat},
		}}

		metadataRecord.Geometry.Coordinates = coordinates

		metadataRecord.Properties.Geocatalogo.Schema = "local"
		metadataRecord.Properties.Geocatalogo.Source = metadataURL

		res, _ := json.Marshal(metadataRecord)
		fmt.Println(string(res))

		result := cat.Index(metadataRecord)
		if !result {
			fmt.Println("ERROR Indexing " + metadataRecord.Properties.Identifier)
		}
	}
	return
}
