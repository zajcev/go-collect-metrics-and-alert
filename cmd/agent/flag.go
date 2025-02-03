package main

import (
	"flag"
	"time"
)

var serverAddress string
var reportInterval time.Duration
var pollInterval time.Duration

func parseFlags() {
	flag.StringVar(&serverAddress, "a", "http://localhost:8080", "address and port destination server [ example usage : -a http://localhost:8080 ]")
	flag.DurationVar(&reportInterval, "r", 10*time.Second, "set the timeout for reports [ example usage : -r 2s ]")
	flag.DurationVar(&pollInterval, "p", 2*time.Second, "set the timeout for collect metrics [ example usage : -p 10s ]")
	flag.Parse()
}
