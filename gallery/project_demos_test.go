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
		"Board":    {press("Add card"), "New card"},
		"TaskRow":  {press("Ship UI"), "done true"},
		"Subtasks": {press("One"), "true"},
		"Priority": {press("Urgent"), "level 0"},
		"Status":   {press("Done"), "status 2"},
		"Assignee": {press("Kai"), "selected 1"},
		"Labels":   {press("UI"), "true"},
		"Estimate": {press("5"), "points 5"},
		"Timeline": {func(tt *ui.Tester) error {
			r, _ := tt.Find("Task")
			x, y := r.X+r.W+50, r.Y+r.H/2
			tt.Press(x, y)
			tt.Move(x+100, y)
			tt.Release(x+100, y)
			return nil
		}, "end 0."},
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
			a, _ := tt.Find("Alpha")
			b, _ := tt.Find("Beta")
			tt.Press(b.X+b.W/2, b.Y+b.H/2)
			tt.Move(a.X+a.W/2, a.Y+a.H/2)
			tt.Release(a.X+a.W/2, a.Y+a.H/2)
			return nil
		}, "[Beta Alpha]"},
	}
	quiet := map[string]bool{"TaskCard": true, "SprintHeader": true, "Burndown": true, "Velocity": true,
		"Roadmap": true, "Milestone": true, "IssueCard": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "project" {
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
