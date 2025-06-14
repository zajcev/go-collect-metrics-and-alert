package handlers

import (
	"context"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"io"
	"log"
	"net/http"
	"time"
)

type UpdateListMetricsJSONStorage interface {
	SetListJSON(ctx context.Context, list []models.Metric) int
	GetAllMetrics(ctx context.Context) *models.MemStorage
}
type UpdateListMetricsJSONConfig interface {
	GetDBHost() string
	GetFilePath() string
	GetHashKey() string
}
type UpdateListJSONHandler struct {
	storage UpdateListMetricsJSONStorage
	config  UpdateListMetricsJSONConfig
}

func NewUpdateListJSONHandler(storage UpdateListMetricsJSONStorage, config UpdateListMetricsJSONConfig) *UpdateListJSONHandler {
	return &UpdateListJSONHandler{
		storage: storage,
		config:  config,
	}
}

// UpdateListJSON add an or update list of metrics from JSON body
func (handler *UpdateListJSONHandler) UpdateListJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()
		var list []models.Metric
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error read body : %v", err)
		}
		metrics := handler.storage
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if handler.config.GetHashKey() != "" {
			if !checkSHA256Hash(body, handler.config.GetHashKey(), r.Header.Get("HashSHA256")) {
				http.Error(w, "Mismatch sha256sum", http.StatusBadRequest)
				return
			}
		}
		if err = json.Unmarshal(body, &list); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		metrics.SetListJSON(ctx, list)
		resp, err := json.Marshal(&list)
		if handler.config.GetHashKey() != "" {
			w.Header().Set("HashSHA256", calculateSHA256Hash(resp, handler.config.GetHashKey()))
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = SaveMetricStorageOnce(handler.config.GetFilePath(), metrics)
		if err != nil {
			log.Printf("Error save metrics to file : %v", err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
