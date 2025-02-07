package listeners

import (
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"math/rand"
	"reflect"
	"runtime"
)

func NewMonitor() {
	var rt runtime.MemStats
	runtime.ReadMemStats(&rt)
	mt := reflect.TypeOf(MemStorage)
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		model.SetFieldValue(&MemStorage, f.Name, model.GetValueByName(rt, f.Name))
		//log.Printf("Monitor: Name: %v = Value: %v", f.Name, getValueByName(m, f.Name))
	}
	AddCustomMetric()
}

func AddCustomMetric() {
	model.SetFieldValue(&MemStorage, "RandomValue", rand.Float64())
	if model.GetValueByName(MemStorage, "PollCount") == nil {
		model.SetFieldValue(MemStorage, "PollCount", int64(1))
	} else {
		counter = model.GetValueByName(&MemStorage, "PollCount").(int64)
		counter++
		model.SetFieldValue(&MemStorage, "PollCount", counter)
	}
}
