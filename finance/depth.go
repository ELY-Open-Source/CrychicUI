package finance

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DepthOptions word the chart.
type DepthOptions struct{ Title, Bids, Asks string }

// Depth charts cumulative bids and asks as meeting steps.
func Depth(c *ui.Context, bids, asks []BookLevel, currency string, o DepthOptions) ui.Element {
	if currency == "" {
		panic("finance: Depth needs a currency")
	}
	bids, asks = sides("Depth", bids, asks)
	// Asks ascending, so both sides walk outward.
	out := make([]BookLevel, len(asks))
	for i, l := range asks {
		out[len(asks)-1-i] = l
	}
	cum := func(levels []BookLevel) []float64 {
		sums, total := make([]float64, len(levels)), 0.0
		for i, l := range levels {
			total += l.Amount
			if !finite(total) {
				panic("finance: Depth total overflows")
			}
			sums[i] = total
		}
		return sums
	}
	bidSum, askSum := cum(bids), cum(out)
	lo, hi := bids[len(bids)-1].Price, out[len(out)-1].Price
	peak := max(bidSum[len(bidSum)-1], askSum[len(askSum)-1])
	def(&o.Title, "Market depth")
	def(&o.Bids, "Bids")
	def(&o.Asks, "Asks")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "chart-line", o.Title, func() {
			kit.Chip(c, o.Bids, t.Success)
			kit.Chip(c, o.Asks, t.Danger)
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			chart := ui.Box(c).FillWidth().Height(t.Space(36))
			grow := chart.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				x := func(price float64) float32 { return r.X + r.W*float32((price-lo)/(hi-lo)) }
				y := func(sum float64) float32 { return r.Y + r.H - r.H*0.92*float32(sum/peak)*grow }
				side := func(levels []BookLevel, sums []float64, ink ui.Color, end float32) {
					var area ui.Path
					area.MoveTo(x(levels[0].Price), r.Y+r.H)
					for i, l := range levels {
						if i > 0 {
							area.LineTo(x(l.Price), y(sums[i-1]))
						}
						area.LineTo(x(l.Price), y(sums[i]))
					}
					area.LineTo(end, y(sums[len(sums)-1])).LineTo(end, r.Y+r.H).Close()
					p.FillPath(&area, ink.Alpha(0.18))
					p.StrokePath(&area, t.Space(0.4), ink)
				}
				side(bids, bidSum, t.Success, r.X)
				side(out, askSum, t.Danger, r.X+r.W)
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, t.Border, 0)
			}))
			ui.Row(c).FillWidth().Children(func() {
				kit.Muted(c, moneyText(c, lo, currency)).FontSize(t.Rem(0.7)).Grow(1).Basis(0)
				ui.Text(c, moneyText(c, lo/2+hi/2, currency)).FontSize(t.Rem(0.75)).Bold().Grow(1).Basis(0).TextAlign(ui.Center)
				kit.Muted(c, moneyText(c, hi, currency)).FontSize(t.Rem(0.7)).Grow(1).Basis(0).TextAlign(ui.End)
			})
		})
	})
	return root
}
