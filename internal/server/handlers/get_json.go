package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"time"
)

// GetMetricHandlerJSON return metric value in JSON Struct
func GetMetricHandlerJSON(metrics *models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m models.Metric
		var buf bytes.Buffer
		var code int
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if config.GetDBHost() != "" {
			ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
			defer cancel()
			m, code = db.GetMetricJSON(ctx, m)
		} else {
			m, code = metrics.GetMetricJSON(m)
		}
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
}
