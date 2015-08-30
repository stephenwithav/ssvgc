package ssvgc_test

import (
	"strings"
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/stephenwithav/ssvgc"
)

func TestTextRendering(t *testing.T) {
	var tests = []StringMap{
		{
			"name": "textredbanana",
			"text-value": `<svg width="200" height="200" fill="white">
<text ttf-font="/tmp/DejaVuSans.ttf" x="15" y="30">You are
    <tspan fill="red">not</tspan>
    a banana
</text>
</svg>`,
		},
		{
			"name": "textbananashiftedright",
			"text-value": `<svg width="200" height="200" fill="white">
<text ttf-font="/tmp/DejaVuSans.ttf" x="45" y="30">You are
    <tspan>not</tspan>
    a banana
</text>
</svg>`,
		},
		{
			"name": "textbananashifteddown",
			"text-value": `<svg width="200" height="200" fill="white">
<text ttf-font="/tmp/DejaVuSans.ttf" font-size="18" x="1" y="70">You are
    <tspan>not</tspan>
    a banana
</text>
</svg>`,
		},
		{
			"name": "textbananaanchormiddle",
			"text-value": `<svg width="500" height="200" fill="white">
<text ttf-font="/tmp/DejaVuSans.ttf" font-size="18" x="250" y="70" text-anchor="middle">You are
    <tspan>not</tspan>
    a banana
</text>
</svg>`,
		},
		{
			"name": "textbananaanchorend",
			"text-value": `<svg width="500" height="200" fill="white">
<text ttf-font="/tmp/DejaVuSans.ttf" font-size="18" x="250" y="70" text-anchor="end">You are
    <tspan>not</tspan>
    a banana
</text>
</svg>`,
		},
	}

	for _, tt := range tests {
		r := strings.NewReader(tt["text-value"])
		p := ssvgc.NewParser(r)
		svg, err := p.ParseSVG()
		if err != nil {
			t.Errorf("ParseSVG() failed => %s", err)
			continue
		}

		draw2dimg.SaveToPngFile("output/"+tt["name"]+".png", svg.Draw())
	}
}
