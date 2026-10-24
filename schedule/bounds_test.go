package schedule

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// swapUnderPress presses name, reorders the data, then releases there.
func swapUnderPress(t *testing.T, tt *ui.Tester, name string, reorder func()) {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("%q missing", name)
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	reorder()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func TestPicksKeepIdentityAcrossReorder(t *testing.T) {
	days := testDays()
	day, slot, confirmed := -1, "", false
	tt := ui.NewTester(padded(func(c *ui.Context) { Booking(c, days, []string{"10:00"}, &day, &slot, &confirmed, BookingOptions{}) }), 600, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Fri 23", func() { days[4], days[5] = days[5], days[4] })
	if day != -1 {
		t.Errorf("Booking picked day %d after a reorder", day)
	}

	calendars := []Calendar{{Name: "Work"}, {Name: "Personal", Tint: 1}}
	visible := []bool{true, true}
	tt = ui.NewTester(padded(func(c *ui.Context) { CalendarList(c, calendars, &visible, CalendarListOptions{}) }), 500, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Work", func() {
		calendars[0], calendars[1] = calendars[1], calendars[0]
		visible[0], visible[1] = visible[1], visible[0]
	})
	if !visible[0] || !visible[1] {
		t.Errorf("CalendarList toggled a moved calendar: %v", visible)
	}

	slots := []SlotDay{{Date: "Mon 19", Available: 2}, {Date: "Tue 20", Available: 2}}
	picked := -1
	tt = ui.NewTester(padded(func(c *ui.Context) { SlotCalendar(c, slots, &picked, SlotCalendarOptions{}) }), 500, 200)
	settles(t, tt)
	swapUnderPress(t, tt, "Mon 19", func() { slots[0], slots[1] = slots[1], slots[0] })
	if picked != -1 {
		t.Errorf("SlotCalendar picked %d after a reorder", picked)
	}

	cells, chosen := []bool{true, true}, []bool{false, false}
	o := AvailabilityOptions{Days: []string{"Mon", "Tue"}, Hours: []string{"09:00"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Availability(c, &cells, &chosen, o) }), 500, 200)
	settles(t, tt)
	swapUnderPress(t, tt, "Mon 09:00", func() { o.Days = []string{"Tue", "Mon"} })
	if slices.Contains(chosen, true) {
		t.Errorf("Availability picked a moved cell: %v", chosen)
	}
}

func TestDaylightSavingKeepsClockHours(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	for _, date := range []time.Time{time.Date(2026, 3, 8, 0, 0, 0, 0, ny), time.Date(2026, 11, 1, 0, 0, 0, 0, ny)} {
		start := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, ny)
		events := []Event{{ID: "e", Title: "Standup", Start: start, End: start.Add(time.Hour)}}
		tt := ui.NewTester(padded(func(c *ui.Context) { DayView(c, date, events, DayViewOptions{Now: date.Add(-time.Hour)}) }), 600, 600)
		settles(t, tt)
		label, _ := tt.Find("09:00")
		block, ok := tt.Find("Standup")
		if !ok || block.Y < label.Y-2 || block.Y > label.Y+20 {
			t.Errorf("%v: 09:00 event at %v, label at %v", date.Format("Jan 2"), block, label)
		}
	}
}

func TestCaptionsUseCalendarZone(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	days := make([]time.Time, 7)
	for i := range days {
		days[i] = time.Date(2026, 10, 19+i, 0, 0, 0, 0, ny)
	}
	events := []Event{ev("Late call", at(24, 2, 0), at(24, 3, 0))}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, days, events, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	if !tt.HasText("22:00–23:00") || tt.HasText("02:00–03:00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestShortEventsDoNotCover(t *testing.T) {
	events := []Event{ev("First", at(23, 9, 0), at(23, 9, 15)), ev("Second", at(23, 9, 15), at(23, 9, 30))}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, testDays(), events, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	a, _ := tt.Find("First")
	b, _ := tt.Find("Second")
	apart := b.X >= a.X+a.W-1 || a.X >= b.X+b.W-1 || b.Y >= a.Y+a.H-1
	if a.H == 0 || !apart {
		t.Errorf("short events cover each other: %v %v", a, b)
	}
}

func TestKeysNeverCollide(t *testing.T) {
	events := []Event{ev("a", at(23, 9, 0), at(23, 10, 0)), ev("day.a", at(24, 9, 0), at(24, 10, 0)),
		ev("2026-10-23", at(23, 11, 0), at(23, 12, 0)), ev("schedule.dayHead:a", at(23, 13, 0), at(23, 14, 0))}
	tt := ui.NewTester(padded(func(c *ui.Context) { Upcoming(c, events, UpcomingOptions{}) }), 500, 400)
	settles(t, tt)
	when := time.Time{}
	options := []time.Time{time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC), time.Date(2584, 7, 21, 8, 34, 33, 709551616, time.UTC)}
	tt = ui.NewTester(padded(func(c *ui.Context) { ScheduleSend(c, &when, options, ScheduleSendOptions{}) }), 500, 300)
	settles(t, tt)
}

func TestRepeatedHourEventsSplit(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 11, 1, 0, 0, 0, 0, ny)
	edt := time.Date(2026, 11, 1, 5, 15, 0, 0, time.UTC)
	est := time.Date(2026, 11, 1, 6, 15, 0, 0, time.UTC)
	events := []Event{{ID: "a", Title: "Early", Start: edt, End: edt.Add(30 * time.Minute)}, {ID: "b", Title: "Late", Start: est, End: est.Add(30 * time.Minute)}}
	tt := ui.NewTester(padded(func(c *ui.Context) { DayView(c, day, events, DayViewOptions{Now: day.AddDate(0, 0, -1)}) }), 600, 600)
	settles(t, tt)
	a, _ := tt.Find("Early")
	b, _ := tt.Find("Late")
	if a.H == 0 || a == b || !(b.X >= a.X+a.W-1 || a.X >= b.X+b.W-1) {
		t.Errorf("events in the repeated hour cover each other: %v %v", a, b)
	}
}

func TestAvailabilityAmbiguousHeadings(t *testing.T) {
	cells, chosen := make([]bool, 4), make([]bool, 4)
	o := AvailabilityOptions{Days: []string{"Mon AM", "Mon"}, Hours: []string{"09:00", "AM 09:00"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Availability(c, &cells, &chosen, o) }), 500, 300)
	settles(t, tt)
}

func TestUpcomingMixedZones(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	late := time.Date(2026, 10, 23, 23, 0, 0, 0, ny)
	events := []Event{ev("a", at(23, 9, 0), at(23, 10, 0)), {ID: "b", Title: "b", Start: late, End: late.Add(time.Hour)}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Upcoming(c, events, UpcomingOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("Saturday, Oct 24") || !tt.HasText("03:00") {
		t.Errorf("texts %q", tt.Texts())
	}
}
