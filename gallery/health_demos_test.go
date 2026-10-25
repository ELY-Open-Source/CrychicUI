package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestHealthDemosLog(t *testing.T) {
	press := func(names ...string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, name := range names {
				if err := tt.Click(name); err != nil {
					return err
				}
				tt.Frame()
			}
			return nil
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Water":      {press("Add a glass"), "glasses 6"},
		"Intervals":  {press("Start", "Reset"), "timer toggled\ntimer reset"},
		"Breathe":    {press("Start"), "breathing true"},
		"Mood":       {press("Good"), "mood 3"},
		"Habits":     {press("Stretch on Sun"), "marks 20"},
		"Medication": {press("Iron at 08:00"), "taken 2"},
		"BodyMap":    {press("Left arm"), "region 4"},
	}
	quiet := map[string]bool{"ActivityRings": true, "Steps": true, "HeartRate": true, "Sleep": true, "Calories": true, "Macros": true,
		"Workout": true, "Weight": true, "Cycle": true, "BMI": true, "Meal": true, "Vitals": true, "GoalStreak": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "health" {
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
		tt.Advance(2 * time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 20 {
		t.Errorf("checked %d health demos", seen)
	}
}
