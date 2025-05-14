package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
)

func TestNewProducer(t *testing.T) {
	fileName := "testfile.json"
	producer, err := NewProducer(fileName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer producer.Close()

	if producer.file == nil {
		t.Fatal("expected file to be opened, got nil")
	}
	if producer.encoder == nil {
		t.Fatal("expected encoder to be initialized, got nil")
	}

	os.Remove(fileName)
}

func TestWriteMetrics(t *testing.T) {
	fileName := "testfile.json"
	producer, err := NewProducer(fileName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer producer.Close()

	metrics := &models.MemStorage{}

	err = producer.WriteMetrics(metrics)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer file.Close()

	var writtenMetrics models.MemStorage
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&writtenMetrics)
	if err != nil {
		t.Fatalf("expected no error decoding metrics, got %v", err)
	}

	os.Remove(fileName)
}

func TestClose(t *testing.T) {
	fileName := "testfile.json"
	producer, err := NewProducer(fileName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = producer.Close()
	if err != nil {
		t.Fatalf("expected no error on close, got %v", err)
	}

	err = producer.WriteMetrics(&models.MemStorage{})
	if err == nil {
		t.Fatal("expected an error when writing to a closed producer, got none")
	}
	os.Remove(fileName)
}
