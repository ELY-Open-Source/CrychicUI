package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// dips are the loop's dim spans: from, to, brightness.
var dips = [...][3]float32{{0.1, 0.12, 0.2}, {0.14, 0.16, 0.5}, {0.6, 0.62, 0.3}}

// Flicker shows text dimming in bursts like a failing sign.
func Flicker(c *ui.Context, text string) ui.Element {
	s := headline(c, text)
	e := filled(c, s)
	k, lit := phase(c, e, 4*time.Second), float32(1)
	for _, d := range dips {
		if k >= d[0] && k < d[1] {
			lit = d[2]
		}
	}
	s.Color = c.Theme().Accent.Alpha(lit)
	return e.Draw(func(p *ui.Painter, r ui.Rect) { p.RichText(r.X, r.Y, 0, s) })
}
