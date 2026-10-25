package devtools

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FPSOptions word the meter; Target is the frame rate goal.
type FPSOptions struct {
	Unit, Low string
	Target    float64
}

// FPS shows the frame rate with its recent history.
func FPS(c *ui.Context, history []float64, o FPSOptions) ui.Element {
	_, hi := bounds("FPS", history)
	for _, v := range history {
		if v < 0 {
			panic(fmt.Sprintf("devtools: FPS rate %v is negative", v))
		}
	}
	if !finite(o.Target) || o.Target < 0 {
		panic(fmt.Sprintf("devtools: FPS target %v is not a positive rate", o.Target))
	}
	if o.Target == 0 {
		o.Target = 60
	}
	def(&o.Unit, "fps")
	def(&o.Low, "%d dropped")
	t := c.Theme()
	ink := func(v float64) ui.Color {
		switch {
		case v >= o.Target*0.9:
			return t.Success
		case v >= o.Target*0.5:
			return t.Warning
		}
		return t.Danger
	}
	now := history[len(history)-1]
	dropped := 0
	for _, v := range history {
		if v < o.Target*0.9 {
			dropped++
		}
	}
	// Divide before the headroom factor so it cannot overflow.
	frac := func(v float64) float32 { return float32(v / max(hi, o.Target) / 1.15) }
	root := card(c).Row().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).Role(ui.RoleStatus)
	root.Children(func() {
		ui.Column(c).Shrink(0).Children(func() {
			n := kit.Numeric(c, strconv.FormatFloat(math.Round(now), 'g', 4, 64)).FontSize(t.Rem(2)).Bold().Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
			n.TextColor(n.AnimateColor("ink", ink(now), ui.Snappy))
			kit.Muted(c, o.Unit).FontSize(t.Rem(0.75)).Bold()
			if dropped > 0 {
				kit.Chip(c, fmt.Sprintf(o.Low, dropped), t.Warning).Margin(t.Space(1), 0, 0, 0)
			}
		})
		chart := ui.Box(c).Grow(1).MinWidth(0).Height(t.Space(18))
		chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			n := len(history)
			w := r.W / float32(n)
			for i, v := range history {
				h := r.H * frac(v)
				gap := min(t.Space(0.5), w/4)
				p.Fill(ui.Rect{X: r.X + float32(i)*w + gap/2, Y: r.Y + r.H - h, W: w - gap, H: h}, ink(v).Alpha(0.75), t.Space(0.5))
			}
			ty := r.Y + r.H - r.H*frac(o.Target)
			p.Line(r.X, ty, r.X+r.W, ty, 1, t.TextMuted.Alpha(0.5))
		}))
	})
	return root
}
