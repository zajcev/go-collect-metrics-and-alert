package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/handlers"
	"log"
	"net/http"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler)
	r.Get("/value/{type}/{name}", handlers.GetMetricHandler)
	r.Get("/", handlers.GetAllMetrics)
	return r
}

func main() {
	config.ParseFlags()
	log.Fatal(http.ListenAndServe(config.ListenAddress, Router()))
}
