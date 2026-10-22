package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Bubbles rises circles behind content, fading near the top.
func Bubbles(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 10*time.Second)
	const n = 24
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for i := range n {
			rise := float32(1 + int(rnd(i, 1)*2))
			f := fract(rnd(i, 2) + rise*k)
			wobble := t.Space(1.5) * float32(math.Sin(2*math.Pi*float64(3*rise*k+rnd(i, 3))))
			var b ui.Path
			p.StrokePath(b.Circle(r.X+rnd(i, 4)*r.W+wobble, r.Y+r.H*(1-f), t.Space(1)+t.Space(2)*rnd(i, 5)), 1,
				t.Accent.Alpha(0.5*(1-f)))
		}
	})).Children(content)
}
