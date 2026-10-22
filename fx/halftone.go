package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Halftone sizes grid dots by a turning gradient behind content.
func Halftone(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	turn := 2 * math.Pi * float64(phase(c, root, 16*time.Second))
	dx, dy := float32(math.Cos(turn)), float32(math.Sin(turn))
	step := t.Space(3)
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for y := r.Y + step/2; y < r.Y+r.H; y += step {
			for x := r.X + step/2; x < r.X+r.W; x += step {
				f := 0.5 + 0.7*((x-r.X)/r.W-0.5)*dx + 0.7*((y-r.Y)/r.H-0.5)*dy
				if f = min(max(f, 0), 1); f > 0 {
					var dot ui.Path
					p.FillPath(dot.Circle(x, y, step*0.45*f), t.Accent.Alpha(0.3))
				}
			}
		}
	})).Children(content)
}
