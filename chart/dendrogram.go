package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Join merges A and B into cluster n+k.
type Join struct {
	A, B   int
	Height float64
}

// DendrogramOptions tune a Dendrogram.
type DendrogramOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
}

// Dendrogram draws cluster joins; at picks a leaf.
func Dendrogram(c *ui.Context, label string, leaves []string, joins []Join, at *int, o DendrogramOptions) ui.Element {
	named("Dendrogram", label, leaves)
	n := len(leaves)
	if n < 2 || len(joins) != n-1 {
		panic(fmt.Sprintf("chart: Dendrogram %q needs two leaves and one join fewer", label))
	}
	heights := make([]float64, 2*n-1)
	used := make([]bool, 2*n-1)
	for k, j := range joins {
		finite("Dendrogram", label, j.Height)
		for _, id := range []int{j.A, j.B} {
			if id < 0 || id >= n+k || used[id] || j.Height < heights[id] {
				panic(fmt.Sprintf("chart: Dendrogram %q joins take each earlier cluster once, no lower", label))
			}
			used[id] = true
		}
		heights[n+k] = j.Height
	}
	chosen("Dendrogram", label, at, n)
	t := c.Theme()
	xs := make([]float32, 2*n-1)
	slot := 0
	var walk func(id int)
	walk = func(id int) {
		if id < n {
			xs[id] = (float32(slot) + 0.5) / float32(n)
			slot++
			return
		}
		j := joins[id-n]
		walk(j.A)
		walk(j.B)
		xs[id] = (xs[j.A] + xs[j.B]) / 2
	}
	walk(2*n - 2)
	top := heights[2*n-2]
	if top == 0 {
		top = 1
	}
	y := mc.NewLinear(0, top).Nice(4)
	ticks := y.Ticks(4)
	bottom := make([]mc.Tick, n)
	for i, name := range leaves {
		bottom[i] = mc.Tick{At: xs[i], Label: name}
	}
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: bottom}, o.Height, func(area ui.Element) {
			g := grow(area)
			hs := springs(area, "h", heights)
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for k, j := range joins {
					ya, yb, yj := mc.Y(r, y.At(hs[j.A])), mc.Y(r, y.At(hs[j.B])), mc.Y(r, y.At(hs[n+k]))
					xa, xb := mc.X(r, xs[j.A]), mc.X(r, xs[j.B])
					var elbow ui.Path
					elbow.MoveTo(xa, ya).LineTo(xa, yj).LineTo(xb, yj).LineTo(xb, yb)
					p.StrokePath(elbow.Trim(0, g), 1.5, tone(t, 0))
				}
			}))
			for i, name := range leaves {
				c.Key(i)
				dot(c, root, at, i, xs[i], 0, t.Space(3)*g, tone(t, 0), name, "Dendrogram")
			}
		})
	})
}
