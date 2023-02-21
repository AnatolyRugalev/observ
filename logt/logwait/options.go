package logwait

import (
	"context"
	"github.com/AnatolyRugalev/observ/logq"
	"time"
)

type Options struct {
	Context        context.Context
	Duration       time.Duration
	ShouldContinue func(record logq.Record) bool
	Filter         logq.FilterFunc
}

type Option func(o *Options)

func For(duration time.Duration) Option {
	return func(o *Options) {
		o.Duration = duration
		o.ShouldContinue = func(record logq.Record) bool {
			return true
		}
	}
}

func Where(filters ...logq.FilterFunc) Option {
	return func(o *Options) {
		o.Filter = logq.And(filters...)
	}
}

func One() Option {
	return N(1)
}

func N(n int) Option {
	return func(o *Options) {
		o.ShouldContinue = waitCount(n)
	}
}

func Ctx(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func waitCount(n int) func(record logq.Record) bool {
	count := 0
	return func(record logq.Record) bool {
		count++
		if count == n {
			return false
		}
		return true
	}
}

func NewOptions(ctx context.Context, duration time.Duration, opts ...Option) Options {
	o := Options{
		Context:        ctx,
		Duration:       duration,
		ShouldContinue: waitCount(1),
		Filter:         logq.True(),
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
