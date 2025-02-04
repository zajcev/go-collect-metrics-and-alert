package main

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Flags struct {
	ADDRESS         string
	REPORT_INTERVAL int
	POLL_INTERVAL   int
}

var serverAddress string
var reportInterval int
var pollInterval int

func parseFlags() {
	var f Flags
	err := envconfig.Process("", &f)
	if err != nil {
		log.Fatal(err.Error())
	}
	flag.StringVar(&serverAddress, "a", "localhost:8080", "address and port destination server [ example usage : -a http://localhost:8080 ]")
	flag.IntVar(&reportInterval, "r", 10, "set the timeout for reports [ example usage : -r 2 ]")
	flag.IntVar(&pollInterval, "p", 2, "set the timeout for collect metrics [ example usage : -p 10s ]")

	if serverAddressOs := f.ADDRESS; serverAddressOs != "" {
		f.ADDRESS = serverAddressOs
	}
	if f.REPORT_INTERVAL != 0 {
		reportInterval = f.REPORT_INTERVAL
	}
	if f.POLL_INTERVAL != 0 {
		pollInterval = f.POLL_INTERVAL
	}
	flag.Parse()
}
