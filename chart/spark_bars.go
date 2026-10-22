package chart

import (
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// SparkBarsOptions tune a SparkBars.
type SparkBarsOptions = SparklineOptions

// SparkBars are tiny bars without axes.
func SparkBars(c *ui.Context, label string, vs []float64, o SparkBarsOptions) ui.Element {
	return spark(c, label, vs, o, "SparkBars", func(p *ui.Painter, r ui.Rect, ys []float32, g float32, hot int, ink ui.Color) {
		w := r.W / float32(len(ys))
		for i, y := range ys {
			k := ink.Alpha(0.6)
			if i == hot {
				k = ink
			}
			top := mc.Y(r, y*g)
			p.Fill(ui.Rect{X: r.X + w*float32(i) + 1, Y: top, W: w - 2, H: r.Y + r.H - top}, k, 1)
		}
	})
}
