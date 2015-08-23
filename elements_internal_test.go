package ssvgc

import "testing"

func TestParseFloatUnits(t *testing.T) {
	var tests = []struct {
		in  string
		out float64
	}{
		{"1", 1.0},
		{"1pt", 1.25},
		{"1pc", 15.0},
		{"1mm", 3.543307},
		{"1cm", 35.43307},
		{"1in", 90.0},
		{"1zz", 1.0},
		{"foobar", 0.0},
	}

	for _, tt := range tests {
		got := parseFloatUnit(tt.in)
		if got != tt.out {
			t.Errorf("parseFloatUnit(%q) => %v, expected %v", tt.in, got, tt.out)
		}
	}
}

func TestParseUnits(t *testing.T) {
	var tests = []struct {
		in  string
		out int
	}{
		{"1", 1},
		{"1pt", 1},
		{"1pc", 15},
		{"1mm", 3},
		{"1cm", 35},
		{"1in", 90},
		{"1zz", 1},
		{"foobar", 0},
	}

	for _, tt := range tests {
		got := parseUnit(tt.in)
		if got != tt.out {
			t.Errorf("parseUnit(%q) => %v, expected %v", tt.in, got, tt.out)
		}
	}
}
