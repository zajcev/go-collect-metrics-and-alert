package main

import (
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"testing"
)

var mtest model.Metrics

func Test_monitor(t *testing.T) {
	tests := []struct {
		name    string
		metric  model.Metrics
		wantErr bool
	}{
		{
			name:    "test",
			metric:  mtest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewMonitor()
		})
	}
}

func Test_reporter(t *testing.T) {
	tests := []struct {
		name    string
		metric  model.Metrics
		wantErr bool
	}{
		{
			name:    "test",
			metric:  mtest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewReporter("http://localhost:8080")
		})
	}
}
