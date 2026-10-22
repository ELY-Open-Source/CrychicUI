package fx

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Spotlight casts a soft light following the pointer.
func Spotlight(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	x, y, a := follow(root)
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		if a > 0 {
			glow(p, r.X+x, r.Y+y, t.Space(30), t.Accent, 0.4*a)
		}
	})).Children(content)
}
