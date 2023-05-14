package logq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
)

// TODO:
// 3. Keep methods order
// 4. Wrap GroupFunc
// 5. Generic assert for Group
// 6. More wrappers and codegen...
// 7. Implement Promise

var _ Promise = Filter{}

// chaingen:"ext(Fn):*,wrap(*)=WithFn"
type Filter struct {
	filter *genq.Filter[Record] `chaingen:"*,wrap(*)=wrap|slice|group"`
}

func (f Filter) Resolve() genq.Slice[Record] {
	return f.filter.Resolve()
}

func NewFilter(fn FilterFunc, source Promise) Filter {
	return Filter{
		filter: genq.NewFilter(fn.unwrap(), source),
	}
}

func (Filter) wrap(filter *genq.Filter[Record]) Filter {
	return Filter{
		filter: filter,
	}
}

func (Filter) slice(slice genq.Slice[Record]) Records {
	return Records(slice)
}

func (f Filter) group(group *genq.Group[string, Record]) Group[string] {
	return Group[string]{
		group: group,
	}
}

func (f Filter) Fn() FilterFunc {
	return FilterFunc(f.filter.Fn())
}

func (f Filter) WithFn(fn FilterFunc) Filter {
	return f.wrap(f.filter.WithFn(fn.unwrap()))
}
