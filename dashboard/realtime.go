package dashboard

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RealtimeOptions word the card; History holds recent counts, oldest first.
type RealtimeOptions struct {
	Label, Note string
	History     []int
}

// Realtime shows a live count; digits tick, the dot pulses.
func Realtime(c *ui.Context, count int, o RealtimeOptions) ui.Element {
	if count < 0 || len(o.History) > 60 {
		panic("dashboard: Realtime needs a count of zero or more and at most 60 history points")
	}
	peak := max(1, count)
	for _, h := range o.History {
		if h < 0 {
			panic("dashboard: Realtime history must not be negative")
		}
		peak = max(peak, h)
	}
	def(&o.Label, "Visitors now")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			dot := ui.Box(c).Size(t.Space(4), t.Space(4)).Shrink(0)
			beat := float32(0.5 + 0.5*math.Sin(2*math.Pi*float64(dot.Loop("pulse", 1400*time.Millisecond, ui.Linear))))
			dot.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, t.Success.Alpha(0.25*beat), 999)
				d := r.W * 0.5
				p.Fill(ui.Rect{X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d}, t.Success, 999)
			}))
			kit.Muted(c, o.Label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
		})
		ui.Column(c).Gap(t.Space(0.5)).Children(func() {
			kit.Numeric(c, number(c, count)).Bold().FontSize(t.Rem(2.4)).Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
			if o.Note != "" {
				kit.Muted(c, o.Note).FontSize(t.Rem(0.75)).Shrink(0)
			}
		})
		if len(o.History) > 0 {
			bars := ui.Box(c).FillWidth().Height(t.Space(12))
			hs := make([]float32, len(o.History))
			for i, h := range o.History {
				hs[i] = bars.AnimateTo(fmt.Sprint("h", i), float32(h)/float32(peak), ui.Smooth)
			}
			bars.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				slot := r.W / float32(len(hs))
				for i, h := range hs {
					bh := max(t.Space(0.5), r.H*h)
					ink := t.Success.Alpha(0.3 + 0.7*unit(i == len(hs)-1))
					p.Fill(ui.Rect{X: r.X + slot*float32(i) + slot*0.15, Y: r.Y + r.H - bh, W: slot * 0.7, H: bh}, ink, t.Space(0.5))
				}
			}))
		}
	})
	return root
}
