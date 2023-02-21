package metrq

type MetricKind int

const (
	KindCounter MetricKind = iota + 1
	KindGauge
	KindHistogram
)

type MetricKey struct {
	Scope         string
	Name          string
	AttributesKey AttributesKey
}

// Metric is a distinct metric, identified by:
// - scope
// - name
// - unique set of attributes
// Each metric could be either:
// 1. Counter / gauge
// 2. Histogram observation
type Metric struct {
	Scope       string
	Name        string
	Description string
	Attributes  Attributes

	Kind MetricKind

	// Value is a counter or gauge value.
	// When the metric is histogram, value is set to the sum of all histogram measurements.
	// This makes histogram measurements easier to work with
	Value float64

	// Histogram properties

	// Count is a number or histogram samples
	Count uint64
	// BucketCounts store how many samples are placed into each bound
	BucketCounts []uint64
	// Bounds is a set of histogram sample bounds.
	// Bound values are absolute, and this slice should be sorted ASC
	Bounds []float64
}

func (m Metric) Key() MetricKey {
	return MetricKey{
		Scope:         m.Scope,
		Name:          m.Name,
		AttributesKey: m.Attributes.Key(),
	}
}

func (m Metric) Clone() Metric {
	// TODO: do we need a deep copy?
	return m
}
