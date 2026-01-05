package logger

import (
	"os"
	"testing"
	"tmobile-stats/internal/models"
)

func TestJSONLogger(t *testing.T) {
	tmpFile := "test-signal.json"
	defer os.Remove(tmpFile)

	l, err := NewJSONLogger(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create JSONLogger: %v", err)
	}

	data := &models.CombinedStats{
		Gateway: models.GatewayResponse{
			Device: models.DeviceInfo{Model: "TEST"},
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

	if len(content) == 0 {
		t.Error("Log file is empty")
	}
}
