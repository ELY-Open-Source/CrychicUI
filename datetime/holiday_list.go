package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Holiday is a named date.
type Holiday struct {
	Name string
	Date time.Time
}

// HolidayListOptions tune a HolidayList.
type HolidayListOptions struct {
	// Today marks the row dated today.
	Today string
}

// HolidayList lists dated rows; today's row is marked.
func HolidayList(c *ui.Context, label string, holidays []Holiday, o HolidayListOptions) ui.Element {
	if label == "" || len(holidays) == 0 {
		panic("datetime: HolidayList needs a label and holidays")
	}
	for i, h := range holidays {
		if h.Name == "" || h.Date.IsZero() {
			panic(fmt.Sprintf("datetime: HolidayList row %d needs a name and a date", i))
		}
	}
	def(&o.Today, "Today")
	t := c.Theme()
	loc := locale(c)
	now := c.Now()
	c.After(tomorrow(now).Sub(now))
	weekdays := i18n.Weekdays(loc, i18n.Abbreviated, true)
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Label(label)
	return root.Children(func() {
		for i, h := range holidays {
			n := days(now, h.Date)
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius)
			bg, fg := t.Accent.Alpha(0), t.Text
			switch {
			case n == 0:
				bg = t.Accent.Alpha(0.12)
			case n < 0:
				fg = t.TextMuted
			}
			row.Background(row.AnimateColor("bg", bg, ui.Snappy)).TextColor(row.AnimateColor("fg", fg, ui.Snappy))
			kit.Arrive(row, ui.Pose{Y: 4, Fade: 1}, time.Duration(min(i, 10))*30*time.Millisecond).Children(func() {
				kit.Numeric(c, i18n.FormatDate(loc, h.Date, i18n.DateOptions{Skeleton: "MMMd"})).Bold().MinWidth(t.Space(16))
				ui.Text(c, h.Name).Grow(1)
				if n == 0 {
					display.Badge(c, o.Today, display.BadgeOptions{Variant: display.BadgeAccent})
				}
				kit.Muted(c, weekdays[h.Date.Weekday()])
			})
		}
	})
}
