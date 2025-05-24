package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"time"
)

// UpdateMetricHandlerJSON add or update metric value from JSON body
func UpdateMetricHandlerJSON(metrics *models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			var m models.Metric
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err = json.Unmarshal(buf.Bytes(), &m); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else if config.GetDBHost() != "" {
				ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
				defer cancel()
				if m.MType == constants.Gauge {
					db.SetValueJSON(ctx, m)
				} else if m.MType == constants.Counter {
					db.SetDeltaJSON(ctx, m)
				}
			} else {
				if m.MType == constants.Gauge {
					metrics.SetDeltaJSON(m)
				} else if m.MType == constants.Counter {
					metrics.SetValueJSON(m)
				}
				syncWriter(metrics)
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
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
