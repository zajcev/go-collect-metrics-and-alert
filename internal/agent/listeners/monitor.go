package listeners

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

var counter = int64(0)

func NewMonitor() {
	var rt runtime.MemStats
	runtime.ReadMemStats(&rt)
	mt := reflect.TypeOf(MemStorage)
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		model.SetFieldValue(&MemStorage, f.Name, model.GetValueByName(rt, f.Name))
	}
	addCustomMetric()
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

func AdditionalMetrics() {
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
