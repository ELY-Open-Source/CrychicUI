package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Waves moves layered sine lines behind content.
func Waves(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := float64(phase(c, root, 8*time.Second))
	cols := palette(t)[:3]
	amp, step := t.Space(3), t.Space(2)
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		n := max(1, int(math.Ceil(float64(r.W/step))))
		for i, col := range cols {
			base := r.Y + r.H*(0.35+0.15*float32(i))
			var line ui.Path
			for j := range n + 1 {
				f := float32(j) / float32(n)
				w := math.Sin(2 * math.Pi * (float64(f)*(1+0.5*float64(i)) + k*float64(i+1)))
				if x, y := r.X+r.W*f, base+amp*float32(w); j == 0 {
					line.MoveTo(x, y)
				} else {
					line.LineTo(x, y)
				}
			}
			p.StrokePath(&line, t.Space(0.375), col.Alpha(0.6))
		}
	})).Children(content)
}
