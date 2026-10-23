package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UpsellOptions are the upsell labels.
type UpsellOptions struct{ Add string }

// Upsell shows related products with add; cue Pop.
func Upsell(c *ui.Context, label string, products []Product, currency string, added *int, o UpsellOptions) ui.Element {
	if label == "" || len(products) == 0 || currency == "" || added == nil || *added < -1 || *added >= len(products) {
		panic("shop: Upsell needs products and bound state")
	}
	def(&o.Add, "Add")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, product := range products {
			validProduct("Upsell", product)
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, product.Name).SingleLine().Grow(1)
				money(c, product.Price, currency)
				b := addButton(c, o.Add).Disabled(product.Stock == 0)
				b.Children(func() { ui.Text(c, o.Add).SingleLine() })
				if b.Clicked() && product.Stock > 0 {
					*added = i
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("Upsell", "added", "index", i)
				}
			})
		}
	})
	return root
}
