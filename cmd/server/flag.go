package main

import (
	"flag"
	"os"
)

var listenAddress string

func parseFlags() {
	flag.StringVar(&listenAddress, "a", ":8080", "address and port to run server [ example usage : -a localhost:8080 ]")
	if serverAddress := os.Getenv("ADDRESS"); serverAddress != "" {
		listenAddress = serverAddress
	}
	flag.Parse()
}
