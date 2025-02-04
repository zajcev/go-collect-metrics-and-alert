package main

import (
	"flag"
)

var serverAddress string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&serverAddress, "a", "http://localhost:8080", "address and port destination server [ example usage : -a http://localhost:8080 ]")
	flag.IntVar(&reportInterval, "r", 10, "set the timeout for reports [ example usage : -r 2 ]")
	flag.IntVar(&pollInterval, "p", 2, "set the timeout for collect metrics [ example usage : -p 10s ]")
	flag.Parse()
}
