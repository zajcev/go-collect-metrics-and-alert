package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateListMetricsJSON(t *testing.T) {
	testMemStorage := models.NewMetricsStorage()
	tests := []struct {
		name       string
		body       []byte
		hashSHA256 string
		wantCode   int
	}{
		{
			name:       "Valid JSON",
			body:       []byte(`[{"id": "metric1", "type": "counter", "value": 10}]`),
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
			handler := UpdateListMetricsJSON(testMemStorage)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}
		})
	}
}

func TestGetMetricHandlerJSON(t *testing.T) {
	testMemStorage := models.NewMetricsStorage()
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
			handler := GetMetricHandlerJSON(testMemStorage)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}
		})
	}
}

func TestGetAllMetricsJSON(t *testing.T) {
	testMemStorage := models.NewMetricsStorage()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler := GetAllMetricsJSON(testMemStorage)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetAllMetrics(t *testing.T) {
	testMemStorage := models.NewMetricsStorage()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler := GetAllMetrics(testMemStorage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func BenchmarkUpdateListMetricsJSON(b *testing.B) {
	testMemStorage := models.NewMetricsStorage()
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
				handler := UpdateListMetricsJSON(testMemStorage)

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

func ExampleUpdateMetricHandlerJSON() {
	testMemStorage := models.NewMetricsStorage()
	handler := UpdateMetricHandlerJSON(testMemStorage)
	jsonBody := `{"id":"test", "type": "gauge", "value": 42}`
	reqBody := bytes.NewReader([]byte(jsonBody))

	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler(w, req)

	m := models.Metric{ID: "test", MType: "gauge"}
	result, status := testMemStorage.GetMetricJSON(m)

	fmt.Printf("Metric details:\n"+
		"ID:    %s\n"+
		"Type:  %s\n"+
		"Delta: %v\n"+
		"Value: %v\n"+
		"Status: %d\n",
		result.ID, result.MType, result.Delta, *result.Value, status)
	//Output:
	// Metric details:
	// ID:    test
	// Type:  gauge
	// Delta: <nil>
	// Value: 42
	// Status: 200
}

func ExampleTestUpdateListMetricsJSON() {
	testMemStorage := models.NewMetricsStorage()
	handler := UpdateListMetricsJSON(testMemStorage)
	jsonBody := `[{"id":"testGauge", "type": "gauge", "value": 42}, {"id":"testCounter", "type": "counter", "delta": 42}]`
	reqBody := bytes.NewReader([]byte(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler(w, req)
	list := testMemStorage.GetAllMetrics()
	jsonData, err := json.Marshal(list)
	if err != nil {
		fmt.Println("Error marshaling metrics:", err)
		return
	}
	fmt.Println(string(jsonData))
	// Output:
	//{"Metrics":{"testCounter":{"id":"testCounter","type":"counter","delta":42},"testGauge":{"id":"testGauge","type":"gauge","value":42}}}
}

func ExampleGetAllMetricsJSON() {
	testMemStorage := models.NewMetricsStorage()
	delta := int64(42)
	value := 3.14
	testMemStorage.SetGauge("testGauge", "gauge", value)
	testMemStorage.SetCounter("testCounter", "counter", delta)

	req := httptest.NewRequest("GET", "/json/all", nil)
	w := httptest.NewRecorder()

	handler := GetAllMetricsJSON(testMemStorage)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "  ")
	fmt.Println(prettyJSON.String())

	// Output:
	// {
	//   "Metrics": {
	//     "testCounter": {
	//       "id": "testCounter",
	//       "type": "counter",
	//       "delta": 42
	//     },
	//     "testGauge": {
	//       "id": "testGauge",
	//       "type": "gauge",
	//       "value": 3.14
	//     }
	//   }
	// }
}

func ExampleGetMetricHandlerJSON() {
	testMemStorage := models.NewMetricsStorage()
	value := 3.14
	testMemStorage.SetGauge("testGauge", "gauge", value)

	jsonBody := `{"id":"testGauge", "type": "gauge"}`
	reqBody := bytes.NewReader([]byte(jsonBody))

	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := GetMetricHandlerJSON(testMemStorage)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "  ")
	fmt.Println(prettyJSON.String())

	// Output:
	//{
	//   "id": "testGauge",
	//   "type": "gauge",
	//   "value": 3.14
	//}
}
