package storage

import (
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"os"
)

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadMetrics() (*models.MemStorage, error) {
	metric := models.MemStorage{}
	if err := c.decoder.Decode(&metric); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
