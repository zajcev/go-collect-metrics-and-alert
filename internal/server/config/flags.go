package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
)

var flags Flags

type Flags struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
	DBHost        string `env:"DATABASE_DSN"`
}

func NewConfig() error {
	flag.StringVar(&flags.Address, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flags.StoreInterval, "i", 300, "interval between stored files")
	flag.StringVar(&flags.FilePath, "f", "/tmp/metrics", "path to store files")
	flag.BoolVar(&flags.Restore, "r", true, "restore files")
	flag.StringVar(&flags.FilePath, "d", "postgres://root:ok@localhost:5432/test?sslmode=disable", "database host")
	flag.Parse()
	if err := env.Parse(&flags); err != nil {
		log.Printf("%+v", err)
		return err
	}
	return nil
}

func GetConfig() *Flags {
	return &flags
}
