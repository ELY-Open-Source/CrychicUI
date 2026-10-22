package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RetroGrid scrolls a perspective grid toward the viewer behind content.
func RetroGrid(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 2*time.Second)
	gap := t.Space(12)
	const rows = 10
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		horizon, cx := r.Y+r.H*0.4, r.X+r.W/2
		depth := r.Y + r.H - horizon
		for i := range rows {
			z := fract((float32(i) + k) / rows)
			y := horizon + depth*z*z
			p.Line(r.X, y, r.X+r.W, y, 1, t.Accent.Alpha(0.5*z))
		}
		for x := float32(0); x <= r.W; x += gap {
			segment(p, cx, horizon, cx+x*2, r.Y+r.H, 1, t.Accent.Alpha(0.5))
			segment(p, cx, horizon, cx-x*2, r.Y+r.H, 1, t.Accent.Alpha(0.5))
		}
	})).Children(content)
}
