package models

import (
	"net/http"

	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
)

type Metric struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
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

// SetDeltaRaw sets the value of a gauge metric
func (ms *MemStorage) SetDeltaRaw(name string, metricType string, value float64) int {
	ms.Metrics[name] = Metric{
		ID:    name,
		MType: metricType,
		Value: &value,
	}
	return http.StatusOK
}

// SetValueRaw sets the value of a counter metric
func (ms *MemStorage) SetValueRaw(name string, metricType string, value int64) int {
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

// SetDeltaJSON sets the value of a gauge metric from JSON
func (ms *MemStorage) SetValueJSON(input Metric) int {
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

// SetDeltaJSON sets the value of a gauge metric from JSON
func (ms *MemStorage) SetDeltaJSON(input Metric) int {
	if input.Value == nil {
		return http.StatusBadRequest
	}
	ms.Metrics[input.ID] = input
	return http.StatusOK
}

// GetMetricRaw returns the value of a metric by name and type
func (ms *MemStorage) GetMetricRaw(name string, metricType string) string {
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

// SetListJSON sets a list of metrics
func (ms *MemStorage) SetListJSON(list []Metric) int {
	for _, v := range list {
		if v.MType == constants.Gauge {
			ms.SetDeltaJSON(v)
		} else if v.MType == constants.Counter {
			ms.SetValueJSON(v)
		} else {
			return http.StatusBadRequest
		}
	}
	return http.StatusOK
}
