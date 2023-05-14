package metrq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
	"github.com/samber/lo"
	"github.com/AnatolyRugalev/observ/stats"
)

type Source = genq.Promise[Metric]

type Metrics []Metric

func (m Metrics) Resolve() []Metric {
	return m
}

func (m Metrics) Add(sets ...Source) Metrics {
	return m.Merge(sets...).Stat(stats.Sum)
}

func (m Metrics) Sub(sets ...Source) Metrics {
	return m.Merge(sets...).Stat(stats.Sub)
}

func (m Metrics) Merge(with ...Source) Group[MetricKey] {
	return NewGroup(Metric.Key, append([]Source{m}, with...)...)
}

func (m Metrics) Where(operands ...FilterFunc) Filter {
	return NewFilter(And(operands...), m)
}

func (m Metrics) Group(fn GroupFunc) Group[string] {
	return NewGroup(fn, m)
}

func (m Metrics) IntValues() []int64 {
	return lo.Map(m, func(m Metric, index int) int64 {
		return int64(m.Value)
	})
}

func (m Metrics) IntSum() int64 {
	return int64(stats.Sum(intsToFloats(m.IntValues())))
}

func (m Metrics) FloatValues() []float64 {
	return lo.Map(m, func(m Metric, index int) float64 {
		return m.Value
	})
}

func (m Metrics) FloatSum() float64 {
	return stats.Sum(m.FloatValues())
}

func intsToFloats(ints []int64) []float64 {
	return lo.Map(ints, func(item int64, index int) float64 {
		return float64(item)
	})
}

func (m Metrics) First() Metric {
	return genq.First(m)
}

func (m Metrics) Last() Metric {
	return genq.Last(m)
}
