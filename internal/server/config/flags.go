package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
)

var flags Flags

type Flags struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
	DBHost        string `env:"DATABASE_DSN"`
	HashKey       string `env:"KEY"`
}

func NewConfig() error {
	flag.StringVar(&flags.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flags.HashKey, "k", "", "key for sha256sum")
	flag.IntVar(&flags.StoreInterval, "i", 300, "interval between stored files")
	flag.StringVar(&flags.FilePath, "f", "/tmp/metrics", "path to store files") ///tmp/metrics
	flag.BoolVar(&flags.Restore, "r", false, "restore files")
	flag.StringVar(&flags.DBHost, "d", "", "database host") //postgres://user:password@localhost:5432/metrics?sslmode=disable
	flag.Parse()
	if err := env.Parse(&flags); err != nil {
		log.Printf("%+v", err)
		return err
	}
	return nil
}

func GetAddress() string { return convert.GetString(&flags.Address) }

func GetStoreInterval() uint64 {
	return convert.GetUint(&flags.StoreInterval)
}
func GetFilePath() string {
	return convert.GetString(&flags.FilePath)
}

func GetRestore() bool { return convert.GetBool(&flags.Restore) }

func GetDBHost() string {
	return convert.GetString(&flags.DBHost)
}

func GetHashKey() string {
	return convert.GetString(&flags.HashKey)
}
