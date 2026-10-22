package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// OrbitRings turns dots on concentric rings behind content.
func OrbitRings(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := float64(phase(c, root, 12*time.Second))
	cols := palette(t)[:3]
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		for i, col := range cols {
			d := min(r.W, r.H) / 2 * float32(i+1) / (float32(len(cols)) + 0.5)
			var ring ui.Path
			p.StrokePath(ring.Circle(cx, cy, d), 1, t.Border)
			turns := float64(len(cols) - i)
			if i%2 == 1 {
				turns = -turns
			}
			a := 2 * math.Pi * (float64(rnd(i, 1)) + k*turns)
			var dot ui.Path
			p.FillPath(dot.Circle(cx+d*float32(math.Cos(a)), cy+d*float32(math.Sin(a)), t.Space(1)), col)
		}
	})).Children(content)
}
