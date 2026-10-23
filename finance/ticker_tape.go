package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TickerTape scrolls a row of tickers.
func TickerTape(c *ui.Context, items []TickerItem, currency string) ui.Element {
	if len(items) == 0 || currency == "" {
		panic("finance: TickerTape needs items and currency")
	}
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(t.Space(8)).Clip()
	phase := root.Loop("tape", 12*time.Second, ui.Linear)
	cycle := ui.Local(root, "cycle", func() float32 { return 0 })
	root.Children(func() {
		row := ui.Row(c).Gap(t.Space(6)).Translate(-float32(phase)*(*cycle), 0)
		var first, second ui.Element
		row.Children(func() {
			for repeat := range 2 {
				for i, item := range items {
					if item.Symbol == "" || item.Price < 0 || !finite(item.Price) || !finite(item.Change) {
						panic("finance: TickerTape has an invalid item")
					}
					c.Key(repeat*len(items) + i)
					itemRow := ui.Row(c).Gap(t.Space(1))
					if repeat == 0 && i == 0 {
						first = itemRow
					}
					if repeat == 1 && i == 0 {
						second = itemRow
					}
					itemRow.Children(func() {
						ui.Text(c, item.Symbol).Bold()
						money(c, item.Price, currency)
						kit.Numeric(c, fmt.Sprintf("%+.2f%%", item.Change)).TextColor(t.TextMuted)
					})
				}
			}
		})
		distance := second.LayoutBounds().X - first.LayoutBounds().X
		kit.Settle(c, root, []float32{distance}, func() []float32 {
			return []float32{second.LayoutBounds().X - first.LayoutBounds().X}
		})
		if distance > 0 {
			*cycle = distance
		}
	})
	return root
}
