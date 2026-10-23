package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// OrderStage is one order tracking stage.
type OrderStage uint8

const (
	Placed OrderStage = iota
	Packed
	Shipped
	Delivered
)

// OrderStatusOptions are the order stage labels.
type OrderStatusOptions struct{ Placed, Packed, Shipped, Delivered string }

// OrderStatus shows an order tracking track.
func OrderStatus(c *ui.Context, current OrderStage, o OrderStatusOptions) ui.Element {
	if current > Delivered {
		panic("shop: OrderStatus needs a valid stage")
	}
	def(&o.Placed, "Placed")
	def(&o.Packed, "Packed")
	def(&o.Shipped, "Shipped")
	def(&o.Delivered, "Delivered")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
		progress := bar.AnimateTo("stage", float32(current)/3, ui.Smooth)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * progress, H: r.H}, t.Success, 999)
		}))
		ui.Row(c).FillWidth().Children(func() {
			for _, stage := range []struct {
				value OrderStage
				name  string
			}{{Placed, o.Placed}, {Packed, o.Packed}, {Shipped, o.Shipped}, {Delivered, o.Delivered}} {
				ui.Text(c, stage.name).FontSize(t.Rem(0.75)).Grow(1).Basis(0).
					TextColor(t.TextMuted.Mix(t.Success, unit(stage.value <= current)))
			}
		})
	})
	return root
}
