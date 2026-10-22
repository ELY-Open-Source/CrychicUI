package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Blob morphs a soft shape behind content.
func Blob(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := 2 * math.Pi * float64(phase(c, root, 10*time.Second))
	const n = 48
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		cx, cy, base := r.X+r.W/2, r.Y+r.H/2, min(r.W, r.H)*0.32
		var blob ui.Path
		for j := range n {
			a := 2 * math.Pi * float64(j) / n
			d := base * float32(1+0.12*math.Sin(3*a+k)+0.08*math.Sin(5*a-2*k))
			if x, y := cx+d*float32(math.Cos(a)), cy+d*float32(math.Sin(a)); j == 0 {
				blob.MoveTo(x, y)
			} else {
				blob.LineTo(x, y)
			}
		}
		p.FillPathGradient(blob.Close(), ui.LinearGradient{From: t.Accent.Alpha(0.35), To: t.Success.Alpha(0.25), Angle: 135})
	})).Children(content)
}
