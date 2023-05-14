package genq

type FilterFunc[V any] func(v V) bool

func (fn FilterFunc[V]) And(operands ...FilterFunc[V]) FilterFunc[V] {
	if fn != nil {
		operands = append([]FilterFunc[V]{fn}, operands...)
	}
	if len(operands) == 0 {
		return nil
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

func (fn FilterFunc[V]) Or(operands ...FilterFunc[V]) FilterFunc[V] {
	if fn != nil {
		operands = append([]FilterFunc[V]{fn}, operands...)
	}
	if len(operands) == 0 {
		return nil
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

func (fn FilterFunc[V]) Invert() FilterFunc[V] {
	if fn == nil {
		return False[V]()
	}
	return func(v V) bool {
		return !fn(v)
	}
}

func (fn FilterFunc[V]) Not(operands ...FilterFunc[V]) FilterFunc[V] {
	return fn.And(new(FilterFunc[V]).And(operands...).Invert())
}

func True[V any]() FilterFunc[V] {
	return nil
}

func False[V any]() FilterFunc[V] {
	return func(v V) bool {
		return false
	}
}
