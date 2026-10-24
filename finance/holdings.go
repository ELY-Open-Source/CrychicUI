package finance

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Holding is one market position; Name is optional.
type Holding struct {
	Symbol, Name  string
	Shares, Price float64
	Change        float64
}

// HoldingsOptions word the card; Shares takes the share count.
type HoldingsOptions struct{ Title, Shares string }

// Holdings lists positions with value and change under a total.
func Holdings(c *ui.Context, positions []Holding, currency string, o HoldingsOptions) ui.Element {
	if len(positions) == 0 || currency == "" {
		panic("finance: Holdings needs positions and a currency")
	}
	total := 0.0
	for _, p := range positions {
		if p.Symbol == "" || p.Shares < 0 || p.Price < 0 || !finite(p.Shares) || !finite(p.Price) || !finite(p.Change) {
			panic(fmt.Sprintf("finance: Holdings position %q needs finite shares, price and change", p.Symbol))
		}
		total += p.Shares * p.Price
	}
	if !finite(total) {
		panic("finance: Holdings total overflows")
	}
	def(&o.Title, "Holdings")
	def(&o.Shares, "%s shares")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "briefcase", o.Title, func() { money(c, total, currency).Bold().Shrink(0) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3))
		list.Children(func() {
			for i, p := range positions {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(positions)-1)))
				kit.Stagger(c, list, row, i).Children(func() {
					symbolTile(c, p.Symbol)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, p.Symbol).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
						meta := fmt.Sprintf(o.Shares, strconv.FormatFloat(p.Shares, 'f', -1, 64))
						if p.Name != "" {
							meta = p.Name + " · " + meta
						}
						kit.Muted(c, meta).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
					})
					ui.Column(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
						money(c, p.Shares*p.Price, currency).Bold().FontSize(t.Rem(0.9))
						changePill(c, p.Change)
					})
				})
			}
		})
	})
	return root
}
