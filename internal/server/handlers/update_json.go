package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"time"
)

type UpdateMetricHandlerJSONStorage interface {
	SetValueJSON(ctx context.Context, input models.Metric) int
	SetDeltaJSON(ctx context.Context, input models.Metric) int
	GetAllMetrics(ctx context.Context) *models.MemStorage
}
type UpdateMetricHandlerJSON struct {
	storage UpdateMetricHandlerJSONStorage
}

func NewUpdateMetricHandlerJSON(storage UpdateMetricHandlerJSONStorage) *UpdateMetricHandlerJSON {
	return &UpdateMetricHandlerJSON{storage: storage}
}

// UpdateJSON add or update metric value from JSON body
func (handler *UpdateMetricHandlerJSON) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()
		var m models.Metric
		var buf bytes.Buffer
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
		if m.MType == constants.Gauge {
			metrics.SetValueJSON(ctx, m)
		} else if m.MType == constants.Counter {
			metrics.SetDeltaJSON(ctx, m)
		}
		resp, err := json.Marshal(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			log.Fatalf("Error wgile write body : %v", err)
		}
		if config.GetDBHost() == "" {
			SaveMetricStorage(config.GetFilePath(), metrics)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
