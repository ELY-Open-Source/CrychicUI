package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HealthRingOptions word the card; Note sits under the band.
type HealthRingOptions struct{ Label, Good, Fair, Poor, Note string }

// HealthRing shows health as a filling ring and a band.
func HealthRing(c *ui.Context, score int, o HealthRingOptions) ui.Element {
	if score < 0 || score > 100 {
		panic(fmt.Sprintf("dashboard: HealthRing score %d must be 0–100", score))
	}
	def(&o.Label, "Health")
	def(&o.Good, "Healthy")
	def(&o.Fair, "Needs attention")
	def(&o.Poor, "Critical")
	t := c.Theme()
	band, ink := o.Poor, t.Danger
	switch {
	case score >= 80:
		band, ink = o.Good, t.Success
	case score >= 50:
		band, ink = o.Fair, t.Warning
	}
	root := card(c).Row().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4))
	root.Children(func() {
		ring := ui.Box(c).Size(t.Space(28), t.Space(28)).Shrink(0).Center()
		fill := ring.AnimateTo("health", float32(score)/100, ui.Smooth)
		tint := ring.AnimateColor("ink", ink, ui.Snappy)
		ring.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			path.Circle(r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(2))
			p.StrokePath(&path, t.Space(3), t.SurfaceHover)
			if fill > 0 {
				p.StrokePath(path.Trim(0, fill), t.Space(3), tint)
			}
		}))
		ring.Children(func() {
			kit.Numeric(c, fmt.Sprint(score)).Bold().FontSize(t.Rem(1.8)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
			kit.Muted(c, o.Label).FontSize(t.Rem(0.8)).Shrink(0)
			kit.Chip(c, band, ink).AlignSelf(ui.Start)
			if o.Note != "" {
				ui.Text(c, o.Note).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).FillWidth().Wrap()
			}
		})
	})
	return root
}
