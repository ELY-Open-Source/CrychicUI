package dashboard

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Status is a service health state.
type Status uint8

const (
	Up Status = iota
	Degraded
	Down
)

// UptimeDay is one day's service status.
type UptimeDay struct {
	Label  string
	Status Status
}

// UptimeOptions word the card; States name each status.
type UptimeOptions struct {
	Title, Uptime string
	States        [3]string
}

// Uptime shows day bars by status; hover names a day.
func Uptime(c *ui.Context, days []UptimeDay, o UptimeOptions) ui.Element {
	if len(days) == 0 || len(days) > 90 {
		panic("dashboard: Uptime needs 1–90 days")
	}
	up := 0
	for _, day := range days {
		if day.Label == "" || day.Status > Down {
			panic(fmt.Sprintf("dashboard: Uptime day %q needs a label and a status", day.Label))
		}
		if day.Status != Down {
			up++
		}
	}
	def(&o.Title, "Uptime")
	def(&o.Uptime, "%s uptime")
	for i, name := range [...]string{"Operational", "Degraded", "Outage"} {
		def(&o.States[i], name)
	}
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "activity", o.Title, func() {
			share := float64(up) / float64(len(days))
			kit.Chip(c, fmt.Sprintf(o.Uptime, fmt.Sprintf("%.2f%%", 100*share)), statusInk(t, Up).Mix(t.Danger, unit(share < 0.99)))
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			hover := -1
			bars := ui.Row(c).FillWidth().Height(t.Space(10)).Gap(t.Space(0.5)).AlignItems(ui.Stretch)
			sweep := bars.Keyframes("crychic.sweep", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			bars.Children(func() {
				for i, day := range days {
					c.Key(i)
					bar := ui.Box(c).Grow(1).Basis(0).MinWidth(0).FillHeight()
					if bar.Hovered() {
						hover = i
					}
					ink := bar.AnimateColor("ink", statusInk(t, day.Status).Alpha(0.75+0.25*unit(bar.Hovered())), ui.Snappy)
					grow := min(1, max(0, sweep*float32(len(days)+6)-float32(i))/6)
					bar.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						h := r.H * grow
						p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - h, W: r.W, H: h}, ink, t.Space(0.5))
					}))
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				if hover < 0 {
					kit.Muted(c, days[0].Label).FontSize(t.Rem(0.7)).Grow(1)
					kit.Muted(c, days[len(days)-1].Label).FontSize(t.Rem(0.7)).Shrink(0)
					return
				}
				day := days[hover]
				ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(statusInk(t, day.Status))
				ui.Text(c, day.Label).FontSize(t.Rem(0.75)).Bold().Shrink(0)
				kit.Muted(c, o.States[day.Status]).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
			})
		})
	})
	return root
}
