package stats

import (
	"math"
	"golang.org/x/exp/slices"
	"sort"
	"github.com/AnatolyRugalev/observ/internal/genq"
)

// StatFunc is an aggregation func for floating point numbers.
type StatFunc genq.AggregationFunc[float64]

// Sum sums all the values.
func Sum(values []float64) float64 {
	sum := 0.
	for _, val := range values {
		sum += val
	}
	return sum
}

// Sub performs subtraction from the first element.
func Sub(values []float64) float64 {
	if len(values) == 0 {
		return math.Inf(-1)
	}
	return values[0] - Sum(values[1:])
}

// Mean evaluates arithmetic mean value.
// NOTE: Arithmetic Mean in statistics is the same as Average in mathematics.
// Learn more: https://medium.com/@seema.singh/average-vs-mean-534b1ac85401
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return math.Inf(1)
	}
	return Sum(values) / float64(len(values))
}

// Median evaluates median value.
func Median(values []float64) float64 {
	length := len(values)
	if length == 0 {
		return 0
	}
	sorted := slices.Clone(values)
	sort.Float64s(sorted)
	if length%2 == 0 {
		return (sorted[length/2-1] + sorted[length/2]) / 2
	}
	return sorted[length/2]
}

// Max evaluates maximum value.
func Max(values []float64) float64 {
	max := math.Inf(-1)
	for _, val := range values {
		if val > max {
			max = val
		}
	}
	return max
}

// Min evaluates minimum value.
func Min(values []float64) float64 {
	min := math.Inf(1)
	for _, val := range values {
		if val < min {
			min = val
		}
	}
	return min
}
