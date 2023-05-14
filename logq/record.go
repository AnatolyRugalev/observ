package logq

import (
	"time"
)

// Nil represents an absence of a record.
var Nil = Record{}

// Record represents a single log record.
// All logger-specific attributes should be put into Attributes map
type Record struct {
	Time       time.Time
	Message    string
	Attributes Attributes
}

// IsNil returns true if the record is Nil
func (r Record) IsNil() bool {
	return r.Time.Equal(Nil.Time)
}
