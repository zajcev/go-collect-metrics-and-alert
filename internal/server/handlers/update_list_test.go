package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpdateListMetricsJSON(t *testing.T) {
	testMemStorage := models.NewMemStorage()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	configuration := config.NewConfig()
	err := configuration.Load()
	if err != nil {
		log.Fatalf("Error init config : %v", err)
	}
	tests := []struct {
		name       string
		hashSHA256 string
		body       []byte
		wantCode   int
	}{
		{
			name:       "Valid JSON",
			body:       []byte(`[{"id": "metric1", "type": "counter", "delta": 10}]`),
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
			handler := NewUpdateListJSONHandler(testMemStorage, configuration)

			handler.UpdateListJSON(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}
		})
	}
}

func BenchmarkUpdateListMetricsJSON(b *testing.B) {
	storage := models.NewMemStorage()
	listSizes := []int{10, 100, 1000, 5000}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	configuration := config.NewConfig()
	err := configuration.Load()
	if err != nil {
		log.Fatalf("Error load config : %v", err)
	}

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
				handler := NewUpdateMetricHandlerJSON(storage, configuration)

				handler.UpdateJSON(rr, req)
			}
		})
	}
}

func ExampleTestUpdateListMetricsJSON() {
	storage := models.NewMemStorage()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	configuration := config.NewConfig()
	err := configuration.Load()
	if err != nil {
		log.Fatalf("Error load config : %v", err)
	}
	handler := NewUpdateListJSONHandler(storage, configuration)
	jsonBody := `[{"id":"testGauge", "type": "gauge", "value": 42}, {"id":"testCounter", "type": "counter", "delta": 42}]`
	reqBody := bytes.NewReader([]byte(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.UpdateListJSON(w, req)
	list := storage.GetAllMetrics(context.Background())
	jsonData, err := json.Marshal(list)
	if err != nil {
		fmt.Println("Error marshaling metrics:", err)
		return
	}
	fmt.Println(string(jsonData))
	// Output:
	//{"Storage":{"testCounter":{"delta":42,"id":"testCounter","type":"counter"},"testGauge":{"value":42,"id":"testGauge","type":"gauge"}}}
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
