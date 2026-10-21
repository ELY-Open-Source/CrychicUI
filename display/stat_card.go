package display

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StatCard pairs a rolling value with a delta.
func StatCard(c *ui.Context, label string, value, delta float64) ui.Element {
	if !finite(delta) {
		panic("display: StatCard delta must be finite")
	}
	t := c.Theme()
	root := kit.Card(c).AlignSelf(ui.Start).Gap(t.Space(2)).Padding(t.Space(4))
	return root.Children(func() {
		kit.Muted(c, label).FontSize(t.Rem(0.85))
		kit.Numeric(c, formatNumber(c, value)).Bold().FontSize(t.Rem(1.5))
		color := t.TextMuted
		if delta > 0 {
			color = t.Success
		} else if delta < 0 {
			color = t.Danger
		}
		kit.Numeric(c, fmt.Sprintf("%+.1f%%", delta)).TextColor(color).FontSize(t.Rem(0.85))
	})
}
