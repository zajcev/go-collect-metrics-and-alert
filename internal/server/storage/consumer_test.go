package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
)

func TestNewConsumer(t *testing.T) {
	// Test with a valid file
	fileName := "test_metrics.json"
	_, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("could not create file: %v", err)
	}
	defer os.Remove(fileName)

	consumer, err := NewConsumer(fileName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if consumer.file == nil {
		t.Fatal("expected file to be opened")
	}
}

func TestNewConsumer_FileNotFound(t *testing.T) {
	consumer, err := NewConsumer("invalid_file.json")
	if err == nil {
		t.Fatal("expected error, got none")
	}
	if consumer != nil {
		t.Fatal("expected consumer to be nil")
	}
}

func TestReadMetrics(t *testing.T) {
	fileName := "test_metrics.json"
	data := &models.MemStorage{ /* populate with test data */ }
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("could not create file: %v", err)
	}
	defer os.Remove(fileName)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		t.Fatalf("could not encode test data: %v", err)
	}
	file.Close()

	consumer, err := NewConsumer(fileName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer consumer.Close()

	metric, err := consumer.ReadMetrics()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if metric == nil {
		t.Fatal("expected metric to be returned")
	}
}

func TestReadMetrics_InvalidJSON(t *testing.T) {
	fileName := "test_invalid_metrics.json"
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("could not create file: %v", err)
	}
	defer os.Remove(fileName)

	if _, err := file.WriteString("{invalid}"); err != nil {
		t.Fatalf("could not write invalid json: %v", err)
	}
	file.Close()

	consumer, err := NewConsumer(fileName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer consumer.Close()

	_, err = consumer.ReadMetrics()
	if err == nil {
		t.Fatal("expected error, got none")
	}
}
