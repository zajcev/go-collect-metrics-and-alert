package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type GetMeticRawStorage interface {
	GetMetricRaw(ctx context.Context, name string, metricType string) any
}
type GetMeticRawHandler struct {
	storage GetMeticRawStorage
}

func NewGetMeticRawHandler(storage GetMeticRawStorage) *GetMeticRawHandler {
	return &GetMeticRawHandler{storage: storage}
}

// GetMetricHandler return metric value by name and type in raw value
func (handler *GetMeticRawHandler) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	mname := chi.URLParam(r, "name")
	mtype := chi.URLParam(r, "type")
	metrics := handler.storage
	value := metrics.GetMetricRaw(ctx, mname, mtype)
	if value != "" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		res, err := w.Write(value.([]byte))
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
