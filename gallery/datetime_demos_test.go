package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDatetimeDemosLog(t *testing.T) {
	click := func(labels ...string) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, l := range labels {
				if err := tt.Click(l); err != nil {
					return err
				}
				tt.Advance(500 * time.Millisecond)
			}
			return nil
		}
	}
	keys := func(ks ...ui.Key) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, k := range ks {
				tt.Key(0, k)
			}
			return nil
		}
	}
	acts := map[string]func(*ui.Tester) error{
		"DatePicker": func(tt *ui.Tester) error {
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeySpace)
			return click("March 20, 2026")(tt)
		},
		"DateRangePicker": click("Trip", "March 20, 2026"),
		"TimePicker":      keys(ui.KeyTab, ui.KeyDown),
		"DateTimePicker":  keys(ui.KeyTab, ui.KeyTab, ui.KeyUp),
		"MonthPicker":     click("Jun"),
		"YearPicker":      click("2027"),
		"WeekPicker":      click("March 20, 2026"),
		"QuarterPicker":   click("Q2"),
		"MonthCalendar":   click("March 20, 2026"),
		"MultiMonth":      click("March 20, 2026"),
		"RangeCalendar":   click("March 20, 2026"),
		"MiniCalendar":    click("March 20, 2026"),
		"Stopwatch":       click("Start"),
		"Countdown":       click("Start"),
		"Pomodoro":        click("Start"),
		"AlarmRow":        click("Monday"),
		"ZoneSelect":      func(tt *ui.Tester) error { tt.Click("Time zone"); tt.Type("tokyo"); tt.Key(0, ui.KeyEnter); return nil },
		"RelativeDate":    click("Tomorrow"),
		"PresetRange":     click("30 days"),
		"TimeRangeSlider": keys(ui.KeyTab, ui.KeyRight),
	}
	for _, e := range Entries() {
		act, ok := acts[e.Name]
		if !ok {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 600)
		tt.Advance(time.Second)
		if err := act(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if len(logs) == 0 {
			t.Errorf("%s: the action logged nothing", e.Name)
		}
	}
}
