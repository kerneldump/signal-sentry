package analysis

import (
	"fmt"
	"time"
)

// TimeFilter defines a time range [Start, End].
// A zero time means no boundary.
type TimeFilter struct {
	Start time.Time
	End   time.Time
}

// Contains returns true if t is within the [Start, End] inclusive range.
// Zero boundaries are treated as open.
func (f *TimeFilter) Contains(t time.Time) bool {
	if f == nil {
		return true
	}
	if !f.Start.IsZero() && t.Before(f.Start) {
		return false
	}
	if !f.End.IsZero() && t.After(f.End) {
		return false
	}
	return true
}

// NewTimeFilter constructs a filter based on user inputs.
// Precedence:
// 1. If rangeDur > 0, Start = Now - rangeDur, End = Now.
// 2. Else, parse start/end strings.
func NewTimeFilter(startStr, endStr string, rangeDur time.Duration) (*TimeFilter, error) {
	// 1. Handle Relative Range
	if rangeDur > 0 {
		now := time.Now()
		return &TimeFilter{
			Start: now.Add(-rangeDur),
			End:   now,
		}, nil
	}

	// 2. Handle Explicit Start/End
	f := &TimeFilter{}
	var err error

	if startStr != "" {
		f.Start, err = ParseISOTime(startStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %w", err)
		}
	}

	if endStr != "" {
		f.End, err = ParseISOTime(endStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
	}

	return f, nil
}

// ParseISOTime attempts to parse standard YYYY-MM-DD or YYYY-MM-DD HH:MM:SS formats.
func ParseISOTime(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse time %q", s)
}
