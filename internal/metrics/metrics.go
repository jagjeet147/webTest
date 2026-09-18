package metrics

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	mu       sync.RWMutex
	requests uint64
	success  uint64
	errors   uint64
	latency  Histogram
}

func (m *Metrics) Record(statusCode int, duration time.Duration, err error) {
	m.mu.Lock()
	m.requests++
	if err != nil || statusCode < 200 || statusCode >= 400 {
		m.errors++
	} else {
		m.success++
	}
	m.mu.Unlock()
	m.latency.Observe(duration.Seconds() * 1000)
}

func (m *Metrics) Snapshot() Snapshot {
	m.mu.RLock()
	snapshot := Snapshot{Requests: m.requests, Success: m.success, Errors: m.errors}
	m.mu.RUnlock()
	snapshot.P50 = m.latency.Percentile(0.50)
	snapshot.P95 = m.latency.Percentile(0.95)
	snapshot.P99 = m.latency.Percentile(0.99)
	return snapshot
}

type Snapshot struct {
	Requests uint64
	Success  uint64
	Errors   uint64
	P50      float64
	P95      float64
	P99      float64
}

func (s Snapshot) ErrorRate() float64 {
	if s.Requests == 0 {
		return 0
	}
	return float64(s.Errors) / float64(s.Requests) * 100
}

func (s Snapshot) String() string {
	return fmt.Sprintf("requests=%d success=%d errors=%d", s.Requests, s.Success, s.Errors)
}
