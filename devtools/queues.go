package devtools

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Queue is one job queue; Oldest is the longest wait.
type Queue struct {
	Name    string
	Depth   int64
	Workers int
	Oldest  time.Duration
}

// QueuesOptions word the panel; Warn is the depth that warns.
type QueuesOptions struct {
	Title, Worker, Workers, Oldest, Waiting string
	Warn                                    int64
}

// Queues lists job queues with their depth.
func Queues(c *ui.Context, queues []Queue, o QueuesOptions) ui.Element {
	def(&o.Title, "Queues")
	if len(queues) == 0 || o.Warn < 0 {
		panic(fmt.Sprintf("devtools: Queues %q needs queues and a non-negative warning depth", o.Title))
	}
	names := make([]string, len(queues))
	var total, deepest int64
	for i, q := range queues {
		if q.Depth < 0 || q.Workers < 0 || q.Oldest < 0 {
			panic(fmt.Sprintf("devtools: Queues %q needs non-negative depth, workers and wait", q.Name))
		}
		if q.Depth > math.MaxInt64-total {
			panic(fmt.Sprintf("devtools: Queues %q total depth overflows", o.Title))
		}
		total += q.Depth
		deepest = max(deepest, q.Depth)
		names[i] = q.Name
	}
	unique("Queues", names)
	def(&o.Worker, "%d worker")
	def(&o.Workers, "%d workers")
	def(&o.Oldest, "oldest %s")
	def(&o.Waiting, "%s waiting")
	if o.Warn == 0 {
		o.Warn = 1000
	}
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "inbox", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Waiting, number(c, total)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, q := range queues {
				ink := t.Accent
				if q.Depth >= o.Warn {
					ink = t.Warning
				}
				if q.Workers == 0 && q.Depth > 0 {
					ink = t.Danger
				}
				c.Key(q.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, q.Name).Font("monospace").FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							format := o.Workers
							if q.Workers == 1 {
								format = o.Worker
							}
							wink := t.TextMuted
							if q.Workers == 0 {
								wink = t.Danger
							}
							ui.Text(c, fmt.Sprintf(format, q.Workers)).FontSize(t.Rem(0.7)).TextColor(wink).Shrink(0)
							if q.Depth > 0 {
								dot(c)
								kit.Muted(c, fmt.Sprintf(o.Oldest, wait(q.Oldest))).FontSize(t.Rem(0.7)).MinWidth(0).SingleLine()
							}
						})
					})
					ui.Column(c).Width(t.Space(24)).Shrink(0).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
						n := kit.Numeric(c, number(c, q.Depth)).FontSize(t.Rem(0.95)).Bold().Shrink(0).MaxWidth(t.Space(24)).SingleLine().
							ContentTransition(ui.ContentTransition{Numeric: true})
						n.TextColor(n.AnimateColor("ink", ink, ui.Snappy))
						share := float32(0)
						if deepest > 0 {
							share = float32(float64(q.Depth) / float64(deepest))
						}
						meter(c, share, ink).Height(t.Space(1.25))
					})
				})
			}
		})
	})
	return root
}

// wait formats a queue wait in fitting units.
func wait(d time.Duration) string {
	if d < time.Minute {
		return took(d)
	}
	return span(int64(d / time.Second))
}
