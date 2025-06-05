package handlers

import (
	"context"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"time"
)

type GetMeticListJSONStorage interface {
	GetAllMetrics(ctx context.Context) *models.MemStorage
}
type GetMetricListJSONHandler struct {
	storage GetMeticListStorage
}

func NewGetMetricListJSONHandler(storage GetMeticListStorage) *GetMetricListJSONHandler {
	return &GetMetricListJSONHandler{storage: storage}
}

// GetAllMetricsJSON return metrics in JSON
func (handler *GetMetricListJSONHandler) GetAllMetricsJSON(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	metrics := handler.storage
	m := metrics.GetAllMetrics(ctx)
	resp, err := json.Marshal(&m)
	if err != nil {
		log.Fatalf("Response: %v \n Error while writing response: %v", resp, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Fatalf("Error wgile write body : %v", err)
	}

}
