package config

import (
	"flag"
	"os"
)

func ParseFlags() string {
	var listenAddress string
	flag.StringVar(&listenAddress, "a", ":8080", "address and port to run server [ example usage : -a localhost:8080 ]")
	if serverAddress := os.Getenv("ADDRESS"); serverAddress != "" {
		listenAddress = serverAddress
	}
	flag.Parse()
	return listenAddress
}
