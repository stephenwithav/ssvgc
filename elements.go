package ssvgc

import (
	"encoding/xml"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"unicode"
)

type Element interface {
	// Draw renders and returns an image.Image based on its attributes.
	Draw() image.Image
	// ParseAttributes sets all relevant attributes on the Element.
	ParseAttributes(token *xml.StartElement)
	// Bounds returns the elements size as an image.Rectangle.
	Bounds() image.Rectangle
	// SetAttribute sets the named attribute to its given value.
	SetAttribute(name string, value string)
}

type commonElement struct {
	strokeColor   color.Color
	strokeOpacity color.Color
	strokeWidth   int

	fillColor   color.Color
	fillOpacity color.Color

	xOffset int
	yOffset int

	width  int
	height int

	elementID string

	canvas   draw.Image
	upToDate bool
}

func (e *commonElement) Draw() image.Image {
	return image.Transparent
}

func (e *commonElement) ParseAttributes(start *xml.StartElement) {
	e.parseCommonAttributes(start)
}

func (e *commonElement) Bounds() image.Rectangle {
	return image.Rect(e.xOffset, e.yOffset, e.width+e.xOffset, e.height+e.yOffset)
}

func (e *commonElement) SetAttribute(attr string, val string) {
	switch attr {
	case "stroke":
		e.strokeColor = colorFromPaintToken(val)
	case "stroke-width":
		e.strokeWidth = parseUnit(val)
	case "stroke-opacity":
		e.strokeOpacity = color.Alpha16{uint16(0xffff * parseFloatUnit(val))}
	case "fill":
		e.fillColor = colorFromPaintToken(val)
	case "fill-opacity":
		e.fillOpacity = color.Alpha16{uint16(0xffff * parseFloatUnit(val))}
	case "width":
		e.width = parseUnit(val)
	case "height":
		e.height = parseUnit(val)
	case "id":
		e.elementID = val
	case "x":
		e.xOffset = parseUnit(val)
	case "y":
		e.yOffset = parseUnit(val)
	}
}

func (e *commonElement) ColorModel() color.Model {
	return color.RGBAModel
}

func (e *commonElement) newRGBA() *image.RGBA {
	return image.NewRGBA(image.Rectangle{image.ZP, e.Bounds().Size()})
}

func (e *commonElement) parseCommonAttributes(start *xml.StartElement) {
	e.strokeColor = color.Transparent
	e.fillColor = color.Transparent

	for _, attr := range start.Attr {
		e.SetAttribute(attr.Name.Local, attr.Value)
	}
}

// utility functions

func parseFloatUnit(s string) float64 {
	var suffix string
	if !unicode.IsDigit(rune(s[len(s)-1])) {
		suffix, s = s[len(s)-2:], s[:len(s)-2]
	}

	scale := 1.0
	switch suffix {
	case "pt":
		scale = 1.25
	case "pc":
		scale = 15.0
	case "mm":
		scale = 3.543307
	case "cm":
		scale = 35.43307
	case "in":
		scale = 90.0
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}

	return scale * f
}

func parseUnit(s string) int {
	return int(parseFloatUnit(s))
}
