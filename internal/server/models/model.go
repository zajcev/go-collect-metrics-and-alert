package models

import (
	"fmt"
	"net/http"
)

type MemStorage struct {
	Mtype string
	Value interface{}
}

func GetMetricValue(m map[string]*MemStorage, k string, t string) string {
	if m[k] != nil && m[k].Mtype == t {
		return fmt.Sprintf("%v", m[k].Value)
	} else {
		return ""
	}
}

func SetCounter(m map[string]*MemStorage, k string, t string, v int64) int {
	if m[k] != nil {
		if m[k].Mtype == t {
			m[k].Value = m[k].Value.(int64) + v
			return http.StatusOK
		} else {
			return http.StatusBadRequest
		}
	} else {
		m[k] = &MemStorage{Mtype: t, Value: v}
		return http.StatusOK
	}
}

func SetGauge(m map[string]*MemStorage, k string, t string, v float64) int {
	if m[k] != nil {
		if m[k].Mtype == t {
			m[k] = &MemStorage{Mtype: t, Value: v}
			return http.StatusOK
		} else {
			return http.StatusBadRequest
		}
	} else {
		m[k] = &MemStorage{Mtype: t, Value: v}
		return http.StatusOK
	}
}
