package files

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestShareDialog(t *testing.T) {
	people := []SharePerson{{Name: "Ana", Role: "Owner"}, {Name: "Kai", Role: "Viewer"}}
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
}

func TestPermissions(t *testing.T) {
	role := Viewer
	tt := ui.NewTester(padded(func(c *ui.Context) { Permissions(c, "Access", &role, PermissionsOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Editor")
	if role != Editor {
		t.Errorf("role %v", role)
	}
	cues(t, before, sound.Select)
}

func TestTrash(t *testing.T) {
	restored := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Trash(c, "Trash", []string{"old.txt"}, &restored, TrashOptions{})
	}), 500, 200)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Restore")
	if restored != 0 {
		t.Errorf("restored %d", restored)
	}
	cues(t, before, sound.Success)
}

func TestProperties(t *testing.T) {
	now := time.Date(2026, 10, 23, 0, 0, 0, 0, time.UTC)
	item := FileItem{Name: "note.txt", Size: 42, Modified: now}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Properties(c, item, "/home/note.txt", "Text", now.Add(-time.Hour), PropertiesOptions{})
	}), 500, 300)
	settles(t, tt)
	for _, text := range []string{"42 B", "Text", "/home/note.txt", "2026-10-23"} {
		if !tt.HasText(text) {
			t.Errorf("missing %q in %q", text, tt.Texts())
		}
	}
}

func TestVersions(t *testing.T) {
	restored := -1
	versions := []VersionItem{
		{ID: "v2", Size: 20, Modified: time.Now(), Current: true},
		{ID: "v1", Size: 10, Modified: time.Now().Add(-time.Hour)},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Versions(c, "Versions", versions, &restored, VersionsOptions{})
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Restore")
	if restored != 1 {
		t.Errorf("restored %d", restored)
	}
	cues(t, before, sound.Success)
}

func TestRecentFiles(t *testing.T) {
	opened := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { RecentFiles(c, "Recent", sampleFiles(), &opened) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "a.go")
	if opened != 1 {
		t.Errorf("opened %d", opened)
	}
	cues(t, before, sound.Tap)
}
