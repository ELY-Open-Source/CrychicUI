package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Charge is one invoice line.
type Charge struct {
	Label  string
	Amount float64
}

// InvoiceOptions are the invoice labels.
type InvoiceOptions struct{ Tax, Total string }

// Invoice shows lines, tax, total, and status.
func Invoice(c *ui.Context, number, status string, lines []Charge, tax float64, currency string, o InvoiceOptions) ui.Element {
	if number == "" || status == "" || len(lines) == 0 || tax < 0 || !finite(tax) || currency == "" {
		panic("finance: Invoice has invalid data")
	}
	def(&o.Tax, "Tax")
	def(&o.Total, "Total")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	total := tax
	root.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, number).Bold().Grow(1)
			ui.Text(c, status).TextColor(t.Accent)
		})
		for i, line := range lines {
			if line.Label == "" || !finite(line.Amount) {
				panic("finance: Invoice has an invalid line")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, line.Label).Grow(1)
				money(c, line.Amount, currency)
			})
			total += line.Amount
		}
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, o.Tax).Grow(1)
			money(c, tax, currency)
		})
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1)
			money(c, total, currency).Bold()
		})
	})
	return root
}
