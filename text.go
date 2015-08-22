package ssvgc

import (
	"encoding/xml"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

type Text struct {
	commonElement

	fontSize  float64
	ttfFont   string
	textValue string

	font *truetype.Font
}

func NewText() *Text {
	t := &Text{}
	t.fillColor = image.Black
	t.strokeColor = image.Transparent
	t.fontSize = 12

	return t
}

func (t *Text) SetAttribute(name, value string) {
	switch name {
	case "font-size":
		t.fontSize, _ = strconv.ParseFloat(value, 64)
	case "ttf-font":
		t.ttfFont = value
		t.font = nil
	case "text-value":
		t.textValue = value
	default:
		t.commonElement.SetAttribute(name, value)
	}
	t.upToDate = false
}

func (t *Text) ParseAttributes(start *xml.StartElement) {
	for _, attr := range start.Attr {
		t.SetAttribute(attr.Name.Local, attr.Value)
	}
}

func (t *Text) Draw() image.Image {
	if t.canvas != nil && t.upToDate {
		return t.canvas
	}
	t.font = t.loadFont(t.ttfFont)

	width, height := t.GetStringSize(t.textValue)
	t.SetAttribute("width", strconv.Itoa(width))
	t.SetAttribute("height", strconv.Itoa(height))

	bounds := t.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, image.Transparent, image.ZP, draw.Src)

	ctx := freetype.NewContext()
	ctx.SetFont(t.font)
	ctx.SetFontSize(t.fontSize)
	ctx.SetSrc(&image.Uniform{t.fillColor})
	ctx.SetDst(img)
	ctx.SetClip(bounds)
	ctx.DrawString(t.textValue, freetype.Pt(bounds.Min.X, bounds.Max.Y))

	return img
}

func (t *Text) loadFont(path string) *truetype.Font {
	if t.font != nil {
		return t.font
	}

	ttfBytes, err := ioutil.ReadFile(t.ttfFont)
	if err != nil {
		log.Fatal(err)
	}

	font, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		log.Fatal(err)
	}

	return font
}

func (t *Text) GetStringSize(s string) (int, int) {
	// Assumes 72 DPI
	fupe := fixed.Int26_6(t.fontSize * 64.0)
	width, height := fixed.I(0), fixed.I(0)

	prev, hasPrev := t.font.Index(0), false
	for _, r := range s {
		idx := t.font.Index(r)
		if hasPrev {
			width += t.font.Kerning(fupe, prev, idx)
		}

		width += t.font.HMetric(fupe, idx).AdvanceWidth
		h := t.font.VMetric(fupe, idx).AdvanceHeight
		if h > height {
			height = h
		}
		prev, hasPrev = idx, true
	}

	return int(width >> 6), int(height >> 6)
}
