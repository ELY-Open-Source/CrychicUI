package schedule

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// panics wants a panic whose message holds want.
func panics(t *testing.T, want string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want %q", r, want)
		}
	}()
	ui.Render(view, 500, 400, 1)
}

func at(day, hour, minute int) time.Time {
	return time.Date(2026, 10, day, hour, minute, 0, 0, time.UTC)
}

func ev(id string, start, end time.Time) Event {
	return Event{ID: id, Title: id, Start: start, End: end}
}

func TestLanesSplitOverlaps(t *testing.T) {
	day := at(23, 0, 0)
	got := lanes([]Event{ev("b", at(23, 10, 30), at(23, 11, 30)), ev("a", at(23, 10, 0), at(23, 11, 0)), ev("c", at(23, 11, 0), at(23, 12, 0))}, day, 0)
	ids := [][]string{}
	for _, lane := range got {
		var row []string
		for _, e := range lane {
			row = append(row, e.ID)
		}
		ids = append(ids, row)
	}
	if fmt.Sprint(ids) != "[[a c] [b]]" {
		t.Errorf("lanes %v, want [[a c] [b]]", ids)
	}
}

func TestWeekViewOverlapsSideBySide(t *testing.T) {
	selected := ""
	events := []Event{ev("Design", at(23, 10, 0), at(23, 11, 0)), ev("Roadmap", at(23, 10, 30), at(23, 11, 30))}
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, testDays(), events, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	a, _ := tt.Find("Design")
	b, _ := tt.Find("Roadmap")
	if a.W == 0 || b.X < a.X+a.W-1 {
		t.Errorf("overlapping events cover each other: %v %v", a, b)
	}
	before := baseline()
	click(t, tt, "Roadmap")
	if selected != "Roadmap" {
		t.Errorf("selected %q", selected)
	}
	cues(t, before, sound.Select)
}

func TestWeekViewSplitsMidnight(t *testing.T) {
	selected := ""
	events := []Event{ev("Night deploy", at(22, 23, 0), at(23, 1, 0))}
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, testDays(), events, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	n := 0
	for _, s := range tt.Texts() {
		if s == "23:00–01:00" {
			n++
		}
	}
	if n != 2 {
		t.Errorf("a midnight event shows %d times, want 2", n)
	}
}

func TestWeekViewOpensAtEarlyEvent(t *testing.T) {
	selected := ""
	events := []Event{ev("Gym", at(20, 6, 0), at(20, 7, 0))}
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, testDays(), events, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	head, _ := tt.Find("Tue 20")
	gym, ok := tt.Find("Gym")
	if !ok || gym.H < 30 || gym.Y < head.Y+head.H {
		t.Errorf("06:00 event at %v is out of view below header %v", gym, head)
	}
}

func TestEventChecks(t *testing.T) {
	selected := ""
	panics(t, `duplicate event "a"`, func(c *ui.Context) {
		WeekView(c, testDays(), []Event{ev("a", at(23, 9, 0), at(23, 10, 0)), ev("a", at(23, 11, 0), at(23, 12, 0))}, &selected, WeekViewOptions{})
	})
	bad := ev("a", at(23, 9, 0), at(23, 10, 0))
	bad.Tint = 4
	panics(t, `invalid event "a"`, func(c *ui.Context) { EventCard(c, bad, EventCardOptions{}) })
}

func TestMonthViewMore(t *testing.T) {
	selected := ""
	var events []Event
	for i := range 4 {
		events = append(events, ev(fmt.Sprint("Event ", i), at(23, 9+i, 0), at(23, 10+i, 0)))
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { MonthView(c, at(1, 0, 0), events, &selected, MonthViewOptions{}) }), 700, 700)
	settles(t, tt)
	if !tt.HasText("+2 more") || tt.HasText("Event 2") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDayViewCount(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		DayView(c, at(23, 0, 0), []Event{testEvent(), ev("Other day", at(24, 9, 0), at(24, 10, 0))}, DayViewOptions{Now: at(23, 9, 0)})
	}), 600, 600)
	tt.Advance(time.Second)
	if !slices.Contains(tt.Texts(), "1 event") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestBookingRebookClears(t *testing.T) {
	day, slot, confirmed := -1, "", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Booking(c, testDays(), []string{"10:00", "11:00"}, &day, &slot, &confirmed, BookingOptions{})
	}), 600, 400)
	settles(t, tt)
	click(t, tt, "Fri 23")
	click(t, tt, "10:00")
	click(t, tt, "Book")
	tt.Advance(time.Second)
	if !tt.HasText("Booked for Fri 23 at 10:00") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Book")
	cues(t, before)
	click(t, tt, "11:00")
	tt.Advance(time.Second)
	if confirmed || tt.HasText("Booked for") {
		t.Errorf("a new slot kept the old confirmation: %v %q", confirmed, tt.Texts())
	}
	panics(t, `slot "10:00"`, func(c *ui.Context) {
		Booking(c, testDays(), []string{"10:00", "10:00"}, &day, &slot, &confirmed, BookingOptions{})
	})
}

func TestAvailabilityGrid(t *testing.T) {
	cells := []bool{true, false, true, true, true, true}
	selected := make([]bool, 6)
	o := AvailabilityOptions{Days: []string{"Mon", "Tue", "Wed"}, Hours: []string{"09:00", "10:00"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Availability(c, &cells, &selected, o) }), 500, 300)
	settles(t, tt)
	head, _ := tt.Find("Tue 09:00")
	cell, _ := tt.Find("Tue 10:00")
	if head.X != cell.X || cell.Y <= head.Y {
		t.Errorf("Tue 10:00 at %v is not under Tue 09:00 at %v", cell, head)
	}
	click(t, tt, "Tue 10:00")
	if !slices.Equal(selected, []bool{false, false, false, true, false, false}) {
		t.Errorf("selected %v", selected)
	}
	click(t, tt, "Mon 10:00")
	if selected[1] {
		t.Error("a busy cell was picked")
	}
	panics(t, "6 cells for 3 days", func(c *ui.Context) {
		Availability(c, &cells, &selected, AvailabilityOptions{Days: []string{"a", "b", "c"}, Hours: []string{"x"}})
	})
}

func TestUpcomingEmptyAndDays(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Upcoming(c, nil, UpcomingOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Nothing scheduled") {
		t.Errorf("texts %q", tt.Texts())
	}
	events := []Event{ev("Late", at(24, 9, 0), at(24, 10, 0)), ev("Early", at(23, 9, 0), at(23, 10, 0))}
	tt = ui.NewTester(padded(func(c *ui.Context) { Upcoming(c, events, UpcomingOptions{}) }), 500, 400)
	settles(t, tt)
	texts := strings.Join(tt.Texts(), "|")
	if !strings.Contains(texts, "Friday, Oct 23|09:00|10:00|Early") || strings.Index(texts, "Early") > strings.Index(texts, "Late") {
		t.Errorf("texts %q", texts)
	}
}

func TestNextEventCountdown(t *testing.T) {
	start := time.Now().Add(26*time.Hour + 5*time.Minute)
	event := Event{ID: "e", Title: "Launch", Start: start, End: start.Add(time.Hour)}
	tt := ui.NewTester(padded(func(c *ui.Context) { NextEvent(c, event, NextEventOptions{}) }), 500, 300)
	tt.Frame()
	if !tt.HasText("in 1d 02:05:00") && !tt.HasText("in 1d 02:04:59") {
		t.Errorf("texts %q", tt.Texts())
	}
	far := Event{ID: "e", Title: "Launch", Start: time.Date(2500, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2500, 1, 1, 1, 0, 0, 0, time.UTC)}
	tt = ui.NewTester(padded(func(c *ui.Context) { NextEvent(c, far, NextEventOptions{}) }), 500, 300)
	tt.Frame()
	days := (far.Start.Unix() - time.Now().Unix()) / 86400
	if !tt.HasText(fmt.Sprintf("in %dd ", days)) && !tt.HasText(fmt.Sprintf("in %dd ", days-1)) {
		t.Errorf("far countdown %q, want about %d days", tt.Texts(), days)
	}
	past := Event{ID: "e", Title: "Launch", Start: time.Now().Add(-time.Minute), End: time.Now().Add(time.Hour)}
	tt = ui.NewTester(padded(func(c *ui.Context) { NextEvent(c, past, NextEventOptions{}) }), 500, 300)
	tt.Frame()
	if !tt.HasText("Happening now") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestQuickEventAdded(t *testing.T) {
	title, created := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) { QuickEvent(c, &title, &created, QuickEventOptions{}) }), 500, 220)
	settles(t, tt)
	click(t, tt, "Event title")
	tt.Type("Dentist")
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	if !created || !tt.HasText("Added to calendar") {
		t.Errorf("created %v texts %q", created, tt.Texts())
	}
	tt.Type("!")
	tt.Advance(time.Second)
	if created || !tt.HasText("Press Enter") {
		t.Errorf("typing kept the added note: %v %q", created, tt.Texts())
	}
}

func TestConflictsClear(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Conflicts(c, []Event{ev("a", at(23, 9, 0), at(23, 10, 0)), ev("b", at(23, 10, 0), at(23, 11, 0))}, ConflictsOptions{})
	}), 500, 400)
	settles(t, tt)
	if !tt.HasText("No conflicts") || tt.HasText("Overlaps") {
		t.Errorf("touching events conflict: %q", tt.Texts())
	}
}

func TestWaitlistNext(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Waitlist(c, 1, WaitlistOptions{Total: 5}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("You're next") || !tt.HasText("of 5") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "position 6 of 5", func(c *ui.Context) { Waitlist(c, 6, WaitlistOptions{Total: 5}) })
}

func TestSlotCalendarFullStaysShut(t *testing.T) {
	days := []SlotDay{{Date: "Mon 19", Available: 0}, {Date: "Tue 20", Available: 1}}
	picked := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { SlotCalendar(c, days, &picked, SlotCalendarOptions{}) }), 500, 200)
	settles(t, tt)
	click(t, tt, "Mon 19")
	if picked != -1 {
		t.Errorf("a full day was picked")
	}
	before := baseline()
	click(t, tt, "Tue 20")
	if picked != 1 || !tt.HasText("1 left") || !tt.HasText("Full") {
		t.Errorf("picked %d texts %q", picked, tt.Texts())
	}
	cues(t, before, sound.Select)
}

func TestZoneCompareDayShift(t *testing.T) {
	zones := []Zone{{Name: "Tokyo", Offset: 9}, {Name: "Honolulu", Offset: -10}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ZoneCompare(c, 20, zones, ZoneCompareOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("05:00") || !tt.HasText("+1 day") || !tt.HasText("10:00") || tt.HasText("−1 day") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestFocusTimeTotal(t *testing.T) {
	blocks := []FocusBlock{{Start: 9.0 / 24, End: 11.0 / 24}, {Start: 13.5 / 24, End: 15.0 / 24}}
	tt := ui.NewTester(padded(func(c *ui.Context) { FocusTime(c, blocks, FocusTimeOptions{}) }), 500, 200)
	settles(t, tt)
	if !tt.HasText("3h 30m of focus") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestListChecks(t *testing.T) {
	choice := ""
	panics(t, `option "1h"`, func(c *ui.Context) { Reminder(c, &choice, []string{"1h", "1h"}, ReminderOptions{}) })
	panics(t, `"Art" is on unknown day "Sat"`, func(c *ui.Context) {
		Classes(c, []ClassSession{{Name: "Art", Day: "Sat", Start: "9", End: "10", Room: "A"}}, ClassesOptions{})
	})
	panics(t, "window 9-9", func(c *ui.Context) {
		ResourceRows(c, []Resource{{Name: "Atlas"}}, ResourceRowsOptions{From: 9, To: 9})
	})
	when := time.Time{}
	options := []time.Time{at(24, 9, 0), at(24, 9, 0)}
	panics(t, "repeats", func(c *ui.Context) { ScheduleSend(c, &when, options, ScheduleSendOptions{}) })
}
