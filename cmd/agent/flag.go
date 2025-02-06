package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
)

type NewConfig struct {
	address        string `env:"ADDRESS"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	pollInterval   int    `env:"POLL_INTERVAL"`
}

func NewParseFlags() (map[string]interface{}, error) {
	var cfg NewConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	var flagAddress string
	var flagReportInterval int
	var flagPollInterval int
	flag.StringVar(&flagAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "interval between report calls")
	flag.IntVar(&flagPollInterval, "p", 2, "interval between polls")
	flag.Parse()

	params := map[string]struct {
		flagValue interface{}
		envValue  interface{}
	}{
		"ADDRESS":         {flagAddress, cfg.address},
		"REPORT_INTERVAL": {flagReportInterval, cfg.reportInterval},
		"POLL_INTERVAL":   {flagPollInterval, cfg.pollInterval},
	}

	result := make(map[string]interface{})
	for k, v := range params {
		if v.envValue != "" && v.envValue != 0 {
			result[k] = v.envValue
		} else {
			result[k] = v.flagValue
		}
	}
	for k, v := range result {
		fmt.Printf("Key: %v, Value: %v\n", k, v)
	}
	return result, nil
}
