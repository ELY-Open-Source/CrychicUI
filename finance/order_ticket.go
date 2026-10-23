package finance

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// OrderTicketOptions are the order ticket labels.
type OrderTicketOptions struct{ Quantity, Decrease, Increase, Buy, Sell string }

// OrderTicket submits buy and sell orders; cue Success.
func OrderTicket(c *ui.Context, quantity *int, price float64, currency string, o OrderTicketOptions) (root, buy, sell ui.Element) {
	if quantity == nil || *quantity < 1 || price <= 0 || !finite(price) || currency == "" {
		panic("finance: OrderTicket has invalid data")
	}
	def(&o.Quantity, "Quantity")
	def(&o.Decrease, "Decrease")
	def(&o.Increase, "Increase")
	def(&o.Buy, "Buy")
	def(&o.Sell, "Sell")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Quantity).Grow(1)
			for _, action := range []struct {
				name, glyph string
				delta       int
			}{{o.Decrease, "minus", -1}, {o.Increase, "plus", 1}} {
				b := kit.Press(ui.ButtonBase(c).Label(action.name).Size(t.Space(6), t.Space(6)).Center().Radius(999).
					Disabled(action.delta < 0 && *quantity == 1 || action.delta > 0 && *quantity == math.MaxInt))
				b.Children(func() { ui.Icon(c, icon.Must(action.glyph)).Size(16, 16) })
				if b.Clicked() && (action.delta < 0 && *quantity > 1 || action.delta > 0 && *quantity < math.MaxInt) {
					*quantity += action.delta
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("OrderTicket", "quantity", "value", *quantity)
				}
			}
			kit.Numeric(c, fmt.Sprint(*quantity))
			money(c, price, currency)
		})
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			buy = kit.Press(ui.ButtonBase(c).Label(o.Buy).Grow(1).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
				Background(t.Success).TextColor(t.AccentText))
			buy.Children(func() { ui.Text(c, o.Buy).SingleLine() })
			sell = kit.Press(ui.ButtonBase(c).Label(o.Sell).Grow(1).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
				Background(t.Danger).TextColor(t.AccentText))
			sell.Children(func() { ui.Text(c, o.Sell).SingleLine() })
		})
	})
	if buy.Clicked() {
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("OrderTicket", "buy", "quantity", *quantity)
	}
	if sell.Clicked() {
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("OrderTicket", "sell", "quantity", *quantity)
	}
	return root, buy, sell
}
