package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TickerItem is one quoted symbol.
type TickerItem struct {
	Symbol string
	Price  float64
	Change float64
}

// Ticker shows price and change; it flashes on change.
func Ticker(c *ui.Context, symbol string, price, change float64, currency string) ui.Element {
	if symbol == "" || price < 0 || !finite(price) || !finite(change) || currency == "" {
		panic("finance: Ticker has invalid data")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
	previous := ui.Local(root, "price", func() float64 { return price })
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	if price != *previous {
		*previous = price
		*until = c.Now().Add(300 * time.Millisecond)
	}
	flash := c.Now().Before(*until)
	if flash {
		c.After(until.Sub(c.Now()))
	}
	root.Background(root.AnimateColor("flash", t.Surface.Mix(t.Accent.Alpha(0.2), unit(flash)), ui.Snappy))
	root.Children(func() {
		ui.Text(c, symbol).Bold().Grow(1)
		money(c, price, currency)
		ink := t.Success
		if change < 0 {
			ink = t.Danger
		}
		kit.Numeric(c, fmt.Sprintf("%+.2f%%", change)).TextColor(ink)
	})
	return root
}
