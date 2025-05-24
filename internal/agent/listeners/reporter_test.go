package listeners

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
)

func TestCalculateSHA256Hash(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
		data     []byte
	}{
		{"Basic test", "91d2330355770ae2a13eb43e62d9ed805aa140d4c7157a7cf69c170d1050fb6c", "key", []byte("test data")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSHA256Hash(tt.data, tt.key)
			if result != tt.expected {
				t.Errorf("Expected %s but got %s", tt.expected, result)
			}
		})
	}
}

func TestRetryFailure(t *testing.T) {
	client := &http.Client{}
	req := httptest.NewRequest("GET", "http://example.com", nil)

	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()

	req.URL, _ = url.Parse(ts.URL)

	resp, err := retry(client, req, 3)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	if resp != nil {
		err = resp.Body.Close()
		if err != nil {
			t.Fatalf("Error while body close %v", err)
		}
		t.Fatalf("Expected nil response, got %v", resp)
	}

}

func TestNewReporter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3)
	defer cancel()
	err := NewReporter(ctx, 2, "http://localhost:8080/update")
	if err != nil {
		t.Fatalf("Error while create reporter %v", err)
	}
}

func BenchmarkSend(b *testing.B) {
	listSizes := []int{10, 100, 500, 1000, 5000}

	for _, size := range listSizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			list := generateRandomMetrics(size)
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(list)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				var buf bytes.Buffer
				gz := gzip.NewWriter(&buf)
				if _, err = gz.Write(resp); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if err = gz.Close(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err = w.Write(buf.Bytes())
				if err != nil {
					b.Fatalf("Error while write body : %v", err)
				}
			}))
			defer testServer.Close()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				send(testServer.URL, &list)
			}
		})
	}
}

func generateRandomMetrics(count int) []model.MetricJSON {
	list := make([]model.MetricJSON, 0, count)

	metricNames := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	}

	for i := 0; i < count; i++ {
		var metric model.MetricJSON
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
