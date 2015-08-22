package ssvgc_test

import (
	"encoding/xml"
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
	}

	for _, tt := range tests {
		r := &ssvgc.Rectangle{}
		r.ParseAttributes(&xml.StartElement{})
		for name, value := range tt {
			r.SetAttribute(name, value)
		}

		draw2dimg.SaveToPngFile("output/"+tt["name"]+".png", r.Draw())
	}
}
