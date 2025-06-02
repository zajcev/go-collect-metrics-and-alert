package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

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
