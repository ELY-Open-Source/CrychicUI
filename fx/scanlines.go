package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Scanlines rolls faint lines and a band down over content.
func Scanlines(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 4*time.Second)
	gap := t.Space(1)
	return root.Children(content).DrawOver(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for y := r.Y - gap + fract(8*k)*gap; y < r.Y+r.H; y += gap {
			p.Fill(ui.Rect{X: r.X, Y: y, W: r.W, H: 1}, t.Text.Alpha(0.08), 0)
		}
		band := r.H / 5
		g := ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{at(t.Text.Alpha(0), 0), at(t.Text.Alpha(0.06), 0.5),
			at(t.Text.Alpha(0), 1)}}
		p.FillGradient(ui.Rect{X: r.X, Y: r.Y - band + (r.H+band)*k, W: r.W, H: band}, g, 0)
	}))
}
