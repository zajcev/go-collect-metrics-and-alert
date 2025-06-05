package models

import (
	"context"
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

type MemStorageService interface {
	SetDeltaRaw(ctx context.Context, name string, metricType string, value int64) int
	SetValueRaw(ctx context.Context, name string, metricType string, value float64) int
	SetValueJSON(ctx context.Context, input Metric) int
	SetDeltaJSON(ctx context.Context, input Metric) int
	GetMetricRaw(ctx context.Context, name string, metricType string) any
	GetMetricJSON(ctx context.Context, input Metric) (Metric, int)
	GetAllMetrics(ctx context.Context, ms *MemStorage) *MemStorage
	SetListJSON(ctx context.Context, list []Metric) int
	Save()
}
type MemStorage struct {
	Storage map[string]Metric
}

// NewMemStorage constructor for MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Storage: make(map[string]Metric),
	}
}

// SetValueRaw sets the value of a gauge metric
func (ms *MemStorage) SetValueRaw(ctx context.Context, name string, metricType string, value float64) int {
	ms.Storage[name] = Metric{
		ID:    name,
		MType: metricType,
		Value: &value,
	}
	return http.StatusOK
}

// SetDeltaRaw sets the value of a counter metric
func (ms *MemStorage) SetDeltaRaw(ctx context.Context, name string, metricType string, value int64) int {
	m, exist := ms.Storage[name]
	if !exist {
		ms.Storage[name] = Metric{
			ID:    name,
			MType: metricType,
			Delta: &value,
		}
	} else {
		*m.Delta += value
		ms.Storage[name] = m
	}
	return http.StatusOK
}

// SetDeltaJSON sets the value of a gauge metric from JSON
func (ms *MemStorage) SetDeltaJSON(ctx context.Context, m Metric) int {
	if m.Delta == nil {
		return http.StatusBadRequest
	}
	currentMetric, exists := ms.Storage[m.ID]
	if !exists {
		ms.Storage[m.ID] = m
	} else {
		if currentMetric.Delta == nil {
			return http.StatusBadRequest
		}
		*currentMetric.Delta += *m.Delta
		ms.Storage[m.ID] = currentMetric
	}
	return http.StatusOK
}

// SetDeltaJSON sets the value of a gauge metric from JSON
func (ms *MemStorage) SetValueJSON(ctx context.Context, m Metric) int {
	if m.Value == nil {
		return http.StatusBadRequest
	}
	ms.Storage[m.ID] = m
	return http.StatusOK
}

// GetMetricRaw returns the value of a metric by name and type
func (ms *MemStorage) GetMetricRaw(ctx context.Context, name string, metricType string) any {
	metric, exists := ms.Storage[name]
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
func (ms *MemStorage) GetMetricJSON(ctx context.Context, m Metric) (Metric, int) {
	result, exist := ms.Storage[m.ID]
	if exist {
		m = result
		return m, http.StatusOK
	}
	return Metric{}, http.StatusNotFound
}

// GetAllMetrics returns all metrics
func (ms *MemStorage) GetAllMetrics(ctx context.Context) *MemStorage {
	return ms
}

// SetListJSON sets a list of metrics
func (ms *MemStorage) SetListJSON(ctx context.Context, list []Metric) int {
	for _, v := range list {
		if v.MType == constants.Gauge {
			ms.SetValueJSON(ctx, v)
		} else if v.MType == constants.Counter {
			ms.SetDeltaJSON(ctx, v)
		} else {
			return http.StatusBadRequest
		}
	}
	return http.StatusOK
}

func (ms *MemStorage) Ping(ctx context.Context) error {
	// Для in-memory хранилища всегда возвращаем nil (всегда доступно)
	return nil
}
