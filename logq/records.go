package logq

import "github.com/AnatolyRugalev/observ/internal/genq"

type Promise = genq.Promise[Record]

// Records is a slice of records
// chaingen:"ext(unwrap):*,unwrap=unwrap,wrap(*)=wrap|filter|group"
type Records genq.Slice[Record]

func (r Records) unwrap() genq.Slice[Record] {
	return genq.Slice[Record](r)
}

func (Records) wrap(slice genq.Slice[Record]) Records {
	return Records(slice)
}

func (r Records) filter(filter *genq.Filter[Record]) Filter {
	return Filter{
		filter: filter,
	}
}

func (r Records) group(group *genq.Group[string, Record]) Group[string] {
	return Group[string]{
		group: group,
	}
}
