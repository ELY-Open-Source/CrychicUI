package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// MonthPickerOptions tune a MonthPicker.
type MonthPickerOptions struct {
	// Previous and Next name the year buttons.
	Previous, Next string
}

// MonthPicker picks a month from a year's twelve.
func MonthPicker(c *ui.Context, date *time.Time, o MonthPickerOptions) ui.Element {
	if date == nil {
		panic("datetime: MonthPicker needs a date")
	}
	def(&o.Previous, "Previous year")
	def(&o.Next, "Next year")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2)).Width(t.Space(72))
	year := ui.Local(root, "year", func() int { return date.Year() })
	return root.Children(func() {
		loc := locale(c)
		shown := time.Date(*year, 1, 1, 0, 0, 0, 0, time.UTC)
		back, ahead := head(c, i18n.FormatDate(loc, shown, i18n.DateOptions{Skeleton: "y"}), o.Previous, o.Next)
		if back.Clicked() {
			(*year)--
		}
		if ahead.Clicked() {
			(*year)++
		}
		names := i18n.Months(loc, i18n.Abbreviated, true)
		slide(c, *year, func() {
			tiles(c, names[:], 3, func(i int, name string) {
				mo := time.Month(i + 1)
				if tile(c, name, date.Year() == *year && date.Month() == mo, false).Clicked() {
					*date = withYearMonth(*date, *year, mo)
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("MonthPicker", "picked", "month", i+1)
				}
			})
		})
	})
}
