package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func NewParseFlags() (map[string]any, error) {
	var NewConfig Config
	if err := env.Parse(&NewConfig); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}
	var flagAddress string
	var flagReportInterval int
	var flagPollInterval int
	flag.StringVar(&flagAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 2, "interval between report calls")
	flag.IntVar(&flagPollInterval, "p", 2, "interval between polls")
	flag.Parse()

	params := map[string]struct {
		flagValue any
		envValue  any
	}{
		"ADDRESS":         {flagAddress, NewConfig.Address},
		"REPORT_INTERVAL": {flagReportInterval, NewConfig.ReportInterval},
		"POLL_INTERVAL":   {flagPollInterval, NewConfig.PollInterval},
	}

	result := make(map[string]any)
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
