package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// QuarterPickerOptions tune a QuarterPicker.
type QuarterPickerOptions struct {
	// Previous and Next name the year buttons; Format writes "Q%d".
	Previous, Next, Format string
}

// QuarterPicker picks a quarter; *date becomes its first day.
func QuarterPicker(c *ui.Context, date *time.Time, o QuarterPickerOptions) ui.Element {
	if date == nil {
		panic("datetime: QuarterPicker needs a date")
	}
	def(&o.Previous, "Previous year")
	def(&o.Next, "Next year")
	def(&o.Format, "Q%d")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2)).Width(t.Space(72))
	year := ui.Local(root, "year", func() int { return date.Year() })
	return root.Children(func() {
		shown := time.Date(*year, 1, 1, 0, 0, 0, 0, time.UTC)
		back, ahead := head(c, i18n.FormatDate(locale(c), shown, i18n.DateOptions{Skeleton: "y"}), o.Previous, o.Next)
		if back.Clicked() {
			(*year)--
		}
		if ahead.Clicked() {
			(*year)++
		}
		labels := make([]string, 4)
		for i := range labels {
			labels[i] = fmt.Sprintf(o.Format, i+1)
		}
		kit.Slide(c, *year, func() {
			tiles(c, labels, 4, func(q int, label string) {
				chosen := date.Year() == *year && (int(date.Month())-1)/3 == q
				if tile(c, label, chosen, false).Clicked() {
					*date = onDay(time.Date(*year, time.Month(q*3+1), 1, 0, 0, 0, 0, date.Location()), *date)
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("QuarterPicker", "picked", "quarter", q+1)
				}
			})
		})
	})
}
