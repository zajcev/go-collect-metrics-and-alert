package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/crypto"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/handlers"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/routes"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printLdFlags()
	configuration := config.NewConfig()
	err := configuration.Load()
	if err != nil {
		log.Fatalf("Error load config : %v", err)
	}
	var memstorage db.Storage
	memstorage = models.NewMemStorage()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	address := configuration.GetAddress()
	cryptoKey := configuration.GetCryptoKey()
	dbHost := configuration.GetDBHost()
	restore := configuration.GetRestore()
	filePath := configuration.GetFilePath()
	interval := configuration.GetStoreInterval()

	if cryptoKey != "" {
		err = crypto.GenKeyPair("tmp")
		if err != nil {
			log.Fatalf("Error generate rsa key pair : %v", err)
		}
	}

	if dbHost == "" {
		if restore {
			memstorage = handlers.RestoreMetricStorage(filePath)
		}
		if configuration.GetStoreInterval() > 0 {
			go func() {
				if err = handlers.SaveMetricStorageSchedule(interval, filePath, memstorage); err != nil {
					log.Printf("Error save metrics to file : %v", err)
				}
			}()
		}
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		pool, poolErr := pgxpool.New(ctx, dbHost)
		if poolErr != nil {
			log.Fatalf("Error while create new PgxPool : %v", poolErr)
		}
		db.Migration(dbHost, "internal/server/db/scripts/")
		memstorage, err = db.NewDatabaseStorage(pool)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := routes.NewRouter(memstorage, configuration)
	go func() {
		log.Fatal(http.ListenAndServe(address, router))
	}()

	sig := <-sigChan
	log.Printf("Received signal %v, shutting down gracefully...", sig)
	time.Sleep(3 * time.Second)
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
