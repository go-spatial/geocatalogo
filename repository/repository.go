///////////////////////////////////////////////////////////////////////////////
//
// Repository interface for geocatalogo backends
//
///////////////////////////////////////////////////////////////////////////////

package repository

import (
	"time"

	"github.com/go-spatial/geocatalogo/metadata"
	"github.com/go-spatial/geocatalogo/search"
)

// Repository defines the interface that all backend implementations must satisfy
type Repository interface {
	Insert(record metadata.Record) error
	Update() bool
	Delete() bool
	Query(collections []string, term string, bbox []float64, timeVal []time.Time, from int, size int, sr *search.Results) error
	Get(identifiers []string, sr *search.Results) error
}
