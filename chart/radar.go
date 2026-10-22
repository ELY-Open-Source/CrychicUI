package chart

import (
	"fmt"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// RadarOptions tune a Radar.
type RadarOptions = PieOptions

// Radar plots series on spokes; at picks a spoke.
func Radar(c *ui.Context, label string, spokes []string, series []Series, at *int, o RadarOptions) ui.Element {
	vals, names := table("Radar", label, spokes, series)
	top := 0.0
	for _, vs := range vals {
		for _, v := range vs {
			if v < 0 {
				panic(fmt.Sprintf("chart: Radar %q values must not be negative", label))
			}
			top = max(top, v)
		}
	}
	if len(spokes) < 3 || top == 0 {
		panic(fmt.Sprintf("chart: Radar %q needs three spokes and a value above zero", label))
	}
	chosen("Radar", label, at, len(spokes))
	def(&o.Format, "%g")
	t := c.Theme()
	size := o.Size
	if size == 0 {
		size = t.Space(56)
	}
	n := len(spokes)
	step := 2 * math.Pi / float64(n)
	outer := size/2 - t.Space(8)
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		shown := make([][]float64, len(vals))
		for s, vs := range vals {
			shown[s] = springs(root, fmt.Sprint("s", s), vs)
		}
		wheel{label: label, comp: "Radar", n: n, size: size,
			hit: func(x, y, _ float32) int {
				if radius(x, y) > outer {
					return -1
				}
				return int(math.Round(turn(x, y)/step)) % n
			},
			describe: func(i int) string {
				var parts []string
				for s, nm := range names {
					parts = append(parts, nm+" "+fmt.Sprintf(o.Format, vals[s][i]))
				}
				return spokes[i] + ": " + strings.Join(parts, ", ")
			},
			draw: func(p *ui.Painter, r ui.Rect, hot int, g float32) {
				cx, cy := r.X+r.W/2, r.Y+r.H/2
				ring := func(k float32) *ui.Path {
					var path ui.Path
					for i := range n {
						x, y := mc.Polar(cx, cy, outer*k, float64(i)*step)
						if i == 0 {
							path.MoveTo(x, y)
						} else {
							path.LineTo(x, y)
						}
					}
					return path.Close()
				}
				for _, k := range []float32{0.25, 0.5, 0.75, 1} {
					p.StrokePath(ring(k), 1, t.Border)
				}
				for i := range n {
					x, y := mc.Polar(cx, cy, outer, float64(i)*step)
					var spoke ui.Path
					ink := t.Border
					if i == hot || i == *at {
						ink = t.TextMuted
					}
					p.StrokePath(spoke.MoveTo(cx, cy).LineTo(x, y), 1, ink)
					fs := t.Rem(0.7)
					w, h := p.MeasureText(0, ui.Span{Text: spokes[i], Size: fs})
					lx, ly := mc.Polar(cx, cy, outer+t.Space(1)+w/2, float64(i)*step)
					p.Text(lx-w/2, ly-h/2, spokes[i], fs, t.TextMuted)
				}
				for s, vs := range shown {
					var shape ui.Path
					for i, v := range vs {
						x, y := mc.Polar(cx, cy, outer*float32(max(v, 0)/top)*g, float64(i)*step)
						if i == 0 {
							shape.MoveTo(x, y)
						} else {
							shape.LineTo(x, y)
						}
					}
					shape.Close()
					p.FillPath(&shape, tone(t, s).Alpha(0.2))
					p.StrokePath(&shape, 2, tone(t, s))
				}
			},
		}.build(c, root, at, nil)
		legend(c, names)
	})
}
