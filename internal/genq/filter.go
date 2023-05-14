package genq

import (
	"sync"
)

func NewFilter[V any](fn FilterFunc[V], source Promise[V]) *Filter[V] {
	return &Filter[V]{
		fn:     fn,
		source: source,
	}
}

// chaingen:"ext(Resolve):*,-Where,ptr(*)"
type Filter[V any] struct {
	source Promise[V]
	fn     FilterFunc[V] `chaingen:"*,ptr(*),wrap(*)=WithFn"`

	once     sync.Once
	resolved Slice[V]
}

func (f *Filter[V]) WithFn(fn FilterFunc[V]) *Filter[V] {
	return &Filter[V]{
		source: f.source,
		fn:     fn,
	}
}

func (f *Filter[V]) Fn() FilterFunc[V] {
	return f.fn
}

func (f *Filter[V]) Resolve() Slice[V] {
	f.once.Do(func() {
		if f.fn == nil {
			f.resolved = f.source.Resolve()
			return
		}
		resolved := f.source.Resolve()
		f.resolved = make([]V, 0, len(resolved))
		for _, v := range resolved {
			if f.fn(v) {
				f.resolved = append(f.resolved, v)
			}
		}
	})
	return f.resolved
}
