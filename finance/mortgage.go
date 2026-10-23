package finance

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// MortgageOptions are the mortgage labels.
type MortgageOptions struct{ Year, Principal, Interest string }

// Mortgage shows principal and interest over years.
func Mortgage(c *ui.Context, principal, annualRate float64, years int, currency string, o MortgageOptions) ui.Element {
	if principal <= 0 || annualRate < 0 || years <= 0 || !finite(principal) || !finite(annualRate) || currency == "" {
		panic("finance: Mortgage has invalid data")
	}
	def(&o.Year, "Year %d")
	def(&o.Principal, "Principal")
	def(&o.Interest, "Interest")
	t := c.Theme()
	payment := loanPayment(principal, annualRate, years)
	interest := payment*float64(years)*12 - principal
	if !finite(interest) {
		panic("finance: Mortgage interest overflow")
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, row := range []struct {
			name  string
			value float64
		}{{o.Principal, principal}, {o.Interest, interest}} {
			ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, row.name).Grow(1)
				money(c, row.value, currency)
			})
		}
		ui.Scroll(c).Height(t.Space(70)).Children(func() {
			remaining, rate := principal, annualRate/12
			for year := 1; year <= years; year++ {
				yearPrincipal, yearInterest := 0.0, 0.0
				for range 12 {
					monthInterest := remaining * rate
					principalPaid := min(remaining, payment-monthInterest)
					yearPrincipal += principalPaid
					yearInterest += monthInterest
					remaining = max(0, remaining-principalPaid)
				}
				c.Key(year)
				ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
					ui.Text(c, fmt.Sprintf(o.Year, year)).Grow(1)
					money(c, yearPrincipal, currency)
					money(c, yearInterest, currency)
				})
			}
		})
	})
	return root
}
