package ssvgc

import (
	"encoding/xml"
	"image"
	"image/draw"
)

type Rectangle struct {
	commonElement

	rX float64
	rY float64
}

func NewRectangle() *Rectangle {
	r := &Rectangle{}
	r.strokeColor = image.Transparent
	r.fillColor = image.Black

	return r
}

func (r *Rectangle) ParseAttributes(start *xml.StartElement) {
	r.parseCommonAttributes(start)

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "rx":
			r.rX = parseFloatUnit(attr.Value)
		case "ry":
			r.rY = parseFloatUnit(attr.Value)
		}
	}
}

func (r *Rectangle) Draw() image.Image {
	if r.upToDate && r.canvas != nil {
		return r.canvas
	}

	fillBounds := r.Bounds()

	switch r.strokeColor {
	case image.Transparent:
		r.canvas = image.NewRGBA(fillBounds)
	default:
		offsetBy := r.strokeWidth >> 1
		strokeBounds := fillBounds.Inset(-offsetBy)
		fillBounds = fillBounds.Inset(offsetBy + (r.strokeWidth & 1))
		r.canvas = image.NewRGBA(strokeBounds)
		draw.Draw(r.canvas, strokeBounds, &image.Uniform{r.strokeColor}, image.ZP, draw.Over)
	}

	draw.Draw(r.canvas, fillBounds, &image.Uniform{r.fillColor}, image.ZP, draw.Over)
	r.upToDate = true

	return r.canvas
}
