package feedback

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GaugeOptions tune a Gauge.
type GaugeOptions struct {
	// Format writes the value; "%.0f" when empty.
	Format string
}

// Gauge is a half circle whose needle springs.
func Gauge(c *ui.Context, label string, value, lo, hi float64, o GaugeOptions) ui.Element {
	if label == "" || !(lo < hi) || math.IsInf(lo, 0) || math.IsInf(hi, 0) || !(value >= lo && value <= hi) {
		panic("feedback: Gauge needs a label and a value within finite lo < hi")
	}
	if o.Format == "" {
		o.Format = "%.0f"
	}
	t := c.Theme()
	w := t.Space(1.5)
	root := ui.Column(c).AlignSelf(ui.Start).AlignItems(ui.Center).Width(t.Space(36)).Role(ui.RoleMeter).Label(label).
		Range(lo, hi, value)
	at := root.AnimateTo("needle", float32((value-lo)/(hi-lo)), ui.Bouncy)
	return root.Children(func() {
		ui.Box(c).FillWidth().Height(t.Space(19)).Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy, rad := r.X+r.W/2, r.Y+r.H-w, r.W/2-w
			var track, fill ui.Path
			p.Transform(ui.Rotation(180), cx, cy, func() {
				p.StrokePath(track.Circle(cx, cy, rad).Trim(0, 0.5), w, t.Border)
				if f := min(max(at, 0), 1); f > 0 {
					p.StrokePath(fill.Circle(cx, cy, rad).Trim(0, 0.5*f), w, t.Accent)
				}
			})
			x, y := polar(cx, cy, rad-w*2, 180+180*float64(at))
			var needle ui.Path
			p.StrokePath(needle.MoveTo(cx, cy).LineTo(x, y), w, t.Text)
			p.Fill(ui.Rect{X: cx - w*1.5, Y: cy - w*1.5, W: w * 3, H: w * 3}, t.Text, w*1.5)
		})
		kit.Numeric(c, fmt.Sprintf(o.Format, value)).Bold().Margin(t.Space(1), 0, 0, 0)
	})
}
