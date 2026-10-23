package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

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

func demoScheduleEvent() schedule.Event {
	return schedule.Event{ID: "e1", Title: "Design review", Place: "Room A",
		Start: time.Date(2026, 10, 23, 10, 0, 0, 0, time.UTC), End: time.Date(2026, 10, 23, 11, 0, 0, 0, time.UTC)}
}

func demoScheduleDays() []time.Time {
	days := make([]time.Time, 7)
	for i := range days {
		days[i] = time.Date(2026, 10, 19+i, 0, 0, 0, 0, time.UTC)
	}
	return days
}

func demoWeekView(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := textState(root, "selected", "")
	root.Children(func() {
		if schedule.WeekView(c, demoScheduleDays(), []schedule.Event{demoScheduleEvent()}, selected, schedule.WeekViewOptions{}).Changed() {
			log("selected " + *selected)
		}
	})
}

func demoDayView(c *ui.Context, log func(string)) {
	schedule.DayView(c, demoScheduleDays()[4], []schedule.Event{demoScheduleEvent()}, schedule.DayViewOptions{})
}

func demoMonthView(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := textState(root, "selected", "")
	month := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	root.Children(func() {
		if schedule.MonthView(c, month, []schedule.Event{demoScheduleEvent()}, selected).Changed() {
			log("selected " + *selected)
		}
	})
}

func demoEventCard(c *ui.Context, log func(string)) {
	schedule.EventCard(c, demoScheduleEvent())
}

func demoEventBlock(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	event := ui.Local(root, "event", func() schedule.Event { return demoScheduleEvent() })
	root.Children(func() {
		if schedule.EventBlock(c, event).Changed() {
			log("moved " + event.Start.Format("15:04"))
		}
	})
}

func demoNowLine(c *ui.Context, log func(string)) {
	schedule.NowLine(c, c.Now())
}

func demoBooking(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	day, slot, confirmed := intState(root, "day", -1), textState(root, "slot", ""), boolState(root, "confirmed", false)
	root.Children(func() {
		root, book := schedule.Booking(c, demoScheduleDays(), []string{"10:00", "11:00"}, day, slot, confirmed, schedule.BookingOptions{})
		if root.Changed() {
			log(fmt.Sprintf("day %d slot %s", *day, *slot))
		}
		if book.Clicked() {
			log("booked " + *slot)
			*confirmed = false
		}
	})
}

func demoAvailability(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	cells := ui.Local(root, "cells", func() []bool { return []bool{true, false} })
	selected := ui.Local(root, "selected", func() []bool { return []bool{false, false} })
	root.Children(func() {
		if schedule.Availability(c, cells, selected, schedule.AvailabilityOptions{}).Changed() {
			log(fmt.Sprint(*selected))
		}
	})
}

func demoMeetingCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	joined := boolState(root, "joined", false)
	root.Children(func() {
		if schedule.MeetingCard(c, "Review", time.Now(), []string{"Ana", "Kai"}, joined, schedule.MeetingCardOptions{}).Changed() {
			log("joined")
		}
	})
}

func demoAttendees(c *ui.Context, log func(string)) {
	people := []schedule.Attendee{{Name: "Ana", Response: "Yes"}, {Name: "Kai", Response: "Maybe"}}
	schedule.Attendees(c, people)
}

func demoRSVP(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	response := ui.Local(root, "response", func() schedule.RSVPResponse { return schedule.RSVPNone })
	root.Children(func() {
		if schedule.RSVP(c, response, schedule.RSVPOptions{}).Changed() {
			log("response " + fmt.Sprint(*response))
		}
	})
}

func demoReminder(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	choice := textState(root, "choice", "")
	root.Children(func() {
		if schedule.Reminder(c, choice, []string{"10m", "1h"}, schedule.ReminderOptions{}).Changed() {
			log("reminder " + *choice)
		}
	})
}

func demoUpcoming(c *ui.Context, log func(string)) {
	schedule.Upcoming(c, []schedule.Event{demoScheduleEvent()})
}

func demoCalendarList(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	visible := ui.Local(root, "visible", func() []bool { return []bool{true} })
	calendars := []schedule.Calendar{{Name: "Work", Tint: 0}}
	root.Children(func() {
		if schedule.CalendarList(c, calendars, visible).Changed() {
			log(fmt.Sprint(*visible))
		}
	})
}

func demoResourceRows(c *ui.Context, log func(string)) {
	resources := []schedule.Resource{{Name: "Room A", Events: []schedule.Event{demoScheduleEvent()}}}
	schedule.ResourceRows(c, resources)
}

func demoShifts(c *ui.Context, log func(string)) {
	shifts := []schedule.Shift{{Name: "Ana", Day: "Mon", Start: "09:00", End: "17:00"}}
	schedule.Shifts(c, shifts)
}

func demoTimeOff(c *ui.Context, log func(string)) {
	requests := []schedule.TimeOffRequest{{Name: "Ana", Range: "Nov 1–3", Status: "Approved"}}
	schedule.TimeOff(c, requests)
}

func demoTimesheet(c *ui.Context, log func(string)) {
	schedule.Timesheet(c, []float64{1, 1, 1, 1, 1, 1, 2}, schedule.TimesheetOptions{})
}

func demoNextEvent(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	event := ui.Local(root, "event", func() schedule.Event {
		e := demoScheduleEvent()
		e.Start = time.Now().Add(time.Hour)
		e.End = e.Start.Add(time.Hour)
		return e
	})
	root.Children(func() {
		schedule.NextEvent(c, *event, schedule.NextEventOptions{})
	})
}

func demoQuickEvent(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	title, created := textState(root, "title", ""), boolState(root, "created", false)
	root.Children(func() {
		root, create := schedule.QuickEvent(c, title, created, schedule.QuickEventOptions{})
		if root.Changed() {
			log("title " + fmt.Sprint(utf8.RuneCountInString(*title)) + " runes")
		}
		if create.Changed() {
			log("created " + *title)
			*created = false
		}
	})
}

func demoConflicts(c *ui.Context, log func(string)) {
	a := demoScheduleEvent()
	b := a
	b.ID, b.Title = "e2", "Overlap"
	schedule.Conflicts(c, []schedule.Event{a, b}, schedule.ConflictsOptions{})
}

func demoYearView(c *ui.Context, log func(string)) {
	schedule.YearView(c, 2026, schedule.YearViewOptions{})
}

func demoItinerary(c *ui.Context, log func(string)) {
	days := []schedule.TripDay{{Label: "Day 1", Stops: []string{"Station", "Museum"}}}
	schedule.Itinerary(c, days)
}

func demoClasses(c *ui.Context, log func(string)) {
	classes := []schedule.ClassSession{{Name: "Math", Day: "Mon", Start: "09:00", End: "10:00", Room: "A1"}}
	schedule.Classes(c, classes)
}

func demoWaitlist(c *ui.Context, log func(string)) {
	schedule.Waitlist(c, 3, schedule.WaitlistOptions{})
}

func demoSlotCalendar(c *ui.Context, log func(string)) {
	days := []schedule.SlotDay{{Date: "Fri 23", Available: 3}}
	schedule.SlotCalendar(c, days)
}

func demoDragCreate(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	start := timeState(root, "start", c)
	end := timeState(root, "end", c)
	created := boolState(root, "created", false)
	root.Children(func() {
		if schedule.DragCreate(c, start, end, created, schedule.DragCreateOptions{}).Changed() {
			log("created")
			*created = false
		}
	})
}

func demoZoneCompare(c *ui.Context, log func(string)) {
	zones := []schedule.Zone{{Name: "UTC", Offset: 0}, {Name: "Tokyo", Offset: 9}}
	schedule.ZoneCompare(c, 10, zones)
}

func demoFocusTime(c *ui.Context, log func(string)) {
	blocks := []schedule.FocusBlock{{Start: 0.25, End: 0.75}}
	schedule.FocusTime(c, blocks)
}

func demoScheduleSend(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	choice := timeState(root, "choice", c)
	options := []time.Time{time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)}
	root.Children(func() {
		if schedule.ScheduleSend(c, choice, options, schedule.ScheduleSendOptions{}).Changed() {
			log("scheduled")
		}
	})
}
