package analysis

import "math"

// CalculateSignalHealth computes a weighted score (1.0 - 5.0) representing signal quality.
// Formula: (RSRP_Score * 0.7) + (SINR_Score * 0.3)
func CalculateSignalHealth(rsrp, sinr int) float64 {
	rsrpScore := scoreRSRP(rsrp)
	sinrScore := scoreSINR(sinr)

	score := (rsrpScore * 0.7) + (sinrScore * 0.3)
	return math.Round(score*10) / 10 // Round to 1 decimal place
}

// scoreRSRP maps RSRP values to a 1.0-5.0 scale.
// Generous mapping for 5G mid-band.
func scoreRSRP(val int) float64 {
	switch {
	case val > -90:
		return 5.0
	case val >= -100:
		return 4.0
	case val >= -110:
		return 3.0
	case val >= -120:
		return 2.0
	default:
		return 1.0
	}
}

// scoreSINR maps SINR values to a 1.0-5.0 scale.
func scoreSINR(val int) float64 {
	switch {
	case val > 20:
		return 5.0
	case val >= 10:
		return 4.0
	case val >= 0:
		return 3.0
	case val >= -10:
		return 2.0
	default:
		return 1.0
	}
}
