package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/handlers"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/middleware"
)

func NewRouter(storage db.Storage) chi.Router {

	updateRow := handlers.NewUpdateMeticRawHandler(storage)
	updateJSON := handlers.NewUpdateMetricHandlerJSON(storage)
	updateList := handlers.NewUpdateListJSONHandler(storage)
	getJSON := handlers.NewGetMetricJSONHandler(storage)

	getRaw := handlers.NewGetMeticRawHandler(storage)
	getAllHTML := handlers.NewGetMetricListHandler(storage)
	getAllJSON := handlers.NewGetMetricListJSONHandler(storage)

	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.ZapMiddleware)

	r.Post("/update/{type}/{name}/{value}", updateRow.UpdateMetric)
	r.Post("/update/", updateJSON.UpdateJSON)
	r.Post("/updates/", updateList.UpdateListJSON)
	r.Post("/value/", getJSON.GetMetricHandlerJSON)

	r.Get("/value/{type}/{name}", getRaw.GetMetricHandler)
	r.Get("/", getAllHTML.GetAllMetrics)
	r.Get("/json", getAllJSON.GetAllMetricsJSON)
	r.Get("/ping", handlers.DatabaseHandler(storage))
  
	return r
}
