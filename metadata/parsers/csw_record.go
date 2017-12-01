package parsers

import (
  "bytes"
  "encoding/xml"
  "strconv"
  "strings"

  "golang.org/x/net/html/charset"

  "github.com/tomkralidis/geocatalogo/metadata"
)

type CSWRecord struct {
    Identifier string `xml:"http://purl.org/dc/elements/1.1/ identifier"`
    Type string `xml:"http://purl.org/dc/elements/1.1/ type"`
    Title string `xml:"http://purl.org/dc/elements/1.1/ title"`
    Modified string `xml:"http://purl.org/dc/terms/ modified"`
    Abstract string `xml:"http://purl.org/dc/terms/ abstract"`
    Subject []string `xml:"http://purl.org/dc/elements/1.1/ subject"`
    Format string `xml:"http://purl.org/dc/elements/1.1/ format"`
    Creator string `xml:"http://purl.org/dc/elements/1.1/ creator"`
    Language string `xml:"http://purl.org/dc/elements/1.1/ language"`
    References []string `xml:"http://purl.org/dc/elements/1.1/ references"`
    WGS84BoundingBox WGS84BoundingBox `xml:http://www.opengis.net/ows WGS84BoundingBox"`
}

type WGS84BoundingBox struct {
    LowerCorner string `xml:"http://www.opengis.net/ows LowerCorner"`
    UpperCorner string `xml:"http://www.opengis.net/ows UpperCorner"`
}

func (e *WGS84BoundingBox) Minx() (float64, error) {
    s := strings.Split(e.LowerCorner, " ")
    minx, err := strconv.ParseFloat(s[0], 64)
    if err != nil {
        return -9999, err
    }
    return minx, nil
}

func (e *WGS84BoundingBox) Miny() (float64, error) {
    s := strings.Split(e.LowerCorner, " ")
    miny, err := strconv.ParseFloat(s[1], 64)
    if err != nil {
        return -9999, err
    }
    return miny, nil
}

func (e *WGS84BoundingBox) Maxx() (float64, error) {
    s := strings.Split(e.UpperCorner, " ")
    maxx, err := strconv.ParseFloat(s[0], 64)
    if err != nil {
        return -9999, err
    }
    return maxx, nil
}

func (e *WGS84BoundingBox) Maxy() (float64, error) {
    s := strings.Split(e.UpperCorner, " ")
    maxy, err := strconv.ParseFloat(s[1], 64)
    if err != nil {
        return -9999, err
    }
    return maxy, nil
}

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
        Type: cswRecord.Type,
        Title: cswRecord.Title,
        Abstract: cswRecord.Abstract,
    }

    return metadataRecord, nil
}
