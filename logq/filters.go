package logq

import "reflect"

func Message(message string) FilterFunc {
	return func(record Record) bool {
		return record.Message == message
	}
}

func Attr(key string, value any) FilterFunc {
	return func(record Record) bool {
		return reflect.DeepEqual(record.Attributes[key], value)
	}
}
