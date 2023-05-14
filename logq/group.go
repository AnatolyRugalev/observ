package logq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
	"github.com/samber/lo"
)

func NewGroup[K comparable](fn genq.GroupFunc[K, Record], sources ...Promise) Group[K] {
	return Group[K]{
		group: genq.NewGroup[K, Record](fn, sources...),
	}
}

type Group[K comparable] struct {
	group *genq.Group[K, Record] `chaingen:"*,unwrap=unwrap,wrap(*)=wrap|slice|sliceMap"`
}

func (Group[K]) slice(slice genq.Slice[Record]) Records {
	return Records(slice)
}

func (g Group[K]) sliceMap(sliceMap map[K]genq.Slice[Record]) map[K]Records {
	return lo.MapValues(sliceMap, func(slice genq.Slice[Record], _ K) Records {
		return g.slice(slice)
	})
}

func (Group[K]) wrap(group *genq.Group[K, Record]) Group[K] {
	return Group[K]{
		group: group,
	}
}

func (g Group[K]) unwrap() *genq.Group[K, Record] {
	return g.group
}
