package listeners

import (
	"context"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"log"
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	err := NewMonitor(ctx, 1)
	if err != nil && err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	if model.GetValueByName(MemStorage, "PollCount") == nil {
		t.Error("PollCount should not be nil after NewMonitor execution")
	}
}

func TestAddCustomMetric(t *testing.T) {
	addCustomMetric()
	pollCount := model.GetValueByName(MemStorage, "PollCount").(int64)

	if pollCount != 2 {
		t.Errorf("expected PollCount to be 1, got %d", pollCount)
	}

	addCustomMetric()
	pollCount = model.GetValueByName(MemStorage, "PollCount").(int64)

	if pollCount != 3 {
		t.Errorf("expected PollCount to be 2, got %d", pollCount)
	}

	randomValue := model.GetValueByName(MemStorage, "RandomValue")
	if randomValue == nil {
		t.Error("RandomValue should not be nil after addCustomMetric execution")
	}
}

func TestAdditionalMetrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	err := AdditionalMetrics(ctx, 1)
	if err != nil && err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	totalMemory := model.GetValueByName(MemStorage, "TotalMemory")
	if totalMemory == nil {
		t.Error("TotalMemory should not be nil after AdditionalMetrics execution")
	}

	freeMemory := model.GetValueByName(MemStorage, "FreeMemory")
	if freeMemory == nil {
		t.Error("FreeMemory should not be nil after AdditionalMetrics execution")
	}

	cpuUtilization := model.GetValueByName(MemStorage, "CPUutilization1")
	if cpuUtilization == nil {
		t.Error("CPUutilization1 should not be nil after AdditionalMetrics execution")
	}
}

func BenchmarkMonitor(*testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), 3)
	defer cancel()
	err := NewMonitor(ctx, 2)
	if err != nil {
		log.Fatalf("Error create monitor : %v", err)
	}
}
