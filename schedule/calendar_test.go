package schedule

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestWeekView(t *testing.T) {
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, testDays(), []Event{testEvent()}, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Design review")
	if selected != "e1" {
		t.Errorf("selected %q", selected)
	}
	cues(t, before, sound.Select)
}

func TestDayView(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { DayView(c, testDays()[4], []Event{testEvent()}, DayViewOptions{}) }), 600, 500)
	settles(t, tt)
	if !tt.HasText("Design review") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMonthView(t *testing.T) {
	selected := ""
	month := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	tt := ui.NewTester(padded(func(c *ui.Context) { MonthView(c, month, []Event{testEvent()}, &selected) }), 700, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Design review")
	if selected != "e1" {
		t.Errorf("selected %q", selected)
	}
	cues(t, before, sound.Select)
}

func TestEventCard(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { EventCard(c, testEvent()) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Design review") || !tt.HasText("Room A") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestEventBlock(t *testing.T) {
	event := testEvent()
	tt := ui.NewTester(padded(func(c *ui.Context) { EventBlock(c, &event) }), 500, 250)
	settles(t, tt)
	r, ok := tt.Find("Design review")
	if !ok {
		t.Fatal("event missing")
	}
	before := baseline()
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Move(x, y+48)
	tt.Release(x, y+48)
	tt.Frame()
	if event.Start.Hour() != 11 {
		t.Errorf("start %v", event.Start)
	}
	cues(t, before, sound.Drop)
}

func TestNowLine(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { NowLine(c, time.Now()) }), 500, 120)
	settles(t, tt)
	if !tt.HasText(":") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestBooking(t *testing.T) {
	day, slot, confirmed := -1, "", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Booking(c, testDays(), []string{"10:00", "11:00"}, &day, &slot, &confirmed, BookingOptions{})
	}), 600, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Fri 23")
	click(t, tt, "10:00")
	click(t, tt, "Book")
	if day != 4 || slot != "10:00" || !confirmed {
		t.Errorf("day %d slot %q confirmed %v", day, slot, confirmed)
	}
	cues(t, before, sound.Select, sound.Select, sound.Success)
}

func TestAvailability(t *testing.T) {
	cells := []bool{true, false}
	selected := []bool{false, false}
	tt := ui.NewTester(padded(func(c *ui.Context) { Availability(c, &cells, &selected, AvailabilityOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Free")
	if !selected[0] {
		t.Errorf("selected %v", selected)
	}
	cues(t, before, sound.Toggle)
}

func TestMeetingCard(t *testing.T) {
	joined := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MeetingCard(c, "Review", time.Now(), []string{"Ana", "Kai"}, &joined, MeetingCardOptions{})
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Join")
	if !joined {
		t.Error("meeting not joined")
	}
	cues(t, before, sound.Tap)
}

func TestAttendees(t *testing.T) {
	people := []Attendee{{Name: "Ana", Response: "Yes"}, {Name: "Kai", Response: "Maybe"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Attendees(c, people) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Yes") || !tt.HasText("Maybe") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestRSVP(t *testing.T) {
	response := RSVPNone
	tt := ui.NewTester(padded(func(c *ui.Context) { RSVP(c, &response, RSVPOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Yes")
	if response != Yes {
		t.Errorf("response %v", response)
	}
	cues(t, before, sound.Select)
}

func TestReminder(t *testing.T) {
	choice := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Reminder(c, &choice, []string{"10m", "1h"}, ReminderOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "10m")
	if choice != "10m" {
		t.Errorf("choice %q", choice)
	}
	cues(t, before, sound.Select)
}

func TestUpcoming(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Upcoming(c, []Event{testEvent()}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Design review") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCalendarList(t *testing.T) {
	visible := []bool{true}
	calendars := []Calendar{{Name: "Work", Tint: 0}}
	tt := ui.NewTester(padded(func(c *ui.Context) { CalendarList(c, calendars, &visible) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Work")
	if visible[0] {
		t.Errorf("visible %v", visible)
	}
	cues(t, before, sound.Toggle)
}

func TestResourceRows(t *testing.T) {
	resources := []Resource{{Name: "Room A", Events: []Event{testEvent()}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ResourceRows(c, resources) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Room A") || !tt.HasText("Design review") {
		t.Errorf("texts %q", tt.Texts())
	}
}
