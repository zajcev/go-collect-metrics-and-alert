package main

import (
	"net/http/httptest"
	"testing"
)

func genPostReq(target string, method string) (int, error) {
	req := httptest.NewRequest(method, target, nil)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	metricCollector(w, req)
	res := w.Result()
	err := res.Body.Close()
	c := res.StatusCode
	return c, err
}

func Test_metricCollector(t *testing.T) {
	tests := []struct {
		name   string
		target string
		method string
		expect int
		err    error
	}{
		{"Metric without name", "/update/counter/123", "POST", 404, nil},
		{"Invalid gauge value", "/update/gauge/test/test", "POST", 400, nil},
		{"Invalid counter value", "/update/counter/test/1.11", "POST", 400, nil},
		{"Invalid metric type", "/update/gauge1/test/123", "POST", 400, nil},
		{"Invalid path", "/updater/gauge1/test/123", "POST", 400, nil},
		{"Invalid request method", "/update/gauge/test1/123", "GET", 405, nil},
		{"Put gauge metric", "/update/gauge/test/123", "POST", 200, nil},
		{"Put counter metric", "/update/counter/test1/123", "POST", 200, nil},
		{"Try to change metric type", "/update/gauge/test1/123", "POST", 400, nil},
	}
	for _, test := range tests {
		c, err := genPostReq(test.target, test.method)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if test.expect != c {
			t.Errorf("Name: %v : Expected %v but got: %v", test.name, test.expect, c)
		}
	}
}
