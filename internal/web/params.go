package web

import (
	"net/http"
	"time"
	"tmobile-stats/internal/analysis"
)

func parseTimeFilter(r *http.Request) (*analysis.TimeFilter, string, error) {
	rangeStr := r.URL.Query().Get("range")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	// Precedence: If start or end is provided, absolute range takes precedence.
	if startStr != "" || endStr != "" {
		filter, err := analysis.NewTimeFilter(startStr, endStr, 0)
		return filter, rangeStr, err
	}

	// Otherwise, relative range
	var rangeDur time.Duration
	var err error

	if rangeStr == "0" || rangeStr == "max" {
		rangeDur = 0
	} else if rangeStr != "" {
		rangeDur, err = time.ParseDuration(rangeStr)
		if err != nil {
			// Default fallback for invalid duration
			rangeDur = 24 * time.Hour
		}
	} else {
		// Default to 24h
		rangeDur = 24 * time.Hour
		rangeStr = "24h"
	}

	filter, err := analysis.NewTimeFilter("", "", rangeDur)
	return filter, rangeStr, err
}
