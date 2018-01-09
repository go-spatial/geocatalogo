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

// Temporal describes temporal bounds
type Temporal struct {
	Begin *time.Time `json:"begin,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

type link struct {
	Name        string
	Description string
	Protocol    string
	URL         string
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
	Identifier     string      `json:"identifier"`
	Title          string      `json:"title,omitempty"`
	Type           string      `json:"type,omitempty"`
	Created        *time.Time  `json:"created,omitempty"`
	Modified       *time.Time  `json:"modified,omitempty"`
	Abstract       string      `json:"abstract,omitempty"`
	KeywordsSets   []keywords  `json:"keywords,omitempty"`
	Contacts       []contact   `json:"contact,omitempty"`
	Dates          []date      `json:"dates,omitempty"`
	License        string      `json:"license,omitempty"`
	Language       string      `json:"language,omitempty"`
	TemporalExtent *Temporal   `json:"temporal_extent,omitempty"`
	Links          []link      `json:"links,omitempty"`
	Geocatalogo    geocatalogo `json:"_geocatalogo,omitempty"`
}

// Record describes a generic metadata record
type Record struct {
	Type       string     `json:"type"`
	Geometry   geometry   `json:"geometry"`
	Properties properties `json:"properties"`
}
