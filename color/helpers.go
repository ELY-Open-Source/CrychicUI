// Package color holds components that choose and inspect colors.
package color

import (
	"fmt"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"
)

// lin turns an sRGB byte into linear light.
func lin(v uint8) float64 {
	c := float64(v) / 255
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

// enc turns linear light into an sRGB byte.
func enc(v float64) uint8 {
	v = max(0, min(v, 1))
	if v <= 0.0031308 {
		v *= 12.92
	} else {
		v = 1.055*math.Pow(v, 1/2.4) - 0.055
	}
	return uint8(v*255 + 0.5)
}

// hsv splits c into hue degrees, saturation and value.
func hsv(c ui.Color) (h, s, v float64) {
	r, g, b := float64(c.R)/255, float64(c.G)/255, float64(c.B)/255
	hi, lo := max(r, g, b), min(r, g, b)
	if hi > 0 {
		s = (hi - lo) / hi
	}
	return hue(r, g, b, hi, lo), s, hi
}

// hue is the hue of r, g, b in degrees.
func hue(r, g, b, hi, lo float64) float64 {
	d := hi - lo
	switch {
	case d == 0:
		return 0
	case hi == r:
		return math.Mod(60*(g-b)/d+360, 360)
	case hi == g:
		return 60*(b-r)/d + 120
	}
	return 60*(r-g)/d + 240
}

// fromHSV builds a color from hue, saturation, value and alpha.
func fromHSV(h, s, v float64, a uint8) ui.Color {
	f := func(n float64) uint8 {
		k := math.Mod(n+h/60, 6)
		return uint8((v-v*s*max(0, min(k, 4-k, 1)))*255 + 0.5)
	}
	return ui.Color{R: f(5), G: f(3), B: f(1), A: a}
}

// hsl splits c into hue degrees, saturation and lightness.
func hsl(c ui.Color) (h, s, l float64) {
	r, g, b := float64(c.R)/255, float64(c.G)/255, float64(c.B)/255
	hi, lo := max(r, g, b), min(r, g, b)
	l = (hi + lo) / 2
	if d := hi - lo; d > 0 {
		s = d / (1 - math.Abs(2*l-1))
	}
	return hue(r, g, b, hi, lo), s, l
}

// fromHSL builds a color from hue, saturation, lightness and alpha.
func fromHSL(h, s, l float64, a uint8) ui.Color {
	k := s * min(l, 1-l)
	f := func(n float64) uint8 {
		m := math.Mod(n+h/30, 12)
		return uint8((l-k*max(-1, min(m-3, 9-m, 1)))*255 + 0.5)
	}
	return ui.Color{R: f(0), G: f(8), B: f(4), A: a}
}

// oklab turns c into Oklab lightness, a and b.
func oklab(c ui.Color) (l, a, b float64) {
	r, g, bl := lin(c.R), lin(c.G), lin(c.B)
	x := math.Cbrt(0.4122214708*r + 0.5363325363*g + 0.0514459929*bl)
	y := math.Cbrt(0.2119034982*r + 0.6806995451*g + 0.1073969566*bl)
	z := math.Cbrt(0.0883024619*r + 0.2817188376*g + 0.6299787005*bl)
	return 0.2104542553*x + 0.7936177850*y - 0.0040720468*z,
		1.9779984951*x - 2.4285922050*y + 0.4505937099*z,
		0.0259040371*x + 0.7827717662*y - 0.8086757660*z
}

// oklch turns c into Oklch lightness, chroma and hue degrees.
func oklch(c ui.Color) (l, ch, h float64) {
	l, a, b := oklab(c)
	return l, math.Hypot(a, b), math.Mod(math.Atan2(b, a)*180/math.Pi+360, 360)
}

// inSRGB reports whether oklch(l c h) lies inside sRGB.
func inSRGB(l, ch, h float64) bool {
	a, b := ch*math.Cos(h*math.Pi/180), ch*math.Sin(h*math.Pi/180)
	x := math.Pow(l+0.3963377774*a+0.2158037573*b, 3)
	y := math.Pow(l-0.1055613458*a-0.0638541728*b, 3)
	z := math.Pow(l-0.0894841775*a-1.2914855480*b, 3)
	const eps = 1e-4
	for _, v := range []float64{4.0767416621*x - 3.3077115913*y + 0.2309699292*z,
		-1.2684380046*x + 2.6097574011*y - 0.3413193965*z, -0.0041960863*x - 0.7034186147*y + 1.7076147010*z} {
		if v < -eps || v > 1+eps {
			return false
		}
	}
	return true
}

// luminance is c's WCAG relative luminance.
func luminance(c ui.Color) float64 { return 0.2126*lin(c.R) + 0.7152*lin(c.G) + 0.0722*lin(c.B) }

// ratio is the WCAG contrast ratio of a and b.
func ratio(a, b ui.Color) float64 {
	x, y := luminance(a), luminance(b)
	return (max(x, y) + 0.05) / (min(x, y) + 0.05)
}

// hexOf writes c as #RRGGBB, or #RRGGBBAA when translucent.
func hexOf(c ui.Color) string {
	if c.A == 255 {
		return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}

// parseHex reads #rgb, #rgba, #rrggbb or #rrggbbaa.
func parseHex(s string) (ui.Color, bool) {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if n := len(h); n != 3 && n != 4 && n != 6 && n != 8 || strings.Trim(h, "0123456789abcdefABCDEF") != "" {
		return ui.Color{}, false
	}
	return ui.Hex("#" + h), true
}

// sim are Machado's full-severity matrices: protan, deutan, tritan.
var sim = [3][9]float64{
	{0.152286, 1.052583, -0.204868, 0.114503, 0.786281, 0.099216, -0.003882, -0.048116, 1.051998},
	{0.367322, 0.860646, -0.227968, 0.280085, 0.672501, 0.047413, -0.011820, 0.042940, 0.968881},
	{1.255528, -0.076749, -0.178779, -0.078411, 0.930809, 0.147602, 0.004733, 0.691367, 0.303900},
}

// seen is c as eyes of kind k see it.
func seen(c ui.Color, k int) ui.Color {
	m := sim[k]
	r, g, b := lin(c.R), lin(c.G), lin(c.B)
	return ui.Color{R: enc(m[0]*r + m[1]*g + m[2]*b), G: enc(m[3]*r + m[4]*g + m[5]*b), B: enc(m[6]*r + m[7]*g + m[8]*b), A: c.A}
}

// checker paints a two-tone check under translucent colors.
func checker(p *ui.Painter, r ui.Rect, t *ui.Theme) {
	s := t.Space(1.5)
	p.Fill(r, t.Background, 0)
	for y, row := r.Y, 0; y < r.Y+r.H; y, row = y+s, row+1 {
		for x, col := r.X, 0; x < r.X+r.W; x, col = x+s, col+1 {
			if (row+col)%2 == 1 {
				p.Fill(ui.Rect{X: x, Y: y, W: min(s, r.X+r.W-x), H: min(s, r.Y+r.H-y)}, t.Border, 0)
			}
		}
	}
}

// swatch is a chip of col over a check.
func swatch(c *ui.Context, col ui.Color, w, h float32) ui.Element {
	t := c.Theme()
	b := ui.Box(c).Size(w, h).Shrink(0).Radius(t.Radius).Border(1, t.Border).Role(ui.RoleImage).Label(hexOf(col))
	shown := b.AnimateColor("col", col, ui.Snappy)
	return b.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Clip(r, t.Radius, func() {
			checker(p, r, t)
			p.Fill(r, shown, 0)
		})
	})
}

// def sets an empty copy field to its English default.
func def(field *string, english string) {
	if *field == "" {
		*field = english
	}
}

func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}
