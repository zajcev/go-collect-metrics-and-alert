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

// SaveMetricStorageSchedule save metrics to file by duration
func SaveMetricStorageSchedule(interval int, file string, metrics interface {
	GetAllMetrics(ctx context.Context) *models.MemStorage
}) error {
	duration := time.Duration(interval) * time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for range ticker.C {
		return save(file, metrics)
	}
	return nil
}

// SaveMetricStorageOnce save metrics to file one time
func SaveMetricStorageOnce(file string, metrics interface {
	GetAllMetrics(ctx context.Context) *models.MemStorage
}) error {
	return save(file, metrics)
}

func save(file string, metrics interface {
	GetAllMetrics(ctx context.Context) *models.MemStorage
}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	producer, err := storage.NewProducer(file)
	m := metrics.GetAllMetrics(ctx)
	if err != nil {
		return err
	}
	err = producer.WriteMetrics(m)
	if err != nil {
		return err
	}
	return nil
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
