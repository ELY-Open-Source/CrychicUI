package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Meteors crosses streaks diagonally behind content.
func Meteors(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 4*time.Second)
	tail := t.Space(16)
	const n, parts = 5, 6
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		run := r.W + r.H + tail
		for i := range n {
			d := fract(rnd(i, 1)+k) * run
			x, y := r.X+rnd(i, 2)*r.W+r.H-d*0.7, r.Y-tail+d*0.7
			for j := range parts {
				a, b := tail*float32(j)/parts, tail*float32(j+1)/parts
				segment(p, x+a*0.7, y-a*0.7, x+b*0.7, y-b*0.7, t.Space(0.375), t.Text.Alpha(0.6*(1-float32(j)/parts)))
			}
		}
	})).Children(content)
}
