package promtest

import (
	"github.com/stretchr/testify/assert"

	"github.com/AnatolyRugalev/observ/prometheus/promquery"
)

type RequireQuery struct {
	q      promquery.Query
	tester *Tester
}

func (r RequireQuery) Labels(labelsKV ...string) RequireQuery {
	r.q = r.q.Labels(labelsKV...)
	return r
}

func (r RequireQuery) Equal(expected float64, msgAndArgs ...any) bool {
	r.tester.t.Helper()
	actual, err := r.q.One()
	if !assert.NoError(r.tester.t, err) {
		return false
	}
	return assert.Equal(r.tester.t, expected, actual, msgAndArgs...)
}

func (r RequireQuery) Greater(than float64, msgAndArgs ...any) bool {
	r.tester.t.Helper()
	actual, err := r.q.One()
	if !assert.NoError(r.tester.t, err) {
		return false
	}
	return assert.Greater(r.tester.t, actual, than, msgAndArgs...)
}

func (r RequireQuery) GreaterOrEqual(than float64, msgAndArgs ...any) bool {
	r.tester.t.Helper()
	actual, err := r.q.One()
	if !assert.NoError(r.tester.t, err) {
		return false
	}
	return assert.GreaterOrEqual(r.tester.t, actual, than, msgAndArgs...)
}

func (r RequireQuery) Less(than float64, msgAndArgs ...any) bool {
	r.tester.t.Helper()
	actual, err := r.q.One()
	if !assert.NoError(r.tester.t, err) {
		return false
	}
	return assert.Less(r.tester.t, actual, than, msgAndArgs...)
}

func (r RequireQuery) LessOrEqual(than float64, msgAndArgs ...any) bool {
	r.tester.t.Helper()
	actual, err := r.q.One()
	if !assert.NoError(r.tester.t, err) {
		return false
	}
	return assert.LessOrEqual(r.tester.t, actual, than, msgAndArgs...)
}

func (r RequireQuery) Positive(msgAndArgs ...any) bool {
	r.tester.t.Helper()
	actual, err := r.q.One()
	if !assert.NoError(r.tester.t, err) {
		return false
	}
	return assert.Positive(r.tester.t, actual, msgAndArgs...)
}

func (r RequireQuery) Negative(msgAndArgs ...any) bool {
	r.tester.t.Helper()
	actual, err := r.q.One()
	if !assert.NoError(r.tester.t, err) {
		return false
	}
	return assert.Negative(r.tester.t, actual, msgAndArgs...)
}

func (r RequireQuery) Group(groupBy string) RequireGroupQuery {
	return RequireGroupQuery{
		GroupQuery: r.q.Group(groupBy),
		tester:     r.tester,
	}
}

type RequireGroupQuery struct {
	promquery.GroupQuery
	tester *Tester
}

func (r RequireGroupQuery) EqualMap(expected M, msgAndArgs ...any) bool {
	r.tester.t.Helper()
	actual, err := r.GroupQuery.Map()
	if !assert.NoError(r.tester.t, err) {
		return false
	}
	return assert.Equal(r.tester.t, expected, actual, msgAndArgs...)
}
