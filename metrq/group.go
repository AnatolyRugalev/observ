package metrq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
	"github.com/samber/lo"
)

type GroupFunc = genq.GroupFunc[string, Metric]

func NewGroup[K comparable](fn genq.GroupFunc[K, Metric], sources ...Source) Group[K] {
	return Group[K]{
		group: genq.NewGroup[K, Metric](fn, sources...),
	}
}

type Group[K comparable] struct {
	group *genq.Group[K, Metric]
}

func (a Group[K]) FloatSum() map[K]float64 {
	return a.Float(Sum())
}

func (a Group[K]) IntSum() map[K]int64 {
	return a.Int(Sum())
}

func (a Group[K]) Aggregate(fn AggregationFunc) map[K]Metric {
	return a.group.Aggregate(MetricAgg(fn))
}

func (a Group[K]) AggregateFlat(fn AggregationFunc) Metrics {
	return a.group.AggregateFlat(MetricAgg(fn))
}

func (a Group[K]) Float(fn AggregationFunc) map[K]float64 {
	metricMap := a.group.Aggregate(MetricAgg(fn))
	return lo.MapValues(metricMap, func(value Metric, key K) float64 {
		return value.Value
	})
}

func (a Group[K]) Int(fn AggregationFunc) map[K]int64 {
	metricMap := a.group.Aggregate(MetricAgg(fn))
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
