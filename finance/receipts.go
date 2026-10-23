package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Receipt is one scanned receipt.
type Receipt struct {
	Merchant, Date string
	Amount         float64
}

// ReceiptsOptions are the receipt labels.
type ReceiptsOptions struct{ Total string }

// Receipts shows scanned receipts and totals.
func Receipts(c *ui.Context, receipts []Receipt, currency string, o ReceiptsOptions) ui.Element {
	if len(receipts) == 0 || currency == "" {
		panic("finance: Receipts needs receipts and currency")
	}
	def(&o.Total, "Total")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	total := 0.0
	root.Children(func() {
		for i, receipt := range receipts {
			if receipt.Merchant == "" || receipt.Date == "" || receipt.Amount < 0 || !finite(receipt.Amount) {
				panic("finance: Receipts has an invalid receipt")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Text(c, receipt.Merchant).Grow(1)
				kit.Muted(c, receipt.Date).FontSize(t.Rem(0.75))
				money(c, receipt.Amount, currency)
			})
			total += receipt.Amount
			if !finite(total) {
				panic("finance: Receipts total overflow")
			}
		}
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1)
			money(c, total, currency).Bold()
		})
	})
	return root
}
