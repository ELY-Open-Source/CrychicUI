package chart

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Task is a named stretch of time.
type Task struct {
	Name       string
	Start, End time.Time
}

// GanttOptions tune a Gantt.
type GanttOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Date is a Go layout; "Jan 2" default.
	Date string
}

// Gantt draws task bars over time, first on top.
func Gantt(c *ui.Context, label string, tasks []Task, at *int, o GanttOptions) ui.Element {
	if len(tasks) == 0 {
		panic(fmt.Sprintf("chart: Gantt %q needs tasks", label))
	}
	names := make([]string, len(tasks))
	lo, hi := tasks[0].Start, tasks[0].End
	for i, k := range tasks {
		if !k.End.After(k.Start) {
			panic(fmt.Sprintf("chart: Gantt %q tasks must end after they start", label))
		}
		names[i] = k.Name
		lo, hi = minTime(lo, k.Start), maxTime(hi, k.End)
	}
	named("Gantt", label, names)
	chosen("Gantt", label, at, len(tasks))
	def(&o.Date, "Jan 2")
	t := c.Theme()
	x := mc.NewTime(lo, hi)
	ticks := x.Ticks(5)
	up := slices.Clone(names)
	slices.Reverse(up)
	y := mc.NewBand(up...).Padding(0.3, 0.1)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: y.Ticks(), Bottom: ticks}, o.Height, func(area ui.Element) {
			g := grow(area)
			area.Draw(func(p *ui.Painter, r ui.Rect) {
				for _, tick := range ticks {
					p.Line(mc.X(r, tick.At), r.Y, mc.X(r, tick.At), r.Y+r.H, 1, t.Border)
				}
			})
			for i, k := range tasks {
				c.Key(k.Name)
				span := x.Span(k.Start, k.End)
				span.Hi = span.Lo + (span.Hi-span.Lo)*g
				tip := fmt.Sprintf("%s: %s – %s", k.Name, k.Start.Format(o.Date), k.End.Format(o.Date))
				mark(c, root, at, i, span, y.Band(len(tasks)-1-i), tone(t, i), tip, "Gantt")
			}
		})
	})
}
