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

// Spatial describes bounding geometry
type Spatial struct {
	Minx float64
	Miny float64
	Maxx float64
	Maxy float64
}

// Temporal describes temporal bounds
type Temporal struct {
	Begin time.Time `json:",omitempty"`
	End   time.Time `json:",omitempty"`
}

type extent struct {
	Spatial  Spatial  `json:",omitempty"`
	Temporal Temporal `json:",omitempty"`
}

type link struct {
	Name        string
	Description string
	Protocol    string
	URL         string
}

// Record describes a generic metadata record
type Record struct {
	Identifier   string
	Title        string     `json:",omitempty"`
	Type         string     `json:",omitempty"`
	DateInserted time.Time  `json:",omitempty"`
	DateModified time.Time  `json:",omitempty"`
	Schema       string     `json:",omitempty"`
	Abstract     string     `json:",omitempty"`
	KeywordsSets []keywords `json:",omitempty"`
	Contacts     []contact  `json:",omitempty"`
	Dates        []date     `json:",omitempty"`
	License      string     `json:",omitempty"`
	Language     string     `json:",omitempty"`
	Extent       extent     `json:",omitempty"`
	Links        []link     `json:",omitempty"`
}
