package schedule

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDragCreateRejectsInvalid(t *testing.T) {
	for _, delta := range []float32{-48, 24} {
		start := time.Now()
		end := start
		created := false
		tt := ui.NewTester(padded(func(c *ui.Context) { DragCreate(c, &start, &end, &created, DragCreateOptions{}) }), 500, 300)
		settles(t, tt)
		r, _ := tt.Find("Drag to create")
		before := baseline()
		x, y := r.X+r.W/2, r.Y+20
		tt.Press(x, y)
		tt.Move(x, y+delta)
		tt.Release(x, y+delta)
		tt.Frame()
		if created {
			t.Errorf("delta %v created invalid interval", delta)
		}
		cues(t, before)
	}
}

func TestDayViewIndependentLine(t *testing.T) {
	day := testDays()[4]
	now := time.Date(2026, 10, 23, 10, 30, 0, 0, time.UTC)
	plain := ui.NewTester(padded(func(c *ui.Context) { DayView(c, day, nil, DayViewOptions{Now: now}) }), 600, 500)
	plain.Advance(time.Second)
	withEvent := ui.NewTester(padded(func(c *ui.Context) {
		DayView(c, day, []Event{testEvent()}, DayViewOptions{Now: now})
	}), 600, 500)
	withEvent.Advance(time.Second)
	if dangerY(plain) != dangerY(withEvent) {
		t.Errorf("event moved now line: plain %d event %d", dangerY(plain), dangerY(withEvent))
	}
}

func TestDayViewTimezone(t *testing.T) {
	zone, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 10, 23, 0, 0, 0, 0, zone)
	now := time.Date(2026, 10, 24, 2, 30, 0, 0, time.UTC)
	tt := ui.NewTester(padded(func(c *ui.Context) { DayView(c, day, nil, DayViewOptions{Now: now}) }), 600, 600)
	tt.Advance(time.Second)
	r, ok := tt.Find("22:00")
	if y := float32(dangerY(tt)); !ok || y < r.Y || y > r.Y+40 {
		t.Errorf("now line at %v is not in the 22:00 row %v", y, r)
	}
}

func TestWeekViewHeaderAlignment(t *testing.T) {
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { WeekView(c, testDays(), []Event{testEvent()}, &selected, WeekViewOptions{}) }), 700, 600)
	settles(t, tt)
	header, _ := tt.Find("Fri 23")
	event, _ := tt.Find("Design review")
	if event.X-header.X > 2 || header.X-event.X > 2 {
		t.Errorf("header and event columns differ: header %v event %v", header.X, event.X)
	}
}
