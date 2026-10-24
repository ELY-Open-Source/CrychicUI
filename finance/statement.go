package finance

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StatementOptions word the card.
type StatementOptions struct{ In, Out, Total string }

// Statement shows a month's rows, money in and out, net.
func Statement(c *ui.Context, month string, rows []Transaction, currency string, o StatementOptions) ui.Element {
	if month == "" || len(rows) == 0 || currency == "" {
		panic("finance: Statement needs a month, rows and a currency")
	}
	def(&o.In, "Money in")
	def(&o.Out, "Money out")
	def(&o.Total, "Total")
	in, out := 0.0, 0.0
	for _, row := range rows {
		if !finite(row.Amount) {
			panic("finance: Statement row amount must be finite")
		}
		if row.Amount > 0 {
			in += row.Amount
		} else {
			out -= row.Amount
		}
	}
	total := in - out
	if !finite(in) || !finite(out) || !finite(total) {
		panic("finance: Statement total overflow")
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "file-text", month, func() { kit.Chip(c, strconv.Itoa(len(rows)), t.TextMuted) })
		ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3), t.Space(3), 0, t.Space(3)).Children(func() {
			fact(c, o.In, moneyText(c, in, currency))
			fact(c, o.Out, moneyText(c, out, currency))
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3))
		list.Children(func() {
			for i, row := range rows {
				c.Key(i)
				kit.Stagger(c, list, txRow(c, row, currency), i)
			}
		})
		ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(3)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).
			Children(func() {
				ui.Text(c, o.Total).Bold().Grow(1).MinWidth(0).SingleLine()
				money(c, total, currency).Bold().FontSize(t.Rem(1.1)).TextColor(changeInk(t, total)).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
	})
	return root
}
