package metrq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
	"github.com/samber/lo"
	"github.com/AnatolyRugalev/observ/stats"
)

// GroupFunc is a grouping function for Metric objects.
type GroupFunc = genq.GroupFunc[string, Metric]

// AggregationFunc is an aggregation func for Metric objects.
type AggregationFunc = genq.AggregationFunc[Metric]

// MetricStat returns a AggregationFunc that applies provided stats.StatFunc to metric values.
func MetricStat(fn stats.StatFunc) AggregationFunc {
	return func(metrics []Metric) Metric {
		metric := metrics[0].Clone()
		metric.Value = fn(lo.Map(metrics, func(m Metric, _ int) float64 {
			return m.Value
		}))
		if metric.Kind == KindHistogram {
			// TODO: sum count, bucket counts
		}
		return metric
	}
}

func NewGroup[K comparable](fn genq.GroupFunc[K, Metric], sources ...Source) Group[K] {
	return Group[K]{
		group: genq.NewGroup[K, Metric](fn, sources...),
	}
}

type Group[K comparable] struct {
	group *genq.Group[K, Metric] `chaingen:"-Resolve,wrap(*)=wrapMetrics|wrapMetricsMap"`
}

func (g Group[K]) wrapMetrics(metrics []Metric) Metrics {
	return metrics
}

func (g Group[K]) wrapMetricsMap(metrics map[K][]Metric) map[K]Metrics {
	return lo.MapValues(metrics, func(value []Metric, _ K) Metrics {
		return value
	})
}

func (g Group[K]) FloatSum() map[K]float64 {
	return g.Float(stats.Sum)
}

func (g Group[K]) IntSum() map[K]int64 {
	return g.Int(stats.Sum)
}

func (g Group[K]) StatMap(fn stats.StatFunc) map[K]Metric {
	return g.group.Aggregate(MetricStat(fn))
}

func (g Group[K]) Stat(fn stats.StatFunc) Metrics {
	return g.group.AggregateFlat(MetricStat(fn))
}

func (g Group[K]) Float(fn stats.StatFunc) map[K]float64 {
	metricMap := g.group.Aggregate(MetricStat(fn))
	return lo.MapValues(metricMap, func(value Metric, key K) float64 {
		return value.Value
	})
}

func (g Group[K]) Int(fn stats.StatFunc) map[K]int64 {
	metricMap := g.group.Aggregate(MetricStat(fn))
	return lo.MapValues(metricMap, func(value Metric, key K) int64 {
		return int64(value.Value)
	})
}

func ByAttr(name string) GroupFunc {
	return func(m Metric) string {
		return m.Attributes.Get(name)
	}
}

func ByName() GroupFunc {
	return func(m Metric) string {
		return m.Name
	}
}
