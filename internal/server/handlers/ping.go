package handlers

import (
	"context"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"net/http"
	"time"
)


func DatabaseHandler(metrics db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()
		err := metrics.Ping(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
