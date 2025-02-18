package model

import (
	"reflect"
)

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
	RandomValue float64
}

type MetricJson struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func GetValueByName(v any, field string) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil
	}
	return f.Interface()
}

func SetFieldValue(obj interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(obj).Elem().FieldByName(fieldName)
	if !v.IsValid() || value == nil {
		return
	}
	v.Set(reflect.ValueOf(value).Convert(v.Type()))
}
