package analysis

import (
	"bytes"
	"strings"
	"testing"
)

func TestAnalyzeEndToEnd(t *testing.T) {
	// 1. Setup Test Data
	jsonInput := `
{"gateway":{"time":{"localTime":1767651600},"signal":{"5g":{"bands":["n41"],"bars":3.0,"rsrp":-100,"sinr":5,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0}}
{"gateway":{"time":{"localTime":1767651660},"signal":{"5g":{"bands":["n41"],"bars":4.0,"rsrp":-90,"sinr":10,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0}}
`
	input := strings.NewReader(strings.TrimSpace(jsonInput))
	var output bytes.Buffer

	// 2. Run Analysis
	err := Analyze(input, &output)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 3. Verify Output
	result := output.String()

	checks := []string{
		"HISTORICAL SIGNAL ANALYSIS",
		"Total Samples: 2",
		"RSRP (dBm)", "-100", "-95.0", "-90", // Min, Avg, Max
		"SINR (dB)", "5", "7.5", "10",
		"BARS SEEN:",
		"3", "1 samples (50.0%)",
		"4", "1 samples (50.0%)",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", check, result)
		}
	}
}

func TestAnalyzeEmptyInput(t *testing.T) {
	input := strings.NewReader("")
	var output bytes.Buffer

	err := Analyze(input, &output)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "No data samples found") {
		t.Errorf("Expected 'No data samples found', got:\n%s", result)
	}
}
