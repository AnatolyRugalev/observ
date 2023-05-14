package logt

import (
	"github.com/AnatolyRugalev/observ/logq"
	"testing"
	"github.com/AnatolyRugalev/observ/internal/gent"
)

// Filter is a slice of records
type Filter struct {
	t      *testing.T
	filter logq.Filter `chaingen:"*,-Resolve,wrap(*)=wrapGroup"`
}

func (f Filter) wrapGroup(group logq.Group[string]) Group[string] {
	return Group[string]{
		t:     f.t,
		group: group,
	}
}

func (f Filter) Assert() Assert {
	return Assert{
		assert: gent.NewAssert[logq.Record](f.t, f.filter),
	}
}

type Group[K comparable] struct {
	t     *testing.T
	group logq.Group[K] `chaingen:"*,wrap(*)=wrapGroup"`
}
