package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
	"os"
)

var flags Flags

type Flags struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
}

func ParseFlags() error {
	flag.StringVar(&flags.Address, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flags.StoreInterval, "i", 5, "interval between stored files")
	flag.StringVar(&flags.FilePath, "f", "/tmp/metrics", "path to store files")
	flag.BoolVar(&flags.Restore, "r", true, "restore files")
	flag.Parse()
	os.Setenv("ADDRESS", "localhost:8080")
	os.Setenv("RESTORE", "true")
	os.Setenv("STORE_INTERVAL", "2")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/metrics")
	if err := env.Parse(&flags); err != nil {
		log.Printf("%+v", err)
		return err
	}
	return nil
}

func GetFlags() *Flags {
	return &flags
}
