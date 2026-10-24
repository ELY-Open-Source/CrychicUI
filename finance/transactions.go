package finance

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Transaction is one merchant transaction; When and Icon are optional.
type Transaction struct {
	Merchant, Category string
	Amount             float64
	When, Icon         string
}

// TransactionsOptions word the card and its empty state.
type TransactionsOptions struct{ Title, Empty string }

// Transactions lists rows under day headings, income in green.
func Transactions(c *ui.Context, rows []Transaction, currency string, o TransactionsOptions) ui.Element {
	if currency == "" {
		panic("finance: Transactions needs a currency")
	}
	def(&o.Title, "Transactions")
	def(&o.Empty, "No transactions yet")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "receipt", o.Title, func() { kit.Chip(c, fmt.Sprint(len(rows)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3))
		list.Children(func() {
			if len(rows) == 0 {
				ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(6), 0).TextColor(t.TextMuted).Children(func() {
					ui.Icon(c, icon.Must("receipt")).Size(24, 24)
					ui.Text(c, o.Empty).FontSize(t.Rem(0.85)).Shrink(0)
				})
				return
			}
			day := ""
			for i, row := range rows {
				if row.When != "" && row.When != day {
					day = row.When
					c.Key(fmt.Sprintf("day %q %d", day, i))
					kit.Muted(c, strings.ToUpper(day)).FontSize(t.Rem(0.7)).Bold().Padding(t.Space(2), 0, t.Space(0.5), 0)
				}
				c.Key(i)
				row.When = ""
				kit.Stagger(c, list, txRow(c, row, currency), i)
			}
		})
	})
	return root
}
