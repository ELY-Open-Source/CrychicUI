package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/schedule"
	"github.com/egoist/mygo/ui"
)

var scheduleEntries = []Entry{
	{ID: 1091, Name: "WeekView", Category: "schedule", Summary: "Days with hour rows and event blocks; cue Select.", Demo: demoWeekView},
	{ID: 1092, Name: "DayView", Category: "schedule", Summary: "Hours with events; now line moves.", Demo: demoDayView},
	{ID: 1093, Name: "MonthView", Category: "schedule", Summary: "Month cells with event chips; cue Select.", Demo: demoMonthView},
	{ID: 1094, Name: "EventCard", Category: "schedule", Summary: "Title, time, place; color bar.", Demo: demoEventCard},
	{ID: 1095, Name: "EventBlock", Category: "schedule", Summary: "Block sized by duration; drag moves; cue Drop.", Demo: demoEventBlock},
	{ID: 1096, Name: "NowLine", Category: "schedule", Summary: "Line at the current time.", Demo: demoNowLine},
	{ID: 1097, Name: "Booking", Category: "schedule", Summary: "Day and slots then confirm; cue Success.", Demo: demoBooking},
	{ID: 1098, Name: "Availability", Category: "schedule", Summary: "Grid of free and busy cells; cue Toggle.", Demo: demoAvailability},
	{ID: 1099, Name: "MeetingCard", Category: "schedule", Summary: "Meeting with join and attendees; cue Tap.", Demo: demoMeetingCard},
	{ID: 1100, Name: "Attendees", Category: "schedule", Summary: "Avatars with responses.", Demo: demoAttendees},
	{ID: 1101, Name: "RSVP", Category: "schedule", Summary: "Yes, maybe, no; cue Select.", Demo: demoRSVP},
	{ID: 1102, Name: "Reminder", Category: "schedule", Summary: "Remind me options; cue Select.", Demo: demoReminder},
	{ID: 1103, Name: "Upcoming", Category: "schedule", Summary: "Next events list.", Demo: demoUpcoming},
	{ID: 1104, Name: "CalendarList", Category: "schedule", Summary: "Calendars with color toggles; cue Toggle.", Demo: demoCalendarList},
	{ID: 1105, Name: "ResourceRows", Category: "schedule", Summary: "Rooms as rows on a time axis.", Demo: demoResourceRows},
	{ID: 1106, Name: "Shifts", Category: "schedule", Summary: "Staff shifts by day.", Demo: demoShifts},
	{ID: 1107, Name: "TimeOff", Category: "schedule", Summary: "Requests with status.", Demo: demoTimeOff},
	{ID: 1108, Name: "Timesheet", Category: "schedule", Summary: "Hours per day and total.", Demo: demoTimesheet},
	{ID: 1109, Name: "NextEvent", Category: "schedule", Summary: "Next event with countdown.", Demo: demoNextEvent},
	{ID: 1110, Name: "QuickEvent", Category: "schedule", Summary: "One-line event entry; cue Pop.", Demo: demoQuickEvent},
	{ID: 1111, Name: "Conflicts", Category: "schedule", Summary: "Overlapping events flagged.", Demo: demoConflicts},
	{ID: 1112, Name: "YearView", Category: "schedule", Summary: "Twelve mini months.", Demo: demoYearView},
	{ID: 1113, Name: "Itinerary", Category: "schedule", Summary: "Trip days with stops.", Demo: demoItinerary},
	{ID: 1114, Name: "Classes", Category: "schedule", Summary: "Weekly class timetable.", Demo: demoClasses},
	{ID: 1115, Name: "Waitlist", Category: "schedule", Summary: "Position with rolls.", Demo: demoWaitlist},
	{ID: 1116, Name: "SlotCalendar", Category: "schedule", Summary: "Dates with available counts.", Demo: demoSlotCalendar},
	{ID: 1117, Name: "DragCreate", Category: "schedule", Summary: "Drag over hours to create an event; cue Pop.", Demo: demoDragCreate},
	{ID: 1118, Name: "ZoneCompare", Category: "schedule", Summary: "Same hour across zones.", Demo: demoZoneCompare},
	{ID: 1119, Name: "FocusTime", Category: "schedule", Summary: "Focus blocks in a day bar.", Demo: demoFocusTime},
	{ID: 1120, Name: "ScheduleSend", Category: "schedule", Summary: "Send later choices; cue Select.", Demo: demoScheduleSend},
}

// at is a time in the demo month, October 2026.
func at(day, hour, minute int) time.Time {
	return time.Date(2026, 10, day, hour, minute, 0, 0, time.UTC)
}

func demoScheduleEvent() schedule.Event {
	return schedule.Event{ID: "design", Title: "Design review", Place: "Atlas room, 4th floor", Start: at(23, 10, 0), End: at(23, 11, 0)}
}

func demoScheduleEvents() []schedule.Event {
	return []schedule.Event{
		{ID: "sprint", Title: "Sprint planning", Place: "Borealis room", Start: at(19, 9, 30), End: at(19, 10, 30), Tint: 1},
		{ID: "maya", Title: "1:1 with Maya Chen", Place: "Video call", Start: at(20, 14, 0), End: at(20, 14, 30), Tint: 2},
		{ID: "northwind", Title: "Northwind customer call", Place: "Video call", Start: at(21, 11, 0), End: at(21, 12, 0), Tint: 3},
		{ID: "learn", Title: "Lunch & learn: accessibility", Place: "Cafeteria", Start: at(22, 12, 30), End: at(22, 13, 30), Tint: 1},
		demoScheduleEvent(),
		{ID: "roadmap", Title: "Roadmap sync", Place: "Cosmos room", Start: at(23, 10, 30), End: at(23, 11, 30), Tint: 2},
		{ID: "retro", Title: "Team retro", Place: "Atlas room", Start: at(23, 15, 0), End: at(23, 16, 0)},
	}
}

func demoScheduleDays() []time.Time {
	days := make([]time.Time, 7)
	for i := range days {
		days[i] = at(19+i, 0, 0)
	}
	return days
}

func demoWeekView(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected := textState(root, "selected", "")
	root.Children(func() {
		o := schedule.WeekViewOptions{Now: at(23, 10, 40)}
		if schedule.WeekView(c, demoScheduleDays(), demoScheduleEvents(), selected, o).Changed() {
			log("selected " + *selected)
		}
	})
}

func demoDayView(c *ui.Context, log func(string)) {
	now := c.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	slot := func(id, title, place string, h, m, mins, tint int) schedule.Event {
		start := day.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
		return schedule.Event{ID: id, Title: title, Place: place, Start: start, End: start.Add(time.Duration(mins) * time.Minute), Tint: tint}
	}
	events := []schedule.Event{
		slot("standup", "Daily standup", "Video call", 9, 0, 15, 1),
		slot("design", "Design review", "Atlas room", 10, 0, 60, 0),
		slot("roadmap", "Roadmap sync", "Cosmos room", 10, 30, 60, 2),
		slot("lunch", "Lunch with Priya Nair", "Cafe Lisboa", 12, 30, 60, 1),
		slot("focus", "Focus: billing export", "", 14, 0, 120, 3),
	}
	ui.Column(c).Width(420).Children(func() { schedule.DayView(c, day, events, schedule.DayViewOptions{}) })
}

func demoMonthView(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected := textState(root, "selected", "")
	root.Children(func() {
		o := schedule.MonthViewOptions{Now: at(23, 9, 0)}
		if schedule.MonthView(c, at(1, 0, 0), demoScheduleEvents(), selected, o).Changed() {
			log("selected " + *selected)
		}
	})
}

func demoEventCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Gap(8).Children(func() {
		schedule.EventCard(c, demoScheduleEvent(), schedule.EventCardOptions{})
		schedule.EventCard(c, demoScheduleEvents()[2], schedule.EventCardOptions{})
	})
}

func demoEventBlock(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	event := ui.Local(root, "event", func() schedule.Event { return demoScheduleEvent() })
	root.Children(func() {
		if schedule.EventBlock(c, event).Changed() {
			log("moved " + event.Start.Format("15:04"))
		}
	})
}

func demoNowLine(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() { schedule.NowLine(c, c.Now()) })
}

func demoBooking(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	day, slot, confirmed := intState(root, "day", -1), textState(root, "slot", ""), boolState(root, "confirmed", false)
	slots := []string{"09:00", "09:30", "10:00", "11:00", "14:00", "15:30"}
	root.Children(func() {
		booking, book := schedule.Booking(c, demoScheduleDays(), slots, day, slot, confirmed, schedule.BookingOptions{Title: "Book a demo with Lena Fischer"})
		if booking.Changed() {
			log(fmt.Sprintf("day %d slot %s", *day, *slot))
		}
		if book.Clicked() {
			log("booked " + *slot)
		}
	})
}

func demoAvailability(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	busy := map[int]bool{1: true, 2: true, 9: true, 13: true, 14: true, 20: true, 25: true, 28: true}
	cells := ui.Local(root, "cells", func() []bool {
		free := make([]bool, 30)
		for i := range free {
			free[i] = !busy[i]
		}
		return free
	})
	selected := ui.Local(root, "selected", func() []bool { return make([]bool, 30) })
	o := schedule.AvailabilityOptions{Days: []string{"Mon", "Tue", "Wed", "Thu", "Fri"}, Hours: []string{"09:00", "10:00", "11:00", "13:00", "14:00", "15:00"}}
	root.Children(func() {
		if schedule.Availability(c, cells, selected, o).Changed() {
			n := 0
			for _, on := range *selected {
				if on {
					n++
				}
			}
			log(fmt.Sprintf("%d picked", n))
		}
	})
}

func demoMeetingCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	joined := boolState(root, "joined", false)
	people := []string{"Maya Chen", "Daniel Okafor", "Priya Nair", "Lucas Martin", "Sofia Rossi"}
	root.Children(func() {
		if schedule.MeetingCard(c, "Design review", at(23, 10, 0), people, joined, schedule.MeetingCardOptions{}).Changed() {
			log("joined")
		}
	})
}

func demoAttendees(c *ui.Context, log func(string)) {
	people := []schedule.Attendee{{Name: "Maya Chen", Response: schedule.Yes}, {Name: "Daniel Okafor", Response: schedule.Yes},
		{Name: "Priya Nair", Response: schedule.Maybe}, {Name: "Lucas Martin", Response: schedule.No}, {Name: "Sofia Rossi"}}
	ui.Column(c).Width(380).Children(func() { schedule.Attendees(c, people, schedule.AttendeesOptions{}) })
}

func demoRSVP(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	response := ui.Local(root, "response", func() schedule.RSVPResponse { return schedule.RSVPNone })
	root.Children(func() {
		if schedule.RSVP(c, response, schedule.RSVPOptions{Prompt: "Going to the team retro?"}).Changed() {
			log("response " + fmt.Sprint(*response))
		}
	})
}

func demoReminder(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	choice := textState(root, "choice", "15 min")
	root.Children(func() {
		options := []string{"At start", "5 min", "15 min", "1 hour", "1 day"}
		if schedule.Reminder(c, choice, options, schedule.ReminderOptions{Label: "Remind me before Design review"}).Changed() {
			log("reminder " + *choice)
		}
	})
}

func demoUpcoming(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		schedule.Upcoming(c, demoScheduleEvents()[2:], schedule.UpcomingOptions{})
	})
}

func demoCalendarList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	visible := ui.Local(root, "visible", func() []bool { return []bool{true, true, false, true} })
	calendars := []schedule.Calendar{{Name: "Work", Tint: 0}, {Name: "Personal", Tint: 1}, {Name: "Team Atlas", Tint: 2}, {Name: "Holidays", Tint: 3}}
	root.Children(func() {
		if schedule.CalendarList(c, calendars, visible, schedule.CalendarListOptions{}).Changed() {
			log(fmt.Sprint(*visible))
		}
	})
}
