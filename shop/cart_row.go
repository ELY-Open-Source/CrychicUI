package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CartRowOptions word the row's controls; Each takes the unit price.
type CartRowOptions struct{ Decrease, Increase, Remove, Each string }

// CartRow edits a line's quantity or collapses out; cue Delete.
func CartRow(c *ui.Context, item *CartItem, currency string, removed *bool, o CartRowOptions) ui.Element {
	if item == nil || removed == nil {
		panic("shop: CartRow needs an item and state")
	}
	validProduct("CartRow", item.Product)
	if item.Product.Image == nil || item.Quantity <= 0 || item.Quantity > item.Product.Stock {
		panic(fmt.Sprintf("shop: CartRow %q needs an image and a quantity within stock", item.Product.Name))
	}
	def(&o.Decrease, "Decrease")
	def(&o.Increase, "Increase")
	def(&o.Remove, "Remove")
	def(&o.Each, "%s each")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2.5)).Role(ui.RoleGroup).Label(item.Product.Name).
		Transition(ui.ElementTransition{Position: true, Exit: &ui.Motion{Collapse: true}})
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Row().AlignItems(ui.Center).Children(func() {
		ui.Image(c, item.Product.Image).Size(t.Space(14), t.Space(14)).Shrink(0).Fit(ui.Cover).Radius(t.Radius)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			if item.Product.Brand != "" {
				kit.Muted(c, item.Product.Brand).FontSize(t.Rem(0.7)).Bold()
			}
			ui.Text(c, item.Product.Name).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
			kit.Muted(c, fmt.Sprintf(o.Each, moneyText(c, item.Product.Price, currency))).FontSize(t.Rem(0.75))
		})
		dec, inc := stepper(c, item.Quantity, item.Product.Stock, o.Decrease, o.Increase)
		delta := 0
		switch {
		case dec.Clicked():
			delta = -1
		case inc.Clicked():
			delta = 1
		}
		if delta != 0 {
			item.Quantity += delta
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("CartRow", "quantity", "value", item.Quantity)
		}
		money(c, item.Product.Price*float64(item.Quantity), currency).Bold().MinWidth(t.Space(20)).TextAlign(ui.End).
			ContentTransition(ui.ContentTransition{Numeric: true})
		remove := kit.Press(ui.ButtonBase(c).Label(o.Remove).Tooltip(o.Remove).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999))
		remove.Background(remove.AnimateColor("bg", t.Danger.Alpha(0.12*unit(remove.Hovered())), ui.Snappy)).
			TextColor(remove.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(remove.Hovered())), ui.Snappy)).
			Children(func() { ui.Icon(c, icon.Must("trash-2")).Size(16, 16) })
		if remove.Clicked() {
			*removed = true
			root.SetChanged()
			sound.Play(sound.Delete)
			kit.Log("CartRow", "removed")
		}
	})
	return root
}
