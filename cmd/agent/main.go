package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/listeners"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/cast"
	"log"
)

func main() {
	f, err := config.NewParseFlags()
	var a = cast.GetString(f["ADDRESS"])
	var pi = cast.GetUint(f["POLL_INTERVAL"])
	var ri = cast.GetUint(f["REPORT_INTERVAL"])
	if err != nil {
		log.Fatal(err)
	}
	gocron.Every(pi).Second().Do(listeners.NewMonitor)
	gocron.Every(ri).Second().Do(listeners.NewReporter, "http://"+a+"/update/")
	<-gocron.Start()
}
