package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestCodeDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"CodeBlock": {press("Copy"), "copied 56 runes"},
		"CodeEditor": {func(tt *ui.Tester) error {
			tt.Key(0, ui.KeyTab)
			tt.Type("x")
			return nil
		}, "edited"},
		"DiffView": {press("Unified"), "layout unified"},
		"Terminal": {press("New line"), "3 lines"},
		"CommandLine": {func(tt *ui.Tester) error {
			if err := tt.Click("Run a command"); err != nil {
				return err
			}
			tt.Type("pwd")
			tt.Key(0, ui.KeyEnter)
			return nil
		}, "ran 3 runes"},
		"CodeTabs":     {press("util.go"), "selected util.go"},
		"FileHeader":   {press("Copy path"), "copied"},
		"Snippet":      {press("Copy"), "copied 33 runes"},
		"StackTrace":   {press("2 more frames"), "open true"},
		"ErrorOverlay": {press("Trigger"), "open true"},
		"Gutter":       {press("3"), "1 breakpoints"},
		"DebugBar":     {press("Step over"), "over"},
		"Variables":    {press("cfg"), "toggled"},
		"Watches": {func(tt *ui.Tester) error {
			if err := tt.Click("Watch an expression"); err != nil {
				return err
			}
			tt.Type("x")
			tt.Key(0, ui.KeyEnter)
			return nil
		}, "2 watches"},
		"CallStack": {press("http.serve"), "selected 1"},
		"Problems":  {press("missing return"), "selected missing return"},
		"Outline":   {press("Serve"), "selected Serve"},
		"Minimap": {func(tt *ui.Tester) error {
			r, ok := tt.Find("Editor map")
			if !ok {
				return errNoText
			}
			tt.ClickAt(r.X+r.W/2, r.Y+r.H-6)
			return nil
		}, "at "},
	}
	// Quiet demos only show; they have no events.
	quiet := map[string]bool{"LineHighlight": true, "Annotated": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "code" {
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
		found := false
		for _, l := range logs {
			found = found || strings.HasPrefix(l, act.want)
		}
		if !found {
			t.Errorf("%s logged %q, want a line starting %q", e.Name, logs, act.want)
		}
	}
	if seen != len(acts)+len(quiet) {
		t.Errorf("checked %d code demos, want %d", seen, len(acts)+len(quiet))
	}
}
