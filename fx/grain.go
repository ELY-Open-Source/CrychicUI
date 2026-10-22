package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Grain flickers faint dither dots over content.
func Grain(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	seed := uint32(48 * phase(c, root, 4*time.Second))
	cell, dot := t.Space(3), t.Space(0.5)
	return root.Children(content).DrawOver(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		n := int(r.W * r.H / (cell * cell))
		for i := range n {
			x, y := r.X+rnd(i, seed)*r.W, r.Y+rnd(i, seed+100)*r.H
			p.Fill(ui.Rect{X: x, Y: y, W: dot, H: dot}, t.Text.Alpha(0.15), 0)
		}
	}))
}
