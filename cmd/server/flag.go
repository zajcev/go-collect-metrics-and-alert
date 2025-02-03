package main

import (
	"flag"
)

var listenAddress string

func parseFlags() {
	flag.StringVar(&listenAddress, "a", ":8080", "address and port to run server [ example usage : -a localhost:8080 ]")
	flag.Parse()
}
