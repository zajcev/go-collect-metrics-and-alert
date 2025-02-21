package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/handlers"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/logging"
	"log"
	"net/http"
)

func Router() chi.Router {
	r := chi.NewRouter()
	//r.Use(compress.GzipMiddleware)
	r.Use(logging.ZapMiddleware)
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler)
	r.Post("/update/", handlers.UpdateMetricHandlerJSON)
	r.Post("/value/", handlers.GetMetricHandlerJSON)
	r.Get("/value/{type}/{name}", handlers.GetMetricHandler)
	r.Get("/", handlers.GetAllMetrics)
	r.Get("/json", handlers.GetAllMetricsJSON)
	return r
}

func main() {
	log.Fatal(http.ListenAndServe(config.ParseFlags(), Router()))
}
