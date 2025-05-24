package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"time"
)

// GetMetricHandler return metric value by name and type in raw value
func GetMetricHandler(metrics *models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mname := chi.URLParam(r, "name")
		mtype := chi.URLParam(r, "type")
		var value string
		if config.GetDBHost() != "" {
			ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
			defer cancel()
			value = convert.GetString(db.GetMetricRaw(ctx, mname, mtype))
		} else {
			value = metrics.GetMetricRaw(mname, mtype)
		}
		if value != "" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/text")
			res, err := w.Write([]byte(value))
			if err != nil {
				log.Fatalf("Response: %v \n Error while writing response: %v", res, err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		err := r.Body.Close()
		if err != nil {
			return
		}
	}
}
