package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetArticleID(t *testing.T) {
	tests := []struct {
		name   string
		target string
		method string
		expect int
	}{
		{"Metric without name", "/update/counter/123", "POST", 404},
		{"Invalid gauge value", "/update/gauge/test/test", "POST", 400},
		{"Invalid counter value", "/update/counter/test/1.11", "POST", 400},
		{"Invalid metric type", "/update/gauge1/test/123", "POST", 400},
		{"Invalid path", "/updater/gauge1/test/123", "POST", 404},
		{"Put gauge metric", "/update/gauge/test/123", "POST", 200},
		{"Put counter metric", "/update/counter/test1/123", "POST", 200},
	}
	testServer := httptest.NewServer(Router())
	testServer.URL = "http://localhost:8080"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//println(test.method, testServer.Listener.Addr().String()+test.target)
			request, err := http.NewRequest(http.MethodPost, "http://"+testServer.Listener.Addr().String()+test.target, nil)
			request.Header.Add("Content-Type", "text/plain")
			if err != nil {
				t.Fatal(err)
			}
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			err = response.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != test.expect {
				t.Fatalf("expect %d, got %d", test.expect, response.StatusCode)
			}
		})
	}
}
