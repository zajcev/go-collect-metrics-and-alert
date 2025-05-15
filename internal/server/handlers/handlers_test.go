package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
)

func TestUpdateListMetricsJSON(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		hashSHA256 string
		wantCode   int
	}{
		{
			name:       "Valid JSON",
			body:       []byte(`[{"name": "metric1", "value": 10}]`),
			hashSHA256: "",
			wantCode:   http.StatusOK,
		},
		{
			name:       "Invalid JSON",
			body:       []byte(`invalid json`),
			hashSHA256: "",
			wantCode:   http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tt.body))
			if tt.hashSHA256 != "" {
				req.Header.Set("HashSHA256", tt.hashSHA256)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(UpdateListMetricsJSON)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}
		})
	}
}

func TestGetMetricHandlerJSON(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		wantCode int
	}{
		{
			name:     "Valid Metric",
			body:     []byte(`{"name": "metric1", "value": 10}`),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Invalid JSON",
			body:     []byte(`invalid json`),
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetMetricHandlerJSON)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}
		})
	}
}

func TestGetAllMetricsJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllMetricsJSON)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetAllMetrics(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllMetrics)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func BenchmarkUpdateListMetricsJSON(b *testing.B) {
	listSizes := []int{10, 100, 1000, 5000}

	for _, size := range listSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			list := generateRandomMetrics(size)
			marshal, err := json.Marshal(&list)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(marshal))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(UpdateListMetricsJSON)

				handler.ServeHTTP(rr, req)
			}
		})
	}
}

func generateRandomMetrics(count int) []models.Metric {
	list := make([]models.Metric, 0, count)

	metricNames := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	}

	for i := 0; i < count; i++ {
		var metric models.Metric
		metric.ID = metricNames[rand.Intn(len(metricNames))]

		if rand.Intn(2) == 0 {
			metric.MType = "gauge"
			val := rand.Float64() * 1000
			metric.Value = &val
		} else {
			metric.MType = "counter"
			delta := rand.Int63n(1000)
			metric.Delta = &delta
		}

		list = append(list, metric)
	}

	return list
}
