package models

import (
	"context"
	"net/http"
	"testing"

	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
)

func TestSetValueRaw(t *testing.T) {
	ms := NewMemStorage()
	ctx := context.Background()
	value := 10.5
	name := "metric1"
	metricType := constants.Gauge

	status := ms.SetValueRaw(ctx, name, metricType, value)
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	metric, exists := ms.Storage[name]
	if !exists || *metric.Value != value {
		t.Fatalf("expected metric with value %f, got %v", value, metric)
	}
}

func TestSetDeltaRaw(t *testing.T) {
	ms := NewMemStorage()
	ctx := context.Background()
	name := "metric2"
	metricType := constants.Counter
	delta := int64(5)

	status := ms.SetDeltaRaw(ctx, name, metricType, delta)
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	metric, exists := ms.Storage[name]
	if !exists || *metric.Delta != delta {
		t.Fatalf("expected metric with delta %d, got %v", delta, metric)
	}

	// Test updating the delta
	status = ms.SetDeltaRaw(ctx, name, metricType, 3)
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	metric, _ = ms.Storage[name]
	if *metric.Delta != 8 {
		t.Fatalf("expected updated delta %d, got %d", 8, *metric.Delta)
	}
}

func TestSetDeltaJSON(t *testing.T) {
	ms := NewMemStorage()
	ctx := context.Background()
	name := "metric3"
	delta := int64(10)

	status := ms.SetDeltaJSON(ctx, Metric{ID: name, Delta: &delta})
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	metric, exists := ms.Storage[name]
	if !exists || *metric.Delta != delta {
		t.Fatalf("expected metric with delta %d, got %v", delta, metric)
	}

	// Test updating delta with existing metric
	deltaUpdate := int64(5)
	status = ms.SetDeltaJSON(ctx, Metric{ID: name, Delta: &deltaUpdate})
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	metric, _ = ms.Storage[name]
	if *metric.Delta != 15 {
		t.Fatalf("expected updated delta %d, got %d", 15, *metric.Delta)
	}
}

func TestSetValueJSON(t *testing.T) {
	ms := NewMemStorage()
	ctx := context.Background()
	name := "metric4"
	value := 20.0

	status := ms.SetValueJSON(ctx, Metric{ID: name, Value: &value})
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	metric, exists := ms.Storage[name]
	if !exists || *metric.Value != value {
		t.Fatalf("expected metric with value %f, got %v", value, metric)
	}
}

func TestGetMetricRaw(t *testing.T) {
	ms := NewMemStorage()
	ctx := context.Background()
	name := "metric5"
	value := 25.0
	metricType := constants.Gauge
	expected := "25"

	ms.SetValueRaw(ctx, name, metricType, value)
	result := ms.GetMetricRaw(ctx, name, metricType)
	if result != expected {
		t.Fatalf("expected value %f, got %v", value, result)
	}

	// Test non-existing metric
	result = ms.GetMetricRaw(ctx, "nonexistent", metricType)
	if result != "" {
		t.Fatalf("expected empty result for nonexistent metric, got %v", result)
	}
}

func TestGetMetricJSON(t *testing.T) {
	ms := NewMemStorage()
	ctx := context.Background()
	name := "metric6"
	value := 30.0
	metricType := constants.Gauge

	ms.SetValueRaw(ctx, name, metricType, value)
	metric, status := ms.GetMetricJSON(ctx, Metric{ID: name})
	if status != http.StatusOK || metric.Value == nil || *metric.Value != value {
		t.Fatalf("expected status %d and value %f, got %d and %v", http.StatusOK, value, status, metric)
	}

	// Test non-existing metric
	metric, status = ms.GetMetricJSON(ctx, Metric{ID: "nonexistent"})
	if status != http.StatusNotFound {
		t.Fatalf("expected status %d for nonexistent metric, got %d", http.StatusNotFound, status)
	}
}
