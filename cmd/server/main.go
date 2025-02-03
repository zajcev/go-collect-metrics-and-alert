package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

var metrics = map[string]*MemStorage{}

func Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", updateMetricHandler)
	r.Get("/value/{type}/{name}", getMetricHandler)
	r.Get("/", getAllMetrics)
	return r
}

func main() {
	parseFlags()
	log.Fatal(http.ListenAndServe(listenAddress, Router()))
}
