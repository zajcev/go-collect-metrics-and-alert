package model

import (
	"reflect"
)

// Metrics represents the structure of the metrics that will be scrape and sent to the server.
type Metrics struct {
	PollCount int64
	Alloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	NumForcedGC,
	NumGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	StackSys,
	Sys,
	TotalAlloc,
	TotalMemory,
	FreeMemory,
	CPUutilization1,
	RandomValue float64
}

// MetricJSON represents the JSON structure
type MetricJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// GetValueByName get metric value by name from a list of metrics
func GetValueByName(v any, field string) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil
	}
	return f.Interface()
}

// SetFieldValue set value of metric by name if isValid
func SetFieldValue(obj interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(obj).Elem().FieldByName(fieldName)
	if !v.IsValid() || value == nil {
		return
	}
	v.Set(reflect.ValueOf(value).Convert(v.Type()))
}
