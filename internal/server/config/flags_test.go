package config

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	if err := os.Setenv("ADDRESS", "127.0.0.1:8080"); err != nil {
		log.Fatalf("failed to set ADDRESS: %v", err)
	}
	if err := os.Setenv("STORE_INTERVAL", "600"); err != nil {
		log.Fatalf("failed to set STORE_INTERVAL: %v", err)
	}
	if err := os.Setenv("FILE_STORAGE_PATH", "/tmp/test_metrics"); err != nil {
		log.Fatalf("failed to set FILE_STORAGE_PATH: %v", err)
	}
	if err := os.Setenv("RESTORE", "true"); err != nil {
		log.Fatalf("failed to set RESTORE: %v", err)
	}
	if err := os.Setenv("DATABASE_DSN", "postgres://user:password@localhost:5432/metrics?sslmode=disable"); err != nil {
		log.Fatalf("failed to set DATABASE_DSN: %v", err)
	}

	err := NewConfig()
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:8080", flags.Address)
	assert.Equal(t, 600, flags.StoreInterval)
	assert.Equal(t, "/tmp/test_metrics", flags.FilePath)
	assert.True(t, flags.Restore)
	assert.Equal(t, "postgres://user:password@localhost:5432/metrics?sslmode=disable", flags.DBHost)
	assert.Equal(t, "testkey", flags.HashKey)
}

func TestGetAddress(t *testing.T) {
	flags.Address = "localhost:9090"
	assert.Equal(t, "localhost:9090", GetAddress())
}

func TestGetStoreInterval(t *testing.T) {
	flags.StoreInterval = 450
	assert.Equal(t, uint64(450), GetStoreInterval())
}

func TestGetFilePath(t *testing.T) {
	flags.FilePath = "/var/data/metrics"
	assert.Equal(t, "/var/data/metrics", GetFilePath())
}

func TestGetRestore(t *testing.T) {
	flags.Restore = false
	assert.False(t, GetRestore())
}

func TestGetDBHost(t *testing.T) {
	flags.DBHost = "localhost:5432"
	assert.Equal(t, "localhost:5432", GetDBHost())
}

func TestGetHashKey(t *testing.T) {
	flags.HashKey = "mysecretkey"
	assert.Equal(t, "mysecretkey", GetHashKey())
}
