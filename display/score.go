package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ScoreOptions tune a Score.
type ScoreOptions struct {
	// Max is the full ring; zero is 100.
	Max float64
}

// Score draws a number ring to its value.
func Score(c *ui.Context, value float64, o ScoreOptions) ui.Element {
	if o.Max == 0 {
		o.Max = 100
	}
	if !finite(value) || !finite(o.Max) || o.Max <= 0 || value < 0 || value > o.Max {
		panic("display: Score value must fit Max")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Size(t.Space(16), t.Space(16)).Center()
	shown := root.AnimateTo("score", float32(value/o.Max), ui.Snappy)
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		radius := (r.W - t.Space(1)) / 2
		ring := (&ui.Path{}).Circle(r.X+r.W/2, r.Y+r.H/2, radius)
		p.StrokePath(ring, t.Space(0.5), t.Border)
		if shown > 0 {
			p.StrokePath(ring.Trim(0, shown), t.Space(0.5), t.Accent)
		}
	})
	return root.Children(func() { kit.Numeric(c, formatNumber(c, value)).Bold() })
}
