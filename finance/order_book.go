package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BookLevel is one order book price and amount.
type BookLevel struct {
	Price, Amount float64
}

// OrderBookOptions are the order book labels.
type OrderBookOptions struct{ Bids, Asks string }

// OrderBook shows bids and asks with depth bars.
func OrderBook(c *ui.Context, bids, asks []BookLevel, currency string, o OrderBookOptions) ui.Element {
	if len(bids) == 0 || len(asks) == 0 || currency == "" {
		panic("finance: OrderBook needs both sides and currency")
	}
	def(&o.Bids, "Bids")
	def(&o.Asks, "Asks")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, side := range []struct {
			name string
			rows []BookLevel
			ink  ui.Color
		}{{o.Bids, bids, t.Success}, {o.Asks, asks, t.Danger}} {
			ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Children(func() {
				ui.Text(c, side.name).Bold()
				maxAmount := 0.0
				for _, row := range side.rows {
					maxAmount = max(maxAmount, row.Amount)
				}
				for i, row := range side.rows {
					if row.Price <= 0 || row.Amount <= 0 || !finite(row.Price) || !finite(row.Amount) {
						panic("finance: OrderBook has an invalid level")
					}
					c.Key(i)
					line := ui.Column(c).FillWidth().Gap(t.Space(0.5))
					previousAmount := ui.Local(line, "amount", func() float64 { return row.Amount })
					previousPrice := ui.Local(line, "price", func() float64 { return row.Price })
					until := ui.Local(line, "until", func() time.Time { return time.Time{} })
					if row.Amount != *previousAmount || row.Price != *previousPrice {
						*previousAmount, *previousPrice = row.Amount, row.Price
						*until = c.Now().Add(300 * time.Millisecond)
					}
					flash := c.Now().Before(*until)
					if flash {
						c.After(until.Sub(c.Now()))
					}
					line.Background(line.AnimateColor("flash", t.Surface.Mix(t.Accent.Alpha(0.2), unit(flash)), ui.Snappy))
					line.Children(func() {
						ui.Row(c).FillWidth().Children(func() {
							money(c, row.Price, currency).Grow(1)
							kit.Numeric(c, fmt.Sprint(row.Amount))
						})
						progress(c, row.Amount/maxAmount, side.ink)
					})
				}
			})
		}
	})
	return root
}
