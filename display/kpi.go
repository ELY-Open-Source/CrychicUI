package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// KPI pairs a big number with a sliding trend arrow.
func KPI(c *ui.Context, label string, value, delta float64) ui.Element {
	if !finite(delta) {
		panic("display: KPI delta must be finite")
	}
	t := c.Theme()
	name, color := "minus", t.TextMuted
	if delta > 0 {
		name, color = "trending-up", t.Success
	} else if delta < 0 {
		name, color = "trending-down", t.Danger
	}
	return ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1)).Children(func() {
		kit.Muted(c, label).FontSize(t.Rem(0.85))
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			kit.Numeric(c, formatNumber(c, value)).Bold().FontSize(t.Rem(2))
			arrow := ui.Icon(c, icon.Must(name)).Size(16, 16).TextColor(color)
			kit.Arrive(arrow, ui.Pose{X: -4, Fade: 1}, 0)
		})
	})
}
