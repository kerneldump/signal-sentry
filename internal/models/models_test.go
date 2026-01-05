package models

import (
	"encoding/json"
	"testing"
)

func TestGatewayResponse_Unmarshal(t *testing.T) {
	jsonStr := `{
		"device": {
			"model": "TMO-G5AR",
			"softwareVersion": "1.00.02"
		},
		"signal": {
			"5g": {
				"bands": ["n41"],
				"bars": 4.0,
				"rsrp": -90,
				"sinr": 15
			},
			"4g": {
				"bands": ["b2"],
				"bars": 3.0
			}
		}
	}`

	var data GatewayResponse
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}


	if data.Device.Model != "TMO-G5AR" {
		t.Errorf("Expected model 'TMO-G5AR', got '%s'", data.Device.Model)
	}
	if len(data.Signal.FiveG.Bands) == 0 || data.Signal.FiveG.Bands[0] != "n41" {
		t.Errorf("Expected 5G band 'n41', got %v", data.Signal.FiveG.Bands)
	}
}
