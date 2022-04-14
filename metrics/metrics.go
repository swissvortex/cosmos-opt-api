package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metric interface {
	Gauges
}

type metric struct {
	gauges map[string]prometheus.Gauge
}

func New() Metric {
	gauges := make(map[string]prometheus.Gauge)
	return &metric{
		gauges: gauges,
	}
}
