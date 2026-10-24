package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Dividend is one dividend payment; Name is optional.
type Dividend struct {
	Symbol, Name, Date string
	Amount             float64
}

// DividendsOptions word the card.
type DividendsOptions struct{ Title string }

// Dividends lists payments by date under their total.
func Dividends(c *ui.Context, payments []Dividend, currency string, o DividendsOptions) ui.Element {
	if len(payments) == 0 || currency == "" {
		panic("finance: Dividends needs payments and a currency")
	}
	total := 0.0
	for _, payment := range payments {
		if payment.Symbol == "" || payment.Date == "" || payment.Amount < 0 || !finite(payment.Amount) {
			panic(fmt.Sprintf("finance: Dividends payment %q needs a date and an amount of zero or more", payment.Symbol))
		}
		total += payment.Amount
	}
	if !finite(total) {
		panic("finance: Dividends total overflows")
	}
	def(&o.Title, "Dividends")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "coins", o.Title, func() {
			money(c, total, currency).Bold().TextColor(t.Success).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3))
		list.Children(func() {
			for i, payment := range payments {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(payments)-1)))
				kit.Stagger(c, list, row, i).Children(func() {
					symbolTile(c, payment.Symbol)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						name := payment.Symbol
						if payment.Name != "" {
							name += " · " + payment.Name
						}
						ui.Text(c, name).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
							ui.Icon(c, icon.Must("calendar")).Size(12, 12)
							kit.Muted(c, payment.Date).FontSize(t.Rem(0.75)).Shrink(0)
						})
					})
					ui.Row(c).Shrink(0).TextColor(t.Success).Children(func() {
						ui.Text(c, "+").Bold().FontSize(t.Rem(0.85)).Shrink(0)
						money(c, payment.Amount, currency).Bold().FontSize(t.Rem(0.85)).Shrink(0)
					})
				})
			}
		})
	})
	return root
}
