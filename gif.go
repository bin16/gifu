package gifu

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"os"
)

func NewGIF(w, h int) *Gif {
	return &Gif{}
}

type Gif struct {
	gif.GIF
}

func (g *Gif) AppendImage(img image.Image, delay int) error {
	if len(g.Image) == 0 {
		g.Config.Width = img.Bounds().Dx()
		g.Config.Height = img.Bounds().Dy()
	}

	tmp := image.NewRGBA(image.Rect(0, 0, g.Config.Width, g.Config.Height))
	draw.Draw(tmp, tmp.Rect, img, image.Pt(0, 0), draw.Src)

	q := &Quantizer{}
	paletted := q.ConvertToPaletted(tmp)

	g.Image = append(g.Image, paletted)
	g.Delay = append(g.Delay, delay)

	return nil
}

func (g *Gif) Encode(w io.Writer) error {
	return gif.EncodeAll(w, &g.GIF)
}

func (g *Gif) SaveGIF(filename string) error {
	gifFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer gifFile.Close()

	return g.Encode(gifFile)
}

func (g *Gif) SaveSheet(filename string, columns int) error {
	w := g.Image[0].Rect.Dx()
	h := g.Image[0].Rect.Dy()
	rows := len(g.Image) / columns
	if len(g.Image)%columns > 0 {
		rows++
	}
	img := image.NewRGBA(image.Rect(0, 0, w*columns, h*rows))
	draw.Draw(img, img.Rect, image.NewUniform(color.White), image.Pt(0, 0), draw.Src)

	for i, frame := range g.Image {
		col := i % columns
		row := i / columns
		draw.Draw(img, image.Rect(col*w, row*h, col*w+w, row*h+h), frame, image.Pt(0, 0), draw.Src)
	}

	imgFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	return png.Encode(imgFile, img)
}
