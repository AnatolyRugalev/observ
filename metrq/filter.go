package metrq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
)

type FilterFunc = genq.FilterFunc[Metric]

func NewFilter(fn FilterFunc, source Source) Filter {
	return Filter{
		filter: genq.NewFilter(fn, source),
	}
}

type Filter struct {
	filter *genq.Filter[Metric] `chaingen:"-Resolve,wrap(*)=wrap|wrapMetrics"`
}

func (f Filter) wrap(filter *genq.Filter[Metric]) Filter {
	return Filter{
		filter: filter,
	}
}

func (f Filter) wrapMetrics(metrics []Metric) Metrics {
	return metrics
}

func (f Filter) Resolve() []Metric {
	return f.filter.Resolve()
}

func (f Filter) Group(fn GroupFunc) Group[string] {
	return NewGroup(fn, f)
}

func And(operands ...FilterFunc) FilterFunc {
	return genq.And(operands...)
}

func Or(operands ...FilterFunc) FilterFunc {
	return genq.Or(operands...)
}

func Not(fn FilterFunc) FilterFunc {
	return genq.Not(fn)
}

func True() FilterFunc {
	return genq.True[Metric]()
}

func False() FilterFunc {
	return genq.False[Metric]()
}
