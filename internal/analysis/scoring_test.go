package analysis

import "testing"

func TestCalculateSignalHealth(t *testing.T) {
	tests := []struct {
		name     string
		rsrp     int
		sinr     int
		expected float64
	}{
		// 1. Perfect Conditions (5.0, 5.0) -> 5.0
		{"Perfect", -70, 25, 5.0},

		// 2. Worst Conditions (1.0, 1.0) -> 1.0
		{"Worst", -120, -15, 1.0},

		// 3. Mixed: Good RSRP (4.0), Poor SINR (2.0)
		// Score = (4.0 * 0.7) + (2.0 * 0.3) = 2.8 + 0.6 = 3.4
		{"GoodRSRP_PoorSINR", -85, -5, 3.4},

		// 4. Mixed: Poor RSRP (2.0), Excellent SINR (5.0)
		// Score = (2.0 * 0.7) + (5.0 * 0.3) = 1.4 + 1.5 = 2.9
		{"PoorRSRP_ExcSINR", -105, 22, 2.9},

		// 5. Boundary Checks
		{"Boundary_RSRP_Minus80", -80, 25, 4.3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSignalHealth(tt.rsrp, tt.sinr)
			if got != tt.expected {
				t.Errorf("CalculateSignalHealth(%d, %d) = %v; want %v", tt.rsrp, tt.sinr, got, tt.expected)
			}
		})
	}
}
