package schedule

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDayViewNowLine(t *testing.T) {
	day := testDays()[4]
	early := ui.NewTester(padded(func(c *ui.Context) {
		DayView(c, day, nil, DayViewOptions{Now: time.Date(2026, 10, 23, 10, 5, 0, 0, time.UTC)})
	}), 600, 500)
	early.Advance(time.Second)
	late := ui.NewTester(padded(func(c *ui.Context) {
		DayView(c, day, nil, DayViewOptions{Now: time.Date(2026, 10, 23, 10, 55, 0, 0, time.UTC)})
	}), 600, 500)
	late.Advance(time.Second)
	if dangerY(early) >= dangerY(late) {
		t.Errorf("now line did not move: early %d late %d", dangerY(early), dangerY(late))
	}
}

func dangerY(tt *ui.Tester) int {
	image := tt.Image()
	for y := image.Rect.Min.Y; y < image.Rect.Max.Y; y++ {
		for x := image.Rect.Min.X; x < image.Rect.Max.X; x++ {
			c := image.RGBAAt(x, y)
			if c.R > 180 && c.G < 100 && c.B < 100 {
				return y
			}
		}
	}
	return image.Rect.Max.Y
}

func TestCopyOverrides(t *testing.T) {
	selected := ""
	week := ui.NewTester(padded(func(c *ui.Context) {
		WeekView(c, testDays(), nil, &selected, WeekViewOptions{DayFormat: "D"})
	}), 700, 500)
	if !week.HasText("D") {
		t.Errorf("week format %q", week.Texts())
	}
	day, slot, confirmed := -1, "", false
	booking := ui.NewTester(padded(func(c *ui.Context) {
		Booking(c, testDays(), []string{"10:00"}, &day, &slot, &confirmed, BookingOptions{DayFormat: "Day"})
	}), 600, 300)
	if !booking.HasText("Day") {
		t.Errorf("booking format %q", booking.Texts())
	}
	joined := false
	meeting := ui.NewTester(padded(func(c *ui.Context) {
		MeetingCard(c, "Review", time.Date(2026, 10, 23, 0, 0, 0, 0, time.UTC), []string{"Ana"}, &joined, MeetingCardOptions{DateFormat: "2006"})
	}), 500, 250)
	if !meeting.HasText("2026") {
		t.Errorf("meeting format %q", meeting.Texts())
	}
	choice := time.Time{}
	options := []time.Time{time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)}
	send := ui.NewTester(padded(func(c *ui.Context) {
		ScheduleSend(c, &choice, options, ScheduleSendOptions{DateFormat: "2006"})
	}), 500, 180)
	if !send.HasText("2026") {
		t.Errorf("send format %q", send.Texts())
	}
	year := ui.NewTester(padded(func(c *ui.Context) {
		YearView(c, 2026, YearViewOptions{Months: [12]string{"M1", "M2", "M3", "M4", "M5", "M6", "M7", "M8", "M9", "M10", "M11", "M12"},
			Weekdays: [7]string{"D1", "D2", "D3", "D4", "D5", "D6", "D7"}})
	}), 600, 500)
	if !year.HasText("M1") || !year.HasText("D1") {
		t.Errorf("year names %q", year.Texts())
	}
}
