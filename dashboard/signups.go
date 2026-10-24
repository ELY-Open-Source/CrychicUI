package dashboard

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SignupDay is one day's signup count.
type SignupDay struct {
	Label string
	Count int
}

// SignupsOptions word the card; Total takes the sum.
type SignupsOptions struct{ Title, Total string }

// Signups charts signups by day; hover shows a count.
func Signups(c *ui.Context, days []SignupDay, o SignupsOptions) ui.Element {
	if len(days) == 0 || len(days) > 14 {
		panic("dashboard: Signups needs 1–14 days")
	}
	peak, total := 1, 0
	for _, day := range days {
		if day.Label == "" || day.Count < 0 || total > math.MaxInt-day.Count {
			panic(fmt.Sprintf("dashboard: Signups day %q needs a label and a count that sums", day.Label))
		}
		peak = max(peak, day.Count)
		total += day.Count
	}
	def(&o.Title, "Signups")
	def(&o.Total, "%s total")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "user-plus", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Total, number(c, total)), t.Accent) })
		ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3)).Children(func() {
			hover := -1
			chart := ui.Row(c).FillWidth().Height(t.Space(26)).Gap(t.Space(1.5)).AlignItems(ui.Stretch)
			grow := chart.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			chart.Children(func() {
				for i, day := range days {
					c.Key(i)
					col := ui.Box(c).Grow(1).Basis(0).MinWidth(0).FillHeight()
					on := col.Hovered()
					if on {
						hover = i
					}
					h := col.AnimateTo("h", float32(day.Count)/float32(peak), ui.Smooth)
					ink := col.AnimateColor("ink", t.Accent.Alpha(0.45+0.55*unit(on)), ui.Snappy)
					col.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						bh := max(t.Space(0.5), r.H*h*grow)
						p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - bh, W: r.W, H: bh}, ink, t.Space(0.75))
					}))
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				for i, day := range days {
					c.Key(i)
					kit.Muted(c, day.Label).FontSize(t.Rem(0.65)).Grow(1).Basis(0).MinWidth(0).TextAlign(ui.Center).SingleLine()
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), 0, 0, 0).Children(func() {
				shown := len(days) - 1
				if hover >= 0 {
					shown = hover
				}
				kit.Numeric(c, number(c, days[shown].Count)).Bold().FontSize(t.Rem(1.3)).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				kit.Muted(c, days[shown].Label).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
			})
		})
	})
	return root
}
