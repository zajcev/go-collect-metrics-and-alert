package storage

import (
	"encoding/json"
	"os"

	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer creates a open file and returns producer instance
func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteMetrics writes metrics to file
func (p *Producer) WriteMetrics(metrics *models.MemStorage) error {
	return p.encoder.Encode(metrics)
}

func (p *Producer) Close() error {
	return p.file.Close()
}
