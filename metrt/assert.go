package metrt

import (
	"fmt"
	assert2 "github.com/AnatolyRugalev/observ/internal/testify/assert"
	"github.com/AnatolyRugalev/observ/metrq"
)

type Assert struct {
	filter Metrics
}

type IntOrFloat any

func sumOfType(t IntOrFloat, metrics metrq.Metrics) IntOrFloat {
	switch t.(type) {
	case int:
		return int(metrics.IntSum())
	case int64:
		return metrics.IntSum()
	default:
		return metrics.FloatSum()
	}
}

func valueOfType(t IntOrFloat, value float64) IntOrFloat {
	switch t.(type) {
	case int:
		return int(value)
	case int64:
		return int64(value)
	default:
		return value
	}
}

func (a Assert) getValue(expected IntOrFloat, msgAndArgs ...any) IntOrFloat {
	a.filter.T.t.Helper()
	if len(a.filter.filter.Metrics()) == 0 {
		return assert2.Fail(a.filter.T.t, "expected exactly one metric in Metrics, got 0", msgAndArgs...)
	}
	if l := len(a.filter.filter.Metrics()); l > 1 {
		return assert2.Fail(a.filter.T.t, fmt.Sprintf("expected exactly one metric in Metrics, got %d", l), msgAndArgs...)
	}
	return valueOfType(expected, a.filter.filter.Metrics().First().Value)
}

func (a Assert) Value(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.Equal(a.filter.T.t, expected, a.getValue(expected, msgAndArgs...), msgAndArgs)
}

func (a Assert) ValueGreater(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.Greater(a.filter.T.t, expected, a.getValue(expected, msgAndArgs...), msgAndArgs)
}

func (a Assert) ValueLess(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.Less(a.filter.T.t, expected, a.getValue(expected, msgAndArgs...), msgAndArgs)
}

func (a Assert) ValueLessOrEqual(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.LessOrEqual(a.filter.T.t, expected, a.getValue(expected, msgAndArgs...), msgAndArgs)
}

func (a Assert) ValueGreaterOrEqual(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.GreaterOrEqual(a.filter.T.t, expected, a.getValue(expected, msgAndArgs...), msgAndArgs)
}

func (a Assert) Sum(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.Equal(a.filter.T.t, expected, sumOfType(expected, a.filter.filter.Metrics()), msgAndArgs...)
}

func (a Assert) SumGreater(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.Greater(a.filter.T.t, expected, sumOfType(expected, a.filter.filter.Metrics()), msgAndArgs...)
}

func (a Assert) SumGreaterOrEqual(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.GreaterOrEqual(a.filter.T.t, expected, sumOfType(expected, a.filter.filter.Metrics()), msgAndArgs...)
}

func (a Assert) SumLess(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.Less(a.filter.T.t, expected, sumOfType(expected, a.filter.filter.Metrics()), msgAndArgs...)
}

func (a Assert) SumLessOrEqual(expected IntOrFloat, msgAndArgs ...any) bool {
	a.filter.T.t.Helper()
	return assert2.LessOrEqual(a.filter.T.t, expected, sumOfType(expected, a.filter.filter.Metrics()), msgAndArgs...)
}

type AssertGroup[K comparable] struct {
	group Group[K]
}

func (a AssertGroup[K]) Sum(expected map[K]int64) bool {
	a.group.T.t.Helper()
	return assert2.Equal(a.group.T.t, expected, a.group.group.IntSum())
}
