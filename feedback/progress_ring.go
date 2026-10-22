package feedback

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ProgressRingOptions tune a ProgressRing.
type ProgressRingOptions struct {
	// Format writes the percent; "%.0f%%" when empty.
	Format string
}

// ProgressRing is a circle whose arc draws to the value.
func ProgressRing(c *ui.Context, label string, value float64, o ProgressRingOptions) ui.Element {
	if label == "" {
		panic("feedback: ProgressRing needs a label")
	}
	v := fraction("ProgressRing", value)
	if o.Format == "" {
		o.Format = "%.0f%%"
	}
	t := c.Theme()
	size, w := t.Space(16), t.Space(1.5)
	root := ui.Box(c).Size(size, size).Center().Role(ui.RoleProgress).Label(label).Range(0, 1, value)
	arc := root.AnimateTo("arc", v, ui.Smooth)
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-w
		var track, fill ui.Path
		p.StrokePath(track.Circle(cx, cy, rad), w, t.Border)
		if arc > 0 {
			p.Transform(ui.Rotation(-90), cx, cy, func() { p.StrokePath(fill.Circle(cx, cy, rad).Trim(0, arc), w, t.Accent) })
		}
	})
	return root.Children(func() {
		kit.Numeric(c, fmt.Sprintf(o.Format, value*100)).Bold()
	})
}
