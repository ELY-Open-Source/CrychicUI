package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ReceiptOptions are the receipt labels.
type ReceiptOptions struct{ Total string }

// Receipt shows an itemized receipt with total.
func Receipt(c *ui.Context, number string, items []CartItem, currency string, o ReceiptOptions) ui.Element {
	if number == "" || len(items) == 0 || currency == "" {
		panic("shop: Receipt needs a number, items, and currency")
	}
	def(&o.Total, "Total")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	total := 0.0
	root.Children(func() {
		ui.Text(c, number).Bold()
		for i, item := range items {
			validProduct("Receipt", item.Product)
			if item.Quantity <= 0 {
				panic("shop: Receipt has an invalid quantity")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, item.Product.Name).Grow(1)
				kit.Numeric(c, "×"+fmt.Sprint(item.Quantity))
				money(c, item.Product.Price*float64(item.Quantity), currency)
			})
			total += item.Product.Price * float64(item.Quantity)
		}
		ui.Row(c).FillWidth().BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1)
			money(c, total, currency).Bold()
		})
	})
	return root
}
