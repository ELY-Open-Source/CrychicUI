package dashboard

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ErrorRateOptions word the card; History holds recent rates, oldest first.
type ErrorRateOptions struct {
	Title, Rate, Limit, Over, Within string
	History                          []float64
}

// ErrorRate shows errors per minute against a limit line.
func ErrorRate(c *ui.Context, rate, limit float64, o ErrorRateOptions) ui.Element {
	if rate < 0 || limit <= 0 || !finite(rate) || !finite(limit) || len(o.History) > 60 {
		panic("dashboard: ErrorRate needs a rate of zero or more, a positive limit and at most 60 history points")
	}
	peak := max(rate, limit)
	for _, h := range o.History {
		if h < 0 || !finite(h) {
			panic("dashboard: ErrorRate history must be finite and not negative")
		}
		peak = max(peak, h)
	}
	def(&o.Title, "Error rate")
	def(&o.Rate, "errors per minute")
	def(&o.Limit, "Limit %.0f")
	def(&o.Over, "Over limit")
	def(&o.Within, "Within limit")
	over := rate > limit
	t := c.Theme()
	ink := t.Success
	if over {
		ink = t.Danger
	}
	history := o.History
	if len(history) == 0 {
		history = []float64{rate}
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "bug", o.Title, func() {
			if over {
				kit.Chip(c, o.Over, t.Danger)
			} else {
				kit.Chip(c, o.Within, t.Success)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
				kit.Numeric(c, fmt.Sprintf("%.1f", rate)).Bold().FontSize(t.Rem(1.8)).TextColor(ink).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				kit.Muted(c, o.Rate).FontSize(t.Rem(0.8)).Padding(0, 0, t.Space(1), 0).Grow(1).MinWidth(0).SingleLine()
				kit.Muted(c, fmt.Sprintf(o.Limit, limit)).FontSize(t.Rem(0.75)).Padding(0, 0, t.Space(1), 0).Shrink(0)
			})
			chart := ui.Box(c).FillWidth().Height(t.Space(20))
			grow := chart.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			hs := make([]float32, len(history))
			for i, h := range history {
				hs[i] = chart.AnimateTo(fmt.Sprint("h", i), float32(h/peak), ui.Smooth)
			}
			chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				slot := r.W / float32(len(hs))
				for i, h := range hs {
					bh := max(t.Space(0.5), r.H*h*grow)
					bar := t.Accent.Alpha(0.5)
					if history[i] > limit {
						bar = t.Danger
					}
					p.Fill(ui.Rect{X: r.X + slot*float32(i) + slot*0.15, Y: r.Y + r.H - bh, W: slot * 0.7, H: bh}, bar, t.Space(0.5))
				}
				y := r.Y + r.H - (r.H-1)*float32(limit/peak)
				dash := t.Space(1.5)
				for x := r.X; x < r.X+r.W; x += 2 * dash {
					p.Fill(ui.Rect{X: x, Y: y - 0.5, W: min(dash, r.X+r.W-x), H: 1}, t.Text, 0)
				}
			}))
		})
	})
	return root
}
