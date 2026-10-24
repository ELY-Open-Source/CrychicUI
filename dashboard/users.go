package dashboard

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// UsersOptions word the card.
type UsersOptions struct{ Title string }

// Users charts active users as bars; hover shows a day.
func Users(c *ui.Context, points []Point, o UsersOptions) ui.Element {
	if len(points) < 2 || len(points) > 31 {
		panic("dashboard: Users needs 2–31 points")
	}
	peak := 0.0
	for _, point := range points {
		if point.Label == "" || point.Value < 0 || !finite(point.Value) {
			panic(fmt.Sprintf("dashboard: Users point %q needs a label and a count of zero or more", point.Label))
		}
		peak = max(peak, point.Value)
	}
	if peak == 0 {
		peak = 1
	}
	def(&o.Title, "Active users")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "users", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			hover := -1
			chart := ui.Row(c).FillWidth().Height(t.Space(28)).Gap(t.Space(1)).AlignItems(ui.Stretch)
			grow := chart.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			chart.Children(func() {
				for i, point := range points {
					c.Key(i)
					col := ui.Box(c).Grow(1).Basis(0).MinWidth(0).FillHeight()
					on := col.Hovered()
					if on {
						hover = i
					}
					h := col.AnimateTo("h", float32(point.Value/peak), ui.Smooth)
					ink := col.AnimateColor("ink", t.Accent.Alpha(0.35+0.65*unit(on || i == len(points)-1)), ui.Snappy)
					col.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						bh := max(t.Space(0.5), r.H*h*grow)
						p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - bh, W: r.W, H: bh}, ink, t.Space(0.75))
					}))
				}
			})
			ui.Row(c).FillWidth().Children(func() {
				kit.Muted(c, points[0].Label).FontSize(t.Rem(0.7)).Grow(1)
				kit.Muted(c, points[len(points)-1].Label).FontSize(t.Rem(0.7)).Shrink(0)
			})
			shown := len(points) - 1
			if hover >= 0 {
				shown = hover
			}
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), 0, 0, 0).
				BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
				kit.Numeric(c, number(c, points[shown].Value)).Bold().FontSize(t.Rem(1.5)).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				kit.Muted(c, points[shown].Label).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine().
					ContentTransition(ui.ContentTransition{})
				if prev := points[max(shown-1, 0)].Value; shown > 0 && prev > 0 {
					change := 100 * ((points[shown].Value - prev) / prev)
					if !finite(change) {
						panic("dashboard: Users change overflows")
					}
					deltaChip(c, fmt.Sprintf("%+.1f%%", change))
				}
			})
		})
	})
	return root
}
