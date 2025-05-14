package listeners

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCalculateSHA256Hash(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		key      string
		expected string
	}{
		{"Basic test", []byte("test data"), "key", "91d2330355770ae2a13eb43e62d9ed805aa140d4c7157a7cf69c170d1050fb6c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSHA256Hash(tt.data, tt.key)
			if result != tt.expected {
				t.Errorf("Expected %s but got %s", tt.expected, result)
			}
		})
	}
}

func TestRetryFailure(t *testing.T) {
	client := &http.Client{}
	req := httptest.NewRequest("GET", "http://example.com", nil)

	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()

	req.URL, _ = url.Parse(ts.URL)

	resp, err := retry(client, req, 3)
	defer resp.Body.Close()
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	if resp != nil {
		t.Fatalf("Expected nil response, got %v", resp)
	}
}
