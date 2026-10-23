package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StockOptions are the stock labels.
type StockOptions struct{ In, Low, Out string }

// Stock shows stock state with a colored dot.
func Stock(c *ui.Context, count int, o StockOptions) ui.Element {
	if count < 0 {
		panic("shop: Stock count must not be negative")
	}
	def(&o.In, "In stock")
	def(&o.Low, "Low stock")
	def(&o.Out, "Out of stock")
	t := c.Theme()
	text, ink := o.In, t.Success
	switch {
	case count == 0:
		text, ink = o.Out, t.Danger
	case count < 5:
		text, ink = o.Low, t.Warning
	}
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	root.Children(func() {
		ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
		ui.Text(c, text).TextColor(ink).FontSize(t.Rem(0.8))
	})
	return root
}
