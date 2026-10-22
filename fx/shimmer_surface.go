package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ShimmerSurface passes a slow sheen over a surface behind content.
func ShimmerSurface(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 4*time.Second)
	return root.Background(t.Surface).Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		w := r.W / 2
		g := ui.LinearGradient{Angle: 110, Stops: []ui.GradientStop{at(t.AccentText.Alpha(0), 0),
			at(t.AccentText.Alpha(0.35), 0.5), at(t.AccentText.Alpha(0), 1)}}
		p.FillGradient(ui.Rect{X: r.X - w + (r.W+2*w)*k, Y: r.Y, W: w, H: r.H}, g, 0)
	})).Children(content)
}
