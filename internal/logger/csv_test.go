package logger

import (
	"os"
	"strings"
	"testing"
	"tmobile-stats/internal/gateway"
)

func TestCSVLogger(t *testing.T) {
	tmpFile := "test-signal.csv"
	defer os.Remove(tmpFile)

	l, err := NewCSVLogger(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create CSVLogger: %v", err)
	}

	data := &gateway.GatewayResponse{
		Signal: gateway.SignalInfo{
			FiveG: gateway.ConnectionStats{Bands: []string{"n41"}, Bars: 4.0, RSRP: -90, SINR: 15},
			FourG: gateway.ConnectionStats{Bands: []string{"b2"}, Bars: 3.0, RSRP: -100, SINR: 5},
		},
	}

	err = l.Log(data)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	l.Close()

	// Verify file content
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines (header + data), got %d", len(lines))
	}

	// Basic header check
	if !strings.Contains(lines[0], "Timestamp") || !strings.Contains(lines[0], "5G_Band") {
		t.Errorf("Header seems incorrect: %s", lines[0])
	}
}
