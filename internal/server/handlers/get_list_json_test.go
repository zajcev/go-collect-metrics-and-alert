package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllMetricsJSON(t *testing.T) {
	storage := models.NewMemStorage()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler := NewGetMetricListJSONHandler(storage)

	handler.GetAllMetricsJSON(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func ExampleGetMetricListJSONHandler() {
	storage := models.NewMemStorage()
	delta := int64(42)
	value := 3.14
	ctx := context.Background()
	storage.SetValueRaw(ctx, "testGauge", "gauge", value)
	storage.SetDeltaRaw(ctx, "testCounter", "counter", delta)

	req := httptest.NewRequest("GET", "/json/all", nil)
	w := httptest.NewRecorder()

	handler := NewGetMetricListJSONHandler(storage)
	handler.GetAllMetricsJSON(w, req)

	resp := w.Result()
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatalf("Error while close result body: %v", err)
		}
	}()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Compact JSON:", string(body))
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		log.Fatalf("Error while formating json : %v", err)
	}

	// Output:
	// Compact JSON: {"Storage":{"testCounter":{"delta":42,"id":"testCounter","type":"counter"},"testGauge":{"value":3.14,"id":"testGauge","type":"gauge"}}}
}
