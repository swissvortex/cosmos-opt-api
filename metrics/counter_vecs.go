package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CounterVecs interface {
	NewCounterVec(counterName string, labelName string, help string)
	IncrementCounterVec(counterName string, labelName string, labelValue string)
	AddToCounterVec(counterName string, labelName string, labelValue string, counterValue float64)
}

func (m *metric) NewCounterVec(counterName string, labelName string, help string) {
	m.counterVecs[counterName] = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: counterName,
			Help: help,
		},
		[]string{labelName},
	)
}

func (m *metric) IncrementCounterVec(counterVecName string, labelName string, labelValue string) {
	if counterVec, ok := m.counterVecs[counterVecName]; ok {
		counterVec.With(prometheus.Labels{labelName: labelValue}).Inc()
	}
}

func (m *metric) AddToCounterVec(counterVecName string, labelName string, labelValue string, counterVecValue float64) {
	if counterVec, ok := m.counterVecs[counterVecName]; ok {
		counterVec.With(prometheus.Labels{labelName: labelValue}).Add(counterVecValue)
	}
}
