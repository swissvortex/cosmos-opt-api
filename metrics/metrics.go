package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metric interface {
	Counters
	CounterVecs
	Gauges
	GaugeVecs
}

type metric struct {
	counters    map[string]prometheus.Counter
	counterVecs map[string]*prometheus.CounterVec
	gauges      map[string]prometheus.Gauge
	gaugeVecs   map[string]*prometheus.GaugeVec
}

func New() Metric {
	counters := make(map[string]prometheus.Counter)
	counterVecs := make(map[string]*prometheus.CounterVec)
	gauges := make(map[string]prometheus.Gauge)
	gaugeVecs := make(map[string]*prometheus.GaugeVec)
	return &metric{
		counters:    counters,
		counterVecs: counterVecs,
		gauges:      gauges,
		gaugeVecs:   gaugeVecs,
	}
}
