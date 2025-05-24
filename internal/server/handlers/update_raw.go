package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

// UpdateMetricHandler add or update metric value by raw value parsed from URI
func UpdateMetricHandler(metrics *models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		mname := chi.URLParam(r, "name")
		mtype := chi.URLParam(r, "type")
		mvalue := chi.URLParam(r, "value")
		if mtype == constants.Gauge {
			v, err := strconv.ParseFloat(mvalue, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				if config.GetDBHost() != "" {
					ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
					defer cancel()
					db.SetValueRaw(ctx, mname, mtype, v)
				} else {
					metrics.SetDeltaRaw(mname, mtype, v)
					syncWriter(metrics)
				}
			}
		} else if mtype == constants.Counter {
			v, err := strconv.ParseInt(mvalue, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				if config.GetDBHost() != "" {
					ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
					defer cancel()
					db.SetDeltaRaw(ctx, mname, mtype, v)
				} else {
					metrics.SetValueRaw(mname, mtype, v)
					syncWriter(metrics)
				}
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		err := r.Body.Close()
		if err != nil {
			log.Fatalf("Error while close body: %v", err)
		}
	}
}
