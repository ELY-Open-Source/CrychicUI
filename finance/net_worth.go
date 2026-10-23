package finance

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NetPoint is one net-worth point in time.
type NetPoint struct {
	Label         string
	Assets, Debts float64
}

// NetWorthOptions are the net-worth labels.
type NetWorthOptions struct{ Value string }

// NetWorth shows assets minus debts over time.
func NetWorth(c *ui.Context, points []NetPoint, currency string, o NetWorthOptions) ui.Element {
	if len(points) < 2 || currency == "" {
		panic("finance: NetWorth needs points and currency")
	}
	def(&o.Value, "Net worth")
	t := c.Theme()
	values := make([]float64, len(points))
	lo, hi := 0.0, 0.0
	for i, point := range points {
		if point.Label == "" || point.Assets < 0 || point.Debts < 0 || !finite(point.Assets) || !finite(point.Debts) {
			panic("finance: NetWorth has an invalid point")
		}
		values[i] = point.Assets - point.Debts
		lo, hi = min(lo, values[i]), max(hi, values[i])
	}
	span := hi - lo
	if !finite(span) {
		panic("finance: NetWorth range overflow")
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Value).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		money(c, values[len(values)-1], currency).Bold()
		box := ui.Box(c).FillWidth().Height(t.Space(32))
		progress := box.Keyframes("net", []ui.Keyframe{{Value: 0}, {At: 400 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			for i, value := range values {
				x := r.X + r.W*float32(i)/float32(len(values)-1)
				y := r.Y + r.H*0.5
				if hi > lo {
					y = r.Y + r.H*(1-float32((value-lo)/span))
				}
				if i == 0 {
					path.MoveTo(x, y)
				} else {
					path.LineTo(x, y)
				}
			}
			p.StrokePath(path.Trim(0, progress), 2, t.Accent)
		}))
		ui.Row(c).FillWidth().Children(func() {
			kit.Muted(c, points[0].Label).FontSize(t.Rem(0.75)).Grow(1)
			kit.Muted(c, points[len(points)-1].Label).FontSize(t.Rem(0.75))
		})
	})
	return root
}
