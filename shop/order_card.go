package shop

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// OrderCardOptions are the order card labels.
type OrderCardOptions struct {
	Open, Placed, Packed, Shipped, Delivered string
}

// OrderCard shows an order and opens it; cue Tap.
func OrderCard(c *ui.Context, number string, when time.Time, total float64, currency string, status OrderStage, opened *bool, o OrderCardOptions) ui.Element {
	if number == "" || total < 0 || !finite(total) || currency == "" || status > Delivered || opened == nil {
		panic("shop: OrderCard has invalid order data")
	}
	def(&o.Open, "View order")
	def(&o.Placed, "Placed")
	def(&o.Packed, "Packed")
	def(&o.Shipped, "Shipped")
	def(&o.Delivered, "Delivered")
	t := c.Theme()
	root := kit.Lift(ui.ButtonBase(c).Label(o.Open).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
		Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Background(t.Surface))
	root.Children(func() {
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, number).Bold().SingleLine()
			kit.Muted(c, when.Format("2006-01-02")).FontSize(t.Rem(0.75))
		})
		money(c, total, currency)
		OrderStatus(c, status, OrderStatusOptions{Placed: o.Placed, Packed: o.Packed, Shipped: o.Shipped, Delivered: o.Delivered})
	})
	if root.Clicked() && !*opened {
		*opened = true
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("OrderCard", "opened")
	}
	return root
}
