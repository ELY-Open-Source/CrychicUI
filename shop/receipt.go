package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ReceiptOptions word the slip; Meta sits under the number.
type ReceiptOptions struct{ Title, Total, Each, Meta string }

// Receipt is an itemized slip with dashed rules and total.
func Receipt(c *ui.Context, number string, items []CartItem, currency string, o ReceiptOptions) ui.Element {
	if number == "" || len(items) == 0 || currency == "" {
		panic(fmt.Sprintf("shop: Receipt %q needs items and a currency", number))
	}
	total := 0.0
	for _, item := range items {
		validProduct("Receipt", item.Product)
		if item.Quantity <= 0 {
			panic(fmt.Sprintf("shop: Receipt %q line %q needs a quantity", number, item.Product.Name))
		}
		total += item.Product.Price * float64(item.Quantity)
	}
	def(&o.Title, "Receipt")
	def(&o.Total, "Total")
	def(&o.Each, "%d × %s")
	t := c.Theme()
	rule := func() {
		ui.Box(c).FillWidth().Height(0).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).BorderStyle(ui.BorderDashed)
	}
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Center().Radius(999).Background(t.Surface).TextColor(t.TextMuted).
				Children(func() { ui.Icon(c, icon.Must("receipt")).Size(18, 18) })
			ui.Text(c, o.Title).Bold().Shrink(0)
			kit.Numeric(c, number).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).Shrink(0)
			if o.Meta != "" {
				kit.Muted(c, o.Meta).FontSize(t.Rem(0.75)).TextAlign(ui.Center).Wrap()
			}
		})
		rule()
		lines := ui.Column(c).FillWidth().Gap(t.Space(2))
		lines.Children(func() {
			for i, item := range items {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start)
				kit.Stagger(c, lines, row, i).Children(func() {
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						ui.Text(c, item.Product.Name).FontSize(t.Rem(0.85)).FillWidth().Wrap()
						kit.Muted(c, fmt.Sprintf(o.Each, item.Quantity, moneyText(c, item.Product.Price, currency))).FontSize(t.Rem(0.75))
					})
					money(c, item.Product.Price*float64(item.Quantity), currency).FontSize(t.Rem(0.85)).Shrink(0)
				})
			}
		})
		rule()
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1).MinWidth(0)
			money(c, total, currency).Bold().FontSize(t.Rem(1.15)).Shrink(0)
		})
	})
	return root
}
