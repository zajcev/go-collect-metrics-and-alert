package main

import (
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MemStorage struct {
	Alloc       float64
	BuckHashSys float64
	RandomValue float64
	PollCount   int64
}

// MetricJSON — структура для JSON
type MetricJSON struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func Test_monitor(t *testing.T) {
	tests := []struct {
		name    string
		metric  model.Metrics
		wantErr bool
	}{
		{
			name:    "test",
			metric:  model.Metrics{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listeners.NewMonitor()
		})
	}
}

func TestNewReporter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got: %s", r.Header.Get("Content-Type"))
		}
		var mj MetricJSON
		err := json.NewDecoder(r.Body).Decode(&mj)
		if err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	listeners.NewReporter(ts.URL)
}
