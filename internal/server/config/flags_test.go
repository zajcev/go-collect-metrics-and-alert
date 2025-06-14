package config

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestLoadFromFlags(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	os.Args = []string{
		"cmd",
		"-d", "test.db.host",
		"-a", "127.0.0.1:8080",
		"-i", "600",
		"-f", "/tmp/test_metrics",
		"-k", "testkey",
		"-crypto-key", "/tmp/test_key.pem",
		"-r", "true",
	}

	config := NewConfig()
	err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := Flags{
		Address:       "127.0.0.1:8080",
		FilePath:      "/tmp/test_metrics",
		DBHost:        "test.db.host",
		HashKey:       "testkey",
		CryptoKey:     "/tmp/test_key.pem",
		StoreInterval: 600,
		Restore:       true,
	}

	if !reflect.DeepEqual(config.Config, expected) {
		t.Errorf("expected %+v,\n got %+v", expected, config.Config)
	}
}

func TestLoadFromJSON(t *testing.T) {
	jsonData := `{
					"address": "127.0.0.1:8080",
					"store_file": "/tmp/test_metrics",
					"database_dsn": "postgres://user:password@localhost:5432/test?sslmode=disable",
					"key": "testkey",
					"crypto_key": "/tmp/test_key.pem",
					"store_interval": 600,
					"restore": true
                  }`
	tmpFile, err := os.Create("test_config.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove("test_config.json")

	if _, err := tmpFile.WriteString(jsonData); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	config := NewConfig()
	err = config.loadFromJSON("test_config.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := Flags{
		Address:       "127.0.0.1:8080",
		FilePath:      "/tmp/test_metrics",
		DBHost:        "postgres://user:password@localhost:5432/test?sslmode=disable",
		HashKey:       "",
		CryptoKey:     "/tmp/test_key.pem",
		StoreInterval: 600,
		Restore:       true,
	}

	if !reflect.DeepEqual(config.Config, expected) {
		t.Errorf("expected %+v,\n got %+v", expected, config.Config)
	}
}

func TestLoadFromENV(t *testing.T) {
	os.Setenv("ADDRESS", "127.0.0.1:8080")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/test_metrics")
	os.Setenv("DATABASE_DSN", "postgres://user:password@localhost:5432/test?sslmode=disable")
	os.Setenv("KEY", "testkey")
	os.Setenv("CRYPTO_KEY", "/tmp/test_key.pem")
	os.Setenv("STORE_INTERVAL", "600")
	os.Setenv("RESTORE", "true")
	defer os.Clearenv()

	config := NewConfig()
	err := config.loadFromENV()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := Flags{
		Address:       "127.0.0.1:8080",
		FilePath:      "/tmp/test_metrics",
		DBHost:        "postgres://user:password@localhost:5432/test?sslmode=disable",
		HashKey:       "testkey",
		CryptoKey:     "/tmp/test_key.pem",
		StoreInterval: 600,
		Restore:       true,
	}

	if !reflect.DeepEqual(config.Config, expected) {
		t.Errorf("expected %+v, got %+v", expected, config.Config)
	}
}
