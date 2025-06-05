package handlers

import (
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllMetrics(t *testing.T) {
	storage := models.NewMemStorage()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler := NewGetMetricListHandler(storage)

	handler.GetAllMetrics(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
