package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestFiles2DemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"DownloadRow":   {press("Open"), "opened download"},
		"ZipList":       {press("main.go"), "open 0"},
		"ImageFiles":    {press("two.png"), "selected 1"},
		"FileSearch":    {press(".go"), "kind .go"},
		"ColumnBrowser": {press("work"), "[1 -1]"},
		"FileCrumbs": {func(tt *ui.Tester) error {
			if err := tt.Click("docs"); err != nil {
				return err
			}
			return tt.ChooseMenuItem("Open")
		}, "opened 2"},
		"QuotaWarning": {press("Upgrade"), "upgraded"},
		"Starred":      {press("Unstar notes.md"), "false"},
		"Conflict":     {press("Keep both"), "choice 3"},
		"SelectBar":    {press("Delete"), "action 2"},
		"FileTags":     {press("urgent"), "true true false"},
		"UploadButtonRow": {func(tt *ui.Tester) error {
			if err := tt.Click("New"); err != nil {
				return err
			}
			return tt.ChooseMenuItem("New folder")
		}, "action 2"},
		"FileTypeFilter": {press("Code 8"), "selected 1"},
		"Thumbs":         {press("main.go"), "selected 1"},
	}
	quiet := map[string]bool{"Storage": true, "SyncRow": true, "Compress": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "files" || e.ID < 959 {
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
	if seen != 17 {
		t.Errorf("checked %d second files demos", seen)
	}
}
