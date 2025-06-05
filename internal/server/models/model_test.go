package models

import (
	"net/http"
	"testing"
)

func TestNewMetricsStorage(t *testing.T) {
	ms := NewMetricsStorage()
	if ms == nil {
		t.Fatal("Expected non-nil MemStorage")
	}
	if len(ms.Metrics) != 0 {
		t.Fatalf("Expected empty metrics map, got %d", len(ms.Metrics))
	}
}

func TestSetGauge(t *testing.T) {
	ms := NewMetricsStorage()
	status := ms.SetDeltaRaw("test_gauge", "gauge", 1.5)
	if status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}
	if metric, exist := ms.Metrics["test_gauge"]; !exist || *metric.Value != 1.5 {
		t.Fatalf("Expected gauge value 1.5, got %v", metric.Value)
	}
}

func TestSetCounter(t *testing.T) {
	ms := NewMetricsStorage()
	status := ms.SetValueRaw("test_counter", "counter", 5)
	if status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}
	if metric, exist := ms.Metrics["test_counter"]; !exist || *metric.Delta != 5 {
		t.Fatalf("Expected counter delta 5, got %v", metric.Delta)
	}
}

func TestSetCounterIncrement(t *testing.T) {
	ms := NewMetricsStorage()
	ms.SetValueRaw("test_counter", "counter", 5)
	status := ms.SetValueRaw("test_counter", "counter", 3)
	if status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}
	if metric, exist := ms.Metrics["test_counter"]; !exist || *metric.Delta != 8 {
		t.Fatalf("Expected counter delta 8, got %v", metric.Delta)
	}
}

func TestSetCounterJSON(t *testing.T) {
	ms := NewMetricsStorage()
	input := Metric{ID: "test_counter", MType: "counter", Delta: int64Ptr(5)}
	status := ms.SetValueJSON(input)
	if status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}

	input.Delta = int64Ptr(3)
	status = ms.SetValueJSON(input)
	if status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}
	if metric, exist := ms.Metrics["test_counter"]; !exist || *metric.Delta != 8 {
		t.Fatalf("Expected counter delta 8, got %v", metric.Delta)
	}
}

func TestSetGaugeJSON(t *testing.T) {
	ms := NewMetricsStorage()
	input := Metric{ID: "test_gauge", MType: "gauge", Value: float64Ptr(2.5)}
	status := ms.SetDeltaJSON(input)
	if status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}
	if metric, exist := ms.Metrics["test_gauge"]; !exist || *metric.Value != 2.5 {
		t.Fatalf("Expected gauge value 2.5, got %v", metric.Value)
	}
}

func TestSetGaugeJSONBadRequest(t *testing.T) {
	ms := NewMetricsStorage()
	input := Metric{ID: "test_gauge", MType: "gauge", Value: nil}
	status := ms.SetDeltaJSON(input)
	if status != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestGetMetric(t *testing.T) {
	ms := NewMetricsStorage()
	ms.SetDeltaRaw("test_gauge", "gauge", 1.5)
	result := ms.GetMetricRaw("test_gauge", "gauge")
	if result != "1.5" {
		t.Fatalf("Expected gauge value string '1.5', got %s", result)
	}
}

func TestGetMetricNotFound(t *testing.T) {
	ms := NewMetricsStorage()
	result := ms.GetMetricRaw("non_existing", "gauge")
	if result != "" {
		t.Fatalf("Expected empty string for non-existing metric, got %s", result)
	}
}

func TestGetMetricJSON(t *testing.T) {
	ms := NewMetricsStorage()
	ms.SetDeltaRaw("test_gauge", "gauge", 1.5)
	input := Metric{ID: "test_gauge"}
	result, status := ms.GetMetricJSON(input)
	if status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}
	if result.ID != "test_gauge" || *result.Value != 1.5 {
		t.Fatalf("Expected result metric with ID 'test_gauge' and value 1.5, got %+v", result)
	}
}

func TestGetAllMetrics(t *testing.T) {
	ms := NewMetricsStorage()
	ms.SetDeltaRaw("test_gauge", "gauge", 1.5)
	allMetrics := ms.GetAllMetrics()
	if len(allMetrics.Metrics) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(allMetrics.Metrics))
	}
}

func TestSetMetricList(t *testing.T) {
	ms := NewMetricsStorage()
	list := []Metric{
		{ID: "gauge1", MType: "gauge", Value: float64Ptr(1.0)},
		{ID: "counter1", MType: "counter", Delta: int64Ptr(10)},
	}
	status := ms.SetListJSON(list)
	if status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}
	if _, exist := ms.Metrics["gauge1"]; !exist {
		t.Fatal("Expected gauge1 to exist in metrics")
	}
	if _, exist := ms.Metrics["counter1"]; !exist {
		t.Fatal("Expected counter1 to exist in metrics")
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
