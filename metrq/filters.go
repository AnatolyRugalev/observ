package metrq

import (
	"strings"
)

func NonZero() FilterFunc {
	return func(m Metric) bool {
		return m.Value != 0
	}
}

func Attr(name string, value string) FilterFunc {
	return func(m Metric) bool {
		return m.Attributes.Get(name) == value
	}
}

func Kind(kind MetricKind) FilterFunc {
	return func(m Metric) bool {
		return m.Kind == kind
	}
}

func Prefix(prefix string) FilterFunc {
	return func(m Metric) bool {
		return strings.HasPrefix(m.Name, prefix)
	}
}

func Name(name string) FilterFunc {
	return func(m Metric) bool {
		return m.Name == name
	}
}

func Scope(scope string) FilterFunc {
	return func(m Metric) bool {
		return m.Scope == scope
	}
}
