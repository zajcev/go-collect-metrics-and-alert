package config

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v11"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"os"
)

type Config interface {
	Load() error
	loadFromFlags() error
	loadFromJSON(path string) error
	loadFromENV() error
	GetAddress() string
	GetStoreInterval() uint64
	GetFilePath() string
	GetRestore() bool
	GetDBHost() string
	GetHashKey() string
	GetCryptoKey() string
}

type Provider struct {
	Config Flags
}
type Flags struct {
	Address       string `env:"ADDRESS"           json:"address"`
	FilePath      string `env:"FILE_STORAGE_PATH" json:"store_file"`
	DBHost        string `env:"DATABASE_DSN"      json:"database_dsn"`
	HashKey       string `env:"KEY"               json:"hash_key"`
	CryptoKey     string `env:"CRYPTO_KEY"        json:"crypto_key"`
	StoreInterval int    `env:"STORE_INTERVAL"    json:"store_interval"`
	Restore       bool   `env:"RESTORE"           json:"restore"`
	ConfigFile    string `env:"CONFIG"            json:"-"`
}

func NewConfig() *Provider {
	return &Provider{
		Config: Flags{},
	}
}

// Load parses the command-line flags, environment variables, json file.
func (c *Provider) Load() error {
	err := c.loadFromFlags()
	if err != nil {
		return err
	}

	configFile := c.Config.ConfigFile
	if configFile == "" {
		configFile = os.Getenv("CONFIG")
	}

	if configFile != "" {
		err = c.loadFromJSON(configFile)
		if err != nil {
			return err
		}
	}

	err = c.loadFromENV()
	if err != nil {
		return err
	}

	return nil
}

func (c *Provider) loadFromFlags() error {
	flag.StringVar(&c.Config.ConfigFile, "c", "", "path to config file")
	flag.StringVar(&c.Config.ConfigFile, "config", "", "path to config file (long form)")
	flag.StringVar(&c.Config.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.Config.HashKey, "k", "", "key for sha256sum")
	flag.IntVar(&c.Config.StoreInterval, "i", 300, "interval between stored files")
	flag.StringVar(&c.Config.FilePath, "f", "/tmp/metrics", "path to store files")                   ///tmp/metrics
	flag.StringVar(&c.Config.CryptoKey, "crypto-key", "/tmp/key.pem", "public key for decrypt data") ///tmp/key.pem
	flag.BoolVar(&c.Config.Restore, "r", false, "restore files")
	flag.StringVar(&c.Config.DBHost, "d", "", "database host") //postgres://user:password@localhost:5432/metrics?sslmode=disable
	flag.Parse()
	return nil
}

func (c *Provider) loadFromJSON(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &c.Config)
}

func (c *Provider) loadFromENV() error {
	return env.Parse(&c.Config)
}

// GetAddress returns the address and port to run the server.
func (c *Provider) GetAddress() string { return convert.GetString(&c.Config.Address) }

// GetStoreInterval returns the interval between stored files.
func (c *Provider) GetStoreInterval() uint64 {
	return convert.GetUint(&c.Config.StoreInterval)
}

// GetFilePath returns the path to store files.
func (c *Provider) GetFilePath() string {
	return convert.GetString(&c.Config.FilePath)
}

// GetRestore returns the flag to restore files.
func (c *Provider) GetRestore() bool { return convert.GetBool(&c.Config.Restore) }

// GetDBHost returns the database host.
func (c *Provider) GetDBHost() string {
	return convert.GetString(&c.Config.DBHost)
}

// GetHashKey returns the key for sha256sum.
func (c *Provider) GetHashKey() string {
	return convert.GetString(&c.Config.HashKey)
}

// GetCryptoKey return path to file with private key
func (c *Provider) GetCryptoKey() string { return convert.GetString(&c.Config.CryptoKey) }
