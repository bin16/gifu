package gifu

import (
	"image"
	"image/color"
	"sort"
)

type Quantizer struct {
	pal        color.Palette
	paletted   *image.Paletted
	ColorDepth int
}

func (q *Quantizer) Quantize(p color.Palette, m image.Image) color.Palette {
	if q.ColorDepth == 0 {
		q.ColorDepth = 8
	}

	q.pal = p
	q.paletted = image.NewPaletted(m.Bounds(), q.pal)
	q.splitIntoBuckets(q.unpackImage(m), q.ColorDepth)
	q.paletted.Palette = q.pal
	return q.pal
}

func (q *Quantizer) PalettedImage() *image.Paletted {
	return q.paletted
}

func (q *Quantizer) ConvertToPaletted(img image.Image) *image.Paletted {
	q.Quantize(color.Palette{}, img)
	return q.PalettedImage()
}

func (q *Quantizer) unpackImage(img image.Image) Pixels {
	pixels := Pixels{}
	for i := 0; i < img.Bounds().Dx(); i++ {
		for r := 0; r < img.Bounds().Dy(); r++ {
			pixels = append(pixels, NewPixel(img.At(i, r), i, r))
		}
	}

	return pixels
}

func (q *Quantizer) medCutQuantize(pixels Pixels) {
	avg := pixels.Average()
	q.pal = append(q.pal, avg)
	colorIndex := uint8(len(q.pal) - 1)
	for _, p := range pixels {
		q.paletted.SetColorIndex(p.X, p.Y, colorIndex)
	}
}

func (q *Quantizer) splitIntoBuckets(pixels Pixels, depth int) {
	if len(pixels) == 0 {
		return
	}

	if depth == 0 {
		q.medCutQuantize(pixels)
		return
	}

	channel := 0
	rRange := pixels.MaxOf(0) - pixels.MinOf(0)
	gRange := pixels.MaxOf(1) - pixels.MinOf(1)
	bRange := pixels.MaxOf(2) - pixels.MinOf(2)

	if rRange >= gRange && rRange >= bRange {
		channel = 0
	} else if gRange >= rRange && gRange >= bRange {
		channel = 1
	} else if bRange >= rRange && bRange >= gRange {
		channel = 2
	}

	sort.SliceStable(pixels, func(i, j int) bool {
		return pixels[i].Get(channel) < pixels[j].Get(channel)
	})
	medianIndex := (len(pixels) + 1) / 2
	q.splitIntoBuckets(pixels[:medianIndex], depth-1)
	q.splitIntoBuckets(pixels[medianIndex:], depth-1)
}
