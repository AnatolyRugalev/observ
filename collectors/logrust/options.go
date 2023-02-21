package logrust

import "github.com/sirupsen/logrus"

type options struct {
	levels []logrus.Level
}

var defaultOptions = options{
	levels: logrus.AllLevels,
}

type Option func(o *options)

// WithLevels specifies at which logrus log levels should be captured.
// If not set, all levels will be used.
func WithLevels(levels []logrus.Level) Option {
	return func(o *options) {
		o.levels = levels
	}
}
