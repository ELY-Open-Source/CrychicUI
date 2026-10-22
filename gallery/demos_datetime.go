package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/datetime"
	"github.com/egoist/mygo/ui"
)

var datetimeEntries = []Entry{
	{491, "DatePicker", "datetime", "A field opening a calendar.", demoDatePicker},
	{492, "DateRangePicker", "datetime", "Two dates with the range shaded.", demoDateRangePicker},
	{493, "TimePicker", "datetime", "Hour and minute wheels.", demoTimePicker},
	{494, "DateTimePicker", "datetime", "A date and a time together.", demoDateTimePicker},
	{495, "MonthPicker", "datetime", "A grid of twelve months.", demoMonthPicker},
	{496, "YearPicker", "datetime", "Years by decade; decades slide.", demoYearPicker},
	{497, "WeekPicker", "datetime", "A calendar choosing whole weeks.", demoWeekPicker},
	{498, "QuarterPicker", "datetime", "Q1 to Q4 by year.", demoQuarterPicker},
	{499, "MonthCalendar", "datetime", "A month that slides as it changes.", demoMonthCalendar},
	{500, "MultiMonth", "datetime", "Two months side by side.", demoMultiMonth},
	{501, "RangeCalendar", "datetime", "A month previewing the range as you point.", demoRangeCalendar},
	{502, "MiniCalendar", "datetime", "A compact month with event dots.", demoMiniCalendar},
	{503, "AnalogClock", "datetime", "Hands that sweep.", demoAnalogClock},
	{504, "DigitalClock", "datetime", "Rolling digits; the colon blinks.", demoDigitalClock},
	{505, "WorldClock", "datetime", "Zones with their local times.", demoWorldClock},
	{506, "Stopwatch", "datetime", "Start, lap and reset.", demoStopwatch},
	{507, "Countdown", "datetime", "A ring counting down to zero.", demoCountdown},
	{508, "Pomodoro", "datetime", "Work and break in turn.", demoPomodoro},
	{509, "AlarmRow", "datetime", "Time, days and a switch.", demoAlarmRow},
	{510, "ZoneSelect", "datetime", "A searchable time zone select.", demoZoneSelect},
	{511, "RelativeDate", "datetime", "In N days, resolved to a date.", demoRelativeDate},
	{512, "PresetRange", "datetime", "Today, 7 days, 30 days or custom.", demoPresetRange},
	{513, "TimeRangeSlider", "datetime", "Two thumbs over a day.", demoTimeRangeSlider},
}

// day is a fixed demo date, the same every day.
func day(root ui.Element) *time.Time {
	return ui.Local(root, "day", func() time.Time { return time.Date(2026, 3, 12, 9, 30, 0, 0, time.Local) })
}

// stamp writes a demo log line for a date.
func stamp(t time.Time) string { return t.Format("Jan 2, 2006 15:04") }

func demoDatePicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.DatePicker(c, "Due date", d).Changed() {
			log(stamp(*d))
		}
	})
}

func demoDateRangePicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	start := day(root)
	end := ui.Local(root, "end", func() time.Time { return start.AddDate(0, 0, 4) })
	root.Children(func() {
		if datetime.DateRangePicker(c, "Trip", start, end, datetime.DateRangePickerOptions{}).Changed() {
			log(stamp(*start) + " to " + stamp(*end))
		}
	})
}

func demoTimePicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.TimePicker(c, "Alarm", d, datetime.TimePickerOptions{}).Changed() {
			log(d.Format("15:04"))
		}
	})
}

func demoDateTimePicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.DateTimePicker(c, "Meeting", d, datetime.DateTimePickerOptions{}).Changed() {
			log(stamp(*d))
		}
	})
}

func demoMonthPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.MonthPicker(c, d, datetime.MonthPickerOptions{}).Changed() {
			log(d.Format("January 2006"))
		}
	})
}

func demoYearPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.YearPicker(c, d, datetime.YearPickerOptions{}).Changed() {
			log(fmt.Sprint(d.Year()))
		}
	})
}

func demoWeekPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.WeekPicker(c, d, datetime.WeekPickerOptions{}).Changed() {
			log("week of " + d.Format("Jan 2"))
		}
	})
}

func demoQuarterPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.QuarterPicker(c, d, datetime.QuarterPickerOptions{}).Changed() {
			log(fmt.Sprintf("Q%d %d", (int(d.Month())-1)/3+1, d.Year()))
		}
	})
}

func demoMonthCalendar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.MonthCalendar(c, d, datetime.MonthCalendarOptions{}).Changed() {
			log(stamp(*d))
		}
	})
}

func demoMultiMonth(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	root.Children(func() {
		if datetime.MultiMonth(c, d, datetime.MultiMonthOptions{}).Changed() {
			log(stamp(*d))
		}
	})
}

func demoRangeCalendar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	start := day(root)
	end := ui.Local(root, "end", func() time.Time { return start.AddDate(0, 0, 3) })
	root.Children(func() {
		if datetime.RangeCalendar(c, start, end, datetime.RangeCalendarOptions{}).Changed() {
			log(stamp(*start) + " to " + stamp(*end))
		}
	})
}

func demoMiniCalendar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := day(root)
	marks := []time.Time{d.AddDate(0, 0, 2), d.AddDate(0, 0, 5), d.AddDate(0, 0, -4)}
	root.Children(func() {
		if datetime.MiniCalendar(c, d, marks, datetime.MiniCalendarOptions{}).Changed() {
			log(stamp(*d))
		}
	})
}

func demoAnalogClock(c *ui.Context, log func(string)) {
	datetime.AnalogClock(c, datetime.AnalogClockOptions{})
}

func demoDigitalClock(c *ui.Context, log func(string)) {
	datetime.DigitalClock(c, datetime.DigitalClockOptions{Seconds: true})
}

func demoWorldClock(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		datetime.WorldClock(c, []datetime.Zone{{Name: "Local", Location: time.Local},
			{Name: "Tokyo", Location: time.FixedZone("JST", 9*3600)},
			{Name: "Reykjavík", Location: time.FixedZone("GMT", 0)},
			{Name: "Bogotá", Location: time.FixedZone("COT", -5*3600)}}, datetime.WorldClockOptions{})
	})
}

func demoStopwatch(c *ui.Context, log func(string)) {
	if datetime.Stopwatch(c, datetime.StopwatchOptions{}).Changed() {
		log("changed")
	}
}

func demoCountdown(c *ui.Context, log func(string)) {
	if datetime.Countdown(c, 10*time.Second, datetime.CountdownOptions{}).Changed() {
		log("changed")
	}
}

func demoPomodoro(c *ui.Context, log func(string)) {
	if datetime.Pomodoro(c, datetime.PomodoroOptions{Work: 20 * time.Second, Break: 10 * time.Second}).Changed() {
		log("changed")
	}
}

func demoAlarmRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	days := ui.Local(root, "days", func() [7]bool { return [7]bool{false, true, true, true, true, true, false} })
	on := boolState(root, "on", true)
	root.Children(func() {
		if datetime.AlarmRow(c, "Wake up", time.Date(2026, 3, 12, 7, 0, 0, 0, time.Local), days, on,
			datetime.AlarmRowOptions{}).Changed() {
			log(fmt.Sprintf("on %v, days %v", *on, *days))
		}
	})
}

func demoZoneSelect(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	zone := ui.Local(root, "zone", func() string { return "Europe/London" })
	root.Children(func() {
		if datetime.ZoneSelect(c, "Time zone", zone, datetime.ZoneSelectOptions{}).Changed() {
			log(*zone)
		}
	})
}

func demoRelativeDate(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	d := ui.Local(root, "d", func() time.Time { return c.Now() })
	root.Children(func() {
		if datetime.RelativeDate(c, "Remind me", d, datetime.RelativeDateOptions{}).Changed() {
			log(d.Format("Mon Jan 2"))
		}
	})
}

func demoPresetRange(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	start := ui.Local(root, "start", func() time.Time { return c.Now().AddDate(0, 0, -6) })
	end := ui.Local(root, "end", func() time.Time { return c.Now() })
	root.Children(func() {
		if datetime.PresetRange(c, "Period", start, end, datetime.PresetRangeOptions{}).Changed() {
			log(start.Format("Jan 2") + " to " + end.Format("Jan 2"))
		}
	})
}

func demoTimeRangeSlider(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	from := ui.Local(root, "from", func() time.Duration { return 9 * time.Hour })
	to := ui.Local(root, "to", func() time.Duration { return 17 * time.Hour })
	root.Children(func() {
		if datetime.TimeRangeSlider(c, "Open hours", from, to, datetime.TimeRangeSliderOptions{}).Changed() {
			log(fmt.Sprintf("%v to %v", *from, *to))
		}
	})
}
