package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DayHours is one weekday, Sunday first; empty Open means shut.
type DayHours struct {
	Day, Open, Close string
}

// StoreHoursOptions word the card; Opens takes day and time.
type StoreHoursOptions struct{ Title, Open, Closed, Shut, Closes, Opens string }

// StoreHours shows open now, the next change and the week.
func StoreHours(c *ui.Context, now time.Time, hours []DayHours, o StoreHoursOptions) ui.Element {
	if len(hours) != 7 || now.IsZero() {
		panic(fmt.Sprintf("shop: StoreHours needs seven days and the time, got %d days", len(hours)))
	}
	for _, day := range hours {
		if day.Day == "" || (day.Open == "") != (day.Close == "") {
			panic(fmt.Sprintf("shop: StoreHours day %q needs a name and both times or neither", day.Day))
		}
		if day.Open != "" {
			clockMinutes(day.Open)
			clockMinutes(day.Close)
		}
	}
	for i, day := range hours {
		next := hours[(i+1)%7]
		if day.Open != "" && next.Open != "" && clockMinutes(day.Open) > clockMinutes(day.Close) &&
			clockMinutes(next.Open) <= clockMinutes(day.Close) {
			panic(fmt.Sprintf("shop: StoreHours %q runs past %q opening", day.Day, next.Day))
		}
	}
	def(&o.Title, "Store hours")
	def(&o.Open, "Open now")
	def(&o.Closed, "Closed now")
	def(&o.Shut, "Closed")
	def(&o.Closes, "Closes at %s")
	def(&o.Opens, "Opens %s at %s")
	t := c.Theme()
	open, next := openAt(hours, now)
	today := int(now.Weekday())
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		status, ink := o.Closed, t.Danger
		if open {
			status, ink = o.Open, t.Success
		}
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(t.Surface).TextColor(t.TextMuted).
				Children(func() { ui.Icon(c, icon.Must("store")).Size(18, 18) })
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, o.Title).Bold().FillWidth().SingleLine()
				if next != "" {
					line := fmt.Sprintf(o.Closes, next)
					if !open {
						line = fmt.Sprintf(o.Opens, hours[nextDay(hours, now)].Day, next)
					}
					kit.Muted(c, line).FontSize(t.Rem(0.75))
				}
			})
			pill := ui.Row(c).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2.5)).Radius(999)
			pill.Background(pill.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).TextColor(ink)
			kit.Arrive(pill, ui.Pose{Scale: 0.8, Fade: 1}, 150*time.Millisecond).Children(func() {
				ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
				ui.Text(c, status).FontSize(t.Rem(0.8)).Bold().Shrink(0)
			})
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5))
		list.Children(func() {
			for i, day := range hours {
				c.Key(i)
				on := i == today
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2)).Radius(t.Space(1.5)).
					Background(t.Accent.Alpha(0.08 * unit(on)))
				kit.Stagger(c, list, row, i).Children(func() {
					name := ui.Text(c, day.Day).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0)
					text := day.Open + "–" + day.Close
					if day.Open == "" {
						text = o.Shut
					}
					value := kit.Numeric(c, text).FontSize(t.Rem(0.85)).Shrink(0)
					switch {
					case on:
						name.Bold().TextColor(t.Accent)
						value.Bold()
					case day.Open == "":
						name.TextColor(t.TextMuted)
						value.TextColor(t.TextMuted)
					}
				})
			}
		})
	})
	return root
}

func clockMinutes(value string) int {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		panic(fmt.Sprintf("shop: StoreHours needs HH:MM times, got %q", value))
	}
	return parsed.Hour()*60 + parsed.Minute()
}

// openAt reports whether it is open and the next change.
func openAt(hours []DayHours, now time.Time) (bool, string) {
	minute := now.Hour()*60 + now.Minute()
	wd := int(now.Weekday())
	today, previous := hours[wd], hours[(wd+6)%7]
	if previous.Open != "" {
		if start, end := clockMinutes(previous.Open), clockMinutes(previous.Close); start > end && minute < end {
			return true, previous.Close
		}
	}
	if today.Open != "" {
		start, end := clockMinutes(today.Open), clockMinutes(today.Close)
		if minute >= start && (start > end || minute < end) {
			return true, today.Close
		}
	}
	if d := nextDay(hours, now); d >= 0 {
		return false, hours[d].Open
	}
	return false, ""
}

// nextDay is the next opening weekday, -1 if never open.
func nextDay(hours []DayHours, now time.Time) int {
	wd := int(now.Weekday())
	if today := hours[wd]; today.Open != "" && now.Hour()*60+now.Minute() < clockMinutes(today.Open) {
		return wd
	}
	for d := 1; d <= 7; d++ {
		if hours[(wd+d)%7].Open != "" {
			return (wd + d) % 7
		}
	}
	return -1
}
