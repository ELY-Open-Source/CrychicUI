package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// MiniCalendarOptions tune a MiniCalendar.
type MiniCalendarOptions struct {
	// Previous and Next name the month buttons.
	Previous, Next string
}

// MiniCalendar is a compact month; marked days carry a dot.
func MiniCalendar(c *ui.Context, date *time.Time, marks []time.Time, o MiniCalendarOptions) ui.Element {
	if date == nil {
		panic("datetime: MiniCalendar needs a date")
	}
	def(&o.Previous, "Previous month")
	def(&o.Next, "Next month")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(1)).AlignSelf(ui.Start)
	cursor := ui.Local(root, "cursor", func() time.Time { return *date })
	return root.Children(func() {
		m := monthOf(*cursor)
		back, ahead := head(c, i18n.FormatDate(locale(c), m, i18n.DateOptions{Skeleton: "yMMM"}), o.Previous, o.Next)
		if back.Clicked() {
			*cursor = addMonths(m, -1)
		}
		if ahead.Clicked() {
			*cursor = addMonths(m, 1)
		}
		day, _ := grid(c, monthOf(*cursor), cursor, true, func(d time.Time) look {
			marked := false
			for _, k := range marks {
				marked = marked || sameDay(d, k)
			}
			return look{chosen: sameDay(d, *date), dot: marked}
		})
		if !day.IsZero() {
			*date = onDay(day, *date)
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("MiniCalendar", "picked", "day", day.Day())
		}
	})
}
