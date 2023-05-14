package genq

import (
	"github.com/samber/lo"
	"sync"
)

func NewGroup[K comparable, V any](fn GroupFunc[K, V], sources ...Promise[V]) *Group[K, V] {
	return &Group[K, V]{
		fn:      fn,
		sources: sources,
	}
}

type Group[K comparable, V any] struct {
	sources []Promise[V]
	fn      GroupFunc[K, V]

	once     sync.Once
	resolved Slice[V]
	index    map[K][]int
}

func (g *Group[K, V]) Resolve() Slice[V] {
	g.once.Do(func() {
		max := 0
		resolved := lo.Map(g.sources, func(s Promise[V], index int) Slice[V] {
			return s.Resolve()
		})
		g.resolved = Slice[V]{}.Merge(resolved...)
		g.index = make(map[K][]int, max)
		var zero K
		for _, vv := range resolved {
			for _, v := range vv {
				key := g.fn(v)
				if key == zero {
					return
				}
				g.index[key] = append(g.index[key], len(g.resolved))
				g.resolved = append(g.resolved, v)
			}
		}
	})
	return g.resolved
}

// Flat returns all items in a group as a slice.
func (g *Group[K, V]) Flat() Slice[V] {
	return g.Resolve()
}

// Aggregate performs an aggregation, and returns results as a map
func (g *Group[K, V]) Aggregate(fn AggregationFunc[V]) map[K]V {
	g.Resolve()
	result := make(map[K]V, len(g.index))
	for k, indexes := range g.index {
		item := fn(lo.Map(indexes, func(idx int, _ int) V {
			return g.resolved[idx]
		}))
		result[k] = item
	}
	return result
}

// AggregateFlat performs an aggregation, and returns results as a slice.
func (g *Group[K, V]) AggregateFlat(fn AggregationFunc[V]) Slice[V] {
	g.Resolve()
	result := make(Slice[V], 0, len(g.index))
	for _, indexes := range g.index {
		item := fn(lo.Map(indexes, func(idx int, _ int) V {
			return g.resolved[idx]
		}))
		result = append(result, item)
	}
	return result
}

// Key returns a list of items grouped by a given key value.
func (g *Group[K, V]) Key(key K) Slice[V] {
	g.Resolve()
	result := make([]V, 0, len(g.index[key]))
	for _, i := range g.index[key] {
		result = append(result, g.resolved[i])
	}
	return result
}

// AsMap returns grouped values as a map
func (g *Group[K, V]) AsMap() map[K]Slice[V] {
	g.Resolve()
	result := make(map[K]Slice[V], len(g.index))
	for k := range g.index {
		result[k] = g.Key(k)
	}
	return result
}

func (g *Group[K, V]) Merge(groups ...*Group[K, V]) *Group[K, V] {
	merged := append([]Promise[V]{g})
	for _, g := range groups {
		merged = append(merged, g)
	}
	return NewGroup[K, V](g.fn, merged...)
}

func (g *Group[K, V]) Count() int {
	g.Resolve()
	return len(g.index)
}
