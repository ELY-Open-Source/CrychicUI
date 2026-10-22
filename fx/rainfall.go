package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Rainfall drops slanted streaks behind content.
func Rainfall(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 3*time.Second)
	long := t.Space(4)
	const n, slant = 60, 0.25
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for i := range n {
			fall := float32(1 + int(rnd(i, 1)*2))
			y := r.Y - long + fract(rnd(i, 2)+fall*k)*(r.H+long)
			x := r.X + rnd(i, 3)*(r.W+r.H*slant) - (y-r.Y)*slant
			segment(p, x, y, x-long*slant, y+long, 1, t.Accent.Alpha(0.35))
		}
	})).Children(content)
}
