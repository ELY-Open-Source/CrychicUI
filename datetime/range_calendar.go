package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// RangeCalendarOptions tune a RangeCalendar.
type RangeCalendarOptions struct {
	// Previous and Next name the month buttons.
	Previous, Next string
}

// RangeCalendar picks start then end, previewing the range.
func RangeCalendar(c *ui.Context, start, end *time.Time, o RangeCalendarOptions) ui.Element {
	if start == nil || end == nil || dayBefore(*end, *start) {
		panic("datetime: RangeCalendar needs a start on or before its end")
	}
	def(&o.Previous, "Previous month")
	def(&o.Next, "Next month")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2)).AlignSelf(ui.Start)
	cursor := ui.Local(root, "cursor", func() time.Time { return *start })
	picking := ui.Local(root, "picking", func() bool { return false })
	hover := ui.Local(root, "hover", func() time.Time { return time.Time{} })
	return root.Children(func() {
		m := monthOf(*cursor)
		back, ahead := head(c, i18n.FormatDate(locale(c), m, i18n.DateOptions{Skeleton: "yMMMM"}), o.Previous, o.Next)
		if back.Clicked() {
			*cursor = addMonths(m, -1)
		}
		if ahead.Clicked() {
			*cursor = addMonths(m, 1)
		}
		lo, hi := *start, *end
		if *picking && !hover.IsZero() {
			lo, hi = *start, *hover
			if hi.Before(lo) {
				lo, hi = hi, lo
			}
		}
		day, over := grid(c, monthOf(*cursor), cursor, false, func(d time.Time) look {
			edge := sameDay(d, lo) || sameDay(d, hi)
			return look{chosen: edge, band: !edge && d.After(lo) && d.Before(hi)}
		})
		if !sameDay(over, *hover) || over.IsZero() != hover.IsZero() {
			*hover = over
			c.AnimationFrame()
		}
		if day.IsZero() {
			return
		}
		if !*picking {
			*start, *end, *picking = onDay(day, *start), onDay(day, *end), true
		} else if dayBefore(day, *start) {
			*start, *picking = onDay(day, *start), false
		} else {
			*end, *picking = onDay(day, *end), false
		}
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("RangeCalendar", "picked", "picking", *picking)
	})
}
