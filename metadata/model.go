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

type Keywords struct {
	Keyword []string
	Type    string
}

type Contact struct {
	Type  string
	Value string
}

type Date struct {
	Type  string
	Value string
}

type Extent struct {
	Spatial  [4]float64 // minx, miny, maxx, maxy
	Temporal [2]time.Time
}

type Link struct {
	Name        string
	Description string
	Protocol    string
	URL         string
}

type Record struct {
	Identifier   string
	Title        string
	Type         string
	DateInserted time.Time
	DateModified time.Time
	Schema       string
	Abstract     string
	KeywordsSets []Keywords
	Contacts     []Contact
	Dates        []Date
	License      string
	Language     string
	Extents      []Extent
	Links        []Link
}
