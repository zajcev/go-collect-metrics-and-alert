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

func TestGetMetricHandlerJSON(t *testing.T) {
	storage := models.NewMemStorage()
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
			handler := NewGetMetricJSONHandler(storage)

			handler.GetMetricHandlerJSON(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
			}
		})
	}
}

func ExampleNewGetMetricJSONHandler() {
	storage := models.NewMemStorage()
	value := 3.14
	ctx := context.Background()
	storage.SetValueRaw(ctx, "testGauge", "gauge", value)

	jsonBody := `{"id":"testGauge", "type": "gauge"}`
	reqBody := bytes.NewReader([]byte(jsonBody))

	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := NewGetMetricJSONHandler(storage)
	handler.GetMetricHandlerJSON(w, req)

	resp := w.Result()
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatalf("Error while close result body: %v", err)
		}
	}()
	body, _ := io.ReadAll(resp.Body)

	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		log.Fatalf("Error while formating json : %v", err)
	}
	fmt.Println(prettyJSON.String())

	// Output:
	//{
	//   "value": 3.14,
	//   "id": "testGauge",
	//   "type": "gauge"
	//}
}
