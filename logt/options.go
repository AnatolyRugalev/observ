package logt

import (
	"context"
	"github.com/AnatolyRugalev/observ/logcollect"
	"github.com/AnatolyRugalev/observ/logq"
	"time"
)

var defaultOptions = options{
	collectFilter: logq.True(),

	waitTimeout: 5 * time.Second,
	waitContext: context.Background(),
}

type options struct {
	collector     logcollect.Collector
	collectFilter logq.FilterFunc

	waitTimeout time.Duration
	waitContext context.Context
}

type Option func(o *options)

// WithCollector sets a collector for LogT.
// Collector is required for LogT to function.
func WithCollector(collector logcollect.Collector) Option {
	return func(o *options) {
		o.collector = collector
	}
}

// WithCollectFilter sets a collector filter which affects which records are visible to LogT
func WithCollectFilter(filters ...logq.FilterFunc) Option {
	return func(o *options) {
		o.collectFilter = logq.And(filters...)
	}
}

// WithCollectFilter sets a collector filter which affects which records are visible to LogT
func (t LogT) WithCollectFilter(filters ...logq.FilterFunc) LogT {
	return t.WithOptions(WithCollectFilter(filters...))
}

// WithWaitTimeout sets a default timeout for LogT.Wait
func WithWaitTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.waitTimeout = timeout
	}
}

// WithWaitTimeout sets a default timeout for LogT.Wait
func (t LogT) WithWaitTimeout(timeout time.Duration) LogT {
	return t.WithOptions(WithWaitTimeout(timeout))
}

// WithWaitContext sets a default context for LogT.Wait
func WithWaitContext(ctx context.Context) Option {
	return func(o *options) {
		o.waitContext = ctx
	}
}

// WithWaitContext sets a default context for LogT.Wait
func (t LogT) WithWaitContext(ctx context.Context) LogT {
	return t.WithOptions(WithWaitContext(ctx))
}

// WithOptions sets new options to LogT
func (t LogT) WithOptions(opts ...Option) LogT {
	for _, opt := range opts {
		opt(&t.options)
	}
	return t
}
