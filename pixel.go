package gifu

import "image/color"

type Pixel struct {
	R, G, B, A uint8
	X, Y       int
}

func NewPixel(c color.Color, x, y int) Pixel {
	r, g, b, a := c.RGBA()
	return Pixel{
		uint8(r),
		uint8(g),
		uint8(b),
		uint8(a),
		x, y,
	}
}

func (p *Pixel) Get(channel int) uint8 {
	if channel == 0 {
		return p.R
	}

	if channel == 1 {
		return p.G
	}

	if channel == 2 {
		return p.B
	}

	if channel == 3 {
		return p.A
	}

	return 0
}

type Pixels []Pixel

func (pl Pixels) MaxOf(channel int) (max uint8) {
	max = 0
	for _, p := range pl {
		val := p.Get(channel)
		if val > max {
			max = val
		}
	}

	return
}

func (pl Pixels) MinOf(channel int) (min uint8) {
	min = 255
	for _, p := range pl {
		val := p.Get(channel)
		if val < min {
			min = val
		}
	}

	return
}

func (pl Pixels) AvgOf(channel int) uint8 {
	var total int64 = 0
	for _, p := range pl {
		total += int64(p.Get(channel))
	}

	return uint8(total / int64(len(pl)))
}

func (pl Pixels) Average() color.Color {
	r := pl.AvgOf(0)
	g := pl.AvgOf(1)
	b := pl.AvgOf(2)
	a := pl.AvgOf(3)
	return color.RGBA{
		r, g, b, a,
	}
}
