package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Dividend is one dividend payment.
type Dividend struct {
	Symbol, Date string
	Amount       float64
}

// Dividends shows payments by date.
func Dividends(c *ui.Context, payments []Dividend, currency string) ui.Element {
	if len(payments) == 0 || currency == "" {
		panic("finance: Dividends needs payments and currency")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, payment := range payments {
			if payment.Symbol == "" || payment.Date == "" || payment.Amount < 0 || !finite(payment.Amount) {
				panic("finance: Dividends has an invalid payment")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Text(c, payment.Symbol).Bold().Grow(1)
				kit.Muted(c, payment.Date).FontSize(t.Rem(0.75))
				money(c, payment.Amount, currency)
			})
		}
	})
	return root
}
