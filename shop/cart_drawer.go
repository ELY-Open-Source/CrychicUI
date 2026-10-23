package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CartDrawerOptions are the cart labels.
type CartDrawerOptions struct{ Empty, Subtotal, Checkout string }

// CartDrawer shows items, subtotal, and checkout; cue Open.
func CartDrawer(c *ui.Context, items []CartItem, currency string, checkout *bool, o CartDrawerOptions) (root, send ui.Element) {
	if checkout == nil || currency == "" {
		panic("shop: CartDrawer needs currency and checkout state")
	}
	def(&o.Empty, "Your cart is empty")
	def(&o.Subtotal, "Subtotal")
	def(&o.Checkout, "Checkout")
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(72)).Gap(t.Space(2)).Padding(t.Space(3)).
		Radius(t.Radius).Border(1, t.Border).Background(t.Surface).Role(ui.RoleDialog)
	total := 0.0
	root.Children(func() {
		if len(items) == 0 {
			kit.Muted(c, o.Empty)
		}
		for i, item := range items {
			validProduct("CartDrawer", item.Product)
			if item.Quantity <= 0 || item.Quantity > item.Product.Stock {
				panic("shop: CartDrawer has an invalid quantity")
			}
			c.Key(i)
			ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
				ui.Text(c, item.Product.Name).SingleLine().Grow(1)
				kit.Numeric(c, "×"+fmt.Sprint(item.Quantity))
				money(c, item.Product.Price*float64(item.Quantity), currency)
			})
			total += item.Product.Price * float64(item.Quantity)
		}
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Subtotal).Bold().Grow(1)
			money(c, total, currency).Bold()
		})
		send = addButton(c, o.Checkout).Disabled(len(items) == 0)
		send.Children(func() { ui.Text(c, o.Checkout).SingleLine() })
	})
	if send.Clicked() && len(items) > 0 {
		*checkout = true
		root.SetChanged()
		sound.Play(sound.Open)
		kit.Log("CartDrawer", "checkout", "items", len(items))
	}
	return root, send
}
