package genq

import (
	"github.com/samber/lo"
	"sync"
)

type GroupFunc[K comparable, V any] func(v V) K

func NewGroup[K comparable, V any](fn GroupFunc[K, V], sources ...Source[V]) *Group[K, V] {
	return &Group[K, V]{
		fn:      fn,
		sources: sources,
	}
}

type Group[K comparable, V any] struct {
	sources []Source[V]
	fn      GroupFunc[K, V]

	once     sync.Once
	resolved []V
	index    map[K][]int
}

func (g *Group[K, V]) Resolve() []V {
	g.once.Do(func() {
		max := 0
		resolved := lo.Map(g.sources, func(s Source[V], index int) []V {
			return s.Resolve()
		})
		capacity := lo.SumBy(resolved, func(vv []V) int {
			l := len(vv)
			if l > max {
				max = l
			}
			return l
		})
		var zero K
		g.resolved = make([]V, 0, capacity)
		g.index = make(map[K][]int, max)
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

func (g *Group[K, V]) Flat() []V {
	return g.Resolve()
}

type AggregationFunc[V any] func(values []V) V

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

func (g *Group[K, V]) AsMap() map[K][]V {
	g.Resolve()
	result := make(map[K][]V, len(g.index))
	for k, indexes := range g.index {
		result[k] = lo.Map(indexes, func(idx int, _ int) V {
			return g.resolved[idx]
		})
	}
	return result
}

func (g *Group[K, V]) AggregateFlat(fn AggregationFunc[V]) []V {
	g.Resolve()
	result := make([]V, 0, len(g.index))
	for _, indexes := range g.index {
		item := fn(lo.Map(indexes, func(idx int, _ int) V {
			return g.resolved[idx]
		}))
		result = append(result, item)
	}
	return result
}
