package metrmulti

import (
	"github.com/AnatolyRugalev/observ/logcollect"
	"github.com/AnatolyRugalev/observ/logt"
)

type MultiCollector struct {
	collectors []logcollect.Collector
}

func (m MultiCollector) CaptureLogs(sink logcollect.Sink) (done func()) {
	stops := make([]func(), 0, len(m.collectors))
	for _, c := range m.collectors {
		stops = append(stops, c.CaptureLogs(sink))
	}
	return func() {
		for _, stop := range stops {
			stop()
		}
	}
}

func New(collectors ...logcollect.Collector) MultiCollector {
	return MultiCollector{
		collectors: collectors,
	}
}

func Multi(collectors ...logcollect.Collector) logt.Option {
	return logt.WithCollector(New(collectors...))
}
