package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// MultiMonthOptions tune a MultiMonth.
type MultiMonthOptions struct {
	// Previous and Next name the month buttons.
	Previous, Next string
}

// MultiMonth picks a day from two months side by side.
func MultiMonth(c *ui.Context, date *time.Time, o MultiMonthOptions) ui.Element {
	if date == nil {
		panic("datetime: MultiMonth needs a date")
	}
	def(&o.Previous, "Previous month")
	def(&o.Next, "Next month")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2)).AlignSelf(ui.Start)
	cursor := ui.Local(root, "cursor", func() time.Time { return *date })
	first := ui.Local(root, "first", func() time.Time { return monthOf(*date) })
	return root.Children(func() {
		loc := locale(c)
		title := i18n.FormatDate(loc, *first, i18n.DateOptions{Skeleton: "yMMMM"}) + " – " +
			i18n.FormatDate(loc, addMonths(*first, 1), i18n.DateOptions{Skeleton: "yMMMM"})
		back, ahead := head(c, title, o.Previous, o.Next)
		if back.Clicked() {
			*first = addMonths(*first, -1)
			*cursor = *first
		}
		if ahead.Clicked() {
			*first = addMonths(*first, 1)
			*cursor = *first
		}
		ui.Row(c).Gap(t.Space(4)).Children(func() {
			for i := range 2 {
				m := addMonths(*first, i)
				day, _ := grid(c, m, cursor, false, func(d time.Time) look {
					return look{chosen: sameDay(d, *date) && d.Month() == m.Month()}
				})
				if !day.IsZero() {
					*date = onDay(day, *date)
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("MultiMonth", "picked", "day", day.Day())
				}
			}
		})
		// Keys moving past the two months bring the view along.
		if at := monthIndex(*cursor); at < monthIndex(*first) {
			*first = monthOf(*cursor)
		} else if at > monthIndex(*first)+1 {
			*first = addMonths(monthOf(*cursor), -1)
		}
	})
}
