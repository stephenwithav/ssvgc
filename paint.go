package ssvgc

import (
	"image/color"
	"strconv"
)

func colorFromPaintToken(s string) color.Color {
	// #rgb or #rrggbb
	if s[:1] == "#" {
		var rgbStr [6]byte
		switch len(s) {
		case 4:
			rgbStr[0], rgbStr[1], rgbStr[2], rgbStr[3], rgbStr[4], rgbStr[5] = s[1], s[1], s[2], s[2], s[3], s[3]
		case 7:
			copy(rgbStr[:], s[1:])
		default:
			return color.Black
		}

		if rgbInt, err := strconv.ParseUint(string(rgbStr[:]), 16, 32); err == nil {
			return color.RGBA{uint8(rgbInt >> 16), uint8((rgbInt >> 8) & 0xFF), uint8(rgbInt & 0xFF), 255}
		}

		return color.Black
	}

	// none
	if s == "none" {
		return color.Transparent
	}

	// officially recognized colors
	if col, ok := recognizedColors[s]; ok {
		return col
	}

	return color.Black
}
