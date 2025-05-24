// package listeners implement collector and reporter for the agent.

package listeners

import (
	"context"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
)

var counter = int64(0)

// NewMonitor start scraping metrics from runtime.MemStats with interval.
// Func variable interval defines the frequency of scraping.
func NewMonitor(ctx context.Context, interval int) error {
	duration := time.Duration(interval) * time.Second
	ticker := time.NewTicker(duration)
	var rt runtime.MemStats
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			runtime.ReadMemStats(&rt)
			mt := reflect.TypeOf(MemStorage)
			for i := 0; i < mt.NumField(); i++ {
				f := mt.Field(i)
				model.SetFieldValue(&MemStorage, f.Name, model.GetValueByName(rt, f.Name))
			}
			addCustomMetric()
		}
	}
}

func addCustomMetric() {
	model.SetFieldValue(&MemStorage, "RandomValue", rand.Float64())
	if model.GetValueByName(MemStorage, "PollCount") == nil {
		model.SetFieldValue(MemStorage, "PollCount", int64(1))
	} else {
		counter = model.GetValueByName(&MemStorage, "PollCount").(int64)
		counter++
		model.SetFieldValue(&MemStorage, "PollCount", counter)
	}
}

// AdditionalMetrics start scraping metrics for monitoring Memory and CPU utilization
// Func variable interval defines the frequency of scraping.
func AdditionalMetrics(ctx context.Context, interval int) error {
	duration := time.Duration(interval) * time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			memoryInfo, err := mem.VirtualMemory()
			if err != nil {
				log.Fatalf("Ошибка при получении информации о памяти: %v", err)
			}
			cpuUtilization, err := cpu.Percent(time.Second, false)
			if err != nil {
				log.Fatalf("Ошибка при получении информации о CPU: %v", err)
			}
			model.SetFieldValue(&MemStorage, "TotalMemory", convert.GetFloat(&memoryInfo.Total))
			model.SetFieldValue(&MemStorage, "FreeMemory", convert.GetFloat(&memoryInfo.Free))
			model.SetFieldValue(&MemStorage, "CPUutilization1", convert.GetFloat(&cpuUtilization[0]))
		}
	}
}
