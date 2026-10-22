package fx

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Sparkles blinks small stars around content.
func Sparkles(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("fx: Sparkles needs content")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Padding(t.Space(4))
	k := float64(phase(c, root, 2400*time.Millisecond))
	big := t.Space(1.5)
	const n = 8
	return root.Children(content).DrawOver(func(p *ui.Painter, r ui.Rect) {
		for i := range n {
			s := float32(math.Max(0, math.Sin(2*math.Pi*(float64(1+i%3)*k+float64(rnd(i, 1))))))
			if s == 0 {
				continue
			}
			x, y := r.X+big+rnd(i, 2)*(r.W-2*big), r.Y+big+rnd(i, 3)*(r.H-2*big)
			o, in := big*s, big*s*0.3
			var star ui.Path
			star.MoveTo(x, y-o).LineTo(x+in, y-in).LineTo(x+o, y).LineTo(x+in, y+in).LineTo(x, y+o).
				LineTo(x-in, y+in).LineTo(x-o, y).LineTo(x-in, y-in).Close()
			p.FillPath(&star, t.Warning)
		}
	})
}
