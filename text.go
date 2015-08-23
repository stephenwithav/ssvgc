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

	// Go needs to know the top-left point to begin drawing at.
	// Freetype and SVG need to know the bottom-left point
	// to begin drawing at.  This adjustment gives everybody
	// what they want.
	t.yOffset -= int(t.fontSize)
}

func (t *Text) Draw() image.Image {
	if t.canvas != nil && t.upToDate {
		return t.canvas
	}
	t.font = t.loadFont(t.ttfFont)

	width, _ := t.GetStringSize(t.textValue)
	t.SetAttribute("width", strconv.Itoa(width))
	t.SetAttribute("height", strconv.Itoa(t.getMaxHeight()))

	bounds := t.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, image.Transparent, image.ZP, draw.Src)

	ctx := freetype.NewContext()
	ctx.SetFont(t.font)
	ctx.SetFontSize(t.fontSize)
	ctx.SetSrc(&image.Uniform{t.fillColor})
	ctx.SetDst(img)
	ctx.SetClip(bounds)
	ctx.DrawString(t.textValue, freetype.Pt(bounds.Min.X, bounds.Min.Y+int(t.fontSize)))

	return img
}

func (t *Text) GetStringSize(s string) (int, int) {
	// Assumes 72 DPI
	fupe := fixed.Int26_6(t.fontSize * 64.0)
	width := fixed.I(0)

	prev, hasPrev := t.font.Index(0), false
	for _, r := range s {
		idx := t.font.Index(r)
		if hasPrev {
			width += t.font.Kerning(fupe, prev, idx)
		}

		width += t.font.HMetric(fupe, idx).AdvanceWidth
		prev, hasPrev = idx, true
	}

	fontBounds := t.font.Bounds(fupe)
	return int(width >> 6), int((fontBounds.YMax - fontBounds.YMin) >> 6)
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

func (t *Text) getMaxHeight() int {
	w, h := t.GetStringSize("|")
	bounds := image.Rect(0, 0, w, h)
	img := image.NewRGBA(bounds)

	ctx := freetype.NewContext()
	ctx.SetFont(t.font)
	ctx.SetFontSize(t.fontSize)
	ctx.SetSrc(&image.Uniform{t.fillColor})
	ctx.SetDst(img)
	ctx.SetClip(bounds)
	ctx.DrawString("|", freetype.Pt(0, int(t.fontSize)))

	var i = len(img.Pix) - 1
	for ; img.Pix[i] == 0; i-- {

	}

	return (i / img.Stride) + 1
}
