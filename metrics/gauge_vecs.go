package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type GaugeVecs interface {
	NewGaugeVec(gaugeName string, labelName string, help string)
	IncrementGaugeVec(gaugeVecName string, labelName string, labelValue string)
	DecrementGaugeVec(gaugeVecName string, labelName string, labelValue string)
	AddToGaugeVec(gaugeVecName string, labelName string, labelValue string, gaugeVecValue float64)
	SubToGaugeVec(gaugeVecName string, labelName string, labelValue string, gaugeVecValue float64)
	SetGaugeVec(gaugeVecName string, labelName string, labelValue string, gaugeVecValue float64)
}

func (m *metric) NewGaugeVec(gaugeName string, labelName string, help string) {
	m.gaugeVecs[gaugeName] = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: gaugeName,
			Help: help,
		},
		[]string{labelName},
	)
}

func (m *metric) IncrementGaugeVec(gaugeVecName string, labelName string, labelValue string) {
	if gaugeVec, ok := m.gaugeVecs[gaugeVecName]; ok {
		gaugeVec.With(prometheus.Labels{labelName: labelValue}).Inc()
	}
}

func (m *metric) DecrementGaugeVec(gaugeVecName string, labelName string, labelValue string) {
	if gaugeVec, ok := m.gaugeVecs[gaugeVecName]; ok {
		gaugeVec.With(prometheus.Labels{labelName: labelValue}).Dec()
	}
}

func (m *metric) AddToGaugeVec(gaugeVecName string, labelName string, labelValue string, gaugeVecValue float64) {
	if gaugeVec, ok := m.gaugeVecs[gaugeVecName]; ok {
		gaugeVec.With(prometheus.Labels{labelName: labelValue}).Add(gaugeVecValue)
	}
}

func (m *metric) SubToGaugeVec(gaugeVecName string, labelName string, labelValue string, gaugeVecValue float64) {
	if gaugeVec, ok := m.gaugeVecs[gaugeVecName]; ok {
		gaugeVec.With(prometheus.Labels{labelName: labelValue}).Sub(gaugeVecValue)
	}
}

func (m *metric) SetGaugeVec(gaugeVecName string, labelName string, labelValue string, gaugeVecValue float64) {
	if gaugeVec, ok := m.gaugeVecs[gaugeVecName]; ok {
		gaugeVec.With(prometheus.Labels{labelName: labelValue}).Set(gaugeVecValue)
	}
}
