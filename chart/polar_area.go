package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// PolarAreaOptions tune a PolarArea.
type PolarAreaOptions = PieOptions

// PolarArea sizes equal wedges by value.
func PolarArea(c *ui.Context, label string, names []string, vs []float64, at *int, o PolarAreaOptions) ui.Element {
	values("PolarArea", label, names, vs)
	top := 0.0
	for _, v := range vs {
		if v < 0 {
			panic(fmt.Sprintf("chart: PolarArea %q values must not be negative", label))
		}
		top = max(top, v)
	}
	chosen("PolarArea", label, at, len(names))
	def(&o.Format, "%g")
	t := c.Theme()
	size := o.Size
	if size == 0 {
		size = t.Space(56)
	}
	n := len(vs)
	step := 2 * math.Pi / float64(n)
	outer := size/2 - t.Space(1)
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		shown := springs(root, "v", vs)
		reach := func(i int) float32 {
			if top == 0 {
				return 0
			}
			return outer * float32(math.Sqrt(max(shown[i], 0)/top))
		}
		wheel{label: label, comp: "PolarArea", n: n, size: size,
			hit: func(x, y, g float32) int {
				i := int(turn(x, y) / step)
				if i < n && radius(x, y) <= reach(i)*g {
					return i
				}
				return -1
			},
			describe: func(i int) string { return fmt.Sprintf("%s: "+o.Format, names[i], vs[i]) },
			draw: func(p *ui.Painter, r ui.Rect, hot int, g float32) {
				cx, cy := r.X+r.W/2, r.Y+r.H/2
				for i := range vs {
					ink := tone(t, i).Alpha(0.85)
					if *at >= 0 && *at != i {
						ink = ink.Alpha(0.4)
					}
					if i == hot {
						ink = ink.Mix(t.Text, 0.15)
					}
					s := mc.Slice{Start: float64(i) * step, End: float64(i+1) * step, Pad: 0.02}
					p.FillPath(mc.Sector(cx, cy, 0, reach(i)*g, s), ink)
				}
			},
		}.build(c, root, at, nil)
		legend(c, names)
	})
}
