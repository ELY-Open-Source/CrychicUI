package datetime

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// BirthdayPickerOptions tune a BirthdayPicker.
type BirthdayPickerOptions struct {
	// Years counts back from this year; 120 when zero.
	Years int
	// Day, Month and Year name the three selects.
	Day, Month, Year string
}

// BirthdayPicker is day, month and year selects.
func BirthdayPicker(c *ui.Context, label string, date *time.Time, o BirthdayPickerOptions) ui.Element {
	if o.Years == 0 {
		o.Years = 120
	}
	def(&o.Day, "Day")
	def(&o.Month, "Month")
	def(&o.Year, "Year")
	last := c.Now().Year()
	first := last - o.Years + 1
	if label == "" || date == nil || o.Years < 1 || date.Year() < first || date.Year() > last {
		panic(fmt.Sprintf("datetime: BirthdayPicker %q needs a date from %d to %d", label, first, last))
	}
	t := c.Theme()
	months := i18n.Months(locale(c), i18n.Wide, true)
	dayNames, years := numbers(1, daysIn(*date)), numbers(last, first)
	day, month, year := strconv.Itoa(date.Day()), months[date.Month()-1], strconv.Itoa(date.Year())
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).Label(label)
	return root.Children(func() {
		c.Key("day")
		changed := choice.Dropdown(c, o.Day, dayNames, &day).Changed()
		c.Key("month")
		changed = choice.Dropdown(c, o.Month, months[:], &month).Changed() || changed
		c.Key("year")
		changed = choice.Dropdown(c, o.Year, years, &year).Changed() || changed
		if changed {
			d := slices.Index(dayNames, day) + 1
			*date = withYearMonth(date.AddDate(0, 0, d-date.Day()), last-slices.Index(years, year),
				time.Month(slices.Index(months[:], month)+1))
			root.SetChanged()
			kit.Log("BirthdayPicker", "picked")
		}
	})
}
