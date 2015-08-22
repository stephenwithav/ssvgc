package ssvgc_test

import (
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/stephenwithav/ssvgc"
)

type StringMap map[string]string

func TestRectangleDrawing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping lengthier tests during short test run.")
	}

	var tests = []StringMap{
		{
			"name":   "solidblue",
			"fill":   "blue",
			"width":  "400",
			"height": "400",
		},
		{
			"name":         "redinsetwithblue",
			"fill":         "blue",
			"stroke":       "red",
			"stroke-width": "30",
			"width":        "400",
			"height":       "400",
		},
		{
			"name":         "redinsetwithblue2",
			"fill":         "blue",
			"stroke":       "red",
			"stroke-width": "1",
			"width":        "400",
			"height":       "400",
		},
		{
			"name":         "solidredstrokedgreen",
			"fill":         "red",
			"width":        "100",
			"height":       "100",
			"x":            "50",
			"y":            "50",
			"stroke":       "green",
			"stroke-width": "15",
		},
	}

	for _, tt := range tests {
		r := ssvgc.NewRectangle()
		for name, value := range tt {
			r.SetAttribute(name, value)
		}

		draw2dimg.SaveToPngFile("output/"+tt["name"]+".png", r.Draw())
	}
}
