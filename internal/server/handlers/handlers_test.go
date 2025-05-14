package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
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
