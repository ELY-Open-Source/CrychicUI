package datetime

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestAlarmRow(t *testing.T) {
	var days [7]bool
	on, changed := true, 0
	var r ui.Rect
	before := baseline()
	at := time.Date(2026, 3, 2, 7, 30, 0, 0, time.UTC)
	tt := ui.NewTester(func(c *ui.Context) {
		row := AlarmRow(c, "Wake up", at, &days, &on, AlarmRowOptions{})
		if row.Changed() {
			changed++
		}
		r = row.LayoutBounds()
	}, 600, 160)
	click(t, tt, "Monday")
	tt.ClickAt(r.X+r.W-ui.LightTheme().Space(4)-12, r.Y+r.H/2)
	tt.Frame()
	if !days[1] || on || changed != 2 {
		t.Fatalf("days %v, on %v, changed %d", days, on, changed)
	}
	settles(t, tt, "AlarmRow dims")
	cues(t, before, sound.Toggle, sound.Toggle)
	panics(t, "nil switch", func() { render(func(c *ui.Context) { AlarmRow(c, "A", at, &days, nil, AlarmRowOptions{}) }) })
}

func TestZoneSelect(t *testing.T) {
	zone := "UTC"
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { ZoneSelect(c, "Time zone", &zone, ZoneSelectOptions{}) }, 400, 500)
	settles(t, tt, "ZoneSelect")
	click(t, tt, "Time zone")
	tt.Type("tokyo")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if zone != "Asia/Tokyo" {
		t.Errorf("zone %q", zone)
	}
	cues(t, before, sound.Select)
}

func TestRelativeDate(t *testing.T) {
	var date time.Time
	tt := ui.NewTester(func(c *ui.Context) { RelativeDate(c, "Remind me", &date, RelativeDateOptions{}) }, 500, 200)
	date = tt.Now()
	tt.Frame()
	before := baseline()
	click(t, tt, "Tomorrow")
	if !sameDay(date, tt.Now().AddDate(0, 0, 1)) {
		t.Errorf("date %v", date)
	}
	cues(t, before, sound.Select)
	panics(t, "unlabeled choice", func() {
		render(func(c *ui.Context) { RelativeDate(c, "R", &date, RelativeDateOptions{Choices: []Relative{{"", 2}}}) })
	})
}

func TestPresetRange(t *testing.T) {
	var start, end time.Time
	tt := ui.NewTester(func(c *ui.Context) { PresetRange(c, "Period", &start, &end, PresetRangeOptions{}) }, 600, 300)
	start, end = tt.Now(), tt.Now()
	tt.Frame()
	if tt.HasText("From") {
		t.Fatal("custom inputs shown for a preset")
	}
	before := baseline()
	click(t, tt, "7 days")
	if !sameDay(start, tt.Now().AddDate(0, 0, -6)) || !sameDay(end, tt.Now()) {
		t.Errorf("range %v to %v", start, end)
	}
	click(t, tt, "Custom")
	settles(t, tt, "PresetRange custom")
	find(t, tt, "From")
	cues(t, before, sound.Select, sound.Select)
}

func TestTimeRangeSlider(t *testing.T) {
	from, to := 9*time.Hour, 17*time.Hour
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		TimeRangeSlider(c, "Open hours", &from, &to, TimeRangeSliderOptions{})
	}, 500, 160)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if from != 9*time.Hour+30*time.Minute || to != 17*time.Hour {
		t.Errorf("from %v, to %v", from, to)
	}
	cues(t, before, sound.Tick)
	from, to = 9*time.Hour, 17*time.Hour
	odd := ui.NewTester(func(c *ui.Context) {
		TimeRangeSlider(c, "Open hours", &from, &to, TimeRangeSliderOptions{Step: 90 * time.Second})
	}, 500, 160)
	odd.Key(0, ui.KeyTab)
	for range 3 {
		odd.Key(0, ui.KeyRight)
		odd.Frame()
	}
	if from != 9*time.Hour+270*time.Second {
		t.Errorf("90 s step: from %v", from)
	}
	off := 9*time.Hour + 10*time.Minute
	panics(t, "off the step", func() {
		render(func(c *ui.Context) { TimeRangeSlider(c, "T", &off, &to, TimeRangeSliderOptions{}) })
	})
}
