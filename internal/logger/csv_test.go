package logger

import (
	"os"
	"strings"
	"testing"
	"tmobile-stats/internal/models"
)

func TestCSVLogger(t *testing.T) {
	tmpFile := "test-signal.csv"
	defer os.Remove(tmpFile)

	l, err := NewCSVLogger(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create CSVLogger: %v", err)
	}

	data := &models.CombinedStats{
		Gateway: models.GatewayResponse{
			Signal: models.SignalInfo{
				FiveG: models.ConnectionStats{Bands: []string{"n41"}, Bars: 4.0, RSRP: -90, SINR: 15},
				FourG: models.ConnectionStats{Bands: []string{"b2"}, Bars: 3.0, RSRP: -100, SINR: 5},
			},
		},
		Ping: models.PingStats{Min: 10.5, Avg: 12.0, Max: 15.0, StdDev: 1.5, Loss: 0.0},
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
	if !strings.Contains(lines[0], "Ping_Min") || !strings.Contains(lines[0], "Ping_Loss") {
		t.Errorf("Header seems incorrect (missing ping columns): %s", lines[0])
	}
}

