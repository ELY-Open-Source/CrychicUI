package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// YearPickerOptions tune a YearPicker.
type YearPickerOptions struct {
	// Previous and Next name the decade buttons.
	Previous, Next string
}

// YearPicker picks a year from a decade; decades slide.
func YearPicker(c *ui.Context, date *time.Time, o YearPickerOptions) ui.Element {
	if date == nil {
		panic("datetime: YearPicker needs a date")
	}
	def(&o.Previous, "Previous decade")
	def(&o.Next, "Next decade")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2)).Width(t.Space(72))
	decade := ui.Local(root, "decade", func() int { return date.Year() - (date.Year()%10+10)%10 })
	return root.Children(func() {
		loc := locale(c)
		year := func(y int) string {
			return i18n.FormatDate(loc, time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC), i18n.DateOptions{Skeleton: "y"})
		}
		back, ahead := head(c, year(*decade)+" – "+year(*decade+9), o.Previous, o.Next)
		if back.Clicked() {
			*decade -= 10
		}
		if ahead.Clicked() {
			*decade += 10
		}
		labels := make([]string, 12)
		for i := range labels {
			labels[i] = year(*decade - 1 + i)
		}
		slide(c, *decade, func() {
			tiles(c, labels, 3, func(i int, label string) {
				y := *decade - 1 + i
				if tile(c, label, date.Year() == y, i == 0 || i == 11).Clicked() {
					*date = withYearMonth(*date, y, date.Month())
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("YearPicker", "picked", "year", y)
				}
			})
		})
	})
}
