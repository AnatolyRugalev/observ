package genq

type Slice[V any] []V

func (s Slice[V]) Resolve() Slice[V] {
	return s
}

func (s Slice[V]) First() V {
	if len(s) == 0 {
		var zero V
		return zero
	}
	return s[0]
}

func (s Slice[V]) Last() V {
	if len(s) == 0 {
		var zero V
		return zero
	}
	return s[len(s)-1]
}

func (s Slice[V]) Count() int {
	return len(s)
}

func (s Slice[V]) Where(operands ...FilterFunc[V]) *Filter[V] {
	return NewFilter[V](True[V]().And(operands...), s)
}

func (s Slice[V]) Group(fn GroupFunc[string, V]) *Group[string, V] {
	return NewGroup[string, V](fn, s)
}

func (s Slice[V]) Merge(slices ...Slice[V]) Slice[V] {
	length := len(s)
	for _, slice := range slices {
		length += len(slice)
	}
	merged := make(Slice[V], 0, length)
	merged = append(merged, s...)
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	return merged
}
