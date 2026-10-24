package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LoyaltyOptions word the card; Next takes the points still needed.
type LoyaltyOptions struct{ Next, NextTier, Top, Points string }

// Loyalty shows points, the tier and progress to the next.
func Loyalty(c *ui.Context, points int, tier string, next int, o LoyaltyOptions) ui.Element {
	if points < 0 || tier == "" || next < 0 {
		panic(fmt.Sprintf("shop: Loyalty tier %q needs points and next at zero or more", tier))
	}
	def(&o.Next, "%d points to the next tier")
	def(&o.Top, "Top tier reached")
	def(&o.Points, "points")
	t := c.Theme()
	share := 1.0
	if next > 0 {
		share = float64(points) / (float64(points) + float64(next))
	}
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			tile := ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).Background(t.Warning.Alpha(0.15)).TextColor(t.Warning)
			spin := tile.Keyframes("crychic.crown", []ui.Keyframe{{Value: -12}, {At: 600 * time.Millisecond, Value: 0, Ease: ui.EaseOut}})
			tile.Children(func() { ui.Icon(c, icon.Must("crown")).Size(18, 18).Rotate(spin) })
			ui.Text(c, tier).Bold().FontSize(t.Rem(1.05)).Grow(1).MinWidth(0).SingleLine()
			ui.Column(c).Shrink(0).AlignItems(ui.End).Children(func() {
				kit.Numeric(c, fmt.Sprint(points)).Bold().FontSize(t.Rem(1.4)).ContentTransition(ui.ContentTransition{Numeric: true})
				kit.Muted(c, o.Points).FontSize(t.Rem(0.7))
			})
		})
		bar := ui.Box(c).FillWidth().Height(t.Space(2)).Radius(999).Background(t.Surface)
		fill := bar.AnimateTo("share", float32(share), ui.Smooth)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			if fill > 0 {
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*fill), H: r.H}, t.Warning, 999)
			}
		}))
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			text := o.Top
			if next > 0 {
				text = fmt.Sprintf(o.Next, next)
			}
			kit.Muted(c, text).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).Wrap()
			if o.NextTier != "" && next > 0 {
				ui.Text(c, o.NextTier).FontSize(t.Rem(0.75)).Bold().Shrink(0)
			}
		})
	})
	return root
}
