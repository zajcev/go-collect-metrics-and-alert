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
	FilePath      string `env:"FILE_STORAGE_PATH"`
	DBHost        string `env:"DATABASE_DSN"`
	HashKey       string `env:"KEY"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	Restore       bool   `env:"RESTORE"`
}

// NewConfig parses the command-line flags and environment variables.
func NewConfig() error {
	flag.StringVar(&flags.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flags.HashKey, "k", "", "key for sha256sum")
	flag.IntVar(&flags.StoreInterval, "i", 300, "interval between stored files")
	flag.StringVar(&flags.FilePath, "f", "/tmp/metrics", "path to store files") ///tmp/metrics
	flag.BoolVar(&flags.Restore, "r", false, "restore files")
	flag.StringVar(&flags.DBHost, "d", "postgres://user:password@localhost:5432/metrics?sslmode=disable", "database host") //postgres://user:password@localhost:5432/metrics?sslmode=disable
	flag.Parse()
	if err := env.Parse(&flags); err != nil {
		log.Printf("%+v", err)
		return err
	}
	return nil
}

// GetAddress returns the address and port to run the server.
func GetAddress() string { return convert.GetString(&flags.Address) }

// GetStoreInterval returns the interval between stored files.
func GetStoreInterval() uint64 {
	return convert.GetUint(&flags.StoreInterval)
}

// GetFilePath returns the path to store files.
func GetFilePath() string {
	return convert.GetString(&flags.FilePath)
}

// GetRestore returns the flag to restore files.
func GetRestore() bool { return convert.GetBool(&flags.Restore) }

// GetDBHost returns the database host.
func GetDBHost() string {
	return convert.GetString(&flags.DBHost)
}

// GetHashKey returns the key for sha256sum.
func GetHashKey() string {
	return convert.GetString(&flags.HashKey)
}
