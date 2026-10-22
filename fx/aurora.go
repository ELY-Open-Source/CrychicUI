package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Aurora moves soft bands of theme colors behind content.
func Aurora(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := float64(phase(c, root, 12*time.Second))
	cols := palette(t)[:3]
	amp, step := t.Space(5), t.Space(2)
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for i, col := range cols {
			base := r.Y + r.H*(0.25+0.15*float32(i))
			var band ui.Path
			band.MoveTo(r.X, r.Y+r.H)
			n := max(1, int(math.Ceil(float64(r.W/step))))
			for j := range n + 1 {
				f := float32(j) / float32(n)
				w := math.Sin(2 * math.Pi * (float64(f)*1.2 + k + float64(i)/3))
				band.LineTo(r.X+r.W*f, base+amp*float32(w))
			}
			band.LineTo(r.X+r.W, r.Y+r.H).Close()
			p.FillPathGradient(&band, ui.LinearGradient{From: col.Alpha(0.3), To: col.Alpha(0), Angle: 180})
		}
	})).Children(content)
}
