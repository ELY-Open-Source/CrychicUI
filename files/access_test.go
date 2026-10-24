package files

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestShareDialog(t *testing.T) {
	people := []SharePerson{{Name: "Ana Reyes", Email: "ana@acme.dev", Role: "Owner"}, {Name: "Kai Mori", Role: "Viewer"}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ShareDialog(c, "Share", people, "https://example.test/f/1", ShareDialogOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Copy link")
	if tt.Clipboard() != "https://example.test/f/1" {
		t.Errorf("clipboard %q", tt.Clipboard())
	}
	cues(t, before, sound.Pop)
	if !tt.HasText("ana@acme.dev") || !tt.HasText("Anyone with the link can view") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestPermissions(t *testing.T) {
	role := Viewer
	tt := ui.NewTester(padded(func(c *ui.Context) { Permissions(c, "Access", &role, PermissionsOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Editor")
	if role != Editor || !tt.HasText("Can edit and share") {
		t.Errorf("role %v texts %q", role, tt.Texts())
	}
	cues(t, before, sound.Select)
}

func TestTrash(t *testing.T) {
	restored := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Trash(c, "Trash", []FileItem{{Name: "old.txt", Size: 3, Modified: time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC), Path: "Documents"}}, &restored, TrashOptions{})
	}), 500, 200)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Restore old.txt")
	if restored != 0 || !tt.HasText("Deleted Oct 20, 2026 · Documents") {
		t.Errorf("restored %d texts %q", restored, tt.Texts())
	}
	cues(t, before, sound.Success)
	empty := ui.NewTester(padded(func(c *ui.Context) { Trash(c, "Trash", nil, &[]int{-1}[0], TrashOptions{}) }), 500, 200)
	settles(t, empty)
	if !empty.HasText("Trash is empty") {
		t.Errorf("empty texts %q", empty.Texts())
	}
}

func TestProperties(t *testing.T) {
	now := time.Date(2026, 10, 23, 0, 0, 0, 0, time.UTC)
	item := FileItem{Name: "note.txt", Size: 42, Modified: now}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Properties(c, item, "/home/note.txt", "Text", now.Add(-time.Hour), PropertiesOptions{})
	}), 500, 300)
	settles(t, tt)
	for _, text := range []string{"42 B", "42 bytes", "Text", "/home/note.txt", "Oct 23, 2026 at 00:00"} {
		if !tt.HasText(text) {
			t.Errorf("missing %q in %q", text, tt.Texts())
		}
	}
}

func TestVersions(t *testing.T) {
	restored := -1
	versions := []VersionItem{
		{ID: "v2", Author: "Ana Reyes", Size: 20, Modified: time.Now(), Current: true},
		{ID: "v1", Author: "Kai Mori", Size: 10, Modified: time.Now().Add(-time.Hour)},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Versions(c, "Versions", versions, &restored, VersionsOptions{})
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Restore v1")
	if restored != 1 || !tt.HasText("Current") || !tt.HasText("Kai Mori") {
		t.Errorf("restored %d texts %q", restored, tt.Texts())
	}
	cues(t, before, sound.Success)
}

func TestRecentFiles(t *testing.T) {
	opened := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { RecentFiles(c, "Recent", sampleFiles(), &opened, RecentFilesOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "a.go")
	if opened != 1 {
		t.Errorf("opened %d", opened)
	}
	cues(t, before, sound.Tap)
}
