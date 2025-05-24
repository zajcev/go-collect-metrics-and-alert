package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/handlers"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/middleware"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
)

func NewRouter(mertics *models.MemStorage) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.ZapMiddleware)
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler(mertics))
	r.Post("/update/", handlers.UpdateMetricHandlerJSON(mertics))
	r.Post("/updates/", handlers.UpdateListMetricsJSON(mertics))
	r.Post("/value/", handlers.GetMetricHandlerJSON(mertics))
	r.Get("/value/{type}/{name}", handlers.GetMetricHandler(mertics))
	r.Get("/", handlers.GetAllMetrics(mertics))
	r.Get("/json", handlers.GetAllMetricsJSON(mertics))
	r.Get("/ping", handlers.DatabaseHandler)
	return r
}
