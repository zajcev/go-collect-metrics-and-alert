package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
)

type Flags struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
}

func ParseFlags() Flags {
	f := Flags{}
	flag.StringVar(&f.Address, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&f.StoreInterval, "i", 300, "interval between stored files")
	flag.StringVar(&f.FilePath, "f", "/tmp/metrics", "path to store files")
	flag.BoolVar(&f.Restore, "r", false, "restore files")
	flag.Parse()

	if err := env.Parse(&f); err != nil {
		log.Printf("%+v", err)
	}

	return f
}
