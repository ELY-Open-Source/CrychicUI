package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ErrorRateOptions are the error rate labels.
type ErrorRateOptions struct{ Rate, Limit string }

// ErrorRate shows errors per minute with a limit line.
func ErrorRate(c *ui.Context, rate, limit float64, o ErrorRateOptions) ui.Element {
	if rate < 0 || limit <= 0 || !finite(rate) || !finite(limit) {
		panic("dashboard: ErrorRate has invalid values")
	}
	def(&o.Rate, "Errors per minute")
	def(&o.Limit, "Limit %.0f")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, o.Rate).Grow(1)
			kit.Numeric(c, fmt.Sprintf("%.1f", rate)).Bold()
			kit.Muted(c, fmt.Sprintf(o.Limit, limit)).FontSize(t.Rem(0.75))
		})
		ink := t.Success
		if rate > limit {
			ink = t.Danger
		}
		maxValue := max(rate, limit)
		bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
		ratio := bar.AnimateTo("rate", float32(rate/maxValue), ui.Smooth)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * ratio, H: r.H}, ink, 999)
			x := r.X + r.W*float32(limit/maxValue)
			p.Fill(ui.Rect{X: x - 0.5, Y: r.Y, W: 1, H: r.H}, t.Text, 0)
		}))
	})
	return root
}
