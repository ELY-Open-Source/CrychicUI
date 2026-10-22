package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// MonthCalendarOptions tune a MonthCalendar.
type MonthCalendarOptions struct {
	// Previous and Next name the month buttons.
	Previous, Next string
}

// MonthCalendar picks a day; months slide as they change.
func MonthCalendar(c *ui.Context, date *time.Time, o MonthCalendarOptions) ui.Element {
	if date == nil {
		panic("datetime: MonthCalendar needs a date")
	}
	def(&o.Previous, "Previous month")
	def(&o.Next, "Next month")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2)).AlignSelf(ui.Start)
	cursor := ui.Local(root, "cursor", func() time.Time { return *date })
	return root.Children(func() {
		m := monthOf(*cursor)
		back, ahead := head(c, i18n.FormatDate(locale(c), m, i18n.DateOptions{Skeleton: "yMMMM"}), o.Previous, o.Next)
		if back.Clicked() {
			*cursor = addMonths(m, -1)
		}
		if ahead.Clicked() {
			*cursor = addMonths(m, 1)
		}
		m = monthOf(*cursor)
		day, _ := grid(c, m, cursor, false, func(d time.Time) look { return look{chosen: sameDay(d, *date)} })
		if !day.IsZero() {
			*date = onDay(day, *date)
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("MonthCalendar", "picked", "day", day.Day())
		}
	})
}
