package datetime

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestMonthCalendar(t *testing.T) {
	date := time.Date(2026, 3, 31, 9, 15, 0, 0, time.UTC)
	changed := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if MonthCalendar(c, &date, MonthCalendarOptions{}).Changed() {
			changed++
		}
	}, 400, 400)
	find(t, tt, "March 2026")
	click(t, tt, "March 12, 2026")
	if date != time.Date(2026, 3, 12, 9, 15, 0, 0, time.UTC) || changed != 1 {
		t.Fatalf("date %v, changed %d", date, changed)
	}
	click(t, tt, "Next month")
	settles(t, tt, "MonthCalendar slide")
	if !tt.HasText("April 2026") || tt.HasText("March 2026") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "April 30, 2026")
	tt.Key(0, ui.KeyPageDown)
	tt.Frame()
	settles(t, tt, "MonthCalendar page")
	find(t, tt, "May 2026")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if date.Month() != time.May || date.Day() != 30 {
		t.Errorf("Enter picked %v", date)
	}
	cues(t, before, sound.Select, sound.Tap, sound.Select, sound.Select)
}

// cellInk sums a day cell's distance from the background.
func cellInk(t *testing.T, tt *ui.Tester, day string) int {
	t.Helper()
	r := find(t, tt, day)
	bg, img, n := ui.LightTheme().Background, tt.Image(), 0
	for y := int(r.Y) + 2; y < int(r.Y+r.H)-2; y++ {
		for x := int(r.X) + 2; x < int(r.X+r.W)-2; x++ {
			p := img.RGBAAt(x, y)
			n += abs(int(p.R)-int(bg.R)) + abs(int(p.G)-int(bg.G)) + abs(int(p.B)-int(bg.B))
		}
	}
	return n
}

func abs(v int) int { return max(v, -v) }

func TestMultiMonth(t *testing.T) {
	date := time.Date(2026, 3, 10, 8, 0, 0, 0, time.UTC)
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { MultiMonth(c, &date, MultiMonthOptions{}) }, 800, 400)
	find(t, tt, "March 2026 – April 2026")
	click(t, tt, "April 20, 2026")
	if date != time.Date(2026, 4, 20, 8, 0, 0, 0, time.UTC) {
		t.Fatalf("date %v", date)
	}
	click(t, tt, "Next month")
	settles(t, tt, "MultiMonth")
	find(t, tt, "April 2026 – May 2026")
	cues(t, before, sound.Select, sound.Tap)
}

func TestRangeCalendar(t *testing.T) {
	start := time.Date(2026, 3, 10, 18, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 11, 9, 0, 0, 0, time.UTC)
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { RangeCalendar(c, &start, &end, RangeCalendarOptions{}) }, 400, 400)
	plain := cellInk(t, tt, "March 25, 2026")
	click(t, tt, "March 12, 2026")
	r := find(t, tt, "March 15, 2026")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	settles(t, tt, "RangeCalendar preview")
	if cellInk(t, tt, "March 13, 2026") <= plain {
		t.Error("no preview band while picking")
	}
	click(t, tt, "March 15, 2026")
	if !sameDay(start, time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC)) || !sameDay(end, time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("range %v to %v", start, end)
	}
	if start.Hour() != 18 || end.Hour() != 9 {
		t.Errorf("times of day lost: %v, %v", start, end)
	}
	click(t, tt, "March 20, 2026")
	click(t, tt, "March 18, 2026")
	if start.Day() != 18 || end.Day() != 20 {
		t.Errorf("reversed pick gave %v to %v", start, end)
	}
	cues(t, before, sound.Select, sound.Select, sound.Select, sound.Select)
	bad := start.AddDate(0, 0, -5)
	panics(t, "end before start", func() { render(func(c *ui.Context) { RangeCalendar(c, &start, &bad, RangeCalendarOptions{}) }) })
}

func TestMiniCalendarDots(t *testing.T) {
	date := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	marks := []time.Time{time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)}
	tt := ui.NewTester(func(c *ui.Context) { MiniCalendar(c, &date, marks, MiniCalendarOptions{}) }, 300, 300)
	settles(t, tt, "MiniCalendar")
	if cellInk(t, tt, "March 18, 2026") <= cellInk(t, tt, "March 19, 2026") {
		t.Error("marked day has no dot")
	}
	click(t, tt, "March 18, 2026")
	if date.Day() != 18 {
		t.Errorf("date %v", date)
	}
}

func TestWeekPicker(t *testing.T) {
	week := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	var first time.Weekday
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		first = c.Region().FirstWeekday
		if WeekPicker(c, &week, WeekPickerOptions{}).Changed() {
			changed++
		}
	}, 400, 400)
	click(t, tt, "March 12, 2026")
	if week.Weekday() != first || week.After(time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC)) || changed != 1 {
		t.Errorf("week %v (%v), first %v", week, week.Weekday(), first)
	}
}

func TestDateRangePicker(t *testing.T) {
	start := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	end := start
	changed := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if DateRangePicker(c, "Trip", &start, &end, DateRangePickerOptions{}).Changed() {
			changed++
		}
	}, 500, 500)
	find(t, tt, "Mar 2, 2026 – Mar 2, 2026")
	click(t, tt, "Trip")
	settles(t, tt, "DateRangePicker open")
	click(t, tt, "March 5, 2026")
	click(t, tt, "March 9, 2026")
	if start.Day() != 5 || end.Day() != 9 || changed != 2 {
		t.Fatalf("range %v to %v, changed %d", start, end, changed)
	}
	find(t, tt, "Mar 5, 2026 – Mar 9, 2026")
	cues(t, before, sound.Open, sound.Select, sound.Select)
}
