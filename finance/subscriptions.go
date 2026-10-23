package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Subscription is one recurring charge.
type Subscription struct {
	Name, Next string
	Amount     float64
}

// Subscriptions shows recurring charges and next dates.
func Subscriptions(c *ui.Context, rows []Subscription, currency string) ui.Element {
	if len(rows) == 0 || currency == "" {
		panic("finance: Subscriptions needs rows and currency")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, row := range rows {
			if row.Name == "" || row.Next == "" || row.Amount < 0 || !finite(row.Amount) {
				panic("finance: Subscriptions has an invalid row")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, row.Name).Grow(1)
				kit.Muted(c, row.Next).FontSize(t.Rem(0.75))
				money(c, row.Amount, currency)
			})
		}
	})
	return root
}
