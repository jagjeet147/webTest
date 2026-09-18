package metrics

import (
	"sort"
	"sync"
)

type Histogram struct {
	mu     sync.RWMutex
	values []float64
}

func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	h.values = append(h.values, value)
	h.mu.Unlock()
}

func (h *Histogram) Percentile(percent float64) float64 {
	h.mu.RLock()
	values := append([]float64(nil), h.values...)
	h.mu.RUnlock()
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values)
	index := int(float64(len(values)-1) * percent)
	return values[index]
}
