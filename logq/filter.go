package logq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
)

type FilterFunc = genq.FilterFunc[Record]

func NewFilter(fn FilterFunc, source Source) Filter {
	return Filter{
		filter: genq.NewFilter(fn, source),
	}
}

type Filter struct {
	filter *genq.Filter[Record] `chaingen:"-Resolve,wrap(*)=wrap"`
}

func (f Filter) wrap(filter *genq.Filter[Record]) Filter {
	return Filter{
		filter: filter,
	}
}

func (f Filter) Records() Records {
	return f.Resolve()
}

func (f Filter) Resolve() []Record {
	return f.filter.Resolve()
}

func (f Filter) Group(fn GroupFunc) Group[string] {
	return NewGroup(fn, f)
}

func (f Filter) Message(msg string) Filter {
	return f.wrap(f.filter.And(Message(msg)))
}

func (f Filter) Attr(key string, value any) Filter {
	return f.wrap(f.filter.And(Attr(key, value)))
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
	return genq.True[Record]()
}

func False() FilterFunc {
	return genq.False[Record]()
}
