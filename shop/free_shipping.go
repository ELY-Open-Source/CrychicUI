package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FreeShippingOptions are the shipping labels.
type FreeShippingOptions struct{ Remaining, Free string }

// FreeShipping shows progress toward free shipping.
func FreeShipping(c *ui.Context, amount, threshold float64, currency string, o FreeShippingOptions) ui.Element {
	if amount < 0 || threshold <= 0 || !finite(amount) || !finite(threshold) || currency == "" {
		panic("shop: FreeShipping has invalid amounts")
	}
	def(&o.Remaining, "%s away from free shipping")
	def(&o.Free, "Free shipping unlocked")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		text := o.Free
		if amount < threshold {
			text = fmt.Sprintf(o.Remaining, moneyText(c, threshold-amount, currency))
		}
		ui.Text(c, text).FontSize(t.Rem(0.85))
		bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
		progress := bar.AnimateTo("shipping", float32(min(1, amount/threshold)), ui.Smooth)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * progress, H: r.H}, t.Success, 999)
		}))
	})
	return root
}
