package logrust

import (
	"github.com/AnatolyRugalev/observ/logcollect"
	"github.com/AnatolyRugalev/observ/logq"
	"github.com/AnatolyRugalev/observ/logt"
	"github.com/mitchellh/copystructure"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"sync/atomic"
)

type Collector struct {
	options options
	logger  *logrus.Logger
	sink    atomic.Pointer[logcollect.Sink]
}

func (c *Collector) Close() {
	for _, level := range c.options.levels {
		c.logger.Hooks[level] = lo.Filter(c.logger.Hooks[level], func(item logrus.Hook, index int) bool {
			return item != c
		})
	}
}

func (c *Collector) Levels() []logrus.Level {
	return c.options.levels
}

func (c *Collector) Fire(entry *logrus.Entry) error {
	sink := c.sink.Load()
	if sink == nil {
		return nil
	}
	data, err := copystructure.Copy(entry.Data)
	if err != nil {
		data = logrus.Fields{
			"error": "loqcap/logruscap: failed to copy logrus entry data: " + err.Error(),
		}
	}
	attributes := data.(logrus.Fields)
	attributes[logq.AttributeLevel] = entry.Level.String()
	if entry.Caller != nil {
		attributes[logq.AttributeFile] = entry.Caller.File
		attributes[logq.AttributeLine] = entry.Caller.Line
	}

	(*sink).Add(logq.Record{
		Time:       entry.Time,
		Message:    entry.Message,
		Attributes: data.(logrus.Fields),
	})
	return nil
}

func (c *Collector) CaptureLogs(sink logcollect.Sink) func() {
	old := c.sink.Swap(&sink)
	return func() {
		c.sink.Store(old)
	}
}

func New(logger *logrus.Logger, opts ...Option) logcollect.Collector {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	c := &Collector{
		options: o,
	}
	logger.AddHook(c)
	return c
}

// Default sets up Logrus collector for logrus.StandardLogger()
func Default(opts ...Option) logt.Option {
	return Logrus(logrus.StandardLogger(), opts...)
}

// Logrus sets up collector for provided Logrus logger.
func Logrus(logger *logrus.Logger, opts ...Option) logt.Option {
	return logt.WithCollector(New(logger, opts...))
}
