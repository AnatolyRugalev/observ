package logq

import (
	"github.com/AnatolyRugalev/observ/internal/genq"
	"fmt"
)

// GroupFunc is a grouping function for Record objects.
type GroupFunc genq.GroupFunc[string, Record]

func ByMessage() GroupFunc {
	return func(v Record) string {
		return v.Message
	}
}

func ByAttr(name string) GroupFunc {
	return func(v Record) string {
		return fmt.Sprintf("%s", v.Attributes[name])
	}
}
