package main

import (
	"context"
	"github.com/jasonlvhit/gocron"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/handlers"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/routes"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	storage := models.NewMetricsStorage()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	err := config.NewConfig()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	if config.GetDBHost() == "" {
		if config.GetRestore() {
			storage = handlers.RestoreMetricStorage(config.GetFilePath())
		}
		if config.GetStoreInterval() > 0 {
			go func() {
				err = startScheduler(convert.GetUint(config.GetStoreInterval()), config.GetFilePath())
				if err != nil {
					log.Printf("Error with startScheduler: %v\n", err)
				}
			}()
		}
	} else {
		ctx := context.Background()
		db.Init(ctx, config.GetDBHost())
	}

	router := routes.NewRouter(storage)
	log.Fatal(http.ListenAndServe(config.GetAddress(), router))
}

func startScheduler(interval uint64, filePath string) error {
	scheduler := gocron.NewScheduler()
	err := scheduler.Every(interval).Seconds().Do(handlers.SaveMetricStorage, filePath)
	if err != nil {
		return err
	}
	scheduler.Start()
	return nil
}
