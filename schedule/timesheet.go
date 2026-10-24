package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TimesheetOptions word the sheet; Target is daily hours (default 8).
type TimesheetOptions struct {
	Title, Total string
	Days         [7]string
	Target       float64
}

// Timesheet shows hours per day and total.
func Timesheet(c *ui.Context, hours []float64, o TimesheetOptions) ui.Element {
	if len(hours) != 7 {
		panic("schedule: Timesheet needs seven days")
	}
	def(&o.Title, "This week")
	def(&o.Total, "Total %.1f h")
	if o.Days == [7]string{} {
		o.Days = [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	}
	if o.Target == 0 {
		o.Target = 8
	}
	if o.Target < 0 || !finite(o.Target) {
		panic("schedule: Timesheet target must be positive")
	}
	total, peak := 0.0, o.Target
	for _, value := range hours {
		if value < 0 || !finite(value) {
			panic("schedule: Timesheet hours must be nonnegative")
		}
		total += value
		if !finite(total) {
			panic("schedule: Timesheet total overflow")
		}
		peak = max(peak, value)
	}
	t := c.Theme()
	barH := t.Space(28)
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "timer", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Total, total), t.Accent) })
		chart := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Padding(t.Space(3))
		chart.DrawOver(func(p *ui.Painter, r ui.Rect) {
			y := r.Y + t.Space(3) + t.Space(5) + barH*float32(1-o.Target/peak)
			for x := r.X + t.Space(3); x < r.X+r.W-t.Space(3); x += t.Space(2) {
				p.Line(x, y, min(x+t.Space(1), r.X+r.W-t.Space(3)), y, 1, t.TextMuted.Alpha(0.6))
			}
		})
		chart.Children(func() {
			for i, value := range hours {
				over := value > o.Target
				ink := t.Accent
				if over {
					ink = t.Warning
				}
				c.Key(i)
				ui.Column(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Center).Gap(t.Space(1)).Children(func() {
					kit.Numeric(c, fmt.Sprintf("%.1f", value)).FontSize(t.Rem(0.7)).Bold().Height(t.Space(4)).TextColor(ink.Mix(t.Text, 0.4))
					well := ui.Box(c).Width(t.Space(6)).Height(barH).Radius(t.Space(1)).Background(t.SurfaceHover.Alpha(0.6))
					delay := time.Duration(i) * 40 * time.Millisecond
					frames := []ui.Keyframe{{Value: 0}}
					if delay > 0 {
						frames = append(frames, ui.Keyframe{At: delay})
					}
					rise := well.Keyframes("rise", append(frames, ui.Keyframe{At: delay + 500*time.Millisecond, Value: 1, Ease: ui.EaseOut}))
					share := well.AnimateTo("share", float32(value/peak), ui.Smooth)
					well.Draw(func(p *ui.Painter, r ui.Rect) {
						h := r.H * share * rise
						p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - h, W: r.W, H: h}, ink, t.Space(1))
					})
					kit.Muted(c, o.Days[i]).FontSize(t.Rem(0.65))
				})
			}
		})
	})
	return root
}
