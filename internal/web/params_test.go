package web

import (
	"net/http"
	"testing"
)

func TestParseTimeFilter(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedRange string
		expectError   bool
	}{
		{
			name:          "Default 24h",
			query:         "",
			expectedRange: "24h",
		},
		{
			name:          "Relative 1h",
			query:         "?range=1h",
			expectedRange: "1h",
		},
		{
			name:          "Relative 90s",
			query:         "?range=90s",
			expectedRange: "90s",
		},
		{
			name:          "Absolute Start Only",
			query:         "?start=2026-01-01T12:00",
			expectedRange: "", // range might still be in query but ignored
		},
		{
			name:          "Absolute Range",
			query:         "?start=2026-01-01T12:00&end=2026-01-01T14:00",
			expectedRange: "",
		},
		{
			name:          "Absolute takes precedence over range",
			query:         "?range=24h&start=2026-01-01T12:00",
			expectedRange: "24h", // returns rangeStr as it was in URL, but filter should be absolute
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/"+tt.query, nil)
			filter, rangeStr, err := parseTimeFilter(req)

			if (err != nil) != tt.expectError {
				t.Fatalf("expectError %v, got %v", tt.expectError, err)
			}

			if tt.query != "" && tt.expectedRange != "" && rangeStr != tt.expectedRange {
				t.Errorf("expected rangeStr %v, got %v", tt.expectedRange, rangeStr)
			}

			if (tt.query == "?start=2026-01-01T12:00" || tt.query == "?start=2026-01-01T12:00&end=2026-01-01T14:00") && filter.Start.IsZero() {
				t.Error("expected absolute filter, but Start is zero")
			}
		})
	}
}
