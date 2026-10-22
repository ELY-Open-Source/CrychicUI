package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// RadialBarsOptions tune a RadialBars.
type RadialBarsOptions = PieOptions

// RadialBars draws concentric arcs, the first outermost.
func RadialBars(c *ui.Context, label string, names []string, vs []float64, at *int, o RadialBarsOptions) ui.Element {
	values("RadialBars", label, names, vs)
	top := 0.0
	for _, v := range vs {
		if v < 0 {
			panic(fmt.Sprintf("chart: RadialBars %q values must not be negative", label))
		}
		top = max(top, v)
	}
	chosen("RadialBars", label, at, len(names))
	def(&o.Format, "%g")
	t := c.Theme()
	size := o.Size
	if size == 0 {
		size = t.Space(56)
	}
	n := len(vs)
	outer := size/2 - t.Space(1)
	width := outer * 0.75 / float32(n)
	ring := func(i int) (float32, float32) {
		hi := outer - width*float32(i)
		return hi - width*0.8, hi
	}
	sweep := 1.5 * math.Pi
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		shown := springs(root, "v", vs)
		arc := func(i int) float64 {
			if top == 0 {
				return 0
			}
			return sweep * max(shown[i], 0) / top
		}
		wheel{label: label, comp: "RadialBars", n: n, size: size,
			hit: func(x, y, g float32) int {
				d := radius(x, y)
				for i := range vs {
					if lo, hi := ring(i); d >= lo && d <= hi && turn(x, y) <= arc(i)*float64(g) {
						return i
					}
				}
				return -1
			},
			describe: func(i int) string { return fmt.Sprintf("%s: "+o.Format, names[i], vs[i]) },
			draw: func(p *ui.Painter, r ui.Rect, hot int, g float32) {
				cx, cy := r.X+r.W/2, r.Y+r.H/2
				for i := range vs {
					lo, hi := ring(i)
					p.FillPath(mc.Sector(cx, cy, lo, hi, mc.Slice{End: sweep}), t.Border.Alpha(0.6))
					ink := tone(t, i)
					if *at >= 0 && *at != i {
						ink = ink.Alpha(0.4)
					}
					if i == hot {
						ink = ink.Mix(t.Text, 0.15)
					}
					p.FillPath(mc.Sector(cx, cy, lo, hi, mc.Slice{End: arc(i) * float64(g)}), ink)
				}
			},
		}.build(c, root, at, nil)
		legend(c, names)
	})
}
