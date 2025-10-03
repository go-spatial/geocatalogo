///////////////////////////////////////////////////////////////////////////////
//
// In-memory repository backend for geocatalogo
// Loads records from JSON file for quick local testing
//
///////////////////////////////////////////////////////////////////////////////

package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-spatial/geocatalogo/config"
	"github.com/go-spatial/geocatalogo/metadata"
	"github.com/go-spatial/geocatalogo/search"
)

// Memory provides an in-memory object model for repository
type Memory struct {
	Type    string
	Records map[string]metadata.Record
	log     *logrus.Logger
}

// NewMemory creates an in-memory repository
func NewMemory(cfg config.Config, log *logrus.Logger) error {
	log.Debug("Creating in-memory repository")
	log.Debug("Type: " + cfg.Repository.Type)
	log.Debug("URL: " + cfg.Repository.URL)

	// For memory backend, URL points to JSON file with records
	return nil
}

// OpenMemory loads an in-memory repository
func OpenMemory(cfg config.Config, log *logrus.Logger) (*Memory, error) {
	log.Debug("Loading in-memory repository from " + cfg.Repository.URL)

	m := &Memory{
		Type:    cfg.Repository.Type,
		Records: make(map[string]metadata.Record),
		log:     log,
	}

	// Load records from JSON file if URL is provided
	if cfg.Repository.URL != "" && cfg.Repository.URL != "memory://" {
		// URL format: file:///path/to/records.json
		filePath := strings.TrimPrefix(cfg.Repository.URL, "file://")

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Warnf("Could not load records from %s: %v", filePath, err)
			return m, nil // Return empty repository, not an error
		}

		var records []metadata.Record
		if err := json.Unmarshal(data, &records); err != nil {
			return nil, fmt.Errorf("failed to parse records JSON: %v", err)
		}

		for _, record := range records {
			m.Records[record.Identifier] = record
		}

		log.Infof("Loaded %d records from %s", len(m.Records), filePath)
	}

	return m, nil
}

// Insert adds a record to the in-memory repository
func (m *Memory) Insert(record metadata.Record) error {
	record.Properties.Geocatalogo.Inserted = time.Now()
	m.Records[record.Identifier] = record
	m.log.Debugf("Inserted record %s", record.Identifier)
	return nil
}

// Update updates a record in the repository
func (m *Memory) Update() bool {
	return true
}

// Delete deletes a record from the repository
func (m *Memory) Delete() bool {
	return true
}

// Get retrieves records by identifier(s)
func (m *Memory) Get(identifiers []string, sr *search.Results) error {
	sr.Records = []metadata.Record{}

	for _, id := range identifiers {
		if record, ok := m.Records[id]; ok {
			sr.Records = append(sr.Records, record)
		}
	}

	sr.Matches = len(sr.Records)
	sr.Returned = sr.Matches
	sr.NextRecord = 0

	return nil
}

// Query performs a search against the in-memory repository
func (m *Memory) Query(collections []string, term string, bbox []float64, timeVal []time.Time, from int, size int, sr *search.Results) error {
	sr.Records = []metadata.Record{}
	matches := []metadata.Record{}

	// Search through all records
	for _, record := range m.Records {
		match := true

		// Collection filter
		if len(collections) > 0 {
			collectionMatch := false
			for _, coll := range collections {
				if record.Properties.Collection == coll {
					collectionMatch = true
					break
				}
			}
			if !collectionMatch {
				match = false
			}
		}

		// Text search (searches in title and abstract)
		if term != "" && match {
			termLower := strings.ToLower(term)
			titleMatch := strings.Contains(strings.ToLower(record.Properties.Title), termLower)
			abstractMatch := strings.Contains(strings.ToLower(record.Properties.Abstract), termLower)
			idMatch := strings.Contains(strings.ToLower(record.Identifier), termLower)

			if !titleMatch && !abstractMatch && !idMatch {
				match = false
			}
		}

		// Bounding box filter (simple overlap check)
		if len(bbox) == 4 && match {
			// bbox format: [minx, miny, maxx, maxy]
			recordBBox := record.BoundingBox

			// Check if bounding boxes overlap
			overlap := !(bbox[2] < recordBBox[0] || // query max_x < record min_x
				bbox[0] > recordBBox[2] || // query min_x > record max_x
				bbox[3] < recordBBox[1] || // query max_y < record min_y
				bbox[1] > recordBBox[3])   // query min_y > record max_y

			if !overlap {
				match = false
			}
		}

		// Time filter
		if len(timeVal) > 0 && match {
			if record.Properties.Datetime != nil {
				// Check if record datetime falls within query time range
				if len(timeVal) == 1 {
					// Exact time match (or close enough - within a day)
					diff := record.Properties.Datetime.Sub(timeVal[0]).Hours()
					if diff < -24 || diff > 24 {
						match = false
					}
				} else if len(timeVal) == 2 {
					// Time range
					if record.Properties.Datetime.Before(timeVal[0]) || record.Properties.Datetime.After(timeVal[1]) {
						match = false
					}
				}
			} else {
				// No datetime in record, doesn't match time query
				match = false
			}
		}

		if match {
			matches = append(matches, record)
		}
	}

	// Pagination
	sr.Matches = len(matches)

	if from >= len(matches) {
		sr.Returned = 0
		sr.NextRecord = 0
		return nil
	}

	end := from + size
	if end > len(matches) {
		end = len(matches)
	}

	sr.Records = matches[from:end]
	sr.Returned = len(sr.Records)

	if end < len(matches) {
		sr.NextRecord = end
	} else {
		sr.NextRecord = 0
	}

	m.log.Debugf("Query found %d matches, returning %d from offset %d", sr.Matches, sr.Returned, from)

	return nil
}

// DeleteAll removes all records (for testing)
func (m *Memory) DeleteAll() error {
	count := len(m.Records)
	m.Records = make(map[string]metadata.Record)
	m.log.Infof("Deleted all %d records", count)
	return nil
}

// Count returns the number of records
func (m *Memory) Count() int {
	return len(m.Records)
}
