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

package metadata

import (
	"time"
)

type keywords struct {
	Keyword []string
	Type    string
}

type contact struct {
	Type  string
	Value string
}

type date struct {
	Type  string
	Value string
}

// ProductInfo describes product specific metadata
// for example EO data
type ProductInfo struct {
	ProductIdentifier string     `json:"product_id,omitempty"`
	SceneIdentifier   string     `json:"scene_id,omitempty"`
	Path              uint64     `json:"path,omitempty"`
	Row               uint64     `json:"row,omitempty"`
	CloudCover        float64    `json:"cloud_cover,omitempty"`
	AcquisitionDate   *time.Time `json:"acquisition_date,omitempty"`
	ProcessingLevel   string     `json:"processing_level,omitempty"`
	SensorIdentifier  string     `json:"sensor_id,omitempty"`
	SensorId          string     `json:"sensor_id,omitempty"`
}

// Temporal describes temporal bounds
type Temporal struct {
	Begin *time.Time `json:"begin,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

type Link struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	URL         string `json:"url,omitempty"`
}

type geometry struct {
	Type        string       `json:"type"`
	Coordinates [][2]float64 `json:"coordinates"`
}

type geocatalogo struct {
	Inserted time.Time `json:"inserted"`
	Source   string    `json:"source"`
	Schema   string    `json:"schema,omitempty"`
}

type properties struct {
	Identifier     string       `json:"identifier"`
	Title          string       `json:"title,omitempty"`
	Type           string       `json:"type,omitempty"`
	Created        *time.Time   `json:"created,omitempty"`
	Modified       *time.Time   `json:"modified,omitempty"`
	Abstract       string       `json:"abstract,omitempty"`
	KeywordsSets   []keywords   `json:"keywords,omitempty"`
	Contacts       []contact    `json:"contact,omitempty"`
	Dates          []date       `json:"dates,omitempty"`
	License        string       `json:"license,omitempty"`
	Language       string       `json:"language,omitempty"`
	TemporalExtent *Temporal    `json:"temporal_extent,omitempty"`
	Links          []Link       `json:"links,omitempty"`
	ProductInfo    *ProductInfo `json:"product_info,omitempty"`
	Geocatalogo    geocatalogo  `json:"_geocatalogo,omitempty"`
}

// Record describes a generic metadata record
type Record struct {
	Type       string     `json:"type"`
	Geometry   geometry   `json:"geometry"`
	Properties properties `json:"properties"`
}
