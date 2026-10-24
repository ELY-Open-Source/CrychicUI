package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Charge is one invoice line.
type Charge struct {
	Label  string
	Amount float64
}

// InvoiceOptions word the invoice; To and Date are caller meta.
type InvoiceOptions struct{ Subtotal, Tax, Total, To, Date string }

// Invoice shows lines, subtotal, tax, total and a status chip.
func Invoice(c *ui.Context, number, status string, lines []Charge, tax float64, currency string, o InvoiceOptions) ui.Element {
	if number == "" || status == "" || len(lines) == 0 || tax < 0 || !finite(tax) || currency == "" {
		panic(fmt.Sprintf("finance: Invoice %q needs a status, lines, tax of zero or more and a currency", number))
	}
	sub := 0.0
	for _, l := range lines {
		if l.Label == "" || !finite(l.Amount) {
			panic(fmt.Sprintf("finance: Invoice %q line %q needs a finite amount", number, l.Label))
		}
		sub += l.Amount
	}
	if !finite(sub + tax) {
		panic(fmt.Sprintf("finance: Invoice %q total overflows", number))
	}
	def(&o.Subtotal, "Subtotal")
	def(&o.Tax, "Tax")
	def(&o.Total, "Total")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	pair := func(label string, value float64) {
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			kit.Muted(c, label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0)
			money(c, value, currency).FontSize(t.Rem(0.85)).Shrink(0)
		})
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			tile(c, "file-text", t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, number).Bold().FillWidth().SingleLine()
				meta := o.To
				if o.Date != "" && meta != "" {
					meta += " · "
				}
				if meta += o.Date; meta != "" {
					kit.Muted(c, meta).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
				}
			})
			kit.Arrive(kit.Chip(c, status, t.Accent), ui.Pose{Scale: 0.8, Fade: 1}, 150*time.Millisecond)
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3), 0).BorderWidth(1, 0, 1, 0).BorderColor(t.Border)
		list.Children(func() {
			for i, l := range lines {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Text(c, l.Label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
					money(c, l.Amount, currency).FontSize(t.Rem(0.85)).Shrink(0)
				})
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
			pair(o.Subtotal, sub)
			pair(o.Tax, tax)
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1).MinWidth(0)
			money(c, sub+tax, currency).Bold().FontSize(t.Rem(1.2)).Shrink(0)
		})
	})
	return root
}
