package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
	"log"
)

func main() {
	f, err := config.NewParseFlags()
	var a = config.GetString(f["ADDRESS"])
	var pi = config.GetUint(f["POLL_INTERVAL"])
	var ri = config.GetUint(f["REPORT_INTERVAL"])
	if err != nil {
		log.Fatal(err)
	}
	gocron.Every(pi).Second().Do(listeners.NewMonitor)
	gocron.Every(ri).Second().Do(listeners.NewReporter, "http://"+a)
	<-gocron.Start()
}
