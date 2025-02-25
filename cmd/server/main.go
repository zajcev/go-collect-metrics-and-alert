package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jasonlvhit/gocron"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/compress"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/handlers"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/logging"
	"log"
	"net/http"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Use(compress.GzipMiddleware)
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
	err := config.ParseFlags()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	env := config.GetFlags()
	if env.Restore {
		handlers.RestoreMetricStorage(env.FilePath)
	}
	if env.StoreInterval > 0 {
		go startScheduler(convert.GetUint(env.StoreInterval), env.FilePath)
	}

	log.Fatal(http.ListenAndServe(env.Address, Router()))
}

func startScheduler(interval uint64, filePath string) {
	scheduler := gocron.NewScheduler()
	scheduler.Every(interval).Seconds().Do(handlers.SaveMetricStorage, filePath)
	scheduler.Start()
}
