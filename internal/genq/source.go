package genq

type Source[V any] interface {
	Resolve() []V
}
