package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestProjectDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Board":    {press("Add card"), "cards [3 1 1]"},
		"TaskRow":  {press("Offline sync for saved trips"), "done true"},
		"Subtasks": {press("Resolve conflicts on reconnect"), "3 done"},
		"Priority": {press("Urgent"), "level 0"},
		"Status":   {press("Done"), "status 2"},
		"Assignee": {press("Priya Nair"), "selected 2"},
		"Labels":   {press("Backend"), "[true true true false false]"},
		"Estimate": {press("5"), "points 5"},
		"Timeline": {func(tt *ui.Tester) error {
			r, _ := tt.Find("Research")
			x, y := r.X+r.W+50, r.Y+r.H/2
			tt.Press(x, y)
			tt.Move(x+100, y)
			tt.Release(x+100, y)
			return nil
		}, "ends [0."},
		"IssueHeader": {press("Close"), "action 2"},
		"Comments": {func(tt *ui.Tester) error {
			if err := tt.Click("Write a comment"); err != nil {
				return err
			}
			tt.Type("Looks good")
			return nil
		}, "draft 10 runes"},
		"TimeTracker": {press("Start"), "running true"},
		"Backlog": {func(tt *ui.Tester) error {
			a, _ := tt.Find("Offline sync for saved trips")
			b, _ := tt.Find("Crash when opening a shared map link")
			tt.Press(b.X+b.W/2, b.Y+b.H/2)
			tt.Move(a.X+a.W/2, a.Y+a.H/2)
			tt.Release(a.X+a.W/2, a.Y+a.H/2)
			return nil
		}, "[ATL-151 ATL-142"},
	}
	quiet := map[string]bool{"TaskCard": true, "SprintHeader": true, "Burndown": true, "Velocity": true,
		"Roadmap": true, "Milestone": true, "IssueCard": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "project" || e.ID >= 1141 {
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
	if seen != 20 {
		t.Errorf("checked %d project demos", seen)
	}
}
