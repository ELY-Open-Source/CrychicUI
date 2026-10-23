package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ProductPriceOptions are the saving badge labels.
type ProductPriceOptions struct{ Save string }

// ProductPrice shows price, old price, and savings.
func ProductPrice(c *ui.Context, current, old float64, currency string, o ProductPriceOptions) ui.Element {
	if current < 0 || old < current || !finite(current) || !finite(old) {
		panic("shop: ProductPrice has invalid prices")
	}
	def(&o.Save, "%.0f%% off")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).AlignItems(ui.Center)
	root.Children(func() {
		display.Price(c, current, old, currency)
		if old > current {
			badge := ui.Text(c, fmt.Sprintf(o.Save, 100*(1-current/old))).Padding(t.Space(0.5), t.Space(1.5)).
				Radius(999).Background(t.Danger.Alpha(0.12)).TextColor(t.Danger).FontSize(t.Rem(0.75)).Bold()
			kit.Arrive(badge, ui.Pose{Scale: 0.8, Fade: 1}, 0)
		}
	})
	return root
}
