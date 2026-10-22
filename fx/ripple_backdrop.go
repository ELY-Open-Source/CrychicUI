package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RippleBackdrop expands rings from the center behind content.
func RippleBackdrop(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 4*time.Second)
	const n = 4
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		far := float32(math.Hypot(float64(r.W), float64(r.H))) / 2
		for i := range n {
			f := fract(k + float32(i)/n)
			var ring ui.Path
			p.StrokePath(ring.Circle(r.X+r.W/2, r.Y+r.H/2, f*far), t.Space(0.375), t.Accent.Alpha(0.4*(1-f)))
		}
	})).Children(content)
}
