package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MiniCartOptions are the mini cart labels.
type MiniCartOptions struct{ Cart, Empty string }

// MiniCart shows a cart popover; cue Open.
func MiniCart(c *ui.Context, items []CartItem, currency string, open *bool, o MiniCartOptions) ui.Element {
	if open == nil || currency == "" {
		panic("shop: MiniCart needs currency and state")
	}
	def(&o.Cart, "Cart")
	def(&o.Empty, "Your cart is empty")
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	root.Children(func() {
		button := kit.Press(ui.ButtonBase(c).Label(o.Cart).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
			Border(1, t.Border).Checked(*open))
		button.Children(func() { ui.Text(c, o.Cart).SingleLine() })
		if button.Clicked() {
			*open = !*open
			root.SetChanged()
			if *open {
				sound.Play(sound.Open)
			} else {
				sound.Play(sound.Close)
			}
			kit.Log("MiniCart", "toggled", "open", *open)
		}
		if ui.PopoverBase(c, button, open, func(panel ui.Element) {
			panel.Width(t.Space(56)).Padding(t.Space(3)).Gap(t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Background(t.Surface).Shadow(0, t.Space(2), t.Space(6), 0, kit.Shade(0.15)).Transition(kit.Enter)
			panel.Children(func() {
				if len(items) == 0 {
					kit.Muted(c, o.Empty)
				}
				for i, item := range items {
					validProduct("MiniCart", item.Product)
					if item.Quantity <= 0 || item.Quantity > item.Product.Stock {
						panic("shop: MiniCart has an invalid quantity")
					}
					c.Key(i)
					ui.Row(c).FillWidth().Children(func() {
						ui.Text(c, item.Product.Name).SingleLine().Grow(1)
						money(c, item.Product.Price*float64(item.Quantity), currency)
					})
				}
			})
		}).Dismissed() {
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("MiniCart", "dismissed")
		}
	})
	return root
}
