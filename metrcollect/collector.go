package metrcollect

import (
	"github.com/AnatolyRugalev/observ/metrq"
)

type Collector interface {
	CollectMetrics(f metrq.FilterFunc) metrq.Metrics
}
