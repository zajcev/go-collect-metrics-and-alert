package main

import (
	"github.com/spf13/cast"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
	"log"
	"time"
)

func main() {
	f, err := config.NewParseFlags()
	var a = cast.ToString(f["ADDRESS"])
	var pi = cast.ToDuration(f["POLL_INTERVAL"])
	var ri = cast.ToDuration(f["REPORT_INTERVAL"])
	if err != nil {
		log.Fatal(err)
	}
	monitorTimer := time.NewTicker(pi * time.Second)
	reporterTimer := time.NewTicker(ri * time.Second)
	for {
		select {
		case <-monitorTimer.C:
			listeners.NewMonitor()
		case <-reporterTimer.C:
			listeners.NewReporter("http://" + a)
		}
	}
}
