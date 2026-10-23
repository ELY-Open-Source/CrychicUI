package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CartRowOptions are the cart row labels.
type CartRowOptions struct{ Decrease, Increase, Remove string }

// CartRow edits quantity and removes; cue Delete.
func CartRow(c *ui.Context, item *CartItem, currency string, removed *bool, o CartRowOptions) ui.Element {
	if item == nil || removed == nil {
		panic("shop: CartRow needs an item and state")
	}
	validProduct("CartRow", item.Product)
	if item.Product.Image == nil || item.Quantity <= 0 || item.Quantity > item.Product.Stock {
		panic("shop: CartRow has an invalid image or quantity")
	}
	def(&o.Decrease, "Decrease")
	def(&o.Increase, "Increase")
	def(&o.Remove, "Remove")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).
		Transition(ui.ElementTransition{Position: true, Exit: &ui.Motion{Collapse: true}})
	root.Children(func() {
		ui.Image(c, item.Product.Image).Size(t.Space(10), t.Space(10)).Fit(ui.Cover).Radius(t.Radius)
		ui.Text(c, item.Product.Name).SingleLine().Grow(1)
		for _, action := range []struct {
			name, glyph string
			delta       int
		}{{o.Decrease, "minus", -1}, {o.Increase, "plus", 1}} {
			b := kit.Press(ui.ButtonBase(c).Label(action.name).Tooltip(action.name).Size(t.Space(6), t.Space(6)).Center().Radius(999).
				Disabled(action.delta < 0 && item.Quantity == 1 || action.delta > 0 && item.Quantity == item.Product.Stock))
			b.Children(func() { ui.Icon(c, icon.Must(action.glyph)).Size(16, 16) })
			if b.Clicked() {
				item.Quantity += action.delta
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("CartRow", "quantity", "value", item.Quantity)
			}
		}
		kit.Numeric(c, fmt.Sprint(item.Quantity)).Width(t.Space(4))
		money(c, item.Product.Price*float64(item.Quantity), currency)
		remove := kit.Press(ui.ButtonBase(c).Label(o.Remove).Tooltip(o.Remove).Size(t.Space(6), t.Space(6)).Center().Radius(999))
		remove.Children(func() { ui.Icon(c, icon.Must("trash-2")).Size(16, 16) })
		if remove.Clicked() {
			*removed = true
			root.SetChanged()
			sound.Play(sound.Delete)
			kit.Log("CartRow", "removed")
		}
	})
	return root
}
