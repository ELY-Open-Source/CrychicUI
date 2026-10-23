package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ExchangeRate shows a pair with rate and change.
func ExchangeRate(c *ui.Context, pair string, rate, change float64) ui.Element {
	if pair == "" || rate <= 0 || !finite(rate) || !finite(change) {
		panic("finance: ExchangeRate has invalid data")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Text(c, pair).Bold().Grow(1)
		kit.Numeric(c, fmt.Sprintf("%.4f", rate)).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
		ink := t.Success
		if change < 0 {
			ink = t.Danger
		}
		kit.Numeric(c, fmt.Sprintf("%+.2f%%", change)).TextColor(ink)
	})
	return root
}
