package ssvgc_test

import (
	"os"
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/stephenwithav/ssvgc"
)

func TestParseSVG(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping lengthier tests during short test run.")
	}

	var tests = []string{
		"rectfull",
		"rectfullwithtext",
		"tworects",
		"threerects",
	}

	for _, tt := range tests {
		f, err := os.Open("testsvgs/" + tt + ".svg")
		if err != nil {
			t.Errorf("Failed to open %s.svg for reading: %s", tt, err)
			continue
		}
		defer f.Close()

		p := ssvgc.NewParser(f)
		svg, err := p.ParseSVG()
		if err != nil {
			t.Errorf("ParseSVG(%s.svg) failed => %s", tt, err)
			continue
		}

		img := svg.Draw()
		draw2dimg.SaveToPngFile("output/"+tt+".png", img)
	}
}
