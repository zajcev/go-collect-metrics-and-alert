package model

import (
	"reflect"
	"testing"
)

func TestGetValueByName(t *testing.T) {
	metrics := Metrics{PollCount: 10, Alloc: 5.5}
	tests := []struct {
		name     string
		field    string
		expected interface{}
	}{
		{"Valid PollCount", "PollCount", int64(10)},
		{"Valid Alloc", "Alloc", float64(5.5)},
		{"Invalid Field", "InvalidField", nil},
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
		name     string
		field    string
		value    interface{}
		expected interface{}
	}{
		{"Set PollCount", "PollCount", int64(20), int64(20)},
		{"Set Alloc", "Alloc", float64(10.5), float64(10.5)},
		{"Invalid Field", "InvalidField", float64(15.5), nil},
		{"Nil Value", "PollCount", nil, int64(20)}, // Should not change
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "Nil Value" {
				SetFieldValue(metrics, "PollCount", int64(20)) // First set to a valid value
			}
			SetFieldValue(metrics, test.field, test.value)
			result := GetValueByName(metrics, test.field)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}
