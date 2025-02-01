package main

import (
	"log"
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
	err := res.Body.Close()
	if err != nil {
		log.Fatalf("Error while close body: %v", err)
	}
}
