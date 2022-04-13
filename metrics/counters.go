package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Counters interface {
	NewCounter(counterName string, help string)
	IncrementCounter(counterName string)
	AddToCounter(counterName string, counterValue float64)
}

func (m *metric) NewCounter(counterName string, help string) {
	m.counters[counterName] = promauto.NewCounter(prometheus.CounterOpts{
		Name: counterName,
		Help: help,
	})
}

func (m *metric) IncrementCounter(counterName string) {
	if counter, ok := m.counters[counterName]; ok {
		counter.Inc()
	}
}

func (m *metric) AddToCounter(counterName string, counterValue float64) {
	if counter, ok := m.counters[counterName]; ok {
		counter.Add(counterValue)
	}
}
