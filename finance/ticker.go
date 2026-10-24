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

// TickerOptions add an optional company name.
type TickerOptions struct{ Name string }

// Ticker shows a quote that flashes as the price moves.
func Ticker(c *ui.Context, symbol string, price, change float64, currency string, o TickerOptions) ui.Element {
	if symbol == "" || price < 0 || !finite(price) || !finite(change) || currency == "" {
		panic(fmt.Sprintf("finance: Ticker %q needs a finite price and change and a currency", symbol))
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup).Label(symbol)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	previous := ui.Local(root, "price", func() float64 { return price })
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	up := ui.Local(root, "up", func() bool { return true })
	if price != *previous {
		*up = price > *previous
		*previous = price
		*until = c.Now().Add(300 * time.Millisecond)
	}
	flash := c.Now().Before(*until)
	if flash {
		c.After(until.Sub(c.Now()))
	}
	ink := t.Danger
	if *up {
		ink = t.Success
	}
	root.Background(root.AnimateColor("flash", t.Background.Mix(ink.Alpha(0.16), unit(flash)), ui.Snappy))
	root.Row().AlignItems(ui.Center).Children(func() {
		symbolTile(c, symbol)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, symbol).Bold().FillWidth().SingleLine()
			if o.Name != "" {
				kit.Muted(c, o.Name).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
			}
		})
		ui.Column(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
			money(c, price, currency).Bold().FontSize(t.Rem(1.1)).ContentTransition(ui.ContentTransition{Numeric: true})
			changePill(c, change)
		})
	})
	return root
}
