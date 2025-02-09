package models

import (
	"fmt"
	"net/http"
)

type Metric struct {
	Type  string
	Value interface{}
}
type MemStorage struct {
	Metrics map[string]Metric
}

func NewMetricsStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]Metric),
	}
}

func (ms *MemStorage) SetGauge(name string, metricType string, value float64) int {
	ms.Metrics[name] = Metric{
		Type:  metricType,
		Value: value,
	}
	return http.StatusOK
}

func (ms *MemStorage) SetCounter(name string, metricType string, value int64) int {
	m, exist := ms.Metrics[name]
	if !exist {
		ms.Metrics[name] = Metric{
			Type:  metricType,
			Value: value,
		}
	} else {
		m.Value = m.Value.(int64) + value
		ms.Metrics[name] = m
	}
	return http.StatusOK
}

func (ms *MemStorage) GetMetric(name string, metricType string) string {
	metric, exists := ms.Metrics[name]
	if !exists || metric.Type != metricType {
		return ""
	}
	return fmt.Sprintf("%v", metric.Value)
}

func (ms *MemStorage) GetAllMetrics() *MemStorage {
	return ms
}
