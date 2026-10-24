package gallery

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestFilesDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"FileList":   {press("Size"), "sort 1"},
		"FileGrid":   {press("billing.go"), "selected billing.go"},
		"UploadList": {press("Cancel office-floorplan.png"), "2 uploads left"},
		"FileCard": {func(tt *ui.Tester) error {
			if err := tt.Click("File actions"); err != nil {
				return err
			}
			return tt.ChooseMenuItem("Rename")
		}, "action 2 on Q3 board report.pdf"},
		"FolderCard": {press("Receipts"), "opened Receipts"},
		"PathBar":    {press("Finance"), "at Finance"},
		"NewFolder": {func(tt *ui.Tester) error {
			if err := tt.Click("Folder name"); err != nil {
				return err
			}
			tt.Type("Assets")
			return tt.Click("Create")
		}, "created 6 runes"},
		"RenameRow": {func(tt *ui.Tester) error {
			r, ok := tt.Find("draft.md")
			if !ok {
				return fmt.Errorf("draft missing")
			}
			tt.ClickAt(r.X+5, r.Y+5)
			tt.ClickAt(r.X+5, r.Y+5)
			tt.Type("final")
			tt.Key(0, ui.KeyTab)
			return nil
		}, "renamed 8 runes"},
		"MoveDialog": {func(tt *ui.Tester) error {
			if err := tt.Click("Acme/Finance/Reports"); err != nil {
				return err
			}
			tt.Advance(time.Second)
			return tt.Click("Move")
		}, "moved to Acme/Finance/Reports"},
		"ShareDialog": {press("Copy link"), "copied link"},
		"Permissions": {press("Editor"), "role 1"},
		"Trash":       {press("Restore old-pricing.xlsx"), "restored old-pricing.xlsx"},
		"Versions":    {press("Restore v3"), "restored v3"},
		"RecentFiles": {press("billing.go"), "opened billing.go"},
	}
	quiet := map[string]bool{"FileIcon": true, "Dropzone": true, "FilePreview": true, "Properties": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "files" || e.ID >= 959 {
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
	if seen != 18 {
		t.Errorf("checked %d files demos", seen)
	}
}
