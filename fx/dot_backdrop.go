package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DotBackdrop draws a dot grid drifting behind content.
func DotBackdrop(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	cell, rad := t.Space(5), t.Space(0.375)
	off := cell * phase(c, root, 10*time.Second)
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for y := r.Y - cell + off; y < r.Y+r.H+cell; y += cell {
			for x := r.X - cell + off; x < r.X+r.W+cell; x += cell {
				var dot ui.Path
				p.FillPath(dot.Circle(x, y, rad), t.TextMuted.Alpha(0.35))
			}
		}
	})).Children(content)
}
