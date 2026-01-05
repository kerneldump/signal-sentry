package main

import "testing"

func TestValidateInterval(t *testing.T) {
	tests := []struct {
		input    int
		wantErr  bool
	}{
		{5, false},   // Default/Valid
		{1, false},   // Valid
		{100, false}, // Valid
		{0, true},    // Invalid (Zero)
		{-1, true},   // Invalid (Negative)
	}

	for _, tt := range tests {
		err := validateInterval(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateInterval(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
	}{
		{"json", false},
		{"csv", false},
		{"", false}, // Default (disabled)
		{"xml", true},
		{"txt", true},
	}

	for _, tt := range tests {
		err := validateFormat(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}