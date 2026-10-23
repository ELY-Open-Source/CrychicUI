// Package dashboard holds analytics components.
package dashboard

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Point is one value in a series.
type Point struct {
	Label string
	Value float64
}

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func progress(c *ui.Context, value float64, color ui.Color) ui.Element {
	if value < 0 || value > 1 || !finite(value) {
		panic("dashboard: progress must be in [0,1]")
	}
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
	w := bar.AnimateTo("progress", float32(value), ui.Smooth)
	bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, color, 999)
	}))
	return bar
}

func series(c *ui.Context, points []Point, ink ui.Color) ui.Element {
	if len(points) < 2 {
		panic("dashboard: series needs two points")
	}
	lo, hi := points[0].Value, points[0].Value
	for _, point := range points {
		if !finite(point.Value) {
			panic("dashboard: series value must be finite")
		}
		lo, hi = min(lo, point.Value), max(hi, point.Value)
	}
	span := hi - lo
	if !finite(span) {
		panic("dashboard: series range overflow")
	}
	t := c.Theme()
	box := ui.Box(c).FillWidth().Height(t.Space(16))
	box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		var path ui.Path
		for i, point := range points {
			x := r.X + r.W*float32(i)/float32(len(points)-1)
			y := r.Y + r.H*0.5
			if span > 0 {
				y = r.Y + r.H*(1-float32((point.Value-lo)/span))
			}
			if i == 0 {
				path.MoveTo(x, y)
			} else {
				path.LineTo(x, y)
			}
		}
		p.StrokePath(&path, 2, ink)
	}))
	return box
}
