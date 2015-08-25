package main

import (
	"bufio"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/stephenwithav/ssvgc"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal(`Error!  ssvgc requires two arguments.

$ ssvgc <in.svg> <out.png>`)
	}

	r, f := LoadSVG(os.Args[1])
	defer f.Close()
	p := ssvgc.NewParser(r)
	svg, err := p.ParseSVG()
	if err != nil {
		log.Fatal("Error parsing SVG.", err)
	}

	SavePNG(os.Args[2], svg.Draw())
}

func LoadSVG(path string) (io.Reader, io.Closer) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error!  Unable to open %s.\n", path)
	}

	return bufio.NewReader(f), f
}

func SavePNG(path string, m image.Image) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error!  Unable to create %s.\n", path)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	err = png.Encode(w, m)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
