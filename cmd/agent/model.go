package main

import (
	"math/rand"
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

func getValueByName(v interface{}, field string) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil
	}
	return f.Interface()
}

func setFieldValue(obj interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(obj).Elem().FieldByName(fieldName)
	if !v.IsValid() || value == nil {
		return
	}
	v.Set(reflect.ValueOf(value).Convert(v.Type()))
}

func addCustomMetric() {
	setFieldValue(&m, "RandomValue", rand.Float64())
	if getValueByName(&m, "PollCount") == nil {
		setFieldValue(&m, "PollCount", int64(1))
	} else {
		counter = getValueByName(&m, "PollCount").(int64)
		counter++
		setFieldValue(&m, "PollCount", counter)
	}
}
