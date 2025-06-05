package db

import (
	"context"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
)

type Storage interface {
	SetDeltaRaw(ctx context.Context, name string, metricType string, value int64) int
	SetValueRaw(ctx context.Context, name string, metricType string, value float64) int
	SetValueJSON(ctx context.Context, input models.Metric) int
	SetDeltaJSON(ctx context.Context, input models.Metric) int
	GetMetricRaw(ctx context.Context, name string, metricType string) any
	GetMetricJSON(ctx context.Context, input models.Metric) (models.Metric, int)
	GetAllMetrics(ctx context.Context) *models.MemStorage
	SetListJSON(ctx context.Context, list []models.Metric) int
	Ping(ctx context.Context) error
}
