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
	textContext
	chunks []*textChunk
	buf    []byte
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
	if len(t.chunks) == 0 {
		return image.Transparent
	}

	if t.canvas != nil && t.upToDate {
		return t.canvas
	}

	width, height, maxFont := fixed.I(0), 0, 0.0

	for _, chunk := range t.chunks {
		t.textContext = chunk.textContext
		t.font = t.loadFont(t.ttfFont)

		s := string(t.buf[chunk.startPos:chunk.endPos])
		cWidth, _ := t.GetStringSize(s)
		width += cWidth

		cHeight := t.getMaxHeight()
		if cHeight > height {
			height = cHeight
		}

		if t.fontSize > maxFont {
			maxFont = t.fontSize
		}
	}

	t.textContext = t.chunks[0].textContext
	glyphOffset := height - int(maxFont)
	bounds := image.Rect(t.xOffset, t.yOffset-int(maxFont), t.xOffset+int(width>>6), t.yOffset+glyphOffset)

	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, image.Transparent, image.ZP, draw.Src)

	ctx := freetype.NewContext()
	ctx.SetDst(img)
	ctx.SetClip(bounds)

	pt := freetype.Pt(bounds.Min.X, bounds.Min.Y+int(t.fontSize))
	t.font = t.loadFont(t.ttfFont)
	for _, chunk := range t.chunks {
		ctx.SetFont(t.font)
		ctx.SetFontSize(chunk.fontSize)
		ctx.SetSrc(&image.Uniform{chunk.fillColor})
		pt, _ = ctx.DrawString(string(t.buf[chunk.startPos:chunk.endPos]), pt)
	}

	return img
}

func (t *Text) GetStringSize(s string) (fixed.Int26_6, fixed.Int26_6) {
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
	return width, (fontBounds.YMax - fontBounds.YMin)
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
	bounds := image.Rect(0, 0, int(w>>6), int(h>>6))
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

type textContext struct {
	commonElement

	fontSize float64
	ttfFont  string

	font *truetype.Font
}

type textChunk struct {
	textContext
	startPos int
	endPos   int
}
