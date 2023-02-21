package logt

import (
	"github.com/AnatolyRugalev/observ/logq"
	"testing"
)

var pkgOptions []Option

func SetOptions(opts ...Option) {
	pkgOptions = opts
}

func AddOptions(opts ...Option) {
	pkgOptions = append(pkgOptions, opts...)
}

func New(t *testing.T, opts ...Option) LogT {
	return LogT{
		t:       t,
		sink:    logq.NewSink(nil, logq.True()),
		options: defaultOptions,
	}.WithOptions(pkgOptions...).WithOptions(opts...)
}

func Start(t *testing.T, opts ...Option) LogT {
	t.Helper()
	return New(t, opts...).Start()
}

func Scope(t *testing.T, f func(lgt LogT), opts ...Option) Records {
	t.Helper()
	return New(t, opts...).Scope(f)
}
