package analysis

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAnalyzeEndToEnd(t *testing.T) {
	// 1. Setup Test Data
	jsonInput := `
{"gateway":{"time":{"localTime":1767651600},"signal":{"5g":{"bands":["n41"],"bars":3.0,"rsrp":-100,"sinr":5,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0,"sent":10,"received":10}}
{"gateway":{"time":{"localTime":1767651660},"signal":{"5g":{"bands":["n41"],"bars":4.0,"rsrp":-90,"sinr":10,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0,"sent":10,"received":10}}
`
	input := strings.NewReader(strings.TrimSpace(jsonInput))
	var output bytes.Buffer

	// 2. Run Analysis
	err := Analyze(input, &output, nil)
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
		"4", "1 samples (50.0%) real-time",
		"Loss (%)", "-", "0.0", "-",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", check, result)
		}
	}
}

func TestAnalyzeLiveTower(t *testing.T) {
	// 1. Setup Test Data with two towers
	// First sample: Tower 100
	// Second sample: Tower 200
	jsonInput := `
{"gateway":{"time":{"localTime":1767651600},"signal":{"5g":{"bands":["n41"],"bars":3.0,"rsrp":-100,"sinr":5,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0}}
{"gateway":{"time":{"localTime":1767651660},"signal":{"5g":{"bands":["n41"],"bars":4.0,"rsrp":-90,"sinr":10,"gNBID":200},"4g":{}}},"ping":{"min":20,"loss":0}}
`
	input := strings.NewReader(strings.TrimSpace(jsonInput))
	var output bytes.Buffer

	// 2. Run Analysis
	err := Analyze(input, &output, nil)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 3. Verify Output
	result := output.String()

	expectedLive := "200        1 samples (50.0%) live"
	if !strings.Contains(result, expectedLive) {
		t.Errorf("Expected output to contain live tower marker %q.\nOutput:\n%s", expectedLive, result)
	}

	unexpectedLive := "100        1 samples (50.0%) live"
	if strings.Contains(result, unexpectedLive) {
		t.Errorf("Expected Tower 100 NOT to be marked live, but it was.\nOutput:\n%s", result)
	}
}

func TestAnalyzeLossCalculation(t *testing.T) {
	// Sample 1: 10 sent, 9 received (10% loss)
	// Sample 2: 10 sent, 10 received (0% loss)
	// Total: 20 sent, 19 received (1 lost). Global Loss = 5.0%
	jsonInput := `
{"gateway":{"time":{"localTime":1767651600},"signal":{"5g":{"bands":["n41"],"bars":3.0,"rsrp":-100,"sinr":5,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":10,"sent":10,"received":9}}
{"gateway":{"time":{"localTime":1767651660},"signal":{"5g":{"bands":["n41"],"bars":4.0,"rsrp":-90,"sinr":10,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0,"sent":10,"received":10}}
`
	input := strings.NewReader(strings.TrimSpace(jsonInput))
	var output bytes.Buffer

	err := Analyze(input, &output, nil)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	result := output.String()
	// Check for "Loss (%)" row with "-" "5.0" "-"
	// Using fields/contains might be safer than exact whitespace matching
	lines := strings.Split(result, "\n")
	foundLoss := false
	for _, line := range lines {
		if strings.Contains(line, "Loss (%)") {
			foundLoss = true
			if !strings.Contains(line, "-") || !strings.Contains(line, "5.0") {
				t.Errorf("Expected Loss line to contain '-' and '5.0', got: %q", line)
			}
		}
	}
	if !foundLoss {
		t.Error("Did not find Loss (%) row in output")
	}
}

func TestAnalyzeEmptyInput(t *testing.T) {
	input := strings.NewReader("")
	var output bytes.Buffer

	err := Analyze(input, &output, nil)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "No data samples found") {
		t.Errorf("Expected 'No data samples found', got:\n%s", result)
	}
}

func TestAnalyzeWithFilter(t *testing.T) {
	// Sample 1: 10:00:00
	// Sample 2: 11:00:00
	// Sample 3: 12:00:00
	jsonInput := `
{"gateway":{"time":{"localTime":1736157600},"signal":{"5g":{"bands":["n41"],"bars":3.0,"rsrp":-100,"sinr":5,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0}}
{"gateway":{"time":{"localTime":1736161200},"signal":{"5g":{"bands":["n41"],"bars":4.0,"rsrp":-90,"sinr":10,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0}}
{"gateway":{"time":{"localTime":1736164800},"signal":{"5g":{"bands":["n41"],"bars":5.0,"rsrp":-80,"sinr":15,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0}}
`
	input := strings.NewReader(strings.TrimSpace(jsonInput))
	
	// Filter for 10:30 to 11:30 (Only Sample 2 should remain)
	filter := &TimeFilter{
		Start: time.Unix(1736159400, 0),
		End:   time.Unix(1736163000, 0),
	}
	
	var output bytes.Buffer
	err := Analyze(input, &output, filter)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "Total Samples: 1") {
		t.Errorf("Expected 1 sample, got:\n%s", result)
	}
	if !strings.Contains(result, "-90") {
		t.Errorf("Expected sample with RSRP -90, but not found.\nOutput:\n%s", result)
	}
}

func TestAnalyzeRealTimeBars(t *testing.T) {
	// 1. Setup Test Data
	// Last sample has 5 bars
	jsonInput := `
{"gateway":{"time":{"localTime":1767651600},"signal":{"5g":{"bands":["n41"],"bars":3.0,"rsrp":-100,"sinr":5,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0}}
{"gateway":{"time":{"localTime":1767651660},"signal":{"5g":{"bands":["n41"],"bars":5.0,"rsrp":-80,"sinr":15,"gNBID":100},"4g":{}}},"ping":{"min":20,"loss":0}}
`
	input := strings.NewReader(strings.TrimSpace(jsonInput))
	var output bytes.Buffer

	// 2. Run Analysis
	err := Analyze(input, &output, nil)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 3. Verify Output
	result := output.String()

	expectedRealTime := "5          1 samples (50.0%) real-time"
	if !strings.Contains(result, expectedRealTime) {
		t.Errorf("Expected output to contain real-time bar marker %q.\nOutput:\n%s", expectedRealTime, result)
	}

	unexpectedRealTime := "3          1 samples (50.0%) real-time"
	if strings.Contains(result, unexpectedRealTime) {
		t.Errorf("Expected Bars 3 NOT to be marked real-time, but it was.\nOutput:\n%s", result)
	}
}