package ssvgc_test

import (
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/stephenwithav/ssvgc"
)

func TestTextRendering(t *testing.T) {
	var tests = []StringMap{
		{
			"name":       "helloworld30red",
			"fill":       "red",
			"font-size":  "30",
			"ttf-font":   "/tmp/DejaVuSans.ttf",
			"text-value": "Hello world!",
			"y":          "12",
		},
		{
			"name":       "helloworld18red",
			"fill":       "red",
			"font-size":  "18",
			"ttf-font":   "/tmp/DejaVuSans.ttf",
			"text-value": "Hello world!",
			"y":          "18",
		},
		{
			"name":       "joytotheworld12red",
			"fill":       "red",
			"font-size":  "12",
			"ttf-font":   "/tmp/DejaVuSans.ttf",
			"text-value": "joy to the world!",
			"y":          "12",
		},
	}

	for _, tt := range tests {
		txt := ssvgc.NewText()
		for name, value := range tt {
			txt.SetAttribute(name, value)
		}

		draw2dimg.SaveToPngFile("output/"+tt["name"]+".png", txt.Draw())
	}
}
