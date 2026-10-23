package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ContextMeterOptions tune a ContextMeter.
type ContextMeterOptions struct {
	Label  string
	Format string
}

// ContextMeter is a ring of context used; warns near full.
func ContextMeter(c *ui.Context, used float64, o ContextMeterOptions) ui.Element {
	if !(used >= 0 && used <= 1) {
		panic(fmt.Sprintf("agent: ContextMeter needs a value from 0 to 1, got %v", used))
	}
	def(&o.Label, "Context used")
	def(&o.Format, "%.0f%%")
	t := c.Theme()
	side, w := t.Space(16), t.Space(1.5)
	root := ui.Box(c).Size(side, side).Center().Role(ui.RoleProgress).Label(o.Label).Range(0, 1, used)
	arc := root.AnimateTo("arc", float32(used), ui.Smooth)
	ink := t.Accent
	switch {
	case used >= 0.95:
		ink = t.Danger
	case used >= 0.8:
		ink = t.Warning
	}
	ink = root.AnimateColor("ink", ink, ui.Snappy)
	root.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, min(r.W, r.H)/2-w
		var track, fill ui.Path
		p.StrokePath(track.Circle(cx, cy, rad), w, t.Border)
		if arc > 0 {
			p.Transform(ui.Rotation(-90), cx, cy, func() { p.StrokePath(fill.Circle(cx, cy, rad).Trim(0, arc), w, ink) })
		}
	}))
	return root.Children(func() {
		kit.Numeric(c, fmt.Sprintf(o.Format, used*100)).Bold()
	})
}
