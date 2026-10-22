package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/datetime"
	"github.com/egoist/mygo/ui"
)

var datetimeEntries2 = []Entry{
	{514, "TimeAgo", "datetime", "How long ago, kept current.", demoTimeAgo},
	{515, "Agenda", "datetime", "Events listed under their days.", demoAgenda},
	{516, "DatePill", "datetime", "Weekday and day in a pill.", demoDatePill},
	{517, "BirthdayPicker", "datetime", "Day, month and year selects.", demoBirthdayPicker},
	{518, "DateWheel", "datetime", "Three wheels that snap.", demoDateWheel},
	{519, "Recurrence", "datetime", "Repeat every n units, in words.", demoRecurrence},
	{520, "BusinessHours", "datetime", "Opening times for each weekday.", demoBusinessHours},
	{521, "HolidayList", "datetime", "Dated rows; today is marked.", demoHolidayList},
	{522, "Age", "datetime", "Whole years since a birth date.", demoAge},
	{523, "MoonPhase", "datetime", "The moon as lit on a date.", demoMoonPhase},
	{524, "SunArc", "datetime", "The sun between sunrise and sunset.", demoSunArc},
	{525, "DayProgress", "datetime", "How much of today has gone.", demoDayProgress},
	{526, "WeekStrip", "datetime", "Days of a week; the mark slides.", demoWeekStrip},
	{527, "DueChip", "datetime", "A due date tinted by nearness.", demoDueChip},
	{528, "DeadlineBadge", "datetime", "Overdue, today or days left.", demoDeadlineBadge},
	{529, "CronField", "datetime", "Five cron fields with a summary.", demoCronField},
	{530, "DateNav", "datetime", "Previous, today and next.", demoDateNav},
	{531, "TimeSlots", "datetime", "A day's free slots as chips.", demoTimeSlots},
	{532, "DurationPicker", "datetime", "Hours and minutes steppers.", demoDurationPicker},
	{533, "FiscalYear", "datetime", "Quarter bars filling with time.", demoFiscalYear},
	{534, "CalendarHeat", "datetime", "Day cells shaded by counts.", demoCalendarHeat},
	{535, "DateBadge", "datetime", "A calendar leaf that flips.", demoDateBadge},
}

func demoTimeAgo(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8)
	opened := ui.Local(root, "opened", func() time.Time { return c.Now() })
	root.Children(func() {
		for _, ago := range []time.Duration{20 * time.Second, 7 * time.Minute, 3 * time.Hour, 50 * time.Hour} {
			datetime.TimeAgo(c, opened.Add(-ago), datetime.TimeAgoOptions{})
		}
	})
}

func demoAgenda(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	d := day(root)
	picked := ui.Local(root, "picked", func() int { return -1 })
	events := []datetime.Event{{Title: "Standup", Start: *d, End: d.Add(15 * time.Minute)},
		{Title: "Design review", Start: d.Add(4 * time.Hour), End: d.Add(5 * time.Hour)},
		{Title: "Retro", Start: d.AddDate(0, 0, 1), End: d.AddDate(0, 0, 1).Add(time.Hour)}}
	root.Children(func() {
		if datetime.Agenda(c, "This week", events, picked, datetime.AgendaOptions{}).Changed() {
			log("picked: " + events[*picked].Title)
		}
	})
}

func demoDatePill(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(8).AlignItems(ui.Center)
	d := day(root)
	root.Children(func() {
		datetime.DatePill(c, *d)
		if button.Outline(c, "Next day").Clicked() {
			*d = d.AddDate(0, 0, 1)
			log(stamp(*d))
		}
	})
}

func demoBirthdayPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := ui.Local(root, "birth", func() time.Time { return time.Date(1994, 5, 17, 0, 0, 0, 0, time.Local) })
	root.Children(func() {
		if datetime.BirthdayPicker(c, "Birthday", d, datetime.BirthdayPickerOptions{}).Changed() {
			log(d.Format("Jan 2, 2006"))
		}
	})
}

func demoDateWheel(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.DateWheel(c, "Date", d, datetime.DateWheelOptions{}).Changed() {
			log(d.Format("Jan 2, 2006"))
		}
	})
}

func demoRecurrence(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	rule := ui.Local(root, "rule", func() datetime.Rule {
		return datetime.Rule{Every: 2, Unit: datetime.UnitWeek, Days: [7]bool{1: true, 3: true}}
	})
	root.Children(func() {
		if datetime.Recurrence(c, "Repeat", rule, datetime.RecurrenceOptions{}).Changed() {
			log(fmt.Sprintf("every %d, unit %d", rule.Every, rule.Unit))
		}
	})
}

func demoBusinessHours(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	week := ui.Local(root, "week", func() [7]datetime.Hours {
		var w [7]datetime.Hours
		for d := 1; d <= 5; d++ {
			w[d] = datetime.Hours{Open: true, From: 9 * time.Hour, To: 17 * time.Hour}
		}
		return w
	})
	root.Children(func() {
		if datetime.BusinessHours(c, "Opening hours", week, datetime.BusinessHoursOptions{}).Changed() {
			open := 0
			for _, h := range week {
				if h.Open {
					open++
				}
			}
			log(fmt.Sprintf("%d days open", open))
		}
	})
}

func demoHolidayList(c *ui.Context, log func(string)) {
	now := c.Now()
	ui.Column(c).Width(420).Children(func() {
		datetime.HolidayList(c, "Holidays", []datetime.Holiday{{Name: "Founders' Day", Date: now.AddDate(0, 0, -12)},
			{Name: "Team picnic", Date: now}, {Name: "Winter break", Date: now.AddDate(0, 0, 40)}}, datetime.HolidayListOptions{})
	})
}

func demoAge(c *ui.Context, log func(string)) {
	datetime.Age(c, time.Date(1994, 5, 17, 0, 0, 0, 0, time.Local), datetime.AgeOptions{})
}

func demoMoonPhase(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(12).AlignItems(ui.Center)
	d := day(root)
	root.Children(func() {
		datetime.MoonPhase(c, *d, datetime.MoonPhaseOptions{})
		if button.Outline(c, "Next week").Clicked() {
			*d = d.AddDate(0, 0, 7)
			log(stamp(*d))
		}
	})
}

func demoSunArc(c *ui.Context, log func(string)) {
	now := c.Now()
	dawn := time.Date(now.Year(), now.Month(), now.Day(), 6, 30, 0, 0, now.Location())
	ui.Column(c).Width(360).Children(func() {
		datetime.SunArc(c, dawn, dawn.Add(12*time.Hour+15*time.Minute))
	})
}

func demoDayProgress(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() { datetime.DayProgress(c, "Today", datetime.DayProgressOptions{}) })
}

func demoWeekStrip(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.WeekStrip(c, "Week", d).Changed() {
			log(stamp(*d))
		}
	})
}

func demoDueChip(c *ui.Context, log func(string)) {
	now := c.Now()
	ui.Row(c).Gap(8).Children(func() {
		for _, n := range []int{-2, 1, 9} {
			datetime.DueChip(c, now.AddDate(0, 0, n), datetime.DueChipOptions{})
		}
	})
}

func demoDeadlineBadge(c *ui.Context, log func(string)) {
	now := c.Now()
	ui.Row(c).Gap(8).Children(func() {
		for _, n := range []int{-1, 0, 3} {
			datetime.DeadlineBadge(c, now.AddDate(0, 0, n), datetime.DeadlineBadgeOptions{})
		}
	})
}

func demoCronField(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	fields := ui.Local(root, "fields", func() [5]string { return [5]string{"30", "9", "*", "*", "1-5"} })
	root.Children(func() {
		if datetime.CronField(c, "Schedule", fields, datetime.CronFieldOptions{}).Changed() {
			log(fmt.Sprintf("edited: %d chars", len(fields[0]+fields[1]+fields[2]+fields[3]+fields[4])))
		}
	})
}

func demoDateNav(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.DateNav(c, d, datetime.DateNavOptions{Unit: datetime.UnitWeek}).Changed() {
			log(stamp(*d))
		}
	})
}

func demoTimeSlots(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	d := day(root)
	picked := ui.Local(root, "picked", func() time.Time { return time.Time{} })
	start := time.Date(d.Year(), d.Month(), d.Day(), 9, 0, 0, 0, d.Location())
	var slots []time.Time
	for i := range 8 {
		slots = append(slots, start.Add(time.Duration(i)*30*time.Minute))
	}
	root.Children(func() {
		if datetime.TimeSlots(c, "Slots", slots, []time.Time{slots[2], slots[5]}, picked).Changed() {
			log(stamp(*picked))
		}
	})
}

func demoDurationPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := ui.Local(root, "length", func() time.Duration { return 90 * time.Minute })
	root.Children(func() {
		if datetime.DurationPicker(c, "Length", d, datetime.DurationPickerOptions{}).Changed() {
			log(d.String())
		}
	})
}

func demoFiscalYear(c *ui.Context, log func(string)) {
	now := c.Now()
	ui.Column(c).Width(420).Children(func() {
		datetime.FiscalYear(c, "FY "+fmt.Sprint(now.Year()), time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()),
			datetime.FiscalYearOptions{})
	})
}

func demoCalendarHeat(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	counts := make([]int, 26*7)
	for i := range counts {
		counts[i] = (i*7 + i/5*3) % 6
	}
	root.Children(func() {
		datetime.CalendarHeat(c, "Commits", d.AddDate(0, 0, -len(counts)+1), counts, datetime.CalendarHeatOptions{})
	})
}

func demoDateBadge(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(12).AlignItems(ui.Center)
	d := day(root)
	root.Children(func() {
		datetime.DateBadge(c, *d)
		if button.Outline(c, "Next day").Clicked() {
			*d = d.AddDate(0, 0, 1)
			log(stamp(*d))
		}
	})
}
