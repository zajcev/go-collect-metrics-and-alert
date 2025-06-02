package main

import (
	"bytes"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/routes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	}
	testStorage := models.NewMemStorage()
	testServer := httptest.NewServer(routes.NewRouter(testStorage))
	testServer.URL = "http://localhost:8080"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

func TestPrintLdFlags(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		date     string
		commit   string
		expected string
	}{
		{
			name:     "All fields set",
			version:  "1.0.0",
			date:     "2025-01-01",
			commit:   "e512351s",
			expected: "Build version: 1.0.0\nBuild date: 2025-01-01\nBuild commit: e512351s\n",
		},
		{
			name:     "All fields missing",
			version:  "",
			date:     "",
			commit:   "",
			expected: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A\n",
		},
	}

	buildVersion = "1.0.0"
	buildDate = "2025-01-01"
	buildCommit = "e512351s"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			defer func() { os.Stdout = oldStdout }()

			r, w, _ := os.Pipe()
			os.Stdout = w

			printLdFlags()

			err := w.Close()
			if err != nil {
				t.Fatalf("Error while close stdout reader : %v", err)
			}
			var buf bytes.Buffer
			io.Copy(&buf, r)

			got := strings.TrimSpace(buf.String())
			want := strings.TrimSpace(tt.expected)
			if got != want {
				t.Errorf("\nExpected:\n%s\nGot:\n%s", want, got)
			}
			buildVersion = ""
			buildDate = ""
			buildCommit = ""
		})
	}
}
