package otelt

import (
	"context"
	"github.com/AnatolyRugalev/observ/metrq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func New() *Collector {
	reader := metric.NewManualReader(metric.WithTemporalitySelector(func(kind metric.InstrumentKind) metricdata.Temporality {
		return metricdata.CumulativeTemporality
	}))
	return &Collector{
		Reader: reader,
	}
}

type Collector struct {
	metric.Reader
}

func (c Collector) CollectMetrics(f metrq.FilterFunc) metrq.Metrics {
	oMetrics, err := c.Reader.Collect(context.TODO())
	if err != nil {
		panic(err)
	}
	metrics := make(metrq.Metrics, 0, len(oMetrics.ScopeMetrics))
	add := func(kind metrq.MetricKind, m metrq.Metric, attributes attribute.Set, value float64) {
		mm := m.Clone()
		mm.Kind = kind
		mm.Attributes = newAttributes(attributes)
		mm.Value = value
		if f(mm) {
			metrics = append(metrics, mm)
		}
	}
	for _, oScope := range oMetrics.ScopeMetrics {
		for _, oMetric := range oScope.Metrics {
			m := metrq.Metric{
				Scope:       oScope.Scope.Name,
				Name:        oMetric.Name,
				Description: oMetric.Description,
			}
			switch t := oMetric.Data.(type) {
			case metricdata.Sum[float64]:
				for _, dp := range t.DataPoints {
					add(metrq.KindCounter, m, dp.Attributes, dp.Value)
				}
			case metricdata.Sum[int64]:
				for _, dp := range t.DataPoints {
					add(metrq.KindCounter, m, dp.Attributes, float64(dp.Value))
				}
			case metricdata.Gauge[float64]:
				for _, dp := range t.DataPoints {
					add(metrq.KindGauge, m, dp.Attributes, dp.Value)
				}
			case metricdata.Gauge[int64]:
				for _, dp := range t.DataPoints {
					add(metrq.KindGauge, m, dp.Attributes, float64(dp.Value))
				}
			case metricdata.Histogram:
				//TODO: support histogram data
				for _, dp := range t.DataPoints {
					add(metrq.KindHistogram, m, dp.Attributes, dp.Sum)
				}
			default:
				continue
			}
		}
	}
	return metrics
}

func newAttributes(attributes attribute.Set) metrq.Attributes {
	attrs := make([]metrq.Attribute, attributes.Len())
	for i := 0; i < attributes.Len(); i++ {
		a, _ := attributes.Get(i)
		attrs[i] = metrq.Attribute{
			Key:   string(a.Key),
			Value: a.Value.AsString(),
		}
	}
	return metrq.NewAttributes(attrs...)
}
