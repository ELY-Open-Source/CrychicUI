package project

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SprintHeaderOptions word the header; Goal is optional.
type SprintHeaderOptions struct{ Left, One, Today, Ended, Day, DateFormat, Goal string }

// SprintHeader shows sprint name, dates, and days left.
func SprintHeader(c *ui.Context, name string, start, end time.Time, o SprintHeaderOptions) ui.Element {
	if name == "" || !end.After(start) {
		panic(fmt.Sprintf("project: SprintHeader %q has invalid dates", name))
	}
	def(&o.Left, "%d days left")
	def(&o.One, "1 day left")
	def(&o.Today, "Ends today")
	def(&o.Ended, "Ended")
	def(&o.Day, "Day %d of %d")
	def(&o.DateFormat, "Jan 02")
	t := c.Theme()
	loc := end.Location()
	now := c.Now().In(loc)
	// Calendar days in the deadline's zone, counted without Sub.
	left := max(0, dayNumber(end)-dayNumber(now))
	days := max(1, dayNumber(end)-dayNumber(start.In(loc)))
	secs := end.Unix() - now.Unix()
	if secs > 0 {
		midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
		next := midnight
		if end.Before(midnight) {
			next = end
		}
		c.After(next.Sub(now) + time.Millisecond)
	}
	day := min(days, max(1, days-left+1))
	text, ink := fmt.Sprintf(o.Left, left), t.Accent
	switch {
	case secs <= 0:
		text, ink = o.Ended, t.TextMuted
	case left == 0:
		text, ink = o.Today, t.Danger
	case left == 1:
		text, ink = o.One, t.Warning
	case left <= 2:
		ink = t.Warning
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(3))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "zap", t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.95)).Bold().SingleLine().MinWidth(0)
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "calendar-range", 12)
					ui.Text(c, start.Format(o.DateFormat)+" – "+end.Format(o.DateFormat)).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
				})
			})
			kit.Chip(c, text, ink)
		})
		if o.Goal != "" {
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Start).Children(func() {
				glyph(c, "target", 14).TextColor(t.TextMuted)
				ui.Text(c, o.Goal).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
			})
		}
		ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
			progress(c, float64(day)/float64(days), ink)
			kit.Muted(c, fmt.Sprintf(o.Day, day, days)).FontSize(t.Rem(0.7))
		})
	})
	return root
}

// dayNumber numbers v's calendar date in its own zone.
func dayNumber(v time.Time) int64 {
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC).Unix() / 86400
}
