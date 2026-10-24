package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TickerTape scrolls a dark strip of quotes in a loop.
func TickerTape(c *ui.Context, items []TickerItem, currency string) ui.Element {
	if len(items) == 0 || currency == "" {
		panic("finance: TickerTape needs items and a currency")
	}
	for _, item := range items {
		if item.Symbol == "" || item.Price < 0 || !finite(item.Price) || !finite(item.Change) {
			panic(fmt.Sprintf("finance: TickerTape item %q needs a finite price and change", item.Symbol))
		}
	}
	t := c.Theme()
	fill, text := inversePair(t)
	root := ui.Box(c).FillWidth().Height(t.Space(10)).Radius(t.Radius).Clip().Background(fill).TextColor(text)
	phase := root.Loop("tape", 12*time.Second, ui.Linear)
	cycle := ui.Local(root, "cycle", func() float32 { return 0 })
	root.Children(func() {
		row := ui.Row(c).Gap(t.Space(6)).AlignItems(ui.Center).Padding(0, t.Space(3)).FillHeight().Translate(-float32(phase)*(*cycle), 0)
		var first, second ui.Element
		row.Children(func() {
			for repeat := range 2 {
				for i, item := range items {
					c.Key(repeat*len(items) + i)
					itemRow := ui.Row(c).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center)
					if i == 0 {
						if repeat == 0 {
							first = itemRow
						} else {
							second = itemRow
						}
					}
					ink := changeInk(t, item.Change)
					glyph := "chevron-up"
					if item.Change < 0 {
						glyph = "chevron-down"
					}
					itemRow.Children(func() {
						ui.Text(c, item.Symbol).Bold().FontSize(t.Rem(0.85)).Shrink(0)
						money(c, item.Price, currency).FontSize(t.Rem(0.85)).TextColor(text.Alpha(0.8)).Shrink(0)
						ui.Icon(c, icon.Must(glyph)).Size(12, 12).TextColor(ink)
						kit.Numeric(c, fmt.Sprintf("%+.2f%%", item.Change)).FontSize(t.Rem(0.8)).Bold().TextColor(ink).Shrink(0)
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
