//go:build !logt

package logt

import (
	"context"
	"github.com/AnatolyRugalev/observ/logq"
	"github.com/AnatolyRugalev/observ/logt/logwait"
	"testing"
	"github.com/AnatolyRugalev/observ/internal/gent"
	"github.com/AnatolyRugalev/observ/internal/genq"
)

type LogT struct {
	t *testing.T
	options

	sink   *logq.Sink
	finish func(LogT) Filter
}

func (t LogT) T() *testing.T {
	return t.t
}

func (t LogT) WithT(tt *testing.T) LogT {
	t.t = tt
	return t
}

func (t LogT) Start() LogT {
	if t.collector == nil {
		panic("observ/logt: collector is not set")
	}
	parent := t.sink
	sink := logq.NewSink(parent, t.collectFilter)
	t.sink = sink
	stop := t.collector.CaptureLogs(t.sink)
	oldFinish := t.finish
	t.finish = func(s LogT) Filter {
		stop()
		s.sink = parent
		s.finish = oldFinish
		return Filter{
			T:      t,
			filter: sink.Records().Where(),
		}
	}
	return t
}

func (t LogT) Finish() Filter {
	if t.finish == nil {
		panic("observ/logt: not started")
	}
	return t.finish(t)
}

func (t LogT) Scope(fn func(lgt LogT)) Filter {
	t.t.Helper()
	scope := t.Start()
	fn(scope)
	return scope.Finish()
}

func (t LogT) Collect(filters ...logq.FilterFunc) Filter {
	records := Filter{
		slice: gent.NewSlice[logq.Record](t.t, genq.Slice[logq.Record](t.sink.Records())),
	}
	return records.Where(filters...)
}

func (t LogT) Wait(opts ...logwait.Option) Filter {
	t.t.Helper()
	o := logwait.NewOptions(t.waitContext, t.waitTimeout, opts...)
	ctx, cancel := context.WithTimeout(o.Context, o.Duration)
	defer cancel()
	match := make(chan logq.Record)
	scope := t.WithCollectFilter(func(r logq.Record) bool {
		if o.Filter == nil || o.Filter(r) {
			match <- r
			return true
		}
		return false
	}).Start()
	var records logq.Records
loop:
	for {
		select {
		case record, ok := <-match:
			if !ok {
				break loop
			}
			records = append(records, record)
			if !o.ShouldContinue(record) {
				break loop
			}
		case <-ctx.Done():
			break loop
		}
	}
	close(match)
	return scope.Finish()
}
