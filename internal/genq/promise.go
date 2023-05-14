package genq

type Promise[V any] interface {
	Resolve() Slice[V]
}
