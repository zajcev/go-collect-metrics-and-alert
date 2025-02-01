package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_metricCollector(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update", nil)
	w := httptest.NewRecorder()
	metricCollector(w, req)
	res := w.Result()
	println(res.StatusCode)
}
