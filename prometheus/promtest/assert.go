package promtest

import (
	"github.com/stretchr/testify/assert"

	"github.com/AnatolyRugalev/observ/prometheus/promquery"
)

type AssertQuery struct {
	q      promquery.Query
	tester *Tester
}

func (a AssertQuery) Labels(labelsKV ...string) AssertQuery {
	a.q = a.q.Labels(labelsKV...)
	return a
}

func (a AssertQuery) Equal(expected float64, msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.q.One()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.Equal(a.tester.t, expected, actual, msgAndArgs...)
}

func (a AssertQuery) Greater(than float64, msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.q.One()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.Greater(a.tester.t, actual, than, msgAndArgs...)
}

func (a AssertQuery) GreaterOrEqual(than float64, msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.q.One()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.GreaterOrEqual(a.tester.t, actual, than, msgAndArgs...)
}

func (a AssertQuery) Less(than float64, msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.q.One()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.Less(a.tester.t, actual, than, msgAndArgs...)
}

func (a AssertQuery) LessOrEqual(than float64, msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.q.One()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.LessOrEqual(a.tester.t, actual, than, msgAndArgs...)
}

func (a AssertQuery) Positive(msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.q.One()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.Positive(a.tester.t, actual, msgAndArgs...)
}

func (a AssertQuery) Negative(msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.q.One()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.Negative(a.tester.t, actual, msgAndArgs...)
}

type AssertGroupQuery struct {
	promquery.GroupQuery
	tester *Tester
}

func (a AssertQuery) Group(groupBy string) AssertGroupQuery {
	return AssertGroupQuery{
		GroupQuery: a.q.Group(groupBy),
		tester:     a.tester,
	}
}

func (a AssertGroupQuery) EqualMap(expected M, msgAndArgs ...any) bool {
	a.tester.t.Helper()
	actual, err := a.GroupQuery.Map()
	if !assert.NoError(a.tester.t, err) {
		return false
	}
	return assert.Equal(a.tester.t, expected, actual, msgAndArgs...)
}
