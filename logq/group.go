package logq

import (
	"fmt"
	"github.com/AnatolyRugalev/observ/internal/genq"
	"github.com/samber/lo"
)

type AggregationFunc = genq.AggregationFunc[Record]

type GroupFunc = genq.GroupFunc[string, Record]

func NewGroup[K comparable](fn genq.GroupFunc[K, Record], sources ...Source) Group[K] {
	return Group[K]{
		group: genq.NewGroup[K, Record](fn, sources...),
	}
}

type Group[K comparable] struct {
	group *genq.Group[K, Record]
}

func (a Group[K]) Aggregate(fn AggregationFunc) map[K]Record {
	return a.group.Aggregate(fn)
}

func (a Group[K]) AggregateFlat(fn AggregationFunc) Records {
	return a.group.AggregateFlat(fn)
}

func ByMessage() GroupFunc {
	return func(v Record) string {
		return v.Message
	}
}

func ByAttr(name string) GroupFunc {
	return func(v Record) string {
		return fmt.Sprintf("%s", v.Attributes[name])
	}
}

func (a Group[K]) Count() map[K]int {
	return lo.MapValues(a.group.AsMap(), func(records []Record, _ K) int {
		return len(records)
	})
}
