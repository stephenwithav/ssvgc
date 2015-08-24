package ssvgc

import (
	"image"
	"image/draw"
)

type SVG struct {
	commonElement

	elements []Element
}

func (s *SVG) Draw() image.Image {
	if s.upToDate && s.canvas != nil {
		return s.canvas
	}

	bounds := s.Bounds()
	s.canvas = image.NewRGBA(bounds)
	draw.Draw(s.canvas, bounds, &image.Uniform{s.fillColor}, image.ZP, draw.Over)

	for _, element := range s.elements {
		elementImage := element.Draw()
		draw.Draw(s.canvas, element.Bounds(), elementImage, elementImage.Bounds().Min, draw.Over)
	}

	return s.canvas
}

func (s *SVG) AddElement(e Element) {
	s.elements = append(s.elements, e)
}
