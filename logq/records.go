package logq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
)

type Source = genq.Source[Record]

type Records []Record

func (r Records) Resolve() []Record {
	return r
}

func (r Records) First() *Record {
	if len(r) == 0 {
		return nil
	}
	return &r[0]
}

func (r Records) Last() *Record {
	if len(r) == 0 {
		return nil
	}
	return &r[len(r)-1]
}

func (r Records) Where(operands ...FilterFunc) Filter {
	return NewFilter(And(operands...), r)
}

func (r Records) Group(fn GroupFunc) Group[string] {
	return NewGroup(fn, r)
}
