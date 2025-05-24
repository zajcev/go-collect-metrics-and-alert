package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
)

var flags Flags

type Flags struct {
	Address        string `env:"ADDRESS"`
	HashKey        string `env:"KEY"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

// NewConfig parses the command-line flags and environment variables.
func NewConfig() error {
	flag.StringVar(&flags.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flags.HashKey, "k", "12h5b12b521b", "key for sha256sum")
	flag.IntVar(&flags.ReportInterval, "r", 1, "interval between report calls")
	flag.IntVar(&flags.PollInterval, "p", 1, "interval between polls")
	flag.IntVar(&flags.RateLimit, "l", 2, "request rate limiter")
	flag.Parse()
	if err := env.Parse(&flags); err != nil {
		log.Printf("%+v", err)
		return err
	}
	return nil
}

// GetAddress returns the address and port to run the server.
func GetAddress() string {
	return convert.GetString(&flags.Address)
}

// GetReportInterval returns the interval between report calls.
func GetReportInterval() int {
	return flags.ReportInterval
}

// GetPollInterval returns the interval between polls.
func GetPollInterval() int {
	return flags.PollInterval
}

// GetHashKey returns the key for sha256sum.
func GetHashKey() string { return convert.GetString(&flags.HashKey) }

// GetRateLimit returns interval for Reporter
func GetRateLimit() int { return flags.RateLimit }
