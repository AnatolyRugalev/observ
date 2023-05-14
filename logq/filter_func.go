package logq

import (
	"reflect"
	"github.com/AnatolyRugalev/observ/internal/genq"
)

// chaingen:"ext(unwrap):*,unwrap=unwrap,wrap(*)=wrap,export(*)=*"
// chaingen:"export(*):*"
type FilterFunc genq.FilterFunc[Record]

func (fn FilterFunc) unwrap() genq.FilterFunc[Record] {
	return genq.FilterFunc[Record](fn)
}

func (FilterFunc) wrap(f genq.FilterFunc[Record]) FilterFunc {
	return FilterFunc(f)
}

func (fn FilterFunc) and(operand FilterFunc) FilterFunc {
	return fn.wrap(fn.unwrap().And(operand.unwrap()))
}

func (fn FilterFunc) Message(message string) FilterFunc {
	return fn.and(func(v Record) bool {
		return v.Message == message
	})
}

func (fn FilterFunc) Attr(key string, value any) FilterFunc {
	return fn.and(func(v Record) bool {
		return reflect.DeepEqual(v.Attributes[key], value)
	})
}
