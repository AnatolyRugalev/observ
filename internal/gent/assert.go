package gent

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
	"testing"
	"github.com/AnatolyRugalev/observ/internal/testify/assert"
)

func NewAssert[V any](t *testing.T, source genq.Promise[V]) Assert[V] {
	return Assert[V]{
		t:      t,
		source: source,
	}
}

type Assert[V any] struct {
	t      *testing.T
	source genq.Promise[V]
}

func (a Assert[V]) NotEmpty(msgAndArgs ...any) bool {
	a.t.Helper()
	return assert.NotEmpty(a.t, a.source.Resolve(), msgAndArgs...)
}

func (a Assert[V]) Empty(msgAndArgs ...any) bool {
	a.t.Helper()
	return assert.Empty(a.t, a.source.Resolve(), msgAndArgs...)
}

func (a Assert[V]) One(msgAndArgs ...any) bool {
	a.t.Helper()
	return assert.Len(a.t, a.source.Resolve(), 1, msgAndArgs...)
}

func (a Assert[V]) Count(count int, msgAndArgs ...any) bool {
	a.t.Helper()
	return assert.Len(a.t, a.source.Resolve(), count, msgAndArgs...)
}

func (a Assert[V]) CountLT(count int, msgAndArgs ...any) bool {
	a.t.Helper()
	return assert.Less(a.t, a.source.Resolve().Count(), count, msgAndArgs...)
}

func (a Assert[V]) CountGT(count int, msgAndArgs ...any) bool {
	a.t.Helper()
	return assert.Greater(a.t, len(a.source.Resolve()), count, msgAndArgs...)
}

func (a Assert[V]) CountLTE(count int, msgAndArgs ...any) bool {
	a.t.Helper()
	return assert.LessOrEqual(a.t, len(a.source.Resolve()), count, msgAndArgs...)
}

func (a Assert[V]) CountGTE(count int, msgAndArgs ...any) bool {
	a.t.Helper()
	return assert.GreaterOrEqual(a.t, len(a.source.Resolve()), count, msgAndArgs...)
}

func (a Assert[V]) Matches(fn genq.FilterFunc[V], msgAndArgs ...any) bool {
	a.t.Helper()
	filter := genq.NewFilter[V](fn, a.source)
	return assert.NotEmpty(a.t, filter.Resolve(), msgAndArgs...)
}

func (a Assert[V]) NotMatches(fn genq.FilterFunc[V], msgAndArgs ...any) bool {
	a.t.Helper()
	filter := genq.NewFilter[V](fn, a.source)
	return assert.Empty(a.t, filter.Resolve(), msgAndArgs...)
}
