package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jasonlvhit/gocron"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/handlers"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/middleware"
	"log"
	"net/http"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.ZapMiddleware)
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler)
	r.Post("/update/", handlers.UpdateMetricHandlerJSON)
	r.Post("/updates/", handlers.UpdateListMetricsJSON)
	r.Post("/value/", handlers.GetMetricHandlerJSON)
	r.Get("/value/{type}/{name}", handlers.GetMetricHandler)
	r.Get("/", handlers.GetAllMetrics)
	r.Get("/json", handlers.GetAllMetricsJSON)
	r.Get("/ping", handlers.DatabaseHandler)
	return r
}

func main() {
	err := config.NewConfig()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	if *config.GetDBHost() == "" {
		if *config.GetRestore() {
			handlers.RestoreMetricStorage(*config.GetFilePath())
		}
		if *config.GetStoreInterval() > 0 {
			go startScheduler(convert.GetUint(*config.GetStoreInterval()), *config.GetFilePath())
		}
	} else {
		db.Init(*config.GetDBHost())
	}

	log.Fatal(http.ListenAndServe(*config.GetAddress(), Router()))
}

func startScheduler(interval uint64, filePath string) {
	scheduler := gocron.NewScheduler()
	scheduler.Every(interval).Seconds().Do(handlers.SaveMetricStorage, filePath)
	scheduler.Start()
}
