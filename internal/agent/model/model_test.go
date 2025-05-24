package model

import (
	"reflect"
	"testing"
)

func TestGetValueByName(t *testing.T) {
	metrics := Metrics{PollCount: 10, Alloc: 5.5}
	tests := []struct {
		expected interface{}
		name     string
		field    string
	}{
		{int64(10), "Valid PollCount", "PollCount"},
		{float64(5.5), "Valid Alloc", "Alloc"},
		{nil, "Invalid Field", "InvalidField"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetValueByName(metrics, test.field)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestSetFieldValue(t *testing.T) {
	metrics := &Metrics{}
	tests := []struct {
		value    interface{}
		expected interface{}
		name     string
		field    string
	}{
		{int64(20), int64(20), "Set PollCount", "PollCount"},
		{float64(10.5), float64(10.5), "Set Alloc", "Alloc"},
		{float64(15.5), nil, "Invalid Field", "InvalidField"},
		{nil, int64(20), "Nil Value", "PollCount"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "Nil Value" {
				SetFieldValue(metrics, "PollCount", int64(20))
			}
			SetFieldValue(metrics, test.field, test.value)
			result := GetValueByName(metrics, test.field)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}
