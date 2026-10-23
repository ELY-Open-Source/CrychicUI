package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestScheduleDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	drag := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			r, _ := tt.Find(name)
			x, y := r.X+r.W/2, r.Y+r.H/2
			tt.Press(x, y)
			tt.Move(x, y+48)
			tt.Release(x, y+48)
			return nil
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"WeekView":   {press("Design review"), "selected e1"},
		"MonthView":  {press("Design review"), "selected e1"},
		"EventBlock": {drag("Design review"), "moved 11:00"},
		"Booking": {func(tt *ui.Tester) error {
			return tt.Click("Fri 23")
		}, "day 4"},
		"Availability": {press("Free"), "true"},
		"MeetingCard":  {press("Join"), "joined"},
		"RSVP":         {press("Yes"), "response 1"},
		"Reminder":     {press("10m"), "reminder 10m"},
		"CalendarList": {press("Work"), "false"},
		"QuickEvent": {func(tt *ui.Tester) error {
			if err := tt.Click("Event title"); err != nil {
				return err
			}
			tt.Type("Lunch")
			return nil
		}, "title 5 runes"},
		"DragCreate":   {drag("Drag to create"), "created"},
		"ScheduleSend": {press("Oct 24 09:00"), "scheduled"},
	}
	quiet := map[string]bool{"DayView": true, "EventCard": true, "NowLine": true, "Attendees": true,
		"Upcoming": true, "ResourceRows": true, "Shifts": true, "TimeOff": true, "Timesheet": true,
		"NextEvent": true, "Conflicts": true, "YearView": true, "Itinerary": true, "Classes": true,
		"Waitlist": true, "SlotCalendar": true, "ZoneCompare": true, "FocusTime": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "schedule" {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 30 {
		t.Errorf("checked %d schedule demos", seen)
	}
}
