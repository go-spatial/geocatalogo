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

package parsers

import (
	"time"

	"github.com/tomkralidis/geocatalogo/metadata"
)

// properties provides OAM Catalog Result properties
type properties struct {
	Sensor    string `json:"sensor"`
	Thumbnail string `json:"thumbnail"`
	TMS       string `json:"tms"`
	WTMS      string `json:"wtms"`
}

// OAMCatalogResult provides an OAM Catalog Result
type OAMCatalogResult struct {
	Identifier       string     `json:"_id"`
	Uuid             string     `json:"uuid"`
	Version          int        `json:"__v"`
	Title            string     `json:"title"`
	Projection       string     `json:"projection"`
	Gsd              int        `json:"gsd"`
	Filesize         int        `json:"file_size"`
	AcquisitionStart *time.Time `json:"acquisition_start"`
	AcquisitionEnd   *time.Time `json:"acquisition_end"`
	Platform         string     `json:"platform"`
	Provider         string     `json:"provider"`
	Contact          string     `json:"contact"`
	MetaUri          string     `json:"meta_uri"`
	Properties       properties `json:"properties"`
	Bbox             [4]float64 `json:"bbox"`
}

// OAMCatalogResults provides OAM Catalog Results
type OAMCatalogResults struct {
	Result []OAMCatalogResult `json:"results"`
}

// ParseOAMCatalogResult parses CSWRecord
func ParseOAMCatalogResult(result OAMCatalogResult) (metadata.Record, error) {
	metadataRecord := metadata.Record{}

	metadataRecord.Type = "Feature"
	metadataRecord.Properties.Identifier = result.Identifier
	metadataRecord.Properties.Type = "dataset"
	metadataRecord.Properties.Title = result.Title
	metadataRecord.Geometry.Type = "Polygon"

	metadataRecord.Properties.Contacts = append(metadataRecord.Properties.Contacts, metadata.Contact{Value: result.Provider})
	metadataRecord.Properties.Contacts = append(metadataRecord.Properties.Contacts, metadata.Contact{Value: result.Contact})

	mpi := metadata.ProductInfo{}
	mpi.Platform = result.Platform
	mpi.SensorIdentifier = result.Properties.Sensor
	mpi.AcquisitionDate = result.AcquisitionStart
	metadataRecord.Properties.ProductInfo = &mpi

	metadataRecord.Properties.Links = append(metadataRecord.Properties.Links, metadata.Link{URL: result.Uuid})
	metadataRecord.Properties.Links = append(metadataRecord.Properties.Links, metadata.Link{URL: result.Properties.Thumbnail, Protocol: "WWW:LINK"})
	metadataRecord.Properties.Links = append(metadataRecord.Properties.Links, metadata.Link{URL: result.Properties.TMS, Protocol: "OSGeo:TMS"})
	metadataRecord.Properties.Links = append(metadataRecord.Properties.Links, metadata.Link{URL: result.Properties.WTMS, Protocol: "OGC:WMTS"})
	metadataRecord.Properties.Links = append(metadataRecord.Properties.Links, metadata.Link{URL: result.MetaUri, Protocol: "WWW:LINK"})

	var bbox = [][][2]float64{{
		{result.Bbox[0], result.Bbox[1]},
		{result.Bbox[0], result.Bbox[3]},
		{result.Bbox[2], result.Bbox[3]},
		{result.Bbox[2], result.Bbox[1]},
		{result.Bbox[0], result.Bbox[1]},
	}}

	metadataRecord.Geometry.Coordinates = bbox

	metadataRecord.Properties.Geocatalogo.Typename = "oam:meta"
	metadataRecord.Properties.Geocatalogo.Schema = "https://api.openaerialmap.org/meta"
	metadataRecord.Properties.Geocatalogo.Source = "local"

	return metadataRecord, nil
}
