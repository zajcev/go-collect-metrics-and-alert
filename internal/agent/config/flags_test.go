package config

import (
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	os.Args = []string{"cmd", "-a", "127.0.0.1:8080", "-k", "testkey", "-r", "5", "-p", "3", "-l", "10"}
	err := NewConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := Flags{
		Address:        "127.0.0.1:8080",
		ReportInterval: 5,
		PollInterval:   3,
		HashKey:        "testkey",
		RateLimit:      10,
	}

	if !reflect.DeepEqual(flags, expected) {
		t.Errorf("Expected %+v, got %+v", expected, flags)
	}
}

func TestGetAddress(t *testing.T) {
	flags.Address = "127.0.0.1:8080"
	got := GetAddress()
	expected := "127.0.0.1:8080"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func TestGetReportInterval(t *testing.T) {
	flags.ReportInterval = 5
	got := GetReportInterval()
	expected := 5
	if got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func TestGetPollInterval(t *testing.T) {
	flags.PollInterval = 3
	got := GetPollInterval()
	expected := 3
	if got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func TestGetHashKey(t *testing.T) {
	flags.HashKey = "testkey"
	got := GetHashKey()
	expected := "testkey"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func TestGetRateLimit(t *testing.T) {
	flags.RateLimit = 10
	got := GetRateLimit()
	expected := 10
	if got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}
