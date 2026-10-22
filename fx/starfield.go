package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Starfield sends points outward from the center behind content.
func Starfield(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 6*time.Second)
	const n = 60
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		far := float32(math.Hypot(float64(r.W), float64(r.H))) / 2
		for i := range n {
			a := 2 * math.Pi * float64(rnd(i, 1))
			d := fract(rnd(i, 2) + k)
			x := r.X + r.W/2 + d*far*float32(math.Cos(a))
			y := r.Y + r.H/2 + d*far*float32(math.Sin(a))
			var star ui.Path
			p.FillPath(star.Circle(x, y, t.Space(0.25)+t.Space(0.5)*d), t.Text.Alpha(d))
		}
	})).Children(content)
}
