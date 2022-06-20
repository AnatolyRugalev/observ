package promquery

import (
	"fmt"

	promcli "github.com/prometheus/client_golang/prometheus"
	prompb "github.com/prometheus/client_model/go"
)

type Querier struct {
	registry *promcli.Registry
}

type M map[string]float64

func New(registry *promcli.Registry) *Querier {
	return &Querier{
		registry: registry,
	}
}

func (t *Querier) Registry() *promcli.Registry {
	return t.registry
}

type FamilyFunc func(family *prompb.MetricFamily) bool

func (t *Querier) Each(eachFunc FamilyFunc) error {
	families, err := t.registry.Gather()
	if err != nil {
		return fmt.Errorf("error gathering metrics: %w", err)
	}
	for _, fam := range families {
		if !eachFunc(fam) {
			return nil
		}
	}
	return nil
}

func (t *Querier) Filter(filterFunc FamilyFunc) ([]*prompb.MetricFamily, error) {
	var filtered []*prompb.MetricFamily
	err := t.Each(func(family *prompb.MetricFamily) bool {
		if filterFunc(family) {
			filtered = append(filtered, family)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating families: %w", err)
	}
	return filtered, nil
}

func (t *Querier) Find(findFunc FamilyFunc) (*prompb.MetricFamily, error) {
	var found *prompb.MetricFamily
	err := t.Each(func(family *prompb.MetricFamily) bool {
		if findFunc(family) {
			found = family
			return false
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating families: %w", err)
	}
	return found, nil
}

func (t *Querier) Query(name string, labelsKV ...string) Query {
	return Query{
		querier: t,
		name:    name,
		labels:  LabelsKV(labelsKV...),
	}
}

type Query struct {
	querier *Querier
	name    string
	labels  L
}

func (q Query) Labels(labelsKV ...string) Query {
	q.labels = q.labels.Merge(LabelsKV(labelsKV...))
	return q
}

func LabelPairsToLabels(pairs []*prompb.LabelPair) L {
	labels := make(L, len(pairs)/2)
	for _, pair := range pairs {
		labels[*pair.Name] = *pair.Value
	}
	return labels
}

func FilterMetricsUsingLabels(metrics []*prompb.Metric, labels L) []*prompb.Metric {
	matched := make([]*prompb.Metric, 0)
	for _, m := range metrics {
		metricLabels := LabelPairsToLabels(m.Label)
		match := true
		for k, v := range labels {
			if metricLabels[k] != v {
				match = false
				break
			}
		}
		if match {
			matched = append(matched, m)
		}
	}
	return matched
}

func (q Query) Match() (*prompb.MetricFamily, error) {
	var found *prompb.MetricFamily
	err := q.querier.Each(func(family *prompb.MetricFamily) bool {
		if *family.Name != q.name {
			return true
		}
		metrics := FilterMetricsUsingLabels(family.Metric, q.labels)
		if len(metrics) == 0 {
			return true
		}
		found = &prompb.MetricFamily{
			Name:   family.Name,
			Help:   family.Help,
			Type:   family.Type,
			Metric: metrics,
		}
		return false
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating over metrics: %w", err)
	}
	return found, nil
}

func (q Query) All() ([]Value, error) {
	fam, err := q.Match()
	if err != nil {
		return nil, fmt.Errorf("error matching metric: %w", err)
	}
	if fam == nil {
		return nil, nil
	}
	values := make([]Value, 0, len(fam.Metric))
	for _, m := range fam.Metric {
		values = append(values, Value{
			Labels: LabelPairsToLabels(m.Label),
			// TODO: support more types
			Value: *m.Counter.Value,
		})
	}
	return values, nil
}

func (q Query) First() (float64, error) {
	all, err := q.All()
	if err != nil {
		return 0, fmt.Errorf("error querying metrics: %w", err)
	}
	if len(all) == 0 {
		return 0, nil
	}
	return all[0].Value, nil
}

func (q Query) One() (float64, error) {
	all, err := q.All()
	if err != nil {
		return 0, fmt.Errorf("error querying metrics: %w", err)
	}
	if len(all) == 0 {
		return 0, fmt.Errorf("metric not found")
	}
	if len(all) > 1 {
		return 0, fmt.Errorf("too many matches: %d, 1 expec21ted", len(all))
	}
	return all[0].Value, nil
}

type Value struct {
	Labels L
	Value  float64
}

func (q Query) Group(groupBy string) GroupQuery {
	return GroupQuery{
		query:   q,
		groupBy: groupBy,
		aggFunc: Sum,
	}
}

type AggFunc func(values []float64) float64

type GroupQuery struct {
	query   Query
	groupBy string
	aggFunc AggFunc
}

func Sum(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

func (g GroupQuery) Map() (M, error) {
	values, err := g.query.All()
	if err != nil {
		return nil, fmt.Errorf("error querying metrics: %w", err)
	}
	grouped := make(map[string][]float64, len(values))
	for _, v := range values {
		key, ok := v.Labels[g.groupBy]
		if !ok {
			continue
		}
		grouped[key] = append(grouped[key], v.Value)
	}
	result := make(M, len(grouped))
	for k, vals := range grouped {
		result[k] = g.aggFunc(vals)
	}
	return result, nil
}
