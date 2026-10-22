package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DotGlobe turns a sphere of dots behind content.
func DotGlobe(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	turn := 2 * math.Pi * float64(phase(c, root, 20*time.Second))
	const n = 160
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, min(r.W, r.H)*0.4
		for i := range n {
			// A Fibonacci lattice spreads points evenly.
			y := 1 - 2*(float64(i)+0.5)/n
			ring, a := math.Sqrt(1-y*y), 2.39996*float64(i)+turn
			x, z := ring*math.Cos(a), ring*math.Sin(a)
			near := float32(z+1) / 2
			var dot ui.Path
			p.FillPath(dot.Circle(cx+rad*float32(x), cy+rad*float32(y), t.Space(0.25)+t.Space(0.25)*near),
				t.Accent.Alpha(0.15+0.6*near))
		}
	})).Children(content)
}
