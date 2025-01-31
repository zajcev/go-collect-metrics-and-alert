package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct { // добавляем слайс тестов
		name  string
		value User
		want  string
	}{
		{
			name:  "user test #1",
			value: User{"biba", "boba"},
			want:  "biba boba",
		},
		{
			name:  "user test #2",
			value: User{"", "last"},
			want:  " last",
		},
		{
			name:  "user test #3",
			value: User{"123", "last"},
			want:  "123 last",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, User.FullName(test.value))
		})
	}
}

func TestUser_FullName(t *testing.T) {
	type fields struct {
		FirstName string
		LastName  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "user test #3",
			fields: fields{"123", "last"},
			want:   "123 last",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
			}
			assert.Equalf(t, tt.want, u.FullName(), "FullName()")
		})
	}
}
