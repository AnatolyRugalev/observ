package metrq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
	"github.com/samber/lo"
)

type Source = genq.Source[Metric]

type Metrics []Metric

func (m Metrics) Resolve() []Metric {
	return m
}

func (m Metrics) First() *Metric {
	if len(m) == 0 {
		return nil
	}
	return &m[0]
}

func (m Metrics) Add(sets ...Source) Metrics {
	return m.Merge(sets...).AggregateFlat(Sum())
}

func (m Metrics) Sub(sets ...Source) Metrics {
	return m.Merge(sets...).AggregateFlat(Sub())
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

func (m Metrics) Floats() []float64 {
	return lo.Map(m, func(m Metric, index int) float64 {
		return m.Value
	})
}

func (m Metrics) Ints() []int64 {
	return lo.Map(m, func(m Metric, index int) int64 {
		return int64(m.Value)
	})
}

func (m Metrics) IntSum() int64 {
	return int64(Sum()(intsToFloats(m.Ints())))
}

func intsToFloats(ints []int64) []float64 {
	return lo.Map(ints, func(item int64, index int) float64 {
		return float64(item)
	})
}

func (m Metrics) FloatSum() float64 {
	return Sum()(m.Floats())
}
