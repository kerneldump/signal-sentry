package gateway

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchStats(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"device": {"model": "TEST-MODEL"}}`)
	}))
	defer ts.Close()

	// For now we test a standalone function similar to the original one
	// but it should be exported as FetchStats
	client := &http.Client{Timeout: 1 * time.Second}
	data, err := FetchStats(client, ts.URL)
	if err != nil {
		t.Fatalf("FetchStats failed: %v", err)
	}

	if data.Device.Model != "TEST-MODEL" {
		t.Errorf("Expected model 'TEST-MODEL', got '%s'", data.Device.Model)
	}
}

func TestFetchStats_Retry(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		fmt.Fprint(w, `{"device": {"model": "SUCCESS"}}`)
	}))
	defer ts.Close()

	client := &http.Client{Timeout: 1 * time.Second}
	// We might need a way to configure retries.
	// Maybe a Client struct is better.
	data, err := FetchStats(client, ts.URL)
	if err != nil {
		t.Fatalf("FetchStats failed after retries: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
	if data.Device.Model != "SUCCESS" {
		t.Errorf("Expected model 'SUCCESS', got '%s'", data.Device.Model)
	}
}
