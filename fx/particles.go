package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Particles drifts dots behind content, wrapping at the edges.
func Particles(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 30*time.Second)
	const n = 40
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for i := range n {
			x := fract(rnd(i, 1) + laps(i, 3)*k)
			y := fract(rnd(i, 2) + laps(i, 4)*k)
			var dot ui.Path
			p.FillPath(dot.Circle(r.X+x*r.W, r.Y+y*r.H, t.Space(0.5)+t.Space(0.5)*rnd(i, 5)), t.Accent.Alpha(0.35))
		}
	})).Children(content)
}

// laps is whole crossings per loop: ±1 or ±2.
func laps(i int, salt uint32) float32 {
	v := float32(1 + int(rnd(i, salt)*2))
	if rnd(i, salt+10) < 0.5 {
		return -v
	}
	return v
}
