package logcollect

import (
	"github.com/AnatolyRugalev/observ/logq"
)

type Sink interface {
	Add(r logq.Record)
}

type Collector interface {
	CaptureLogs(sink Sink) (done func())
}
