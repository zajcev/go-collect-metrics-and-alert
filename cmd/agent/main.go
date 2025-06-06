package main

import (
	"context"
	"fmt"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
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
	go func() {
		log.Println(http.ListenAndServe("localhost:5050", nil))
	}()

	err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error while config initialization: %v", err)
	}

	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err = listeners.NewMonitor(ctx, config.GetPollInterval()); err != nil {
			errChan <- fmt.Errorf("monitor start failed: %w", err)
		}
	}()

	go func() {
		if err = listeners.NewReporter(ctx, config.GetReportInterval(), "http://"+config.GetAddress()+"/updates/"); err != nil {
			errChan <- fmt.Errorf("HTTP server failed: %w", err)
		}
	}()

	go func() {
		if err = listeners.AdditionalMetrics(ctx, config.GetPollInterval()); err != nil {
			errChan <- fmt.Errorf("HTTP server failed: %w", err)
		}
	}()

	err = <-errChan
	log.Printf("Fatal error: %v", err)
	cancel()
  
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
