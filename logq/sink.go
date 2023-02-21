package logq

import "sync"

func NewSink(parent *Sink, filter FilterFunc) *Sink {
	return &Sink{
		parent: parent,
		filter: filter,
	}
}

// Sink is a shared store of records.
type Sink struct {
	parent *Sink
	filter FilterFunc

	mu          sync.RWMutex
	subscribers []Subscriber
	records     []Record
}

func (s *Sink) Add(r Record) {
	if s.parent != nil {
		s.parent.Add(r)
	}
	if !s.filter(r) {
		return
	}
	s.mu.Lock()
	s.records = append(s.records, r)
	s.mu.Unlock()
	s.mu.RLock()
	subs := s.subscribers
	s.mu.RUnlock()
	for _, sub := range subs {
		sub(r)
	}
}

func (s *Sink) Records() Records {
	s.mu.RLock()
	records := s.records
	s.mu.RUnlock()
	return records
}

type Subscriber func(record Record)

func (s *Sink) Subscribe(sub Subscriber) func() {
	s.mu.Lock()
	old := s.subscribers
	s.subscribers = append(s.subscribers, sub)
	s.mu.Unlock()
	return func() {
		s.mu.Lock()
		s.subscribers = old
		s.mu.Unlock()
	}
}
