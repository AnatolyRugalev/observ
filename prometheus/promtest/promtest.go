package promtest

import (
	"testing"
	"time"

	promcli "github.com/prometheus/client_golang/prometheus"

	"github.com/AnatolyRugalev/observ/prometheus/promquery"
)

type L = promquery.L
type M = promquery.M

type Tester struct {
	*promquery.Querier
	t       testing.TB
	options promOptions
}

type Option func(o *promOptions)

type promOptions struct {
	registry    *promcli.Registry
	pollingRate time.Duration
	prefix      string
}

var defaultOptions = promOptions{
	pollingRate: 100 * time.Millisecond,
}

func WithPollingRate(pollingRate time.Duration) Option {
	return func(o *promOptions) {
		o.pollingRate = pollingRate
	}
}

func WithRegistry(registry *promcli.Registry) Option {
	return func(o *promOptions) {
		o.registry = registry
	}
}

func WithPrefix(prefix string) Option {
	return func(o *promOptions) {
		o.prefix = prefix
	}
}

func New(t testing.TB, opts ...Option) *Tester {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}
	if options.registry == nil {
		options.registry = promcli.NewRegistry()
	}
	return &Tester{
		Querier: promquery.New(options.registry),
		t:       t,
		options: options,
	}
}

func (t *Tester) Registry() *promcli.Registry {
	return t.options.registry
}

func (t *Tester) Register(collector promcli.Collector) error {
	return t.options.registry.Register(collector)
}

func (t *Tester) MustRegister(collector ...promcli.Collector) {
	t.options.registry.MustRegister(collector...)
}

func (t *Tester) Unregister(collector promcli.Collector) bool {
	return t.options.registry.Unregister(collector)
}

func (t *Tester) Assert(name string, labelsKV ...string) AssertQuery {
	return AssertQuery{
		q:      t.Querier.Query(t.options.prefix+name, labelsKV...),
		tester: t,
	}
}

func (t *Tester) Require(name string, labelsKV ...string) RequireQuery {
	return RequireQuery{
		q:      t.Querier.Query(t.options.prefix+name, labelsKV...),
		tester: t,
	}
}
