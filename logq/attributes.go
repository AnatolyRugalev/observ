package logq

import (
	"github.com/spf13/cast"
)

// These constants are used as attribute keys for well-known use-cases and applications.
const (
	// AttributeFile is an attribute key name that is used to store logger-supplied file name of the caller.
	AttributeFile = "_file"
	// AttributeLine is an attribute key name that is used to store logger-supplied line number of the caller.
	AttributeLine = "_line"
	// AttributeLevel is an attribute key name that is used to store logger-supplied level.
	// Some loggers use numeric log levels, some - strings. By default, log collectors should put the value that is aligned
	// with their log level semantics.
	AttributeLevel = "_level"
)

// Attributes represents a set of attributes of a log record.
type Attributes map[string]any

func (a Attributes) Get(name string) any {
	val, ok := a[name]
	if !ok {
		return nil
	}
	return val
}

func (a Attributes) String(name string) string {
	return cast.ToString(a.Get(name))
}

func (a Attributes) Float64(name string) float64 {
	return cast.ToFloat64(a.Get(name))
}

func (a Attributes) Int(name string) int {
	return cast.ToInt(a.Get(name))
}

func (a Attributes) Int64(name string) int64 {
	return cast.ToInt64(a.Get(name))
}
