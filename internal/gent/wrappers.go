package gent

import (
	"testing"
	"github.com/AnatolyRugalev/observ/internal/genq"
)

func NewSlice[V any](t *testing.T, slice genq.Slice[V]) Slice[V] {
	return Slice[V]{
		t:     t,
		slice: slice,
	}
}

type Slice[V any] struct {
	t     *testing.T
	slice genq.Slice[V] `chaingen:"*,wrap(*)=wrap|filter|group,unwrap=unwrap"`
}

func (s Slice[V]) wrap(slice genq.Slice[V]) Slice[V] {
	s.slice = slice
	return s
}

func (s Slice[V]) unwrap() genq.Slice[V] {
	return s.slice
}

func (s Slice[V]) group(group *genq.Group[string, V]) Group[string, V] {
	return Group[string, V]{
		t:     s.t,
		group: group,
	}
}

func (s Slice[V]) filter(filter *genq.Filter[V]) Filter[V] {
	return Filter[V]{
		t:      s.t,
		filter: filter,
	}
}

func (s Slice[V]) Assert() Assert[V] {
	return Assert[V]{
		t:      s.t,
		source: s.slice,
	}
}

func NewFilter[V any](t *testing.T, filter *genq.Filter[V]) Filter[V] {
	return Filter[V]{
		t:      t,
		filter: filter,
	}
}

type Filter[V any] struct {
	t      *testing.T
	filter *genq.Filter[V] `chaingen:"*,wrap(*)=wrap|slice,unwrap=unwrap"`
}

func (f Filter[V]) wrap(slice *genq.Filter[V]) Filter[V] {
	f.filter = slice
	return f
}

func (f Filter[V]) unwrap() *genq.Filter[V] {
	return f.filter
}

func (f Filter[V]) slice(slice genq.Slice[V]) Slice[V] {
	return Slice[V]{
		t:     f.t,
		slice: slice,
	}
}

func (f Filter[V]) Assert() Assert[V] {
	return Assert[V]{
		t:      f.t,
		source: f.filter,
	}
}

func NewGroup[K comparable, V any](t *testing.T, group *genq.Group[K, V]) Group[K, V] {
	return Group[K, V]{
		t:     t,
		group: group,
	}
}

type Group[K comparable, V any] struct {
	t     *testing.T
	group *genq.Group[K, V] `chaingen:"*,wrap(*)=wrap|slice,unwrap=unwrap"`
}

func (g Group[K, V]) wrap(group *genq.Group[K, V]) Group[K, V] {
	g.group = group
	return g
}

func (g Group[K, V]) slice(slice genq.Slice[V]) Slice[V] {
	return Slice[V]{
		t:     g.t,
		slice: slice,
	}
}

func (g Group[K, V]) unwrap() *genq.Group[K, V] {
	return g.group
}
