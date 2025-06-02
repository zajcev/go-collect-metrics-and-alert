package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printLdFlags()
	var memstorage db.Storage
	memstorage = models.NewMemStorage()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	err := config.NewConfig()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	if config.GetDBHost() == "" {
		if config.GetRestore() {
			memstorage = handlers.RestoreMetricStorage(config.GetFilePath())
		}
		if config.GetStoreInterval() > 0 {
			go func() {
				err = startScheduler(convert.GetUint(config.GetStoreInterval()), config.GetFilePath(), memstorage)
				if err != nil {
					log.Printf("Error with startScheduler: %v\n", err)
				}
			}()
		}
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		pool, poolErr := pgxpool.New(ctx, config.GetDBHost())
		if err != nil {
			log.Fatalf("Error while create new PgxPool : %v", poolErr)
		}
		db.Migration(config.GetDBHost(), "internal/server/db/scripts/")
		memstorage, err = db.NewDatabaseStorage(pool)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := routes.NewRouter(memstorage)
	log.Fatal(http.ListenAndServe(config.GetAddress(), router))
}

func startScheduler(interval uint64, filePath string, storage db.Storage) error {
	scheduler := gocron.NewScheduler()
	err := scheduler.Every(interval).Seconds().Do(handlers.SaveMetricStorage, filePath, storage)
	if err != nil {
		return err
	}
	scheduler.Start()
	return nil
}

func printLdFlags() {
	getValue := func(v string) string {
		if v == "" {
			return "N/A"
		}
		return v
	}
	fmt.Printf("Build version: %s\n", getValue(buildVersion))
	fmt.Printf("Build date: %s\n", getValue(buildDate))
	fmt.Printf("Build commit: %s\n", getValue(buildCommit))
}
