package ssvgc

import "encoding/xml"

type Image struct {
	commonElement

	href string
}

func (i *Image) SetAttribute(name string, value string) {
	switch name {
	case "href":
		i.href = value
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
