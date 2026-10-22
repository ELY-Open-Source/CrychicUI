package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// WaveText bobs the letters of text in a travelling wave.
func WaveText(c *ui.Context, text string) ui.Element {
	if text == "" {
		panic("motion: WaveText needs text")
	}
	gs, w, h := letters(c, text)
	amp, rm := c.Theme().Space(1), still(c)
	root := ui.Box(c).Size(w, h+2*amp).Role(ui.RoleText).Label(text)
	start := ui.Local(root, "start", func() time.Time { return c.Now() })
	return root.Draw(func(p *ui.Painter, r ui.Rect) {
		phase := p.Now().Sub(*start).Seconds() / 1.2
		for i, g := range gs {
			var dy float32
			if !rm {
				dy = -amp * float32(math.Sin(2*math.Pi*(phase-float64(i)*0.08)))
			}
			p.RichText(r.X+g.x, r.Y+amp+dy, 0, g.span)
		}
		if !rm {
			p.AnimationFrame()
		}
	})
}
