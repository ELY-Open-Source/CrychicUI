package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// WeekPickerOptions tune a WeekPicker.
type WeekPickerOptions struct {
	// Previous and Next name the month buttons.
	Previous, Next string
}

// WeekPicker picks a whole week; *week becomes its first day.
func WeekPicker(c *ui.Context, week *time.Time, o WeekPickerOptions) ui.Element {
	if week == nil {
		panic("datetime: WeekPicker needs a week")
	}
	def(&o.Previous, "Previous month")
	def(&o.Next, "Next month")
	t := c.Theme()
	first := c.Region().FirstWeekday
	startOf := func(d time.Time) time.Time {
		d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
		return d.AddDate(0, 0, -((int(d.Weekday()) - int(first) + 7) % 7))
	}
	root := ui.Column(c).Gap(t.Space(2)).AlignSelf(ui.Start)
	cursor := ui.Local(root, "cursor", func() time.Time { return *week })
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
		chosen := startOf(*week)
		day, over := grid(c, monthOf(*cursor), cursor, false, func(d time.Time) look {
			return look{chosen: sameDay(startOf(d), chosen), band: !hover.IsZero() && sameDay(startOf(d), startOf(*hover))}
		})
		if !sameDay(over, *hover) || over.IsZero() != hover.IsZero() {
			*hover = over
			c.AnimationFrame()
		}
		if !day.IsZero() {
			*week = startOf(day)
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("WeekPicker", "picked")
		}
	})
}
