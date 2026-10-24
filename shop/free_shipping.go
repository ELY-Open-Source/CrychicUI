package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FreeShippingOptions word the bar; Remaining takes the amount left.
type FreeShippingOptions struct{ Remaining, Free string }

// FreeShipping fills a bar toward free shipping, green when reached.
func FreeShipping(c *ui.Context, amount, threshold float64, currency string, o FreeShippingOptions) ui.Element {
	if amount < 0 || threshold <= 0 || !finite(amount) || !finite(threshold) || currency == "" {
		panic(fmt.Sprintf("shop: FreeShipping needs amounts and a positive threshold, got %v of %v", amount, threshold))
	}
	def(&o.Remaining, "%s away from free shipping")
	def(&o.Free, "Free shipping unlocked")
	t := c.Theme()
	done := amount >= threshold
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).
		Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	ink := root.AnimateColor("ink", t.Accent.Mix(t.Success, unit(done)), ui.Snappy)
	root.Children(func() {
		tile := ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink)
		glyph := "truck"
		if done {
			glyph = "package-check"
		}
		tile.Children(func() { ui.Icon(c, icon.Must(glyph)).Size(18, 18).ContentTransition(ui.ContentTransition{}) })
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(2)).Children(func() {
			text := o.Free
			if !done {
				text = fmt.Sprintf(o.Remaining, moneyText(c, threshold-amount, currency))
			}
			label := ui.Text(c, text).FontSize(t.Rem(0.85)).FillWidth().Wrap().ContentTransition(ui.ContentTransition{})
			if done {
				label.Bold().TextColor(ink)
			}
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				bar := ui.Box(c).Grow(1).Height(t.Space(1.5)).Radius(999).Background(t.Surface)
				share := bar.AnimateTo("share", float32(min(1, amount/threshold)), ui.Smooth)
				bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					if share > 0 {
						p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*share), H: r.H}, ink, 999)
					}
				}))
				kit.Muted(c, moneyText(c, threshold, currency)).FontSize(t.Rem(0.7)).Shrink(0)
			})
		})
	})
	return root
}
