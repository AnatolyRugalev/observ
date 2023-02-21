package logt

import (
	"github.com/AnatolyRugalev/observ/logq"
)

type Records struct {
	T      LogT        `chaingen:"-"`
	filter logq.Filter `chaingen:"wrap(*)=wrapFilter,wrap(*)=wrapGroup,-Resolve"`
}

func (f Records) wrapGroup(group logq.Group[string]) Group[string] {
	return Group[string]{
		T:     f.T,
		group: group,
	}
}

func (f Records) wrapFilter(filter logq.Filter) Records {
	return Records{
		T:      f.T,
		filter: filter,
	}
}

func (f Records) Assert() Assert {
	return Assert{
		filter: f,
	}
}

func (f Records) Require() Require {
	return Require{
		assert: f.Assert(),
	}
}

type Group[K comparable] struct {
	T     LogT          `chaingen:"-"`
	group logq.Group[K] `chaingen:"wrap(*)=wrap(*)=wrapFilter,wrap(*)=wrapGroup"`
}

func (g Group[K]) wrapGroup(group logq.Group[string]) Group[string] {
	return Group[string]{
		T:     g.T,
		group: group,
	}
}

func (g Group[K]) wrapFilter(filter logq.Filter) Records {
	return Records{
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
