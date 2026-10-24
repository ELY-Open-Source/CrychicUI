package schedule

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FocusBlock is one focus interval as day fractions.
type FocusBlock struct {
	Start, End float64
}

// FocusTimeOptions word the card; Total takes "3h 30m".
type FocusTimeOptions struct{ Title, Total string }

// FocusTime shows focus blocks in a day bar.
func FocusTime(c *ui.Context, blocks []FocusBlock, o FocusTimeOptions) ui.Element {
	if len(blocks) == 0 {
		panic("schedule: FocusTime needs blocks")
	}
	def(&o.Title, "Focus time")
	def(&o.Total, "%s of focus")
	sum := 0.0
	for i, block := range blocks {
		if block.Start < 0 || block.End <= block.Start || block.End > 1 || !finite(block.Start) || !finite(block.End) {
			panic(fmt.Sprintf("schedule: FocusTime block %d is invalid", i))
		}
		sum += block.End - block.Start
	}
	if sum > 1 {
		panic("schedule: FocusTime blocks overlap past a day")
	}
	t := c.Theme()
	minutes := int(math.Round(sum * 24 * 60))
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "target", o.Title, func() {
			kit.Chip(c, fmt.Sprintf(o.Total, fmt.Sprintf("%dh %02dm", minutes/60, minutes%60)), t.Accent)
		})
		ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3)).Children(func() {
			bar := ui.Box(c).FillWidth().Height(t.Space(8)).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
			grow := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			bar.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
				for h := 1; h < 24; h++ {
					x := r.X + r.W*float32(h)/24
					p.Line(x, r.Y+r.H*0.7, x, r.Y+r.H, 1, t.Border)
				}
				for _, block := range blocks {
					x := r.X + r.W*float32(block.Start)
					w := r.W * float32(block.End-block.Start) * grow
					p.Fill(ui.Rect{X: x, Y: r.Y + t.Space(1), W: w, H: r.H - t.Space(2)}, t.Accent.Alpha(0.85), t.Space(1))
				}
			}))
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
				for _, h := range []int{0, 6, 12, 18, 24} {
					kit.Muted(c, fmt.Sprintf("%02d", h)).FontSize(t.Rem(0.6)).Shrink(0)
				}
			})
		})
	})
	return root
}
