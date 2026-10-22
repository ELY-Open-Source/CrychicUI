package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// DateWheelOptions tune a DateWheel.
type DateWheelOptions struct {
	// First and Last bound the years; 100 back, 10 ahead.
	First, Last int
	// Day, Month and Year name the three wheels.
	Day, Month, Year string
}

// DateWheel is day, month and year wheels that snap.
func DateWheel(c *ui.Context, label string, date *time.Time, o DateWheelOptions) ui.Element {
	if o.First == 0 && o.Last == 0 {
		now := c.Now().Year()
		o.First, o.Last = now-100, now+10
	}
	def(&o.Day, "Day")
	def(&o.Month, "Month")
	def(&o.Year, "Year")
	if label == "" || date == nil || o.First > o.Last || date.Year() < o.First || date.Year() > o.Last {
		panic(fmt.Sprintf("datetime: DateWheel %q needs a date from %d to %d", label, o.First, o.Last))
	}
	t := c.Theme()
	months := i18n.Months(locale(c), i18n.Abbreviated, true)
	d, m, y := date.Day()-1, int(date.Month())-1, date.Year()-o.First
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).Label(label)
	return root.Children(func() {
		c.Key("day")
		changed := choice.Wheel(c, o.Day, numbers(1, daysIn(*date)), &d).Changed()
		c.Key("month")
		changed = choice.Wheel(c, o.Month, months[:], &m).Changed() || changed
		c.Key("year")
		changed = choice.Wheel(c, o.Year, numbers(o.First, o.Last), &y).Changed() || changed
		if changed {
			*date = withYearMonth(date.AddDate(0, 0, d+1-date.Day()), o.First+y, time.Month(m+1))
			root.SetChanged()
			kit.Log("DateWheel", "turned")
		}
	})
}
