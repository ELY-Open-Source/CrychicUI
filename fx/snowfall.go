package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Snowfall drops swaying flakes behind content.
func Snowfall(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 12*time.Second)
	const n = 50
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for i := range n {
			fall := float32(1 + int(rnd(i, 1)*3))
			y := fract(rnd(i, 2) + fall*k)
			sway := t.Space(2) * float32(math.Sin(2*math.Pi*float64(2*fall*k+rnd(i, 3))))
			var flake ui.Path
			p.FillPath(flake.Circle(r.X+rnd(i, 4)*r.W+sway, r.Y+y*r.H, t.Space(0.5)+t.Space(0.75)*rnd(i, 5)),
				t.TextMuted.Alpha(0.5))
		}
	})).Children(content)
}
