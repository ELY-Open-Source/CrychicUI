package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DayViewOptions word the day; Now overrides the clock.
type DayViewOptions struct {
	Now                time.Time
	Format, Count, One string
}

// DayView shows hours with events and a moving now line.
func DayView(c *ui.Context, day time.Time, events []Event, o DayViewOptions) ui.Element {
	validEvents("DayView", events)
	def(&o.Format, "Monday, January 2")
	def(&o.Count, "%d events")
	def(&o.One, "1 event")
	t := c.Theme()
	now := o.Now
	if now.IsZero() {
		now = c.Now()
		c.After(now.Truncate(time.Minute).Add(time.Minute).Sub(now))
	}
	now = now.In(day.Location())
	hourH := t.Space(10)
	count := len(lanesOf(events, day))
	root := card(c)
	scroll := ui.Local(root, "scroll", func() *ui.ScrollState {
		open := firstHour([]time.Time{day}, events)
		if sameDate(now, day) {
			open = max(0, float32(now.Hour()-1))
		}
		return &ui.ScrollState{Y: open * hourH}
	})
	root.Children(func() {
		kit.CardHead(c, "calendar-days", day.Format(o.Format), func() {
			text := fmt.Sprintf(o.Count, count)
			if count == 1 {
				text = o.One
			}
			kit.Chip(c, text, t.Accent)
		})
		ui.Scroll(c).TrackScroll(*scroll).FillWidth().Height(t.Space(80)).Children(func() {
			ui.Row(c).FillWidth().Padding(t.Space(1), t.Space(2), t.Space(1), 0).Children(func() {
				gutter(c, hourH)
				dayColumn(c, day, events, now, hourH, func(e Event, top, height float32) {
					b := ui.Column(c).Absolute().Left(t.Space(1)).Right(t.Space(1)).Top(top).Height(height)
					kit.Arrive(eventBox(c, e, day.Location(), b, false), ui.Pose{X: t.Space(2), Fade: 1}, 0)
				})
			})
		})
	})
	return root
}

// lanesOf lists the events that touch day.
func lanesOf(events []Event, day time.Time) []Event {
	var out []Event
	for _, lane := range lanes(events, day, 0) {
		out = append(out, lane...)
	}
	return out
}
