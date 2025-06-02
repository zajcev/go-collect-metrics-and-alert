package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	_ "github.com/lib/pq"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/storage"
	"log"
	"time"
)

// SaveMetricStorage save metrics to file
func SaveMetricStorage(file string, metrics interface {
	GetAllMetrics(ctx context.Context) *models.MemStorage
}) {
	producer, err := storage.NewProducer(file)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	m := metrics.GetAllMetrics(ctx)
	if err != nil {
		log.Fatalf("Error create NewProducer bla bla bla : %v", err)
	}
	err = producer.WriteMetrics(m)
	if err != nil {
		log.Fatalf("Error write metrics : %v", err)
	}
}

func calculateSHA256Hash(data []byte, key string) string {
	k := []byte(key)
	signedData := append(k, data...)
	hash := sha256.Sum256(signedData)
	return hex.EncodeToString(hash[:])
}

func checkSHA256Hash(data []byte, key string, sum string) bool {
	hash := calculateSHA256Hash(data, key)
	if hash == sum {
		return true
	} else {
		return false
	}
}

// RestoreMetricStorage restore metrics from a file
func RestoreMetricStorage(file string) *models.MemStorage {
	consumer, err := storage.NewConsumer(file)
	if err != nil {
		log.Printf("Error while init file consumer %v", err)
		return &models.MemStorage{}
	}
	metrics, err := consumer.ReadMetrics()
	if err != nil {
		log.Printf("Error while read metric %v", err)
	}
	return metrics
}
