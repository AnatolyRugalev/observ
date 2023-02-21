package metrt

import (
	"github.com/AnatolyRugalev/observ/metrq"
	"testing"
)

// MetrT stands for Metrics Testing.
// In your test code, use `mt` variable name for it of a `MetrT` type (no pointer!).
// All modifier functions of MetrT return a new copy of MetrT, making it thread safe, and easy to reuse.
// If you want to persist changes, reassign the result of modifier function back to the variable.
// Example: mt = mt.WithCollectFilter(mcqold.Prefix("my_awesome_app")).Start()
// You can store and access multiple copies of MetrT, which will contain their own snapshots, so you can mix and match
// different sets of metrics.
type MetrT struct {
	t *testing.T
	options

	snapshot metrq.Metrics
}

// WithT sets new T to MetrT.
func (t MetrT) WithT(tt *testing.T) MetrT {
	t.t = tt
	return t
}

// T returns *testing.T
func (t MetrT) T() *testing.T {
	return t.t
}

// WithoutSnapshot sets empty snapshot.
func (t MetrT) WithoutSnapshot() MetrT {
	return t.WithSnapshot(metrq.Metrics{})
}

// WithSnapshot sets new snapshot to MetricsT
func (t MetrT) WithSnapshot(metrics metrq.Metrics) MetrT {
	t.snapshot = metrics
	return t
}

// Start collects metrics, and puts the result as current snapshot
func (t MetrT) Start() MetrT {
	return t.WithSnapshot(t.collect())
}

// Snapshot returns metrics snapshot, taken earlier using Start.
func (t MetrT) Snapshot() Metrics {
	return Metrics{
		T:      t,
		filter: metrq.NewFilter(metrq.True(), t.snapshot),
	}
}

// Scope runs provided function, and returns metrics changes that happened during its execution.
// This is a shorthand for:
// mt = mt.Start()
// .. do something
// delta := mt.Finish()
func (t MetrT) Scope(fn func(mt MetrT)) Metrics {
	snap := t.Start()
	fn(snap)
	return snap.Finish()
}

// Finish returns difference between metrics snapshot and newly collected metrics.
func (t MetrT) Finish() Metrics {
	before := t.snapshot
	after := t.collect()
	return Metrics{
		T:      t,
		filter: metrq.NewFilter(metrq.True(), after.Sub(before)),
	}
}

// Collect collects all metrics from the collector.
// If collection filter is set, it will ignore metrics that don't match the filter.
func (t MetrT) Collect(filters ...metrq.FilterFunc) Metrics {
	return Metrics{
		T:      t,
		filter: metrq.NewFilter(metrq.And(filters...), t.collect()),
	}
}

func (t MetrT) collect() metrq.Metrics {
	if t.collector == nil {
		panic("collect/metrt: collector is not set")
	}
	return t.collector.CollectMetrics(t.collectFilter)
}
