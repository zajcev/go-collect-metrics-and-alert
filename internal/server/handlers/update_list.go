package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"net/http"
	"time"
)

// UpdateListMetricsJSON add an or update list of metrics from JSON body
func UpdateListMetricsJSON(metrics *models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			var list []models.Metric
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if config.GetHashKey() != "" {
				if !checkSHA256Hash(buf.Bytes(), config.GetHashKey(), r.Header.Get("HashSHA256")) {
					http.Error(w, "Mismatch sha256sum", http.StatusBadRequest)
					return
				}
			}
			if err = json.Unmarshal(buf.Bytes(), &list); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else if config.GetDBHost() != "" {
				ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
				defer cancel()
				db.SetListJSON(ctx, list)
			} else {
				metrics.SetListJSON(list)
				syncWriter(metrics)
			}
			resp, err := json.Marshal(&list)
			if config.GetHashKey() != "" {
				w.Header().Set("HashSHA256", calculateSHA256Hash(resp, config.GetHashKey()))
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
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
