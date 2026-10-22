package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// DateNavOptions tune a DateNav.
type DateNavOptions struct {
	// Unit is one step's size; a day when zero.
	Unit Unit
	// Previous, Next and Today are button copy.
	Previous, Next, Today string
}

// DateNav steps a date back, to today, or forward.
func DateNav(c *ui.Context, day *time.Time, o DateNavOptions) ui.Element {
	if day == nil || day.IsZero() || o.Unit < UnitDay || o.Unit > UnitYear {
		panic(fmt.Sprintf("datetime: DateNav needs a day and a known unit, not %d", o.Unit))
	}
	def(&o.Previous, "Previous")
	def(&o.Next, "Next")
	def(&o.Today, "Today")
	t := c.Theme()
	loc := locale(c)
	now := c.Now()
	shown := *day
	write := func(skeleton string, x time.Time) string {
		return i18n.FormatDate(loc, x, i18n.DateOptions{Skeleton: skeleton})
	}
	start := shown.AddDate(0, 0, -((int(shown.Weekday()) - int(c.Region().FirstWeekday) + 7) % 7))
	epoch := time.Date(2000, 1, 2, 0, 0, 0, 0, shown.Location())
	title, index := write("yMMMEd", shown), days(epoch, shown)
	switch o.Unit {
	case UnitWeek:
		title, index = write("MMMd", start)+" – "+write("yMMMd", start.AddDate(0, 0, 6)), days(epoch, start)/7
	case UnitMonth:
		title, index = write("yMMMM", shown), monthIndex(shown)
	case UnitYear:
		title, index = write("y", shown), shown.Year()
	}
	step := func(n int) time.Time {
		switch o.Unit {
		case UnitWeek:
			return day.AddDate(0, 0, 7*n)
		case UnitMonth:
			return addMonths(*day, n)
		case UnitYear:
			return withYearMonth(*day, day.Year()+n, day.Month())
		}
		return day.AddDate(0, 0, n)
	}
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center)
	move := func(to time.Time, event string) {
		*day = to
		root.SetChanged()
		kit.Log("DateNav", event)
	}
	return root.Children(func() {
		if button.IconButton(c, o.Previous, "chevron-left").Clicked() {
			move(step(-1), "previous")
		}
		ui.Box(c).MinWidth(t.Space(52)).Center().Children(func() {
			slide(c, index, func() { ui.Text(c, title).Bold().SingleLine() })
		})
		if button.IconButton(c, o.Next, "chevron-right").Clicked() {
			move(step(1), "next")
		}
		if button.Outline(c, o.Today).Disabled(sameDay(*day, now)).Clicked() {
			move(onDay(now, *day), "today")
		}
	})
}
