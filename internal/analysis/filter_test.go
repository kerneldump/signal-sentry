package analysis

import (
	"testing"
	"time"
)

func TestNewTimeFilter_RangePrecedence(t *testing.T) {
	// If range is provided, it should ignore start/end strings
	dur := 1 * time.Hour
	f, err := NewTimeFilter("2020-01-01", "2020-01-02", dur)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify Start is approx Now - 1h
	// We allow a small delta for execution time
	now := time.Now()
	expectedStart := now.Add(-dur)
	
	if f.Start.Sub(expectedStart).Abs() > time.Second {
		t.Errorf("Expected Start ~ %v, got %v", expectedStart, f.Start)
	}
	if f.End.Sub(now).Abs() > time.Second {
		t.Errorf("Expected End ~ %v, got %v", now, f.End)
	}
}

func TestNewTimeFilter_NegativeDuration(t *testing.T) {
	// If negative range is provided, it should be treated as positive magnitude
	dur := -1 * time.Hour
	f, err := NewTimeFilter("", "", dur)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify Start is approx Now - 1h (magnitude of -1h)
	now := time.Now()
	expectedDur := 1 * time.Hour
	expectedStart := now.Add(-expectedDur)
	
	if f.Start.Sub(expectedStart).Abs() > time.Second {
		t.Errorf("Expected Start ~ %v, got %v", expectedStart, f.Start)
	}
	if f.End.Sub(now).Abs() > time.Second {
		t.Errorf("Expected End ~ %v, got %v", now, f.End)
	}
}

func TestNewTimeFilter_ExplicitDates(t *testing.T) {
	startStr := "2025-01-01"
	endStr := "2025-01-02 12:00:00"
	
	f, err := NewTimeFilter(startStr, endStr, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify Start
	expectedStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
	if !f.Start.Equal(expectedStart) {
		t.Errorf("Expected Start %v, got %v", expectedStart, f.Start)
	}

	// Verify End
	expectedEnd := time.Date(2025, 1, 2, 12, 0, 0, 0, time.Local)
	if !f.End.Equal(expectedEnd) {
		t.Errorf("Expected End %v, got %v", expectedEnd, f.End)
	}
}

func TestNewTimeFilter_InvalidFormat(t *testing.T) {
	_, err := NewTimeFilter("bad-date", "", 0)
	if err == nil {
		t.Error("Expected error for invalid date, got nil")
	}
}

func TestTimeFilter_Contains(t *testing.T) {
	start := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	f := &TimeFilter{Start: start, End: end}

	tests := []struct {
		input    time.Time
		expected bool
	}{
		{time.Date(2025, 1, 1, 9, 59, 0, 0, time.UTC), false},  // Before
		{time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC), true},   // Exact Start
		{time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC), true},   // Middle
		{time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC), true},   // Exact End
		{time.Date(2025, 1, 1, 12, 0, 1, 0, time.UTC), false},  // After
	}

	for _, tt := range tests {
		if got := f.Contains(tt.input); got != tt.expected {
			t.Errorf("Contains(%v) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestTimeFilter_OpenEnded(t *testing.T) {
	// Start set, End zero (open)
	start := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	f := &TimeFilter{Start: start}

	if !f.Contains(start.Add(24 * time.Hour)) {
		t.Error("Open-ended filter should contain future time")
	}
	if f.Contains(start.Add(-1 * time.Hour)) {
		t.Error("Open-ended filter should NOT contain past time")
	}
}