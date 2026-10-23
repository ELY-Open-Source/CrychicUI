package shop

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DayHours is one weekday's opening hours.
type DayHours struct {
	Day, Open, Close string
}

// StoreHoursOptions are the store status labels.
type StoreHoursOptions struct{ Open, Closed string }

// StoreHours shows open-now status and weekly hours.
func StoreHours(c *ui.Context, now time.Time, hours []DayHours, o StoreHoursOptions) ui.Element {
	if len(hours) != 7 || now.IsZero() {
		panic("shop: StoreHours needs seven days and current time")
	}
	def(&o.Open, "Open now")
	def(&o.Closed, "Closed now")
	for _, day := range hours {
		if day.Day == "" || (day.Open == "") != (day.Close == "") {
			panic("shop: StoreHours has invalid day hours")
		}
		if day.Open != "" {
			clockMinutes(day.Open)
			clockMinutes(day.Close)
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		status, ink := o.Closed, t.Danger
		if openAt(hours, now) {
			status, ink = o.Open, t.Success
		}
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
			ui.Text(c, status).TextColor(ink).Bold()
		})
		for i, day := range hours {
			c.Key(i)
			ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, day.Day).Grow(1)
				text := day.Open + "–" + day.Close
				if day.Open == "" {
					text = o.Closed
				}
				kit.Muted(c, text).FontSize(t.Rem(0.8))
			})
		}
	})
	return root
}

func clockMinutes(value string) int {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		panic("shop: StoreHours invalid clock time")
	}
	return parsed.Hour()*60 + parsed.Minute()
}

func openAt(hours []DayHours, now time.Time) bool {
	minute := now.Hour()*60 + now.Minute()
	today := hours[int(now.Weekday())]
	if today.Open != "" {
		start, end := clockMinutes(today.Open), clockMinutes(today.Close)
		if start <= end && minute >= start && minute < end || start > end && minute >= start {
			return true
		}
	}
	previous := hours[(int(now.Weekday())+6)%7]
	if previous.Open != "" {
		start, end := clockMinutes(previous.Open), clockMinutes(previous.Close)
		if start > end && minute < end {
			return true
		}
	}
	return false
}
