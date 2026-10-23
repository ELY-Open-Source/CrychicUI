package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LoyaltyOptions are the loyalty labels.
type LoyaltyOptions struct{ Next string }

// Loyalty shows points, tier, and next-tier progress.
func Loyalty(c *ui.Context, points int, tier string, next int, o LoyaltyOptions) ui.Element {
	if points < 0 || tier == "" || next < 0 {
		panic("shop: Loyalty has invalid points")
	}
	def(&o.Next, "%d points to the next tier")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			kit.Numeric(c, fmt.Sprint(points)).Bold().FontSize(t.Rem(1.4))
			ui.Text(c, tier).Grow(1)
			kit.Muted(c, fmt.Sprintf(o.Next, next)).FontSize(t.Rem(0.75))
		})
		progress := float64(1)
		if next > 0 {
			progress = float64(points) / (float64(points) + float64(next))
		}
		bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
		value := bar.AnimateTo("points", float32(progress), ui.Smooth)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * value, H: r.H}, t.Accent, 999)
		}))
	})
	return root
}
