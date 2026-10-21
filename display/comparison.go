package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Comparison shows before and after values with an arrow.
func Comparison(c *ui.Context, before, after float64) ui.Element {
	t := c.Theme()
	color := t.TextMuted
	if after > before {
		color = t.Success
	} else if after < before {
		color = t.Danger
	}
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		kit.Muted(c, formatNumber(c, before))
		ui.Icon(c, icon.Must("arrow-right")).Size(16, 16).TextColor(t.TextMuted)
		kit.Numeric(c, formatNumber(c, after)).TextColor(color).Bold()
	})
}
