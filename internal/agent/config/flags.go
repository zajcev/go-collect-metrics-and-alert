package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"log"
)

var flags Flags

type Flags struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	HashKey        string `env:"KEY"`
}

func NewConfig() error {
	flag.StringVar(&flags.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flags.HashKey, "k", "12h5b12b521b", "key for sha256sum")
	flag.IntVar(&flags.ReportInterval, "r", 2, "interval between report calls")
	flag.IntVar(&flags.PollInterval, "p", 1, "interval between polls")
	flag.Parse()
	if err := env.Parse(&flags); err != nil {
		log.Printf("%+v", err)
		return err
	}
	return nil
}

func GetAddress() string {
	return convert.GetString(&flags.Address)
}
func GetReportInterval() uint64 {
	return convert.GetUint(&flags.ReportInterval)
}
func GetPollInterval() uint64 {
	return convert.GetUint(&flags.PollInterval)
}
func GetHashKey() string {
	return convert.GetString(&flags.HashKey)
}
