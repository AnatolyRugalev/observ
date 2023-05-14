package genq

type GroupFunc[K comparable, V any] func(v V) K

type AggregationFunc[V any] func(values Slice[V]) V
