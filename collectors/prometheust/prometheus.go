package prometheust

import (
	"github.com/AnatolyRugalev/observ/metrq"
	"github.com/AnatolyRugalev/observ/metrt"
	"github.com/prometheus/client_golang/prometheus"
	prompb "github.com/prometheus/client_model/go"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/sdk/metric"
)

var defaultOptions = options{
	scope: "prometheus",
}

type options struct {
	scope string
}

type Option func(o *options)

// WithScope will set the scope of Prometheus metrics.
// Prometheus doesn't provide scoping on its own, so it is possible to override the value using this option.
// If not set, "prometheus" scope will be assigned to all metrics.
func WithScope(scope string) Option {
	return func(o *options) {
		o.scope = scope
	}
}

type Collector struct {
	options
	gatherer prometheus.Gatherer

	reader metric.Reader
}

func New(gatherer prometheus.Gatherer, opts ...Option) *Collector {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return &Collector{
		options:  o,
		gatherer: gatherer,
	}
}

func (c *Collector) CollectMetrics(f metrq.FilterFunc) metrq.Metrics {
	pFamilies := lo.Must(c.gatherer.Gather())
	metrics := make(metrq.Metrics, 0, len(pFamilies))
	for _, family := range pFamilies {
		m := metrq.Metric{
			Scope:       c.scope,
			Name:        family.GetName(),
			Description: family.GetHelp(),
		}
		for _, pMetric := range family.Metric {
			mm := m.Clone()
			mm.Attributes = newAttributes(pMetric.Label)
			switch family.GetType() {
			case prompb.MetricType_COUNTER:
				mm.Kind = metrq.KindCounter
				mm.Value = pMetric.GetCounter().GetValue()
			case prompb.MetricType_GAUGE:
				mm.Kind = metrq.KindGauge
				mm.Value = pMetric.GetGauge().GetValue()
			case prompb.MetricType_HISTOGRAM:
				mm.Kind = metrq.KindHistogram
				mm.Value = pMetric.GetHistogram().GetSampleSum()
				// TODO: support histogram data
			default:
				continue
			}
			if f(mm) {
				metrics = append(metrics, mm)
			}
		}
	}
	return metrics
}

func newAttributes(labels []*prompb.LabelPair) metrq.Attributes {
	attrs := make([]metrq.Attribute, len(labels))
	for i, label := range labels {
		attrs[i] = metrq.Attribute{
			Key:   label.GetName(),
			Value: label.GetValue(),
		}
	}
	return metrq.NewAttributes(attrs...)
}

func Gatherer(gatherer prometheus.Gatherer, opts ...Option) metrt.Option {
	return metrt.WithCollector(New(gatherer, opts...))
}

// Default returns an option for metrt.New that is using prometheus.DefaultGatherer.
// If your metrics are registered using `promauto` or `prometheus.DefaultRegistry`, they will be picked up by this
// collector.
func Default(opts ...Option) metrt.Option {
	return Gatherer(prometheus.DefaultGatherer, opts...)
}
