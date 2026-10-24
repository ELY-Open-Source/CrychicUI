package finance

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BookLevel is one order book price and amount.
type BookLevel struct {
	Price, Amount float64
}

// OrderBookOptions word the book; Spread takes the gap.
type OrderBookOptions struct{ Title, Bids, Asks, Price, Amount, Spread string }

// sides validates a book and sorts each side descending.
func sides(component string, bids, asks []BookLevel) ([]BookLevel, []BookLevel) {
	if len(bids) == 0 || len(asks) == 0 {
		panic(fmt.Sprintf("finance: %s needs bids and asks", component))
	}
	for _, l := range slices.Concat(bids, asks) {
		if l.Price <= 0 || l.Amount <= 0 || !finite(l.Price) || !finite(l.Amount) {
			panic(fmt.Sprintf("finance: %s level at %v needs a positive price and amount", component, l.Price))
		}
	}
	down := func(a, b BookLevel) int { return cmp.Compare(b.Price, a.Price) }
	bids, asks = slices.SortedFunc(slices.Values(bids), down), slices.SortedFunc(slices.Values(asks), down)
	if bids[0].Price >= asks[len(asks)-1].Price {
		panic(fmt.Sprintf("finance: %s is crossed: best bid %v meets best ask %v", component, bids[0].Price, asks[len(asks)-1].Price))
	}
	return bids, asks
}

// OrderBook shows asks over bids with depth; changed rows flash.
func OrderBook(c *ui.Context, bids, asks []BookLevel, currency string, o OrderBookOptions) ui.Element {
	if currency == "" {
		panic("finance: OrderBook needs a currency")
	}
	bids, asks = sides("OrderBook", bids, asks)
	def(&o.Title, "Order book")
	def(&o.Bids, "Bids")
	def(&o.Asks, "Asks")
	def(&o.Price, "Price")
	def(&o.Amount, "Amount")
	def(&o.Spread, "Spread %s")
	peak := 0.0
	for _, l := range slices.Concat(bids, asks) {
		peak = max(peak, l.Amount)
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	level := func(list ui.Element, i int, key string, l BookLevel, ink ui.Color) {
		c.Key(key)
		row := ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1), t.Space(3))
		kit.Stagger(c, list, row, i)
		was := ui.Local(row, "level", func() BookLevel { return l })
		until := ui.Local(row, "until", func() time.Time { return time.Time{} })
		if *was != l {
			*was = l
			*until = c.Now().Add(400 * time.Millisecond)
		}
		flash := c.Now().Before(*until)
		if flash {
			c.After(until.Sub(c.Now()))
		}
		row.Background(row.AnimateColor("flash", ink.Alpha(0.22*unit(flash)), ui.Snappy))
		depth := row.AnimateTo("depth", float32(l.Amount/peak), ui.Smooth)
		row.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			w := r.W * depth
			p.Fill(ui.Rect{X: r.X + r.W - w, Y: r.Y, W: w, H: r.H}, ink.Alpha(0.12), 0)
		}))
		row.Children(func() {
			money(c, l.Price, currency).FontSize(t.Rem(0.8)).TextColor(ink).Grow(1).MinWidth(0)
			kit.Numeric(c, strconv.FormatFloat(l.Amount, 'f', -1, 64)).FontSize(t.Rem(0.8)).Shrink(0)
		})
	}
	root.Children(func() {
		kit.CardHead(c, "chart-no-axes-column", o.Title, func() {
			kit.Chip(c, o.Bids, t.Success)
			kit.Chip(c, o.Asks, t.Danger)
		})
		ui.Row(c).FillWidth().Padding(t.Space(2), t.Space(3), t.Space(1), t.Space(3)).Children(func() {
			kit.Muted(c, o.Price).FontSize(t.Rem(0.7)).Bold().Grow(1)
			kit.Muted(c, o.Amount).FontSize(t.Rem(0.7)).Bold().Shrink(0)
		})
		for i, l := range asks {
			level(root, i, fmt.Sprintf("ask %d", i), l, t.Danger)
		}
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
			BorderWidth(1, 0, 1, 0).BorderColor(t.Border).Background(t.Surface).Children(func() {
			money(c, bids[0].Price/2+asks[len(asks)-1].Price/2, currency).Bold().Grow(1).MinWidth(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
			kit.Muted(c, fmt.Sprintf(o.Spread, moneyText(c, asks[len(asks)-1].Price-bids[0].Price, currency))).FontSize(t.Rem(0.75)).Shrink(0)
		})
		for i, l := range bids {
			level(root, len(asks)+1+i, fmt.Sprintf("bid %d", i), l, t.Success)
		}
		ui.Box(c).Height(t.Space(1))
	})
	return root
}
