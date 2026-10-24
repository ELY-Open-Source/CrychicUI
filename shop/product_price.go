package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ProductPriceOptions word savings: Save the percent, Saving the amount.
type ProductPriceOptions struct{ Save, Saving string }

// ProductPrice shows a price, the struck old price, the saving.
func ProductPrice(c *ui.Context, current, old float64, currency string, o ProductPriceOptions) ui.Element {
	if current < 0 || old != 0 && old < current || !finite(current) || !finite(old) {
		panic(fmt.Sprintf("shop: ProductPrice needs old 0 or at least current, got %v and %v", current, old))
	}
	def(&o.Save, "%d%% off")
	def(&o.Saving, "You save %s")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1)).Role(ui.RoleGroup).Label(moneyText(c, current, currency))
	root.Children(func() {
		ui.Row(c).Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			money(c, current, currency).FontSize(t.Rem(1.6)).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
			if old > current {
				money(c, old, currency).FontSize(t.Rem(0.9)).TextColor(t.TextMuted).Strikethrough()
				if off := percentOff(current, old); off > 0 {
					badge := ui.Row(c).Shrink(0).Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.Danger.Alpha(0.12))
					kit.Arrive(badge, ui.Pose{Scale: 0.8, Fade: 1}, 120*time.Millisecond).Children(func() {
						ui.Text(c, fmt.Sprintf(o.Save, off)).FontSize(t.Rem(0.75)).Bold().TextColor(t.Danger).Shrink(0)
					})
				}
			}
		})
		if old > current {
			ui.Text(c, fmt.Sprintf(o.Saving, moneyText(c, old-current, currency))).FontSize(t.Rem(0.8)).TextColor(t.Success)
		}
	})
	return root
}
