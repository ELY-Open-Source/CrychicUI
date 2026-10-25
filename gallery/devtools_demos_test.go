package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDevtoolsDemosLog(t *testing.T) {
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
		"LogView":       {press("Follow"), "follow false"},
		"LevelFilter":   {press("Debug"), "levels [true true true true]"},
		"Traces":        {press("cart.load"), "span c"},
		"Requests":      {press("POST /api/checkout"), "request 5"},
		"RequestDetail": {press("Body"), "tab 1"},
		"FlameGraph":    {press("checkout"), "frame checkout"},
		"Events": {func(tt *ui.Tester) error {
			if err := tt.Click("Filter events"); err != nil {
				return err
			}
			tt.Type("route")
			return nil
		}, "filterRunes 5"},
		"StateTree": {press("Add a mug"), "items 3"},
		"Flags":     {press("ai-descriptions"), "on 3"},
		"Deploys":   {press("Roll back to v2.13.1"), "rollback d3"},
	}
	quiet := map[string]bool{"Metrics": true, "SpanDetail": true, "Memory": true, "CPU": true, "FPS": true,
		"EnvBadge": true, "Deploy": true, "Pods": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "devtools" || e.ID > 1283 {
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
	if seen != 18 {
		t.Errorf("checked %d devtools demos", seen)
	}
}
