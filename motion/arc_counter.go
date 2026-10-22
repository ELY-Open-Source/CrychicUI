package motion

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// ArcCounterOptions tune an ArcCounter.
type ArcCounterOptions struct {
	// Format writes the number; "%.0f" when empty.
	Format string
}

var climb = ui.Animation{Duration: 800 * time.Millisecond, Ease: ui.EaseOut}

// ArcCounter draws an arc to value of total, counting inside.
func ArcCounter(c *ui.Context, label string, value, total float64, o ArcCounterOptions) ui.Element {
	finite("ArcCounter", value)
	finite("ArcCounter", total)
	if label == "" || total <= 0 || value < 0 || value > total {
		panic("motion: ArcCounter needs a label and 0 ≤ value ≤ total")
	}
	if o.Format == "" {
		o.Format = "%.0f"
	}
	t := c.Theme()
	d, w := t.Space(24), t.Space(2)
	root := ui.Box(c).Size(d, d).Shrink(0).Center().Role(ui.RoleMeter).Label(label).Range(0, total, value)
	fresh := ui.Local(root, "fresh", func() bool { return true })
	if *fresh {
		*fresh = false
		root.Jump("v", 0)
	}
	v := root.AnimateTo("v", float32(value/total), climb)
	shown := float64(v) * total
	if root.Settled("v") {
		shown = value
	}
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, (r.W-w)/2
		var track, arc ui.Path
		p.StrokePath(track.Circle(cx, cy, rad), w, t.Border)
		if v > 0 {
			p.Transform(ui.Rotation(-90), cx, cy, func() { p.StrokePath(arc.Circle(cx, cy, rad).Trim(0, v), w, t.Accent) })
		}
	})
	return root.Children(func() {
		ui.Text(c, fmt.Sprintf(o.Format, shown)).Bold().FontSize(t.Rem(1.25)).FontFeatures("tnum")
	})
}
