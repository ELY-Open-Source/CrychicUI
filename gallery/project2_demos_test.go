package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestProject2DemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Watchers":      {press("Watch"), "watching true"},
		"Dependencies":  {press("A"), "opened A"},
		"ProjectHeader": {press("Docs"), "tab Docs"},
		"Checklist":     {press("Add item"), "items 2"},
		"Retro":         {press("Add card"), "Add card 1"},
		"Vote":          {press("Vote"), "1"},
		"Whiteboard": {func(tt *ui.Tester) error {
			r, _ := tt.Find("Idea")
			x, y := r.X+r.W/2, r.Y+r.H/2
			tt.Press(x, y)
			tt.Move(x+40, y+20)
			tt.Release(x+40, y+20)
			return nil
		}, "moved"},
		"MindMap":   {press("Add branch"), "branches 2"},
		"Release":   {press("Ship"), "shipped"},
		"Templates": {press("Kanban"), "selected 0"},
		"QuickAdd":  {press("Add"), "created"},
	}
	quiet := map[string]bool{"Epic": true, "Workload": true, "ProjectCard": true,
		"StandUp": true, "DocCard": true, "OKR": true, "Risks": true, "Changes": true, "TeamLoad": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "project" || e.ID < 1141 {
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
		t.Errorf("checked %d second project demos", seen)
	}
}
