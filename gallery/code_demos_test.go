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
		"CodeBlock": {press("Copy"), "copied 349 runes"},
		"CodeEditor": {func(tt *ui.Tester) error {
			tt.Key(0, ui.KeyTab)
			tt.Type("x")
			return nil
		}, "edited"},
		"DiffView": {press("Split"), "layout split"},
		"Terminal": {press("Run vet"), "6 lines"},
		"CommandLine": {func(tt *ui.Tester) error {
			if err := tt.Click("Run a command"); err != nil {
				return err
			}
			tt.Type("pwd")
			tt.Key(0, ui.KeyEnter)
			return nil
		}, "ran 3 runes"},
		"CodeTabs":     {press("price.go"), "selected price.go"},
		"FileHeader":   {press("Copy path"), "copied"},
		"Snippet":      {press("Copy"), "copied 40 runes"},
		"StackTrace":   {press("Show 4 library frames"), "open true"},
		"ErrorOverlay": {press("Close"), "open false"},
		"Gutter":       {press("Breakpoint on line 3"), "2 breakpoints"},
		"DebugBar":     {press("Step over"), "over"},
		"Variables":    {press("ws"), "toggled"},
		"Watches": {func(tt *ui.Tester) error {
			if err := tt.Click("Watch an expression"); err != nil {
				return err
			}
			tt.Type("x")
			tt.Key(0, ui.KeyEnter)
			return nil
		}, "4 watches"},
		"CallStack":    {press("net/http.HandlerFunc.ServeHTTP"), "selected 2"},
		"Problems":     {press("integer division truncates the discount"), "selected integer division truncates the discount"},
		"Outline":      {press("Close"), "selected Close"},
		"Breadcrumb":   {press("gallery"), "at gallery"},
		"SymbolSearch": {press("runServer"), "selected 0"},
		"RegexTester": {func(tt *ui.Tester) error {
			if err := tt.Click("Pattern"); err != nil {
				return err
			}
			tt.Key(0, ui.KeyHome)
			tt.TypeKey(0, ui.KeyX, "x")
			return nil
		}, "edited"},
		"JSONFormat": {press("Format"), "action on"},
		"Base64Tool": {press("Decode"), "encode false"},
		"HashTool": {func(tt *ui.Tester) error {
			if err := tt.Click("SHA-256"); err != nil {
				return err
			}
			tt.Type("x")
			return nil
		}, "edited"},
		"UnixTime": {func(tt *ui.Tester) error {
			if err := tt.Click("Unix seconds"); err != nil {
				return err
			}
			tt.Key(ui.Ctrl, ui.KeyA)
			tt.Type("1704067200")
			return nil
		}, "epoch 1704067200"},
		"Commit":       {press("Copy hash"), "copied hash"},
		"BranchPicker": {press("dev"), "selected dev"},
		"Console":      {press("Error 1"), "filter 3"},
		"REPL": {func(tt *ui.Tester) error {
			if err := tt.Click("1 + 2 * 3"); err != nil {
				return err
			}
			tt.Type("2+2")
			tt.Key(0, ui.KeyEnter)
			return nil
		}, "2 lines"},
		"PackageCard": {press("Copy install"), "copied install"},
		"Env":         {press("Show values"), "values shown"},
		"Playground":  {press("Run"), "ran 32 runes"},
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
	quiet := map[string]bool{"LineHighlight": true, "Annotated": true, "PullRequest": true, "Checks": true, "Blame": true, "CommitGraph": true, "Language": true, "KeyTable": true}
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
