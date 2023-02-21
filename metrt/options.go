package metrt

import (
	"github.com/AnatolyRugalev/observ/metrcollect"
	"github.com/AnatolyRugalev/observ/metrq"
)

var defaultOptions = options{
	collectFilter: metrq.True(),
}

type options struct {
	collector     metrcollect.Collector
	collectFilter metrq.FilterFunc
}

type Option func(o *options)

// WithCollector sets a collector for MetrT.
// Collector is required for MetrT to function.
func WithCollector(collector metrcollect.Collector) Option {
	return func(o *options) {
		o.collector = collector
	}
}

// WithCollectFilter sets collector filter, which affects which metrics are visible toe MetrT.
func WithCollectFilter(fn ...metrq.FilterFunc) Option {
	return func(o *options) {
		o.collectFilter = metrq.And(fn...)
	}
}

// WithCollectFilter sets collector filter, which affects which metrics are visible to MetrT.
func (t MetrT) WithCollectFilter(filter ...metrq.FilterFunc) MetrT {
	t.collectFilter = metrq.And(filter...)
	return t
}

// WithOptions sets new options to MetrT.
func (t MetrT) WithOptions(opts ...Option) MetrT {
	for _, opt := range opts {
		opt(&t.options)
	}
	return t
}
