package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"net/http/httptest"
)

func ExampleUpdateMetricHandlerJSON() {
	storage := models.NewMemStorage()
	handler := NewUpdateMetricHandlerJSON(storage)
	jsonBody := `{"id":"test", "type": "gauge", "value": 42}`
	reqBody := bytes.NewReader([]byte(jsonBody))

	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.UpdateJSON(w, req)

	m := models.Metric{ID: "test", MType: "gauge"}
	ctx := context.Background()
	result, status := storage.GetMetricJSON(ctx, m)

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
