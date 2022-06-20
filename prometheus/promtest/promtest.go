package promtest

import (
	"testing"
	"time"

	promcli "github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"github.com/AnatolyRugalev/observ/prometheus/promquery"
)

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

func (t *Tester) Assert(name string, labelsKV ...string) *AssertQuery {
	return &AssertQuery{
		Query:  t.Querier.Query(t.options.prefix+name, labelsKV...),
		tester: t,
	}
}

type AssertQuery struct {
	promquery.Query
	tester *Tester
}

func (a AssertQuery) Equal(expected float64, msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.Query.One()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.Equal(a.tester.t, expected, actual, msgAndArgs...)
}

type AssertGroupQuery struct {
	promquery.GroupQuery
	tester *Tester
}

func (a AssertQuery) Group(groupBy string) AssertGroupQuery {
	return AssertGroupQuery{
		GroupQuery: a.Query.Group(groupBy),
		tester:     a.tester,
	}
}

type L = promquery.L
type M = promquery.M

func (a AssertGroupQuery) EqualMap(expected M, msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.GroupQuery.Map()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.Equal(a.tester.t, expected, actual, msgAndArgs...)
}
