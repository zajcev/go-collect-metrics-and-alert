package config

import (
	"flag"
	"os"
)

var ListenAddress string

func ParseFlags() {
	flag.StringVar(&ListenAddress, "a", ":8080", "address and port to run server [ example usage : -a localhost:8080 ]")
	if serverAddress := os.Getenv("ADDRESS"); serverAddress != "" {
		ListenAddress = serverAddress
	}
	flag.Parse()
}
