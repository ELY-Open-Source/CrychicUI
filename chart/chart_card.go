package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ChartCardOptions tune a ChartCard.
type ChartCardOptions struct {
	// Delta prints the change; "%+.1f%%" when empty.
	Delta string
	// Down marks a fall as good, as for costs.
	Down bool
}

// ChartCard frames a chart under a title, value and change.
func ChartCard(c *ui.Context, title, value string, delta float64, o ChartCardOptions, content func()) ui.Element {
	if title == "" || content == nil {
		panic("chart: ChartCard needs a title and content")
	}
	finite("ChartCard", title, delta)
	def(&o.Delta, "%+.1f%%")
	t := c.Theme()
	good := (delta >= 0) != o.Down
	ink, glyph := t.Success, "trending-up"
	if delta < 0 {
		glyph = "trending-down"
	}
	if !good {
		ink = t.Danger
	}
	card := kit.Card(c).FillWidth().Padding(t.Space(4)).Gap(t.Space(3))
	kit.Arrive(card, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return card.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
			kit.Muted(c, title).Grow(1).FontSize(t.Rem(0.85))
			badge := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999)
			badge.Background(badge.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).Children(func() {
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(ink).ContentTransition(ui.ContentTransition{})
				kit.Numeric(c, fmt.Sprintf(o.Delta, delta)).FontSize(t.Rem(0.8)).Bold().TextColor(ink).
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
		})
		kit.Numeric(c, value).FontSize(t.Rem(1.6)).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
		content()
	})
}
