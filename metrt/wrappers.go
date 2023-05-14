package metrt

import (
	"github.com/AnatolyRugalev/observ/metrq"
)

type Metrics struct {
	T      MetrT        `chaingen:"-"`
	filter metrq.Filter `chaingen:"-Resolve,wrap(*)=wrapFilter|wrapGroup"`
}

func (f Metrics) wrapGroup(group metrq.Group[string]) Group[string] {
	return Group[string]{
		T:     f.T,
		group: group,
	}
}

func (f Metrics) wrapFilter(filter metrq.Filter) Metrics {
	return Metrics{
		T:      f.T,
		filter: filter,
	}
}

func (f Metrics) Assert() Assert {
	return Assert{
		filter: f,
	}
}

func (f Metrics) Require() Require {
	return Require{
		assert: f.Assert(),
	}
}

type Group[K comparable] struct {
	T     MetrT          `chaingen:"-"`
	group metrq.Group[K] `chaingen:"-Resolve,wrap(*)=wrapFilter|wrapGroup"`
}

func (g Group[K]) wrapFilter(filter metrq.Filter) Metrics {
	return Metrics{
		T:      g.T,
		filter: filter,
	}
}

func (g Group[K]) Assert() AssertGroup[K] {
	return AssertGroup[K]{
		group: Group[K]{
			T:     g.T,
			group: g.group,
		},
	}
}

func (g Group[K]) Require() RequireGroup[K] {
	return RequireGroup[K]{
		assert: g.Assert(),
	}
}
