package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LightRays fans swaying rays from the top left behind content.
func LightRays(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := float64(phase(c, root, 8*time.Second))
	const n = 7
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		far := float32(math.Hypot(float64(r.W), float64(r.H)))
		for i := range n {
			mid := (10 + 70*float64(i)/(n-1) + 3*math.Sin(2*math.Pi*(k+float64(i)/n))) * math.Pi / 180
			const half = 3 * math.Pi / 180
			a := 0.08 + 0.06*float32(math.Sin(2*math.Pi*(2*k+float64(rnd(i, 1)))))
			var ray ui.Path
			ray.MoveTo(r.X, r.Y).
				LineTo(r.X+far*float32(math.Cos(mid-half)), r.Y+far*float32(math.Sin(mid-half))).
				LineTo(r.X+far*float32(math.Cos(mid+half)), r.Y+far*float32(math.Sin(mid+half)))
			p.FillPathGradient(ray.Close(), ui.LinearGradient{From: t.Accent.Alpha(a), To: t.Accent.Alpha(0), Angle: 135})
		}
	})).Children(content)
}
