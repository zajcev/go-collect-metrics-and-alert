package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
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
			ctx, cancel := context.WithTimeout(context.Background(), 3)
			defer cancel()
			listeners.NewMonitor(ctx, 2)
		})
	}
}

func TestNewReporter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3)
	defer cancel()
	listeners.NewReporter(ctx, 2, "http://localhost:8080/update")
}

func BenchmarkMonitor(*testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), 3)
	defer cancel()
	listeners.NewMonitor(ctx, 2)
}

func BenchmarkSend(*testing.B) {
	var list []model.MetricJSON
	delta := convert.GetInt64(123)
	list = append(list, model.MetricJSON{ID: "PollCount", MType: "counter", Delta: &delta, Value: nil})
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(list)
		var buf bytes.Buffer
		g := gzip.NewWriter(&buf)
		g.Write(resp)
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}))

	defer testServer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3)
	listeners.NewReporter(ctx, 2, testServer.URL)
	cancel()
}
