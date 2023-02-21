package metrq

import (
	"bytes"
	"github.com/cespare/xxhash/v2"
)

// AttributesKey can be used as a map key.
type AttributesKey struct {
	hash uint64
}

type Attributes struct {
	// attrs are sorted at all times
	attrs []Attribute
	index map[string]int
	key   AttributesKey
}

func (a Attributes) Get(name string) string {
	idx, ok := a.index[name]
	if !ok {
		return ""
	}
	return a.attrs[idx].Value
}

// NewAttributes creates Attributes object.
// The order of attrs should be deterministic (e.g. sorted alphabetically by key).
func NewAttributes(attrs ...Attribute) Attributes {
	a := Attributes{
		attrs: attrs,
		index: make(map[string]int, len(attrs)),
		key:   NewAttributesKey(attrs...),
	}
	for i, attr := range attrs {
		a.index[attr.Key] = i
	}
	return a
}

type Attribute struct {
	Key   string
	Value string
}

func (a Attributes) Key() AttributesKey {
	return a.key
}

func NewAttributesKey(attrs ...Attribute) AttributesKey {
	if len(attrs) == 0 {
		return AttributesKey{}
	}
	b := bytes.Buffer{}
	for _, attr := range attrs {
		b.WriteString(attr.Key)
		b.WriteRune('=')
		b.WriteString(attr.Value)
		b.WriteRune(';')
	}
	return AttributesKey{hash: xxhash.Sum64(b.Bytes())}
}

// Access Patterns:
// 1. Set of metrics with different names
// 2. Set of metrics with different labels
// 3. Metric with one label

// Deltas work for each scope, metric, and the set of labels
// Filter contains a set of metrics that we need
// And / Or can be applied per-attribute, per-name, per-scope etc.
// Then, we CollectMetrics metrics, apply deltas and perform assertions
// Temporality can be controlled in MetricQ
// Deltas for histograms are per-bucket

// Query levels:
// 1. Metrics (any name, any attributes, any scope)
// 2. Group - based on set of metrics retrieved from (1), has grouping func and aggregating func

// Features:
// 1. Metrics math: add sets, subtract sets

// So:
// We need Metrics which is fully materialized set of metrics, with or without applied deltas
// Metrics can add other Metrics to itself
// Metrics can filter itself
// Metrics can subtract deltas
// Metrics can then me transfromed into MetricsGroup, which is non-meterialized Metrics with classifier and aggregation functions
// MetricsGroup then can be materialized into maps, slices, single number, etc.
// MetricsQ creates Metrics via:
// 1. Collecting sets from collectors
// 2. Subtracting deltas per-collector (since last collection)
// 3. Merging sets of all collectors
// 4. Returning one Metrics as a result
