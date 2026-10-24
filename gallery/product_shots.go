package gallery

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/egoist/mygo/ui"
)

// Shapes of the generated product photos.
const (
	shotLamp = iota
	shotMug
	shotGlobe
	shotTray
	shotCushion
	shotVase
)

type shotSpec struct {
	shape        int
	back, object color.RGBA
}

// productShots draws studio-style product photos for the shop demos.
var productShots = sync.OnceValue(func() []*ui.Bitmap {
	specs := []shotSpec{
		{shotLamp, color.RGBA{236, 230, 220, 255}, color.RGBA{196, 150, 80, 255}},
		{shotMug, color.RGBA{222, 232, 236, 255}, color.RGBA{196, 104, 72, 255}},
		{shotGlobe, color.RGBA{52, 66, 90, 255}, color.RGBA{236, 234, 226, 255}},
		{shotTray, color.RGBA{232, 226, 214, 255}, color.RGBA{176, 128, 80, 255}},
		{shotCushion, color.RGBA{226, 232, 222, 255}, color.RGBA{150, 170, 140, 255}},
		{shotVase, color.RGBA{238, 226, 226, 255}, color.RGBA{90, 110, 150, 255}},
		{shotLamp, color.RGBA{228, 228, 226, 255}, color.RGBA{56, 56, 60, 255}},
		{shotLamp, color.RGBA{210, 220, 228, 255}, color.RGBA{240, 238, 230, 255}},
	}
	const side = 192
	out := make([]*ui.Bitmap, 0, len(specs))
	for _, s := range specs {
		img := image.NewRGBA(image.Rect(0, 0, side, side))
		for y := range side {
			for x := range side {
				var r, g, b float64
				// Four samples a pixel smooth the shape edges.
				for _, d := range [4][2]float64{{0.25, 0.25}, {0.75, 0.25}, {0.25, 0.75}, {0.75, 0.75}} {
					sr, sg, sb := shotPixel(s, (float64(x)+d[0])/side*2-1, (float64(y)+d[1])/side*2-1)
					r, g, b = r+sr/4, g+sg/4, b+sb/4
				}
				img.SetRGBA(x, y, color.RGBA{byteOf(r), byteOf(g), byteOf(b), 255})
			}
		}
		out = append(out, ui.NewBitmap(img))
	}
	return out
})

func byteOf(v float64) uint8 { return uint8(math.Max(0, math.Min(255, v))) }

func sq(v float64) float64 { return v * v }

// shotPixel is a backdrop, floor shadow and lit object.
func shotPixel(s shotSpec, x, y float64) (r, g, b float64) {
	k := 1.04 - 0.07*(y+1)
	if y > 0.62 {
		k -= 0.04
	}
	if d := sq(x/0.55) + sq((y-0.66)/0.07); d < 1 {
		k *= 1 - 0.28*(1-d)
	}
	r, g, b = float64(s.back.R)*k, float64(s.back.G)*k, float64(s.back.B)*k
	if lit, ok := shotObject(s.shape, x, y); ok {
		r, g, b = float64(s.object.R)*lit, float64(s.object.G)*lit, float64(s.object.B)*lit
	}
	return r, g, b
}

// shotObject reports whether (x, y) is lit object, and how.
func shotObject(shape int, x, y float64) (float64, bool) {
	side := func(w float64) float64 { return 0.78 + 0.32*(-x/w) }
	switch shape {
	case shotLamp:
		if y >= -0.62 && y <= -0.12 {
			if w := 0.2 + 0.24*(y+0.62)/0.5; math.Abs(x) < w {
				return side(w), true
			}
		}
		if math.Abs(x) < 0.03 && y > -0.12 && y < 0.6 {
			return 0.8, true
		}
		if sq(x/0.26)+sq((y-0.6)/0.06) < 1 {
			return side(0.26) * 0.9, true
		}
	case shotMug:
		if math.Abs(x) < 0.3 && y > -0.3 && y < 0.62 {
			return side(0.3), true
		}
		if d := math.Hypot(x-0.32, y-0.14); x > 0.3 && d > 0.1 && d < 0.18 {
			return 0.75, true
		}
	case shotGlobe:
		if math.Hypot(x, y-0.08) < 0.5 {
			return math.Max(0.55, 1.12-0.6*math.Hypot(x+0.18, y+0.12)), true
		}
		if math.Abs(x) < 0.012 && y < -0.41 {
			return 0.5, true
		}
	case shotTray:
		if math.Abs(x) < 0.62 && y > 0.4 && y < 0.6 {
			if y < 0.46 {
				return 1.1, true
			}
			return side(0.62) * 0.85, true
		}
	case shotCushion:
		if sq(x/0.62)+sq((y-0.3)/0.3) < 1 {
			return side(0.62), true
		}
	case shotVase:
		if y > -0.62 && y < 0.62 {
			if w := 0.1 + 0.22*math.Sin(math.Pi*(0.15+0.75*(y+0.62)/1.24)); math.Abs(x) < w {
				return side(w), true
			}
		}
	}
	return 0, false
}
