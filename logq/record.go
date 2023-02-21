package logq

import (
	"time"
)

const (
	AttributeFile  = "_file"
	AttributeLine  = "_line"
	AttributeLevel = "_level"
)

type Record struct {
	Time       time.Time
	Message    string
	Attributes map[string]any
}
