package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Price shows a current price and an old struck price.
func Price(c *ui.Context, current, old float64, currency string) ui.Element {
	if currency == "" || !finite(current) || !finite(old) || current < 0 || old < current {
		panic("display: Price needs nonnegative values and an old price")
	}
	t := c.Theme()
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		kit.Numeric(c, currency+" "+formatMoney(c, current)).Bold().FontSize(t.Rem(1.25))
		if old > current {
			kit.Muted(c, currency+" "+formatMoney(c, old)).Strikethrough()
		}
	})
}
