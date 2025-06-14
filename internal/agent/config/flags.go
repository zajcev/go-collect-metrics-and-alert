package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
)

type Config interface {
	Load() error
	loadFromFlags() error
	loadFromJSON(path string) error
	loadFromENV() error
	GetAddress() string
	GetReportInterval() int
	GetPollInterval() int
	GetHashKey() string
	GetRateLimit() int
	GetCryptoKey() string
}
type Provider struct {
	Config Flags
}

type Flags struct {
	Address        string `env:"ADDRESS"             json:"address"`
	HashKey        string `env:"KEY"                 json:"hashkey"`
	CryptoKey      string `env:"CRYPTO_KEY"          json:"crypto_key"`
	ReportInterval int    `env:"REPORT_INTERVAL"     json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL"       json:"poll_interval"`
	RateLimit      int    `env:"RATE_LIMIT"          json:"rate_limit"`
	ConfigFile     string `env:"CONFIG"              json:"-"`
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

func (c *Provider) loadFromFlags() error {
	flag.StringVar(&c.Config.ConfigFile, "c", "", "path to config file")
	flag.StringVar(&c.Config.ConfigFile, "config", "", "path to config file (long form)")
	flag.StringVar(&c.Config.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.Config.HashKey, "k", "12h5b12b521b", "key for sha256sum")
	flag.StringVar(&c.Config.CryptoKey, "crypto-key", "/tmp/cert.pem", "public key for encrypt data") ///tmp/cert.pem
	flag.IntVar(&c.Config.ReportInterval, "r", 1, "interval between report calls")
	flag.IntVar(&c.Config.PollInterval, "p", 1, "interval between polls")
	flag.IntVar(&c.Config.RateLimit, "l", 2, "request rate limiter")
	flag.Parse()
	return nil
}

// GetAddress returns the address and port to run the server.
func (c *Provider) GetAddress() string {
	return convert.GetString(&c.Config.Address)
}

// GetReportInterval returns the interval between report calls.
func (c *Provider) GetReportInterval() int {
	return c.Config.ReportInterval
}

// GetPollInterval returns the interval between polls.
func (c *Provider) GetPollInterval() int {
	return c.Config.PollInterval
}

// GetHashKey returns the key for sha256sum.
func (c *Provider) GetHashKey() string { return convert.GetString(&c.Config.HashKey) }

// GetRateLimit returns interval for Reporter
func (c *Provider) GetRateLimit() int { return c.Config.RateLimit }

// GetCryptoKey return path to file with public key
func (c *Provider) GetCryptoKey() string { return convert.GetString(&c.Config.CryptoKey) }
