package chart

import (
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// SparklineOptions tune a Sparkline or SparkBars.
type SparklineOptions struct {
	// Width and Height; zero gives the defaults.
	Width, Height float32
	// Format prints a hovered value; "%g" when empty.
	Format string
}

// Sparkline is a tiny line without axes.
func Sparkline(c *ui.Context, label string, vs []float64, o SparklineOptions) ui.Element {
	return spark(c, label, vs, o, "Sparkline", func(p *ui.Painter, r ui.Rect, ys []float32, g float32, hot int, ink ui.Color) {
		xs := make([]float32, len(ys))
		for i := range ys {
			xs[i] = float32(i) / float32(max(len(ys)-1, 1))
		}
		p.StrokePath(mc.Line(r, xs, ys, mc.CurveMonotone).Trim(0, g), 1.5, ink)
		if hot >= 0 {
			p.FillPath(new(ui.Path).Circle(mc.X(r, xs[hot]), mc.Y(r, ys[hot]), 2.5), ink)
		}
	})
}
