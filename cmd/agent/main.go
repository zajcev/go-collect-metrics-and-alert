package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"log"
)

func main() {
	f, err := config.NewParseFlags()
	var a = convert.GetString(f["ADDRESS"])
	var pi = convert.GetUint(f["POLL_INTERVAL"])
	var ri = convert.GetUint(f["REPORT_INTERVAL"])
	if err != nil {
		log.Fatal(err)
	}

	scheduler := gocron.NewScheduler()
	monitorDone := make(chan bool, 1)
	reporterDone := make(chan bool, 1)
	err = scheduler.Every(pi).Second().Do(func() {
		<-monitorDone
		listeners.NewMonitor()
		monitorDone <- true
	})
	if err != nil {
		log.Fatal(err)
	}
	err = scheduler.Every(ri).Second().Do(func(url string) {
		<-reporterDone
		listeners.NewReporter(url)
		reporterDone <- true
	}, "http://"+a+"/updates/")
	if err != nil {
		log.Fatal(err)
	}
	monitorDone <- true
	reporterDone <- true
	sc := scheduler.Start()
	<-sc
}
