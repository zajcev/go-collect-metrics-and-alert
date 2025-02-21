package main

import (
	"compress/gzip"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Content-Encoding") != "gzip" {
			t.Errorf("Expected Content-Encoding: gzip, got %s", r.Header.Get("Content-Encoding"))
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			t.Fatalf("Error creating gzip reader: %v", err)
		}
		defer gz.Close()

		body, err := io.ReadAll(gz)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}
		var mj model.MetricJSON
		if err := json.Unmarshal(body, &mj); err != nil {
			t.Fatalf("Error unmarshalling JSON: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	listeners.NewReporter(server.URL + "/update")
}
