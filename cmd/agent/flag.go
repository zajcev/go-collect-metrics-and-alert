package main

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

type Flags struct {
	ServerAddress  string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

type config struct {
	Home string `env:"HOME"`
}

func NewParseFlags() (*Flags, error) {
	var f Flags
	err := env.Parse(&f)
	if err != nil {
		return nil, err
	}
	flag.StringVar(&f.ServerAddress, "a", f.ServerAddress, "address and port destination server [ example usage : -a http://localhost:8080 ]")
	flag.IntVar(&f.ReportInterval, "r", f.ReportInterval, "set the timeout for reports [ example usage : -r 2 ]")
	flag.IntVar(&f.PollInterval, "p", f.PollInterval, "set the timeout for collect metrics [ example usage : -p 10s ]")
	flag.Parse()

	return &f, nil
}
