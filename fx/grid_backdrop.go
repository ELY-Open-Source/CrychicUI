package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GridBackdrop draws faint grid lines drifting behind content.
func GridBackdrop(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	cell := t.Space(6)
	off := cell * phase(c, root, 8*time.Second)
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for x := r.X - cell + off; x < r.X+r.W; x += cell {
			p.Line(x, r.Y, x, r.Y+r.H, 1, t.Border)
		}
		for y := r.Y - cell + off; y < r.Y+r.H; y += cell {
			p.Line(r.X, y, r.X+r.W, y, 1, t.Border)
		}
	})).Children(content)
}
