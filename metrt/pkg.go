package metrt

import "testing"

var pkgOptions []Option

func SetOptions(opts ...Option) {
	pkgOptions = opts
}

func AddOptions(opts ...Option) {
	pkgOptions = append(pkgOptions, opts...)
}

// New creates new MetrT.
func New(t *testing.T, opts ...Option) MetrT {
	return MetrT{
		options: defaultOptions,
		t:       t,
	}.WithOptions(pkgOptions...).WithOptions(opts...)
}

// Start creates MetrT and starts it immediately
func Start(t *testing.T, opts ...Option) MetrT {
	return New(t, opts...).Start()
}

func Scope(t *testing.T, fn func(mt MetrT), opts ...Option) Metrics {
	return New(t, opts...).Scope(fn)
}
