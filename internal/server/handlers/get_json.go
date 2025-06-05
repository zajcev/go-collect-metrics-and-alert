package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"time"
)

type GetMeticJSONStorage interface {
	GetMetricJSON(ctx context.Context, input models.Metric) (models.Metric, int)
}
type GetMetricJSONHandler struct {
	storage GetMeticJSONStorage
}

func NewGetMetricJSONHandler(storage GetMeticJSONStorage) *GetMetricJSONHandler {
	return &GetMetricJSONHandler{storage: storage}
}

// GetMetricHandlerJSON return metric value in JSON Struct
func (handler *GetMetricJSONHandler) GetMetricHandlerJSON(w http.ResponseWriter, r *http.Request) {
	var m models.Metric
	var buf bytes.Buffer
	var code int
	metrics := handler.storage
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	m, code = metrics.GetMetricJSON(ctx, m)
	resp, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(resp)
	if err != nil {
		log.Fatalf("Error write response : %v", err)
	}
}
