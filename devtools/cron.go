package devtools

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CronJob is one scheduled job; Failed marks its last run.
type CronJob struct {
	Name, Schedule string
	Next           time.Time
	Paused, Failed bool
}

// CronOptions word the panel; Next takes a time left.
type CronOptions struct{ Title, Next, Due, Paused, Failed string }

// Cron lists scheduled jobs counting down to their next run.
func Cron(c *ui.Context, jobs []CronJob, o CronOptions) ui.Element {
	def(&o.Title, "Scheduled jobs")
	if len(jobs) == 0 {
		panic(fmt.Sprintf("devtools: Cron %q needs jobs", o.Title))
	}
	names := make([]string, len(jobs))
	for i, j := range jobs {
		if strings.TrimSpace(j.Schedule) == "" || !j.Paused && j.Next.IsZero() {
			panic(fmt.Sprintf("devtools: Cron %q needs a schedule, and a next run unless paused", j.Name))
		}
		names[i] = j.Name
	}
	unique("Cron", names)
	def(&o.Next, "in %s")
	def(&o.Due, "Due now")
	def(&o.Paused, "Paused")
	def(&o.Failed, "Last run failed")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "calendar-clock", o.Title, func() { kit.Chip(c, fmt.Sprint(len(jobs)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, j := range jobs {
				ink := t.Accent
				switch {
				case j.Paused:
					ink = t.TextMuted
				case j.Failed:
					ink = t.Danger
				}
				c.Key(j.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					tile(c, "clock", ink, t.Space(8))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, j.Name).Font("monospace").FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							kit.Muted(c, j.Schedule).Font("monospace").FontSize(t.Rem(0.68)).Shrink(0).MaxWidth(t.Space(28)).SingleLine()
							if j.Failed {
								dot(c)
								ui.Text(c, o.Failed).FontSize(t.Rem(0.68)).TextColor(t.Danger).MinWidth(0).SingleLine()
							}
						})
					})
					if j.Paused {
						kit.Chip(c, o.Paused, t.TextMuted)
						return
					}
					left := countdown(c, j.Next)
					text, tone := fmt.Sprintf(o.Next, span(left)), t.Text
					if left == 0 {
						text, tone = o.Due, t.Warning
					}
					n := kit.Numeric(c, text).Font("monospace").FontSize(t.Rem(0.75)).Bold().Shrink(0).
						ContentTransition(ui.ContentTransition{Numeric: true})
					n.TextColor(n.AnimateColor("ink", tone, ui.Snappy))
				})
			}
		})
	})
	return root
}
