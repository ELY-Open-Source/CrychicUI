package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WaterRipple expands a ring from each click over content.
func WaterRipple(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	const life = 800 * time.Millisecond
	list, now := taps(c, root, "WaterRipple", sound.Drop, life), c.Now()
	return root.Children(content).DrawOver(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for _, tp := range list {
			k := float32(now.Sub(tp.at)) / float32(life)
			x, y := r.X+tp.x, r.Y+tp.y
			var ring ui.Path
			p.StrokePath(ring.Circle(x, y, t.Space(16)*(1-(1-k)*(1-k))), t.Space(0.5), t.Accent.Alpha(0.6*(1-k)))
		}
	}))
}
