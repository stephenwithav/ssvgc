package ssvgc

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"image"
	"io"
)

type Entity uint8

type Parser struct {
	d *xml.Decoder
}

func NewParser(r io.Reader) *Parser {
	p := &Parser{d: xml.NewDecoder(r)}

	return p
}

func (p *Parser) ParseSVG() (*SVG, error) {
	s := &SVG{}
	s.fillColor = image.Transparent

	tok, err := p.d.Token()
	for ; err == nil; tok, err = p.d.Token() {
		switch el := tok.(type) {
		case xml.StartElement:
			p.dispatchStartElement(s, el.Name.Local, &el)
		}
	}

	if err != io.EOF {
		return nil, err
	}

	return s, nil
}

func (p *Parser) dispatchStartElement(s *SVG, tagname string, tok *xml.StartElement) {
	switch svgType[tagname] {
	case SVG_BASE_ELEMENT:
		s.ParseAttributes(tok)
	case SVG_CIRCLE_ELEMENT:
	case SVG_ELLIPSE_ELEMENT:
	case SVG_GROUP_ELEMENT:
	case SVG_IMAGE_ELEMENT:
	case SVG_LINE_ELEMENT:
	case SVG_PATH_ELEMENT:
	case SVG_POLYGON_ELEMENT:
	case SVG_POLYLINE_ELEMENT:
	case SVG_RECT_ELEMENT:
		r := NewRectangle()
		r.ParseAttributes(tok)
		s.AddElement(r)
	case SVG_TEXT_ELEMENT:
		s.AddElement(p.consumeText(tok))
	case SVG_TSPAN_ELEMENT:
	case SVG_UNKNOWN_ELEMENT:
	}
}

func (p *Parser) consumeText(tok *xml.StartElement) *Text {
	t := NewText()
	t.ParseAttributes(tok)
	buf := new(bytes.Buffer)
	cs := &contextStack{}
	cs.Push(t.textContext)

	hasPrev := false
	for tok, err := p.d.Token(); err == nil; tok, err = p.d.Token() {
		switch el := tok.(type) {
		case xml.CharData:
			startPos := buf.Len()

			s := bufio.NewScanner(bytes.NewReader(el))
			s.Split(bufio.ScanWords)
			for s.Scan() {
				switch hasPrev {
				case true:
					buf.WriteString(" " + s.Text())
				default:
					buf.WriteString(s.Text())

				}
				hasPrev = true
			}

			endPos := buf.Len()
			t.chunks = append(t.chunks, &textChunk{cs.Top(), startPos, endPos})
		case xml.StartElement:
			t.ParseAttributes(&el)
			cs.Push(t.textContext)
		case xml.EndElement:
			cs.Pop()
			if el.Name.Local == "text" {
				t.buf = buf.Bytes()
				return t
			}
		}
	}

	return t
}

var svgType map[string]Entity

func init() {
	svgType = map[string]Entity{
		"svg":      SVG_BASE_ELEMENT,
		"circle":   SVG_CIRCLE_ELEMENT,
		"ellipse":  SVG_ELLIPSE_ELEMENT,
		"group":    SVG_GROUP_ELEMENT,
		"image":    SVG_IMAGE_ELEMENT,
		"line":     SVG_LINE_ELEMENT,
		"path":     SVG_PATH_ELEMENT,
		"polygon":  SVG_POLYGON_ELEMENT,
		"polyline": SVG_POLYLINE_ELEMENT,
		"rect":     SVG_RECT_ELEMENT,
		"text":     SVG_TEXT_ELEMENT,
		"tspan":    SVG_TSPAN_ELEMENT,
	}
}

const (
	SVG_BASE_ELEMENT Entity = iota
	SVG_CIRCLE_ELEMENT
	SVG_ELLIPSE_ELEMENT
	SVG_GROUP_ELEMENT
	SVG_IMAGE_ELEMENT
	SVG_LINE_ELEMENT
	SVG_PATH_ELEMENT
	SVG_POLYGON_ELEMENT
	SVG_POLYLINE_ELEMENT
	SVG_RECT_ELEMENT
	SVG_TEXT_ELEMENT
	SVG_TSPAN_ELEMENT
	SVG_UNKNOWN_ELEMENT
)

type contextStack struct {
	contexts []textContext
}

func (c *contextStack) Push(tc textContext) {
	c.contexts = append(c.contexts, tc)
}

func (c *contextStack) Pop() textContext {
	i := len(c.contexts) - 1
	ctx := c.contexts[i]
	c.contexts = c.contexts[0:i]
	return ctx
}

func (c contextStack) Top() textContext {
	return c.contexts[len(c.contexts)-1]
}
