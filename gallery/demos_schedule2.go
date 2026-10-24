package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/schedule"
	"github.com/egoist/mygo/ui"
)

func demoResourceRows(c *ui.Context, log func(string)) {
	events := demoScheduleEvents()
	rooms := []schedule.Resource{
		{Name: "Atlas", Events: []schedule.Event{events[4], events[6]}},
		{Name: "Borealis", Events: []schedule.Event{{ID: "hiring", Title: "Hiring panel", Start: at(23, 9, 0), End: at(23, 11, 0), Tint: 3}}},
		{Name: "Cosmos", Events: []schedule.Event{events[5], {ID: "board", Title: "Board prep", Start: at(23, 13, 0), End: at(23, 14, 30), Tint: 1}}},
		{Name: "Dune", Events: nil},
	}
	ui.Column(c).Width(520).Children(func() { schedule.ResourceRows(c, rooms, schedule.ResourceRowsOptions{Title: "Rooms on Friday"}) })
}

func demoShifts(c *ui.Context, log func(string)) {
	shifts := []schedule.Shift{
		{Name: "Ana Lima", Role: "Barista", Day: "Monday", Start: "06:30", End: "14:30"},
		{Name: "Kai Moreno", Role: "Shift lead", Day: "Monday", Start: "12:00", End: "20:00"},
		{Name: "Noor Haddad", Role: "Barista", Day: "Tuesday", Start: "06:30", End: "14:30"},
		{Name: "Tom Becker", Role: "Kitchen", Day: "Tuesday", Start: "10:00", End: "18:00"},
	}
	ui.Column(c).Width(400).Children(func() { schedule.Shifts(c, shifts, schedule.ShiftsOptions{Title: "Harbor Cafe roster"}) })
}

func demoTimeOff(c *ui.Context, log func(string)) {
	requests := []schedule.TimeOffRequest{
		{Name: "Maya Chen", Kind: "Vacation", Range: "Nov 2–6", Status: schedule.Pending},
		{Name: "Daniel Okafor", Kind: "Sick leave", Range: "Oct 21", Status: schedule.Approved},
		{Name: "Priya Nair", Kind: "Conference", Range: "Nov 12–14", Status: schedule.Approved},
		{Name: "Lucas Martin", Kind: "Vacation", Range: "Dec 21–Jan 2", Status: schedule.Declined},
	}
	ui.Column(c).Width(400).Children(func() { schedule.TimeOff(c, requests, schedule.TimeOffOptions{}) })
}

func demoTimesheet(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		schedule.Timesheet(c, []float64{8, 7.5, 9.5, 8, 6.5, 0, 0}, schedule.TimesheetOptions{})
	})
}

func demoNextEvent(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	event := ui.Local(root, "event", func() schedule.Event {
		e := demoScheduleEvent()
		e.Start = c.Now().Add(65 * time.Minute).Truncate(time.Minute)
		e.End = e.Start.Add(time.Hour)
		return e
	})
	root.Children(func() { schedule.NextEvent(c, *event, schedule.NextEventOptions{}) })
}

func demoQuickEvent(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	title, created := textState(root, "title", ""), boolState(root, "created", false)
	root.Children(func() {
		entry, create := schedule.QuickEvent(c, title, created, schedule.QuickEventOptions{Placeholder: "Coffee with Maya on Friday at 3pm"})
		if entry.Changed() {
			log("title " + fmt.Sprint(utf8.RuneCountInString(*title)) + " runes")
		}
		if create.Changed() {
			log("created")
			*title = ""
		}
	})
}

func demoConflicts(c *ui.Context, log func(string)) {
	events := demoScheduleEvents()
	ui.Column(c).Width(400).Children(func() {
		schedule.Conflicts(c, []schedule.Event{events[4], events[5], events[6]}, schedule.ConflictsOptions{})
	})
}

func demoYearView(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() { schedule.YearView(c, 2026, schedule.YearViewOptions{Now: at(23, 9, 0)}) })
}

func demoItinerary(c *ui.Context, log func(string)) {
	days := []schedule.TripDay{
		{Label: "Day 1", Date: "Fri, Nov 6", Stops: []schedule.Stop{
			{Time: "08:40", Title: "Flight TP 1351 to Lisbon", Place: "Terminal 1", Icon: "plane"},
			{Time: "13:00", Title: "Check in at Casa do Rio", Place: "Rua das Flores 12", Icon: "bed-double"},
			{Time: "19:30", Title: "Dinner at Taberna Alfama", Place: "Alfama", Icon: "utensils"},
		}},
		{Label: "Day 2", Date: "Sat, Nov 7", Stops: []schedule.Stop{
			{Time: "10:00", Title: "Tram 28 to Graça", Place: "Praça Martim Moniz", Icon: "tram-front"},
			{Time: "14:30", Title: "MAAT museum", Place: "Belém", Icon: "landmark"},
		}},
	}
	ui.Column(c).Width(420).Children(func() { schedule.Itinerary(c, days) })
}

func demoClasses(c *ui.Context, log func(string)) {
	classes := []schedule.ClassSession{
		{Name: "Calculus II", Day: "Mon", Start: "09:00", End: "10:30", Room: "B204"},
		{Name: "Organic Chemistry", Day: "Mon", Start: "13:00", End: "14:30", Room: "Lab 3"},
		{Name: "Modern History", Day: "Tue", Start: "10:00", End: "11:30", Room: "A110"},
		{Name: "Calculus II", Day: "Wed", Start: "09:00", End: "10:30", Room: "B204"},
		{Name: "Data Structures", Day: "Thu", Start: "11:00", End: "12:30", Room: "C301"},
		{Name: "Organic Chemistry", Day: "Fri", Start: "13:00", End: "14:30", Room: "Lab 3"},
	}
	ui.Column(c).Width(520).Children(func() { schedule.Classes(c, classes, schedule.ClassesOptions{}) })
}

func demoWaitlist(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(8)
	position := intState(root, "position", 7)
	root.Children(func() {
		schedule.Waitlist(c, *position, schedule.WaitlistOptions{Label: "Pottery class waitlist", Total: 18})
		if ui.Button(c, "Someone left").Disabled(*position == 1).Clicked() {
			*position--
			log(fmt.Sprintf("position %d", *position))
		}
	})
}

func demoSlotCalendar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	picked := intState(root, "picked", -1)
	counts := []int{4, 2, 0, 6, 1, 0, 0, 5, 3, 3, 0, 2, 0, 0}
	days := make([]schedule.SlotDay, len(counts))
	for i, n := range counts {
		days[i] = schedule.SlotDay{Date: at(19+i, 0, 0).Format("Mon 2"), Available: n}
	}
	root.Children(func() {
		if schedule.SlotCalendar(c, days, picked, schedule.SlotCalendarOptions{}).Changed() {
			log("picked " + days[*picked].Date)
		}
	})
}

func demoDragCreate(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	start := ui.Local(root, "start", func() time.Time { return at(23, 13, 0) })
	end := ui.Local(root, "end", func() time.Time { return at(23, 13, 0) })
	created := boolState(root, "created", false)
	root.Children(func() {
		if schedule.DragCreate(c, start, end, created, schedule.DragCreateOptions{}).Changed() {
			log(fmt.Sprintf("created %d min", int(end.Sub(*start).Minutes())))
			*created = false
		}
	})
}

func demoZoneCompare(c *ui.Context, log func(string)) {
	zones := []schedule.Zone{{Name: "San Francisco", Offset: -7}, {Name: "New York", Offset: -4}, {Name: "London", Offset: 1},
		{Name: "Tokyo", Offset: 9}, {Name: "Sydney", Offset: 11}}
	ui.Column(c).Width(420).Children(func() { schedule.ZoneCompare(c, 15, zones, schedule.ZoneCompareOptions{}) })
}

func demoFocusTime(c *ui.Context, log func(string)) {
	hour := func(h float64) float64 { return h / 24 }
	blocks := []schedule.FocusBlock{{Start: hour(9), End: hour(11)}, {Start: hour(13.5), End: hour(15)}, {Start: hour(16), End: hour(17)}}
	ui.Column(c).Width(440).Children(func() { schedule.FocusTime(c, blocks, schedule.FocusTimeOptions{}) })
}

func demoScheduleSend(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	choice := ui.Local(root, "choice", func() time.Time { return time.Time{} })
	options := []time.Time{at(24, 8, 0), at(24, 13, 0), at(26, 9, 0), at(26, 19, 30)}
	root.Children(func() {
		if schedule.ScheduleSend(c, choice, options, schedule.ScheduleSendOptions{DateFormat: "Mon, Jan 2 · 15:04"}).Changed() {
			log("scheduled " + choice.Format("Jan 2 15:04"))
		}
	})
}
