package datetime

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestDatePicker(t *testing.T) {
	date := time.Date(2026, 3, 2, 14, 30, 0, 0, time.UTC)
	changed := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if DatePicker(c, "Due", &date).Changed() {
			changed++
		}
	}, 400, 500)
	if !tt.FramePending() {
		t.Error("DatePicker did not rise in")
	}
	settles(t, tt, "DatePicker")
	click(t, tt, "Due")
	tt.Key(0, ui.KeySpace)
	settles(t, tt, "DatePicker open")
	click(t, tt, "March 20, 2026")
	if date != time.Date(2026, 3, 20, 14, 30, 0, 0, time.UTC) || changed != 1 {
		t.Errorf("date %v, changed %d", date, changed)
	}
	cues(t, before, sound.Select)
}

func TestDateTimePicker(t *testing.T) {
	at := time.Date(2026, 3, 2, 9, 30, 0, 0, time.UTC)
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { DateTimePicker(c, "Starts", &at, DateTimePickerOptions{}) }, 500, 200)
	settles(t, tt, "DateTimePicker")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if at.Hour() != 10 || at.Day() != 2 {
		t.Errorf("at %v", at)
	}
	cues(t, before, sound.Select)
}

func TestTimePicker(t *testing.T) {
	at := time.Date(2026, 3, 2, 9, 30, 0, 0, time.UTC)
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { TimePicker(c, "Alarm", &at, TimePickerOptions{Step: 15}) }, 400, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	tt.Advance(time.Second)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	tt.Advance(time.Second)
	if at.Hour() != 10 || at.Minute() != 45 || at.Day() != 2 {
		t.Errorf("at %v", at)
	}
	cues(t, before, sound.Tick, sound.Tick)
	panics(t, "step 7", func() { render(func(c *ui.Context) { TimePicker(c, "A", &at, TimePickerOptions{Step: 7}) }) })
}

func TestMonthPicker(t *testing.T) {
	date := time.Date(2026, 3, 31, 9, 0, 0, 0, time.UTC)
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { MonthPicker(c, &date, MonthPickerOptions{}) }, 400, 300)
	click(t, tt, "Jun")
	if date != time.Date(2026, 6, 30, 9, 0, 0, 0, time.UTC) {
		t.Fatalf("date %v", date)
	}
	click(t, tt, "Next year")
	settles(t, tt, "MonthPicker")
	find(t, tt, "2027")
	click(t, tt, "Feb")
	if date != time.Date(2027, 2, 28, 9, 0, 0, 0, time.UTC) {
		t.Errorf("date %v", date)
	}
	cues(t, before, sound.Select, sound.Tap, sound.Select)
}

func TestYearPicker(t *testing.T) {
	date := time.Date(2024, 2, 29, 9, 0, 0, 0, time.UTC)
	tt := ui.NewTester(func(c *ui.Context) { YearPicker(c, &date, YearPickerOptions{}) }, 400, 300)
	find(t, tt, "2020 – 2029")
	click(t, tt, "2025")
	if date != time.Date(2025, 2, 28, 9, 0, 0, 0, time.UTC) {
		t.Fatalf("date %v", date)
	}
	click(t, tt, "Next decade")
	settles(t, tt, "YearPicker")
	find(t, tt, "2030 – 2039")
}

func TestQuarterPicker(t *testing.T) {
	date := time.Date(2026, 2, 14, 9, 0, 0, 0, time.UTC)
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { QuarterPicker(c, &date, QuarterPickerOptions{}) }, 400, 200)
	click(t, tt, "Q3")
	if date != time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC) {
		t.Errorf("date %v", date)
	}
	cues(t, before, sound.Select)
}
