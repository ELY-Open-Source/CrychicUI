package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Stat is a muted label over a rolling value.
func Stat(c *ui.Context, label string, value float64) ui.Element {
	t := c.Theme()
	return ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(0.5)).Children(func() {
		kit.Muted(c, label).FontSize(t.Rem(0.85))
		kit.Numeric(c, formatNumber(c, value)).Bold().FontSize(t.Rem(1.25))
	})
}
