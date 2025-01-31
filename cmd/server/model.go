package main

import (
	"net/http"
)

type MemStorage struct {
	mtype string
	value interface{}
}

func setCounter(m map[string]*MemStorage, k string, t string, v int64) int {
	if m[k] != nil {
		if m[k].mtype == t {
			m[k].value = m[k].value.(int64) + v
			return http.StatusOK
		} else {
			return http.StatusBadRequest
		}
	} else {
		m[k] = &MemStorage{mtype: t, value: v}
		return http.StatusOK
	}
}

func setGauge(m map[string]*MemStorage, k string, t string, v float64) int {
	if m[k] != nil {
		if m[k].mtype == t {
			m[k] = &MemStorage{mtype: t, value: v}
			return http.StatusOK
		} else {
			return http.StatusBadRequest
		}
	} else {
		m[k] = &MemStorage{mtype: t, value: v}
		return http.StatusOK
	}
}
