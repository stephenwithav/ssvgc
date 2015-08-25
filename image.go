package ssvgc

import (
	"encoding/xml"
	"image"
	"strings"

	"github.com/disintegration/imaging"
)

type Image struct {
	commonElement

	href   string
	filter imaging.ResampleFilter
}

var resampleFilters map[string]imaging.ResampleFilter = map[string]imaging.ResampleFilter{
	"nearestneighbor":   imaging.NearestNeighbor,
	"box":               imaging.Box,
	"linear":            imaging.Linear,
	"hermite":           imaging.Hermite,
	"mitchellnetravali": imaging.MitchellNetravali,
	"catmullrom":        imaging.CatmullRom,
	"bspline":           imaging.BSpline,
	"gaussian":          imaging.Gaussian,
	"bartlett":          imaging.Bartlett,
	"lanczos":           imaging.Lanczos,
	"hann":              imaging.Hann,
	"hamming":           imaging.Hamming,
	"blackman":          imaging.Blackman,
	"welch":             imaging.Welch,
	"cosine":            imaging.Cosine,
}

func NewImage() *Image {
	i := &Image{}
	i.filter = imaging.NearestNeighbor

	return i
}

func (i *Image) SetAttribute(name string, value string) {
	switch name {
	case "href":
		i.href = value
	case "resample-filter":
		if filter, ok := resampleFilters[strings.ToLower(value)]; ok {
			i.filter = filter
		}
	default:
		i.commonElement.SetAttribute(name, value)
	}
	i.upToDate = false
}

func (i *Image) ParseAttributes(start *xml.StartElement) {
	for _, attr := range start.Attr {
		i.SetAttribute(attr.Name.Local, attr.Value)
	}
}

func (i *Image) Draw() image.Image {
	if i.canvas != nil && i.upToDate {
		return i.canvas
	}

	m, err := imaging.Open(i.href)
	if err != nil {
		return image.Transparent
	}

	if i.width|i.height != 0 {
		m = imaging.Resize(m, i.width, i.height, i.filter)
	}
	size := m.Bounds().Size()
	i.setDimensions(size.X, size.Y)

	return m
}
