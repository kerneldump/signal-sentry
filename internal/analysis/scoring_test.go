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

		// 2. Worst Conditions
		// RSRP -120 maps to >= -120 (2.0)
		// SINR -15 maps to < -10 (1.0)
		// Score = 1.4 + 0.3 = 1.7
		{"Worst", -120, -15, 1.7},

		// 3. Mixed: Good RSRP (4.0), Poor SINR (2.0)
		// RSRP -85 maps to > -90 (5.0) in new generous mapping! 
		// Wait, -85 is > -90. Correct.
		// SINR -5 maps to (-10 to 0) -> 2.0. (Was 2.0 before too? < 0 -> 3.0? No, 0 to 10 is 3. -10 to 0 is 2.)
		// Calculation: (5.0 * 0.7) + (2.0 * 0.3) = 3.5 + 0.6 = 4.1
		{"GoodRSRP_PoorSINR", -85, -5, 4.1},

		// 4. Mixed: Poor RSRP (2.0), Excellent SINR (5.0)
		// RSRP -105 maps to (-100 to -110) -> 3.0
		// SINR 22 maps to > 20 -> 5.0
		// Score = (3.0 * 0.7) + (5.0 * 0.3) = 2.1 + 1.5 = 3.6
		{"PoorRSRP_ExcSINR", -105, 22, 3.6},

		// 5. Boundary Checks
		// RSRP -80 is > -90 -> 5.0
		// SINR 25 is > 20 -> 5.0
		// Score = 5.0
		{"Boundary_RSRP_Minus80", -80, 25, 5.0},
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
