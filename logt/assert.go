package logt

import (
	"github.com/AnatolyRugalev/observ/internal/gent"
	"github.com/AnatolyRugalev/observ/logq"
)

type Assert struct {
	assert gent.Assert[logq.Record]
}

type AssertGroup[K comparable] struct {
	group gent.Group[string, K]
}
