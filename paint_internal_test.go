package ssvgc

import (
	"image/color"
	"testing"
)

func TestColorFromPaintToken(t *testing.T) {
	var tests = []struct {
		in  string
		out color.Color
	}{
		{"red", color.RGBA{0xff, 0, 0, 0xff}},
		{"#fff", color.RGBA{0xff, 0xff, 0xff, 0xff}},
		{"#fe0000", color.RGBA{0xfe, 00, 00, 0xff}},
		{"none", color.Transparent},
		{"qwerty", color.Black},
	}

	for _, tt := range tests {
		col := colorFromPaintToken(tt.in)
		r, g, b, a := col.RGBA()
		rr, gg, bb, aa := tt.out.RGBA()
		if r != rr || g != gg || b != bb || a != aa {
			t.Errorf("colorFromPaintToken(%q) => %v, expected %v", tt.in, col, tt.out)
		}
	}
}
