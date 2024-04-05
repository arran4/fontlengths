package main

import (
	"fmt"
	"github.com/arran4/golang-wordwrap"
	"github.com/flopp/go-findfont"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func GetTTFontFace(fontsize float64, dpi float64, gr *truetype.Font) (font.Face, error) {
	return truetype.NewFace(gr, &truetype.Options{
		Size: fontsize,
		DPI:  dpi,
	}), nil
}

func GetOTFontFace(fontsize float64, dpi float64, gr *opentype.Font) (font.Face, error) {
	return opentype.NewFace(gr, &opentype.FaceOptions{
		Size: fontsize,
		DPI:  dpi,
	})
}

func main() {
	_ = os.MkdirAll("out", 0755)
	log.Printf("Starting")
	for _, fn := range findfont.List() {
		b, err := os.ReadFile(fn)
		if err != nil {
			log.Printf("Failed to open %s because %s so skipping", fn, err)
			continue
		}
		if err := ProcessFont(b, fn); err != nil {
			log.Printf("Error: %s", err)
		}
	}
	log.Printf("Done")
}

func ProcessFont(b []byte, fn string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s panic'd: %s", fn, err)
		}
	}()
	ext := filepath.Ext(fn)
	switch strings.ToLower(ext) {
	case ".ttf":
		tt, err := truetype.Parse(b)
		if err != nil {
			return fmt.Errorf("reading font %s because %w", fn, err)
		}
		name := tt.Name(truetype.NameIDFontFullName)
		log.Printf("Generarting for %s", name)
		face, _ := GetTTFontFace(16, 300, tt)
		if err := CreateImage(face, name); err != nil {
			return fmt.Errorf("creating for %s because %w", name, err)
		}
	case ".otf":
		ot, err := opentype.Parse(b)
		if err != nil {
			return fmt.Errorf("reading font %s because %w", fn, err)
		}
		name, err := ot.Name(nil, sfnt.NameIDFull)
		if err != nil {
			return fmt.Errorf("reading font name of %s because %w", fn, err)
		}
		log.Printf("Generarting for %s", name)
		face, err := GetOTFontFace(16, 300, ot)
		if err != nil {
			return fmt.Errorf("creating font face of %s because %w", fn, err)
		}
		if err := CreateImage(face, name); err != nil {
			return fmt.Errorf("creating for %s because %w", name, err)
		}
	default:
		return fmt.Errorf("unknown font extension: %s on %s", ext, fn)
	}
	return nil
}

func CreateImage(face font.Face, name string) error {
	d := &font.Drawer{
		//Dst:  i,
		Src:  image.White,
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	type Line struct {
		wordwrap.Box
		Length fixed.Int26_6
		Height fixed.Int26_6
	}
	lines := make([]*Line, 26)
	var height fixed.Int26_6
	var length fixed.Int26_6
	for i := int('A'); i <= int('Z'); i++ {
		stb, err := wordwrap.NewSimpleTextBox(d, strings.Repeat(fmt.Sprintf("%c", i), 20))
		if err != nil {
			return fmt.Errorf("box creation failed; %w", err)
		}
		line := &Line{
			Box:    stb,
			Length: stb.AdvanceRect(),
			Height: stb.MetricsRect().Height,
		}
		lines[i-int('A')] = line
		height += line.Height
		if line.Length.Ceil() > length.Ceil() {
			length = line.Length
		}
	}
	slices.SortFunc(lines, func(a, b *Line) int {
		return a.Length.Ceil() - b.Length.Ceil()
	})
	i := image.NewRGBA(image.Rect(0, 0, length.Ceil(), height.Ceil()))
	draw.Draw(i, i.Rect, image.Black, image.Pt(0, 0), draw.Over)
	var ypos fixed.Int26_6
	for _, line := range lines {
		ypos += line.Height
		line.DrawBox(i, ypos, &wordwrap.DrawConfig{})
	}
	f, err := os.Create(fmt.Sprintf("out/%s.png", name))
	if err != nil {
		return fmt.Errorf("file create failed; %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	if err := png.Encode(f, i); err != nil {
		return fmt.Errorf("PNG encode failed; %w", err)
	}
	return nil
}
