package metrmulti

import (
	"github.com/AnatolyRugalev/observ/metrcollect"
	"github.com/AnatolyRugalev/observ/metrq"
	"github.com/AnatolyRugalev/observ/metrt"
)

type MultiCollector struct {
	collectors []metrcollect.Collector
}

func (m MultiCollector) CollectMetrics(f metrq.FilterFunc) metrq.Metrics {
	metrics := metrq.Metrics{}
	for _, c := range m.collectors {
		metrics = metrics.Add(c.CollectMetrics(f))
	}
	return metrics
}

func New(collectors ...metrcollect.Collector) MultiCollector {
	return MultiCollector{
		collectors: collectors,
	}
}

func Multi(collectors ...metrcollect.Collector) metrt.Option {
	return metrt.WithCollector(New(collectors...))
}
