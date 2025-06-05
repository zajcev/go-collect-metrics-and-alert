package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UpdateMeticRawStorage interface {
	SetDeltaRaw(ctx context.Context, name string, metricType string, value int64) int
	SetValueRaw(ctx context.Context, name string, metricType string, value float64) int
}
type UpdateMeticRawHandler struct {
	storage UpdateMeticRawStorage
}

func NewUpdateMeticRawHandler(storage UpdateMeticRawStorage) *UpdateMeticRawHandler {
	return &UpdateMeticRawHandler{storage: storage}
}

// UpdateMetric add or update metric value by raw value parsed from URI
func (handler *UpdateMeticRawHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	metrics := handler.storage
	mname := chi.URLParam(r, "name")
	mtype := chi.URLParam(r, "type")
	mvalue := chi.URLParam(r, "value")
	if mtype == constants.Gauge {
		v, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		metrics.SetValueRaw(ctx, mname, mtype, v)
	} else if mtype == constants.Counter {
		v, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			metrics.SetDeltaRaw(ctx, mname, mtype, v)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	err := r.Body.Close()
	if err != nil {
		log.Fatalf("Error while close body: %v", err)
	}
}
