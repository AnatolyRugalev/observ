package metrq

import (
	"bytes"
	"github.com/cespare/xxhash/v2"
)

// AttributesKey can be used as a map key.
type AttributesKey struct {
	hash uint64
}

type Attributes struct {
	// attrs are sorted at all times
	attrs []Attribute
	index map[string]int
	key   AttributesKey
}

func (a Attributes) Get(name string) string {
	idx, ok := a.index[name]
	if !ok {
		return ""
	}
	return a.attrs[idx].Value
}

// NewAttributes creates Attributes object.
// The order of attrs should be deterministic (e.g. sorted alphabetically by key).
func NewAttributes(attrs ...Attribute) Attributes {
	a := Attributes{
		attrs: attrs,
		index: make(map[string]int, len(attrs)),
		key:   NewAttributesKey(attrs...),
	}
	for i, attr := range attrs {
		a.index[attr.Key] = i
	}
	return a
}

type Attribute struct {
	Key   string
	Value string
}

func (a Attributes) Key() AttributesKey {
	return a.key
}

func NewAttributesKey(attrs ...Attribute) AttributesKey {
	if len(attrs) == 0 {
		return AttributesKey{}
	}
	b := bytes.Buffer{}
	for _, attr := range attrs {
		b.WriteString(attr.Key)
		b.WriteRune('=')
		b.WriteString(attr.Value)
		b.WriteRune(';')
	}
	return AttributesKey{hash: xxhash.Sum64(b.Bytes())}
}
