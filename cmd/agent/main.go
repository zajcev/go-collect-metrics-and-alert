package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
	"log"
)

func main() {
	err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error while config initialization: %v", err)
	}
	scheduler := gocron.NewScheduler()
	monitorDone := make(chan bool, 1)
	reporterDone := make(chan bool, 1)
	err = scheduler.Every(config.GetPollInterval()).Second().Do(func() {
		<-monitorDone
		listeners.NewMonitor()
		monitorDone <- true
	})
	if err != nil {
		log.Fatal(err)
	}
	err = scheduler.Every(config.GetReportInterval()).Second().Do(func(url string) {
		<-reporterDone
		listeners.NewReporter(url)
		reporterDone <- true
	}, "http://"+config.GetAddress()+"/updates/")
	if err != nil {
		log.Fatal(err)
	}
	monitorDone <- true
	reporterDone <- true
	sc := scheduler.Start()
	<-sc
}
