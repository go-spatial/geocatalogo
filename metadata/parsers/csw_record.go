package parsers

import (
	"bytes"
	"encoding/xml"
	"fmt"
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
	References       []string    `xml:"http://purl.org/dc/terms/ references"`
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
func (e *boundingBox) BBox() [][2]float64 {
	minx, _ := e.Minx()
	miny, _ := e.Miny()
	maxx, _ := e.Maxx()
	maxy, _ := e.Maxy()
	var a = [][2]float64{
		{minx, miny},
		{minx, maxy},
		{maxx, maxy},
		{maxx, miny},
		{minx, miny},
	}
	return a
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

	metadataRecord = metadata.Record{}
	metadataRecord.Type = "Feature"
	metadataRecord.Properties.Identifier = cswRecord.Identifier
	metadataRecord.Properties.Type = cswRecord.Type
	metadataRecord.Properties.Title = cswRecord.Title
	metadataRecord.Properties.Abstract = cswRecord.Abstract
	metadataRecord.Geometry.Type = "Polygon"

	fmt.Println(cswRecord)
	for _, ref := range cswRecord.References {
		metadataRecord.Properties.Links = append(metadataRecord.Properties.Links, metadata.Link{URL: ref})
	}

	if (cswRecord.WGS84BoundingBox != boundingBox{}) {
		metadataRecord.Geometry.Coordinates = cswRecord.WGS84BoundingBox.BBox()
	} else if (cswRecord.BoundingBox != boundingBox{}) {
		metadataRecord.Geometry.Coordinates = cswRecord.BoundingBox.BBox()
	}

	metadataRecord.Properties.Geocatalogo.Schema = "local"
	metadataRecord.Properties.Geocatalogo.Source = "local"

	return metadataRecord, nil
}
