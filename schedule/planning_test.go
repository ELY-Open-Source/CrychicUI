package schedule

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestShifts(t *testing.T) {
	shifts := []Shift{{Name: "Ana Lima", Day: "Mon", Start: "09:00", End: "17:00"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Shifts(c, shifts, ShiftsOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Ana Lima") || !tt.HasText("Mon") || !tt.HasText("09:00–17:00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTimeOff(t *testing.T) {
	requests := []TimeOffRequest{{Name: "Ana Lima", Range: "Nov 1–3", Status: Approved}}
	tt := ui.NewTester(padded(func(c *ui.Context) { TimeOff(c, requests, TimeOffOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Approved") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTimesheet(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Timesheet(c, []float64{1, 1, 1, 1, 1, 1, 2}, TimesheetOptions{}) }), 500, 320)
	settles(t, tt)
	if !tt.HasText("Total 8.0 h") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestNextEvent(t *testing.T) {
	event := testEvent()
	event.Start = time.Now().Add(time.Hour)
	event.End = event.Start.Add(time.Hour)
	tt := ui.NewTester(padded(func(c *ui.Context) { NextEvent(c, event, NextEventOptions{}) }), 500, 180)
	tt.Frame()
	if !tt.HasText("Design review") || !tt.HasText("Next event") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestQuickEvent(t *testing.T) {
	title, created := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		QuickEvent(c, &title, &created, QuickEventOptions{})
	}), 500, 180)
	settles(t, tt)
	click(t, tt, "Event title")
	tt.Type("Lunch")
	before := baseline()
	click(t, tt, "Create")
	if title != "Lunch" || !created {
		t.Errorf("title %q created %v", title, created)
	}
	cues(t, before, sound.Pop)
}

func TestConflicts(t *testing.T) {
	a := testEvent()
	b := a
	b.ID = "e2"
	b.Title = "Overlap"
	tt := ui.NewTester(padded(func(c *ui.Context) { Conflicts(c, []Event{a, b}, ConflictsOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("1 conflict") || tt.HasText("1 conflicts") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestYearView(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { YearView(c, 2026, YearViewOptions{}) }), 600, 500)
	settles(t, tt)
	if !tt.HasText("January") || !tt.HasText("December") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestItinerary(t *testing.T) {
	days := []TripDay{{Label: "Day 1", Stops: []Stop{{Time: "09:00", Title: "Station", Icon: "train-front"}, {Time: "11:00", Title: "Museum", Icon: "landmark"}}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Itinerary(c, days) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Day 1") || !tt.HasText("Museum") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestClasses(t *testing.T) {
	classes := []ClassSession{{Name: "Math", Day: "Mon", Start: "09:00", End: "10:00", Room: "A1"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Classes(c, classes, ClassesOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Math") || !tt.HasText("A1") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestWaitlist(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Waitlist(c, 3, WaitlistOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("#3") || !tt.HasText("2 ahead of you") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSlotCalendar(t *testing.T) {
	days := []SlotDay{{Date: "Fri 23", Available: 3}}
	picked := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { SlotCalendar(c, days, &picked, SlotCalendarOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Fri 23") || !tt.HasText("3 left") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDragCreate(t *testing.T) {
	start := time.Now()
	end := start
	created := false
	tt := ui.NewTester(padded(func(c *ui.Context) { DragCreate(c, &start, &end, &created, DragCreateOptions{}) }), 500, 250)
	settles(t, tt)
	r, ok := tt.Find("Drag to create")
	if !ok {
		t.Fatal("drop area missing")
	}
	before := baseline()
	x, y := r.X+r.W/2, r.Y+10
	tt.Press(x, y)
	tt.Move(x, y+48)
	tt.Release(x, y+48)
	tt.Frame()
	if !created {
		t.Error("event was not created")
	}
	cues(t, before, sound.Pop)
}

func TestZoneCompare(t *testing.T) {
	zones := []Zone{{Name: "UTC", Offset: 0}, {Name: "Tokyo", Offset: 9}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ZoneCompare(c, 10, zones, ZoneCompareOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("10:00") || !tt.HasText("19:00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestFocusTime(t *testing.T) {
	blocks := []FocusBlock{{Start: 0.25, End: 0.75}}
	tt := ui.NewTester(padded(func(c *ui.Context) { FocusTime(c, blocks, FocusTimeOptions{}) }), 500, 180)
	settles(t, tt)
}

func TestScheduleSend(t *testing.T) {
	choice := time.Time{}
	options := []time.Time{time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)}
	tt := ui.NewTester(padded(func(c *ui.Context) { ScheduleSend(c, &choice, options, ScheduleSendOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Oct 24 09:00")
	if !choice.Equal(options[0]) {
		t.Errorf("choice %v", choice)
	}
	cues(t, before, sound.Select)
}
