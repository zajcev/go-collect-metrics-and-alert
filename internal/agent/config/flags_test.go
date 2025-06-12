package config

import (
	"flag"
	"os"
	"testing"
)

func TestLoadFromFlags(t *testing.T) {
	os.Args = []string{"cmd", "-a", "localhost:9090", "-k", "mysecretkey", "-crypto-key", "/tmp/mycert.pem", "-r", "5", "-p", "10", "-l", "3"}
	provider := NewConfig()
	err := provider.Load()
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if provider.GetAddress() != "localhost:9090" {
		t.Errorf("expected address to be 'localhost:9090', got '%s'", provider.GetAddress())
	}
	if provider.GetHashKey() != "mysecretkey" {
		t.Errorf("expected hash key to be 'mysecretkey', got '%s'", provider.GetHashKey())
	}
	if provider.GetCryptoKey() != "/tmp/mycert.pem" {
		t.Errorf("expected crypto key to be '/tmp/mycert.pem', got '%s'", provider.GetCryptoKey())
	}
	if provider.GetReportInterval() != 5 {
		t.Errorf("expected report interval to be 5, got %d", provider.GetReportInterval())
	}
	if provider.GetPollInterval() != 10 {
		t.Errorf("expected poll interval to be 10, got %d", provider.GetPollInterval())
	}
	if provider.GetRateLimit() != 3 {
		t.Errorf("expected rate limit to be 3, got %d", provider.GetRateLimit())
	}
}

func TestLoadFromJSON(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fileContent := `{
"address": "localhost:8081",
"hashkey": "testkey",
"crypto_key": "/tmp/testcert.pem",
"report_interval": 2,
"poll_interval": 2,
"rate_limit": 5
}`

	err := os.WriteFile("config.json", []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("could not write config file: %v", err)
	}
	defer os.Remove("config.json")

	provider := NewConfig()
	err = provider.loadFromJSON("config.json")
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if provider.GetAddress() != "localhost:8081" {
		t.Errorf("expected address to be 'localhost:8081', got '%s'", provider.GetAddress())
	}
	if provider.GetHashKey() != "testkey" {
		t.Errorf("expected hash key to be 'testkey', got '%s'", provider.GetHashKey())
	}
	if provider.GetCryptoKey() != "/tmp/testcert.pem" {
		t.Errorf("expected crypto key to be '/tmp/testcert.pem', got '%s'", provider.GetCryptoKey())
	}
	if provider.GetReportInterval() != 2 {
		t.Errorf("expected report interval to be 2, got %d", provider.GetReportInterval())
	}
	if provider.GetPollInterval() != 2 {
		t.Errorf("expected poll interval to be 2, got %d", provider.GetPollInterval())
	}
	if provider.GetRateLimit() != 5 {
		t.Errorf("expected rate limit to be 5, got %d", provider.GetRateLimit())
	}
}

func TestLoadFromENV(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Setenv("ADDRESS", "localhost:7070")
	os.Setenv("KEY", "envkey")
	os.Setenv("CRYPTO_KEY", "/tmp/envcert.pem")
	os.Setenv("REPORT_INTERVAL", "15")
	os.Setenv("POLL_INTERVAL", "5")
	os.Setenv("RATE_LIMIT", "10")
	defer os.Unsetenv("ADDRESS")
	defer os.Unsetenv("KEY")
	defer os.Unsetenv("CRYPTO_KEY")
	defer os.Unsetenv("REPORT_INTERVAL")
	defer os.Unsetenv("POLL_INTERVAL")
	defer os.Unsetenv("RATE_LIMIT")

	provider := NewConfig()
	err := provider.Load()
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if provider.GetAddress() != "localhost:7070" {
		t.Errorf("expected address to be 'localhost:7070', got '%s'", provider.GetAddress())
	}
	if provider.GetHashKey() != "envkey" {
		t.Errorf("expected hash key to be 'envkey', got '%s'", provider.GetHashKey())
	}
	if provider.GetCryptoKey() != "/tmp/envcert.pem" {
		t.Errorf("expected crypto key to be '/tmp/envcert.pem', got '%s'", provider.GetCryptoKey())
	}
	if provider.GetReportInterval() != 15 {
		t.Errorf("expected report interval to be 15, got %d", provider.GetReportInterval())
	}
	if provider.GetPollInterval() != 5 {
		t.Errorf("expected poll interval to be 5, got %d", provider.GetPollInterval())
	}
	if provider.GetRateLimit() != 10 {
		t.Errorf("expected rate limit to be 10, got %d", provider.GetRateLimit())
	}
}
