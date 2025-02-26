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
	gocron.Every(pi).Second().Do(listeners.NewMonitor)
	gocron.Every(ri).Second().Do(listeners.NewReporter, "http://"+a+"/update/")
	<-gocron.Start()
}
