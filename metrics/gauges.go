package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Gauges interface {
	NewGauge(gaugeName string, help string)
	IncrementGauge(gaugeName string)
	DecrementGauge(gaugeName string)
	AddToGauge(gaugeName string, gaugeValue float64)
	SubToGauge(gaugeName string, gaugeValue float64)
	SetGauge(gaugeName string, gaugeValue float64)
}

func (m *metric) NewGauge(gaugeName string, help string) {
	m.gauges[gaugeName] = promauto.NewGauge(prometheus.GaugeOpts{
		Name: gaugeName,
		Help: help,
	})
}

func (m *metric) IncrementGauge(gaugeName string) {
	if gauge, ok := m.gauges[gaugeName]; ok {
		gauge.Inc()
	}
}

func (m *metric) DecrementGauge(gaugeName string) {
	if gauge, ok := m.gauges[gaugeName]; ok {
		gauge.Dec()
	}
}

func (m *metric) AddToGauge(gaugeName string, gaugeValue float64) {
	if gauge, ok := m.gauges[gaugeName]; ok {
		gauge.Add(gaugeValue)
	}
}

func (m *metric) SubToGauge(gaugeName string, gaugeValue float64) {
	if gauge, ok := m.gauges[gaugeName]; ok {
		gauge.Sub(gaugeValue)
	}
}

func (m *metric) SetGauge(gaugeName string, gaugeValue float64) {
	if gauge, ok := m.gauges[gaugeName]; ok {
		gauge.Set(gaugeValue)
	}
}
