package schedule

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestWeekViewColumns(t *testing.T) {
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, testDays(), []Event{testEvent()}, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	r, ok := tt.Find("Design review")
	if !ok || r.X < 300 {
		t.Errorf("event fills row instead of Friday column: %v", r)
	}
}

func TestWeekViewTimezone(t *testing.T) {
	zone, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	days := testDays()
	for i := range days {
		days[i] = time.Date(2026, 10, 19+i, 0, 0, 0, 0, zone)
	}
	event := testEvent()
	event.Start = time.Date(2026, 10, 24, 2, 0, 0, 0, time.UTC)
	event.End = event.Start.Add(time.Hour)
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, days, []Event{event}, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	r, ok := tt.Find("Design review")
	if !ok || r.X < 300 || r.X > 450 {
		t.Errorf("event is not in New York Friday: %v", r)
	}
}

func TestEventBlockReset(t *testing.T) {
	event := testEvent()
	tt := ui.NewTester(padded(func(c *ui.Context) { EventBlock(c, &event) }), 500, 300)
	settles(t, tt)
	r, _ := tt.Find("Design review")
	for range 2 {
		x, y := r.X+r.W/2, r.Y+r.H/2
		tt.Press(x, y)
		tt.Move(x, y+30)
		tt.Release(x, y+30)
		tt.Frame()
	}
	if event.Start.Hour() != 10 {
		t.Errorf("separate drags accumulated to %v", event.Start)
	}
}

func TestDragCreateInterval(t *testing.T) {
	start := time.Now()
	end := start
	created := false
	tt := ui.NewTester(padded(func(c *ui.Context) { DragCreate(c, &start, &end, &created, DragCreateOptions{}) }), 500, 250)
	settles(t, tt)
	r, _ := tt.Find("Drag to create")
	x, y := r.X+r.W/2, r.Y+10
	before := baseline()
	tt.Press(x, y)
	tt.Move(x, y+48)
	tt.Release(x, y+48)
	tt.Frame()
	if !created || end.Sub(start) != time.Hour {
		t.Errorf("created %v interval %v", created, end.Sub(start))
	}
	cues(t, before, sound.Pop)
}

func TestDragMotion(t *testing.T) {
	event := testEvent()
	tt := ui.NewTester(padded(func(c *ui.Context) { EventBlock(c, &event) }), 500, 300)
	settles(t, tt)
	r, _ := tt.Find("Design review")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W/2, r.Y+r.H/2+20)
	if !tt.FramePending() {
		t.Error("drag has no visual motion")
	}
	tt.Release(r.X+r.W/2, r.Y+r.H/2+20)
}

func TestBookingInvalidSlot(t *testing.T) {
	day, slot, confirmed := 0, "10:00", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Booking(c, testDays(), []string{"11:00"}, &day, &slot, &confirmed, BookingOptions{})
	}), 600, 350)
	settles(t, tt)
	click(t, tt, "Book")
	if confirmed {
		t.Error("unavailable slot confirmed")
	}
}

func TestResourceRowsAxis(t *testing.T) {
	event := testEvent()
	event.Start = event.Start.Add(-time.Hour)
	event.End = event.Start.Add(time.Hour)
	resources := []Resource{{Name: "Room A", Events: []Event{event}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ResourceRows(c, resources, ResourceRowsOptions{}) }), 700, 220)
	settles(t, tt)
	early, _ := tt.Find("Design review")
	resources[0].Events[0].Start = resources[0].Events[0].Start.Add(6 * time.Hour)
	resources[0].Events[0].End = resources[0].Events[0].Start.Add(time.Hour)
	tt.Frame()
	late, _ := tt.Find("Design review")
	if late.X <= early.X {
		t.Errorf("event did not move on time axis: early %v late %v", early.X, late.X)
	}
}

func TestConflictsFlagged(t *testing.T) {
	a := testEvent()
	b := a
	b.ID, b.Title = "e2", "Overlap"
	other := a
	other.ID, other.Title = "e3", "Later"
	other.Start = other.Start.Add(2 * time.Hour)
	other.End = other.Start.Add(time.Hour)
	tt := ui.NewTester(padded(func(c *ui.Context) { Conflicts(c, []Event{a, b, other}, ConflictsOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("Overlaps Overlap") || !tt.HasText("Overlaps Design review") {
		t.Errorf("conflicting events not flagged: %q", tt.Texts())
	}
	n := 0
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, "Overlaps ") {
			n++
		}
	}
	if n != 2 {
		t.Errorf("%d events flagged, want 2", n)
	}
}

func TestYearViewWeekdays(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { YearView(c, 2025, YearViewOptions{}) }), 600, 500)
	settles(t, tt)
	first, _ := tt.Find("1")
	other := ui.NewTester(padded(func(c *ui.Context) { YearView(c, 2026, YearViewOptions{}) }), 600, 500)
	second, _ := other.Find("1")
	if first.X == second.X {
		t.Errorf("weekday offsets match: %v %v", first.X, second.X)
	}
}

func TestFocusTimeValidation(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Height(1000)
		FocusTime(c, []FocusBlock{{Start: -1, End: 2}}, FocusTimeOptions{})
	}, 500, 100)
}

func TestTimesheetOverflow(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Timesheet(c, []float64{1e308, 1e308, 1, 1, 1, 1, 1}, TimesheetOptions{})
	}, 500, 200, 1)
}

func TestCalendarListDuplicate(t *testing.T) {
	visible := []bool{true, false}
	calendars := []Calendar{{Name: "Work", Tint: 0}, {Name: "Work", Tint: 1}}
	panics(t, `repeated calendar "Work"`, func(c *ui.Context) { CalendarList(c, calendars, &visible, CalendarListOptions{}) })
}
