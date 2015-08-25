package ssvgc_test

import (
	"strings"
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/stephenwithav/ssvgc"
)

func TestImageRendering(t *testing.T) {
	var tests = []StringMap{
		{
			"name": "tuxcentered",
			"xml": `<svg width="200" height="200" fill="white">
<image href="testimages/tux-construction.png" width="150" x="25" y="25" />
</svg>`,
		},
		{
			"name": "tuxbottomright",
			"xml": `<svg width="200" height="200" fill="white">
<image href="testimages/tux-construction.png" width="150" x="50" y="50" />
</svg>`,
		},
		{
			"name": "tuxfullsize",
			"xml": `<svg width="200" height="200" fill="white">
<image href="testimages/tux-construction.png" />
</svg>`,
		},
	}

	for _, tt := range tests {
		r := strings.NewReader(tt["xml"])
		p := ssvgc.NewParser(r)
		svg, err := p.ParseSVG()
		if err != nil {
			t.Errorf("ParseSVG() failed => %s", err)
			continue
		}

		draw2dimg.SaveToPngFile("output/"+tt["name"]+".png", svg.Draw())
	}
}
