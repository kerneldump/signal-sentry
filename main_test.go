package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"tmobile-stats/internal/gateway"
)

// TestColorizeRSRP tests the threshold logic for RSRP colors
func TestColorizeRSRP(t *testing.T) {
	tests := []struct {
		val      int
		expected string // We look for the color code substring
	}{
		{-70, ColorGreen},   // Excellent > -80
		{-85, ColorYellow},  // Good/Fair -80 to -100
		{-105, ColorRed},    // Poor < -100
		{-80, ColorYellow},  // Edge case
		{-100, ColorYellow}, // Edge case
	}

	for _, tt := range tests {
		result := colorizeRSRP(tt.val)
		if !strings.Contains(result, tt.expected) {
			t.Errorf("colorizeRSRP(%d) = %q; want containing %q", tt.val, result, tt.expected)
		}
	}
}

// TestColorizeSINR tests the threshold logic for SINR colors
func TestColorizeSINR(t *testing.T) {
	tests := []struct {
		val      int
		expected string
	}{
		{25, ColorGreen},  // Excellent > 20
		{10, ColorYellow}, // Good 0 to 20
		{-5, ColorRed},    // Poor < 0
		{20, ColorYellow}, // Edge case
		{0, ColorYellow},  // Edge case
	}

	for _, tt := range tests {
		result := colorizeSINR(tt.val)
		if !strings.Contains(result, tt.expected) {
			t.Errorf("colorizeSINR(%d) = %q; want containing %q", tt.val, result, tt.expected)
		}
	}
}

// TestColorizeBars tests the threshold logic for Bars colors
func TestColorizeBars(t *testing.T) {
	tests := []struct {
		val      float64
		expected string
	}{
		{5.0, ColorGreen},  // Excellent >= 4
		{4.0, ColorGreen},  // Edge case >= 4
		{3.0, ColorYellow}, // Good >= 2
		{2.0, ColorYellow}, // Edge case >= 2
		{1.0, ColorRed},    // Poor < 2
	}

	for _, tt := range tests {
		result := colorizeBars(tt.val)
		if !strings.Contains(result, tt.expected) {
			t.Errorf("colorizeBars(%.1f) = %q; want containing %q", tt.val, result, tt.expected)
		}
	}
}

// TestFetchStats mocks the T-Mobile gateway response and verifies parsing
func TestFetchStats(t *testing.T) {
	// Mock JSON Response
	mockJSON := `{
		"device": {
			"model": "TMO-G5AR",
			"serial": "***REMOVED***",
			"hardwareVersion": "R01",
			"softwareVersion": "1.00.02"
		},
		"signal": {
			"5g": {
				"bands": ["n41"],
				"bars": 3.0,
				"rsrp": -104,
				"sinr": 3,
				"gnbid": 1870191
			},
			"4g": {
				"bands": ["b66"],
				"bars": 4.0,
				"rsrp": -95,
				"sinr": 8,
				"eid": 12345
			},
			"generic": {
				"apn": "FBB.HOME"
			}
		}
	}`

	// Create Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	// Create Client
	client := &http.Client{Timeout: 2 * time.Second}

	// Call gateway.FetchStats using the mock server URL
	data, err := gateway.FetchStats(client, server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify Data
	if data.Device.Model != "TMO-G5AR" {
		t.Errorf("Expected model TMO-G5AR, got %s", data.Device.Model)
	}
	if len(data.Signal.FiveG.Bands) != 1 || data.Signal.FiveG.Bands[0] != "n41" {
		t.Errorf("Expected 5G band n41, got %v", data.Signal.FiveG.Bands)
	}
	if data.Signal.FiveG.RSRP != -104 {
		t.Errorf("Expected 5G RSRP -104, got %d", data.Signal.FiveG.RSRP)
	}
	if data.Signal.FourG.RSRP != -95 {
		t.Errorf("Expected 4G RSRP -95, got %d", data.Signal.FourG.RSRP)
	}
}