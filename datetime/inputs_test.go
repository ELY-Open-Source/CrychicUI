package datetime

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func TestAgenda(t *testing.T) {
	day := time.Date(2026, 3, 4, 9, 0, 0, 0, time.Local)
	events := []Event{{"Standup", day, day.Add(15 * time.Minute)}, {"Review", day.Add(5 * time.Hour), day.Add(6 * time.Hour)},
		{"Retro", day.AddDate(0, 0, 1), day.AddDate(0, 0, 1).Add(time.Hour)}}
	picked, changed := -1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Agenda(c, "Week", events, &picked, AgendaOptions{}).Changed() {
			changed++
		}
	}, 500, 400)
	if !tt.HasText("Wed, Mar 4, 2026") || !tt.HasText("Thu, Mar 5, 2026") {
		t.Fatalf("day headers %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Review")
	if picked != 1 || changed != 1 {
		t.Errorf("picked %d, changed %d", picked, changed)
	}
	cues(t, before, sound.Tap)
	events = events[:2]
	tt.Frame()
	if !tt.FramePending() {
		t.Error("removed row did not transition out")
	}
	settles(t, tt, "Agenda")
	events, picked = nil, -1
	tt.Frame()
	if !tt.HasText("Nothing scheduled") {
		t.Errorf("empty %q", tt.Texts())
	}
	twins := []Event{{"Sync", day, day.Add(time.Hour)}, {"Sync", day, day.Add(2 * time.Hour)}, {"Sync", day, day.Add(2 * time.Hour)}}
	p := -1
	tw := ui.NewTester(func(c *ui.Context) { Agenda(c, "Twins", twins, &p, AgendaOptions{}) }, 500, 300)
	tw.Advance(time.Second)
	if n := len(tw.Texts()); n < 7 {
		t.Errorf("twin events: %d texts %q", n, tw.Texts())
	}
	panics(t, "unsorted", func() {
		render(func(c *ui.Context) {
			p := -1
			Agenda(c, "x", []Event{{"B", day.Add(time.Hour), day.Add(time.Hour)}, {"A", day, day}}, &p, AgendaOptions{})
		})
	})
}

func TestBirthdayPickerClamps(t *testing.T) {
	date := time.Date(2000, 1, 31, 8, 0, 0, 0, time.Local)
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if BirthdayPicker(c, "Birthday", &date, BirthdayPickerOptions{}).Changed() {
			changed++
		}
	}, 600, 500)
	before := baseline()
	click(t, tt, "Month")
	click(t, tt, "February")
	if !date.Equal(time.Date(2000, 2, 29, 8, 0, 0, 0, time.Local)) || changed != 1 {
		t.Errorf("date %v, changed %d", date, changed)
	}
	cues(t, before, sound.Select)
	if tt.HasText("31") {
		t.Error("February still offers day 31")
	}
	panics(t, "year out of range", func() {
		render(func(c *ui.Context) {
			d := time.Date(1800, 1, 1, 0, 0, 0, 0, time.Local)
			BirthdayPicker(c, "x", &d, BirthdayPickerOptions{})
		})
	})
}

func TestDateWheelClamps(t *testing.T) {
	date := time.Date(2026, 1, 31, 8, 0, 0, 0, time.Local)
	tt := ui.NewTester(func(c *ui.Context) { DateWheel(c, "Date", &date, DateWheelOptions{}) }, 600, 400)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Advance(time.Second)
	if !date.Equal(time.Date(2026, 2, 28, 8, 0, 0, 0, time.Local)) {
		t.Errorf("date %v", date)
	}
	cues(t, before, sound.Tick)
	panics(t, "year out of range", func() {
		render(func(c *ui.Context) {
			d := time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)
			DateWheel(c, "x", &d, DateWheelOptions{First: 1990, Last: 2000})
		})
	})
}

func TestRecurrence(t *testing.T) {
	rule := Rule{Every: 1, Unit: UnitDay}
	tt := ui.NewTester(func(c *ui.Context) { Recurrence(c, "Repeat", &rule, RecurrenceOptions{}) }, 600, 400)
	if !tt.HasText("Every day") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Increase Interval")
	click(t, tt, "Unit")
	click(t, tt, "Weeks")
	tt.Advance(time.Second)
	click(t, tt, "Monday")
	if rule.Every != 2 || rule.Unit != UnitWeek || !rule.Days[time.Monday] {
		t.Errorf("rule %+v", rule)
	}
	if !tt.HasText("Every 2 weeks on Mon") {
		t.Errorf("summary %q", tt.Texts())
	}
	cues(t, before, sound.Tick, sound.Select, sound.Toggle)
	panics(t, "zero interval", func() {
		render(func(c *ui.Context) { Recurrence(c, "x", &Rule{Unit: UnitDay}, RecurrenceOptions{}) })
	})
}

func TestBusinessHours(t *testing.T) {
	var week [7]Hours
	week[time.Monday] = Hours{true, 9 * time.Hour, 17 * time.Hour}
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if BusinessHours(c, "Hours", &week, BusinessHoursOptions{}).Changed() {
			changed++
		}
	}, 600, 600)
	before := baseline()
	click(t, tt, "Tuesday")
	if !week[time.Tuesday].Open {
		t.Error("Tuesday did not open")
	}
	click(t, tt, "Monday opens")
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if from := week[time.Monday].From; from != 10*time.Hour && from != 9*time.Hour+time.Minute || changed != 2 {
		t.Errorf("Monday opens %v, changed %d", week[time.Monday].From, changed)
	}
	cues(t, before, sound.Toggle, sound.Tick)
	panics(t, "a day past 24 h", func() {
		render(func(c *ui.Context) {
			var w [7]Hours
			w[0].To = 25 * time.Hour
			BusinessHours(c, "x", &w, BusinessHoursOptions{})
		})
	})
}

func TestHolidayListMarksToday(t *testing.T) {
	var list []Holiday
	tt := ui.NewTester(func(c *ui.Context) {
		if list != nil {
			HolidayList(c, "Holidays", list, HolidayListOptions{})
		}
	}, 500, 300)
	next := tomorrow(tt.Now())
	list = []Holiday{{"Founders", tt.Now().AddDate(0, 0, -3)}, {"Picnic", next.Add(9 * time.Hour)}}
	tt.Frame()
	settles(t, tt, "HolidayList")
	if tt.HasText("Today") {
		t.Error("no holiday is today yet")
	}
	tt.Advance(next.Sub(tt.Now()) + time.Second)
	if !tt.HasText("Today") {
		t.Error("the mark did not move to the new day")
	}
	panics(t, "no holidays", func() { render(func(c *ui.Context) { HolidayList(c, "x", nil, HolidayListOptions{}) }) })
}

func TestWeekStrip(t *testing.T) {
	day := time.Date(2026, 3, 4, 9, 30, 0, 0, time.Local)
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if WeekStrip(c, "Week", &day).Changed() {
			changed++
		}
	}, 600, 120)
	before := baseline()
	click(t, tt, "Fri, Mar 6, 2026")
	if !day.Equal(time.Date(2026, 3, 6, 9, 30, 0, 0, time.Local)) || changed != 1 {
		t.Errorf("day %v, changed %d", day, changed)
	}
	if !tt.FramePending() {
		t.Error("the mark did not slide")
	}
	cues(t, before, sound.Select)
	settles(t, tt, "WeekStrip")
}

func TestCronField(t *testing.T) {
	fields := [5]string{"30", "9", "*", "*", "1-5"}
	tt := ui.NewTester(func(c *ui.Context) { CronField(c, "Schedule", &fields, CronFieldOptions{}) }, 600, 200)
	if !tt.HasText("Minute 30, Hour 9, Weekday 1-5") {
		t.Fatalf("summary %q", tt.Texts())
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Type("x")
	tt.Frame()
	if fields[1] == "9" || !tt.HasText("Check the Hour field") {
		t.Errorf("hour %q, texts %q", fields[1], tt.Texts())
	}
	cues(t, before, sound.Tick)
	for _, c := range []struct {
		f      string
		lo, hi int
		ok     bool
	}{{"*", 0, 59, true}, {"*/15", 0, 59, true}, {"1,2,3", 0, 59, true}, {"1-5/2", 0, 6, true},
		{"5-1", 0, 6, false}, {"60", 0, 59, false}, {"", 0, 59, false}, {"*/0", 0, 59, false}, {"a", 0, 59, false}, {"0", 1, 31, false}} {
		if cronValid(c.f, c.lo, c.hi) != c.ok {
			t.Errorf("cronValid(%q, %d, %d) = %v", c.f, c.lo, c.hi, !c.ok)
		}
	}
}

func TestDateNav(t *testing.T) {
	day := time.Date(2026, 3, 31, 9, 0, 0, 0, time.Local)
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if DateNav(c, &day, DateNavOptions{Unit: UnitMonth}).Changed() {
			changed++
		}
	}, 600, 120)
	if !tt.HasText("March 2026") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Next")
	if !day.Equal(time.Date(2026, 4, 30, 9, 0, 0, 0, time.Local)) || !tt.FramePending() {
		t.Errorf("next month: %v, sliding %v", day, tt.FramePending())
	}
	settles(t, tt, "DateNav")
	if !tt.HasText("April 2026") {
		t.Errorf("label %q", tt.Texts())
	}
	click(t, tt, "Today")
	if !sameDay(day, tt.Now()) || changed != 2 {
		t.Errorf("today: %v, changed %d", day, changed)
	}
	cues(t, before, sound.Tap, sound.Tap)
	panics(t, "unknown unit", func() {
		render(func(c *ui.Context) { DateNav(c, &day, DateNavOptions{Unit: 9}) })
	})
}

func TestTimeSlots(t *testing.T) {
	day := time.Date(2026, 3, 4, 9, 0, 0, 0, time.Local)
	slots := []time.Time{day, day.Add(30 * time.Minute), day.Add(time.Hour)}
	var selected time.Time
	var ten string
	tt := ui.NewTester(func(c *ui.Context) {
		ten = i18n.FormatDate(locale(c), slots[2], i18n.DateOptions{Skeleton: "jm"})
		TimeSlots(c, "Slots", slots, slots[1:2], &selected)
	}, 500, 120)
	before := baseline()
	click(t, tt, ten)
	if !selected.Equal(slots[2]) {
		t.Errorf("selected %v", selected)
	}
	cues(t, before, sound.Select)
	panics(t, "unlisted selection", func() {
		render(func(c *ui.Context) { s := day.Add(time.Minute); TimeSlots(c, "x", slots, nil, &s) })
	})
	panics(t, "no slots", func() {
		render(func(c *ui.Context) { var s time.Time; TimeSlots(c, "x", nil, nil, &s) })
	})
	panics(t, "taken yesterday", func() {
		render(func(c *ui.Context) { var s time.Time; TimeSlots(c, "x", slots, []time.Time{day.AddDate(0, 0, -1)}, &s) })
	})
	panics(t, "repeated slot", func() {
		render(func(c *ui.Context) { var s time.Time; TimeSlots(c, "x", []time.Time{day, day.In(time.UTC)}, nil, &s) })
	})
	early := time.Date(2026, 11, 1, 1, 30, 0, 0, time.FixedZone("EDT", -4*3600))
	twice := []time.Time{early, early.Add(time.Hour).In(time.FixedZone("EST", -5*3600))}
	var second time.Time
	tw := ui.NewTester(func(c *ui.Context) { TimeSlots(c, "Night", twice, nil, &second) }, 500, 120)
	tw.Key(0, ui.KeyTab)
	tw.Key(0, ui.KeyRight)
	tw.Frame()
	if !second.Equal(twice[1]) {
		t.Errorf("same-label slots: picked %v, want %v", second, twice[1])
	}
}

func TestDurationPicker(t *testing.T) {
	d := 90 * time.Minute
	tt := ui.NewTester(func(c *ui.Context) { DurationPicker(c, "Length", &d, DurationPickerOptions{}) }, 600, 120)
	if !tt.HasText("1 h") || !tt.HasText("30 min") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Increase Minutes")
	click(t, tt, "Decrease Hours")
	if d != 35*time.Minute {
		t.Errorf("duration %v", d)
	}
	cues(t, before, sound.Tick, sound.Tick)
	panics(t, "off the step", func() {
		render(func(c *ui.Context) { x := 93 * time.Minute; DurationPicker(c, "x", &x, DurationPickerOptions{}) })
	})
}
