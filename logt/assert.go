package logt

import (
	"github.com/AnatolyRugalev/observ/internal/testify/assert"
)

type Assert struct {
	filter Records
}

func (a Assert) NotEmpty(msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert.NotEmpty(a.filter.T.t, a.filter.filter.Records(), msgAndArgs...)
}

func (a Assert) Empty(msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert.Empty(a.filter.T.t, a.filter.filter.Records(), msgAndArgs...)
}

func (a Assert) Count(count int, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert.Len(a.filter.T.t, a.filter.filter.Records(), count, msgAndArgs...)
}

func (a Assert) Message(expected string, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	record := a.filter.filter.Records().First()
	if record == nil {
		return assert.Fail(a.filter.T.t, "expected at least one record, got none", msgAndArgs...)
	}
	return assert.Equal(a.filter.T.t, expected, record.Message, msgAndArgs...)
}

func (a Assert) Attr(name string, expected any, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	record := a.filter.filter.Records().First()
	if record == nil {
		return assert.Fail(a.filter.T.t, "expected at least one record, got none")
	}
	return assert.Equal(a.filter.T.t, expected, record.Attributes[name], msgAndArgs...)
}

type AssertGroup[K comparable] struct {
	group Group[K]
}

func (a AssertGroup[K]) Count(expected map[K]int, msgAndArgs ...any) bool {
	a.group.T.t.Helper()
	return assert.Equal(a.group.T.t, expected, a.group.group.Count(), msgAndArgs...)
}
