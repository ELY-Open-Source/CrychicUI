package shop

import "github.com/egoist/mygo/ui"

// SummaryLine is one order amount.
type SummaryLine struct {
	Label  string
	Amount float64
}

// OrderSummaryOptions are the summary labels.
type OrderSummaryOptions struct{ Total string }

// OrderSummary shows lines, discounts, tax, and total.
func OrderSummary(c *ui.Context, lines []SummaryLine, currency string, o OrderSummaryOptions) ui.Element {
	if len(lines) == 0 || currency == "" {
		panic("shop: OrderSummary needs lines and currency")
	}
	def(&o.Total, "Total")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	total := 0.0
	root.Children(func() {
		for i, line := range lines {
			if line.Label == "" || !finite(line.Amount) {
				panic("shop: OrderSummary has an invalid line")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, line.Label).Grow(1)
				signedMoney(c, line.Amount, currency)
			})
			total += line.Amount
		}
		ui.Row(c).FillWidth().Padding(t.Space(1), 0).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1)
			signedMoney(c, total, currency).Bold()
		})
	})
	return root
}
