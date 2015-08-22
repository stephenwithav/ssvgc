package ssvgc_test

import (
	"strings"
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/stephenwithav/ssvgc"
)

func TestSVGDrawing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping lengthier tests during short test run.")
	}

	var rectangleDefinitions = []StringMap{
		{
			"name":   "solidblue",
			"fill":   "blue",
			"width":  "200",
			"height": "200",
			"x":      "0",
			"y":      "0",
		},
		{
			"name":   "solidred",
			"fill":   "red",
			"width":  "100",
			"height": "100",
			"x":      "50",
			"y":      "50",
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

	var tests = []StringMap{
		{
			"name":       "solidbluesvg",
			"rectangles": "solidblue",
			"width":      "200",
			"height":     "200",
			"fill":       "none",
		},
		{
			"name":       "bluewithredsvg",
			"rectangles": "solidblue,solidred",
			"width":      "200",
			"height":     "200",
			"fill":       "none",
		},
		{
			"name":       "bluewithredstrokedgreensvg",
			"rectangles": "solidblue,solidredstrokedgreen",
			"width":      "200",
			"height":     "200",
			"fill":       "none",
		},
	}

	rectangles := make(map[string]*ssvgc.Rectangle)

	for _, def := range rectangleDefinitions {
		r := ssvgc.NewRectangle()
		for name, value := range def {
			r.SetAttribute(name, value)
		}

		rectangles[def["name"]] = r
	}

	for _, tt := range tests {
		s := &ssvgc.SVG{}
		for name, value := range tt {
			s.SetAttribute(name, value)
		}

		for _, rect := range strings.Split(tt["rectangles"], ",") {
			s.AddElement(rectangles[rect])
		}

		draw2dimg.SaveToPngFile("output/"+tt["name"]+".png", s.Draw())
	}
}
