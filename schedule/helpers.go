// Package schedule holds calendar components.
package schedule

import (
	"fmt"
	"math"
	"time"
)

// Event is one calendar event.
type Event struct {
	ID, Title, Place string
	Start, End       time.Time
}

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func validEvent(component string, event Event) {
	if event.ID == "" || event.Title == "" || !event.End.After(event.Start) {
		panic(fmt.Sprintf("schedule: %s has an invalid event", component))
	}
}

func validEvents(component string, events []Event) {
	for _, event := range events {
		validEvent(component, event)
	}
}

func eventText(event Event) string {
	return event.Start.Format("15:04") + "–" + event.End.Format("15:04")
}

func sameDate(a, b time.Time) bool {
	a = a.In(b.Location())
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}
