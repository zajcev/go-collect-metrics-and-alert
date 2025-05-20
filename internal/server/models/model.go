package models

import (
	"net/http"

	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MemStorage struct {
	Metrics map[string]Metric
}

// NewMetricsStorage creates a new instance of MemStorage
func NewMetricsStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]Metric, 100),
	}
}

// SetGauge sets the value of a gauge metric
func (ms *MemStorage) SetGauge(name string, metricType string, value float64) int {
	ms.Metrics[name] = Metric{
		ID:    name,
		MType: metricType,
		Value: &value,
	}
	return http.StatusOK
}

// SetCounter sets the value of a counter metric
func (ms *MemStorage) SetCounter(name string, metricType string, value int64) int {
	m, exist := ms.Metrics[name]
	if !exist {
		ms.Metrics[name] = Metric{
			ID:    name,
			MType: metricType,
			Delta: &value,
		}
	} else {
		*m.Delta += value
		ms.Metrics[name] = m
	}
	return http.StatusOK
}

// SetGaugeJSON sets the value of a gauge metric from JSON
func (ms *MemStorage) SetCounterJSON(input Metric) int {
	m, exist := ms.Metrics[input.ID]
	if !exist {
		ms.Metrics[input.ID] = input
	} else {
		if m.Delta == nil || input.Delta == nil {
			return http.StatusBadRequest
		}
		*m.Delta += *input.Delta
		ms.Metrics[input.ID] = m
	}
	return http.StatusOK
}

// SetGaugeJSON sets the value of a gauge metric from JSON
func (ms *MemStorage) SetGaugeJSON(input Metric) int {
	if input.Value == nil {
		return http.StatusBadRequest
	}
	ms.Metrics[input.ID] = input
	return http.StatusOK
}

// GetMetric returns the value of a metric by name and type
func (ms *MemStorage) GetMetric(name string, metricType string) string {
	metric, exists := ms.Metrics[name]
	if !exists || metric.MType != metricType {
		return ""
	}
	if metricType == constants.Gauge {
		return convert.GetString(metric.Value)
	} else if metricType == constants.Counter {
		return convert.GetString(metric.Delta)
	} else {
		return ""
	}
}

// GetMetricJSON returns the value of a metric by name and type from JSON
func (ms *MemStorage) GetMetricJSON(input Metric) (Metric, int) {
	m, exist := ms.Metrics[input.ID]
	if exist {
		input = m
		return input, http.StatusOK
	}
	return Metric{}, http.StatusNotFound
}

// GetAllMetrics returns all metrics
func (ms *MemStorage) GetAllMetrics() *MemStorage {
	return ms
}

// SetMetricList sets a list of metrics
func (ms *MemStorage) SetMetricList(list []Metric) int {
	for _, v := range list {
		if v.MType == constants.Gauge {
			ms.SetGaugeJSON(v)
		} else if v.MType == constants.Counter {
			ms.SetCounterJSON(v)
		} else {
			return http.StatusBadRequest
		}
	}
	return http.StatusOK
}
