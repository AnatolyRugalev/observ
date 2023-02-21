package genq

import (
	"sync"
)

type FilterFunc[V any] func(v V) bool

func NewFilter[V any](fn FilterFunc[V], source Source[V]) *Filter[V] {
	return &Filter[V]{
		fn:     fn,
		source: source,
	}
}

type Filter[V any] struct {
	source Source[V]
	fn     FilterFunc[V]

	once     sync.Once
	resolved []V
}

func (f *Filter[V]) Resolve() []V {
	f.once.Do(func() {
		resolved := f.source.Resolve()
		// TODO: how to predict capacity better?
		f.resolved = make([]V, 0, len(resolved))
		for _, v := range resolved {
			if f.fn(v) {
				f.resolved = append(f.resolved, v)
			}
		}
	})
	return f.resolved
}

func (f *Filter[V]) Where(operands ...FilterFunc[V]) *Filter[V] {
	return f.And(operands...)
}

func (f *Filter[V]) And(operands ...FilterFunc[V]) *Filter[V] {
	return &Filter[V]{
		source: f.source,
		fn:     f.fn.And(operands...),
	}
}

func (f *Filter[V]) Or(operands ...FilterFunc[V]) *Filter[V] {
	return &Filter[V]{
		source: f.source,
		fn:     f.fn.Or(operands...),
	}
}

func (f FilterFunc[V]) Or(operands ...FilterFunc[V]) FilterFunc[V] {
	return Or(append([]FilterFunc[V]{f}, operands...)...)
}

func (f FilterFunc[V]) And(operands ...FilterFunc[V]) FilterFunc[V] {
	return And(append([]FilterFunc[V]{f}, operands...)...)
}

func (f FilterFunc[V]) Not() FilterFunc[V] {
	return Not(f)
}

func Or[V any](operands ...FilterFunc[V]) FilterFunc[V] {
	if len(operands) == 0 {
		return True[V]()
	}
	if len(operands) == 1 {
		return operands[0]
	}
	return func(v V) bool {
		for _, o := range operands {
			if o(v) {
				return true
			}
		}
		return false
	}
}

func And[V any](operands ...FilterFunc[V]) FilterFunc[V] {
	if len(operands) == 0 {
		return True[V]()
	}
	if len(operands) == 1 {
		return operands[0]
	}
	return func(v V) bool {
		for _, o := range operands {
			if !o(v) {
				return false
			}
		}
		return true
	}
}

func Not[V any](f FilterFunc[V]) FilterFunc[V] {
	return func(v V) bool {
		return !f(v)
	}
}

func True[V any]() FilterFunc[V] {
	return func(v V) bool {
		return true
	}
}

func False[V any]() FilterFunc[V] {
	return Not[V](True[V]())
}
