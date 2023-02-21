package metrq

import (
	"fmt"
	"github.com/AnatolyRugalev/observ/internal/genq"
	"github.com/montanaflynn/stats"
	"github.com/samber/lo"
)

type AggregationFunc = genq.AggregationFunc[float64]

type MetricAggregationFunc = genq.AggregationFunc[Metric]

func MetricAgg(fn AggregationFunc) MetricAggregationFunc {
	return func(metrics []Metric) Metric {
		metric := metrics[0].Clone()
		metric.Value = fn(lo.Map(metrics, func(m Metric, _ int) float64 {
			return m.Value
		}))
		if metric.Kind == KindHistogram {
			// TODO: sum count, bucket counts
		}
		return metric
	}
}

func StatsAggregation(f func(input stats.Float64Data) (float64, error)) AggregationFunc {
	return func(values []float64) float64 {
		result, err := f(values)
		if err != nil {
			panic(fmt.Sprintf("metricq: aggregation function error: %s", err.Error()))
		}
		return result
	}
}

func Sum() AggregationFunc {
	return StatsAggregation(stats.Sum)
}

// Sub subtracts sum of all sets after the first one from the first set.
func Sub() AggregationFunc {
	return func(values []float64) float64 {
		if len(values) < 2 {
			panic("metricq: Subtract requires at least two metric sets")
		}
		return values[0] - Sum()(values[1:])
	}
}

func Mean() AggregationFunc {
	return StatsAggregation(stats.Mean)
}

func Median() AggregationFunc {
	return StatsAggregation(stats.Median)
}

func Max() AggregationFunc {
	return StatsAggregation(stats.Max)
}

func Min() AggregationFunc {
	return StatsAggregation(stats.Min)
}
