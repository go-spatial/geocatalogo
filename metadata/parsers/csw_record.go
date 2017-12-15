package parsers

import (
	"bytes"
	"encoding/xml"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/tomkralidis/geocatalogo/metadata"
)

// CSWRecord provides a CSW 2.0.2 Record model
type CSWRecord struct {
	Identifier       string      `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	Type             string      `xml:"http://purl.org/dc/elements/1.1/ type"`
	Title            string      `xml:"http://purl.org/dc/elements/1.1/ title"`
	Modified         string      `xml:"http://purl.org/dc/terms/ modified"`
	Abstract         string      `xml:"http://purl.org/dc/terms/ abstract"`
	Subject          []string    `xml:"http://purl.org/dc/elements/1.1/ subject"`
	Format           string      `xml:"http://purl.org/dc/elements/1.1/ format"`
	Creator          string      `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Language         string      `xml:"http://purl.org/dc/elements/1.1/ language"`
	References       []string    `xml:"http://purl.org/dc/elements/1.1/ references"`
	WGS84BoundingBox boundingBox `xml:"http://www.opengis.net/ows WGS84BoundingBox"`
	BoundingBox      boundingBox `xml:"http://www.opengis.net/ows BoundingBox"`
}

type boundingBox struct {
	Crs         string `xml:"crs,attr"`
	Dimensions  uint   `xml:"dimensions,attr"`
	LowerCorner string `xml:"http://www.opengis.net/ows LowerCorner"`
	UpperCorner string `xml:"http://www.opengis.net/ows UpperCorner"`
}

func (e *boundingBox) Minx() (float64, error) {
	s := strings.Split(e.LowerCorner, " ")
	minx, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return 0, err
	}
	return minx, nil
}

func (e *boundingBox) Miny() (float64, error) {
	s := strings.Split(e.LowerCorner, " ")
	miny, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		return 0, err
	}
	return miny, nil
}

func (e *boundingBox) Maxx() (float64, error) {
	s := strings.Split(e.UpperCorner, " ")
	maxx, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return 0, err
	}
	return maxx, nil
}

func (e *boundingBox) Maxy() (float64, error) {
	s := strings.Split(e.UpperCorner, " ")
	maxy, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		return 0, err
	}
	return maxy, nil
}

// BBox generates a list of minx,miny,maxx,maxy
func (e *boundingBox) BBox() metadata.Spatial {
	bbox := metadata.Spatial{}
	bbox.Minx, _ = e.Minx()
	bbox.Miny, _ = e.Miny()
	bbox.Maxx, _ = e.Maxx()
	bbox.Maxy, _ = e.Maxy()
	return bbox
}

// ParseCSWRecord parses CSWRecord
func ParseCSWRecord(xmlBuffer []byte) (metadata.Record, error) {
	var cswRecord CSWRecord
	var metadataRecord metadata.Record
	reader := bytes.NewReader(xmlBuffer)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	err := decoder.Decode(&cswRecord)

	if err != nil {
		return metadataRecord, err
	}

	metadataRecord = metadata.Record{
		Identifier: cswRecord.Identifier,
		Type:       cswRecord.Type,
		Title:      cswRecord.Title,
		Abstract:   cswRecord.Abstract,
	}
	if (cswRecord.WGS84BoundingBox != boundingBox{}) {
		metadataRecord.Extent.Spatial = cswRecord.WGS84BoundingBox.BBox()
	} else if (cswRecord.BoundingBox != boundingBox{}) {
		metadataRecord.Extent.Spatial = cswRecord.BoundingBox.BBox()
	}

	return metadataRecord, nil
}
