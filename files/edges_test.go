package files

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestTrashTellsFoldersApart(t *testing.T) {
	at := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	items := []FileItem{
		{Name: "report.pdf", Size: 1, Modified: at, Path: "Finance"},
		{Name: "report.pdf", Size: 2, Modified: at, Path: "Marketing"},
	}
	restored := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Trash(c, "Trash", items, &restored, TrashOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Finance") || !tt.HasText("Marketing") {
		t.Fatalf("texts %q", tt.Texts())
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) { Trash(c, "x", []FileItem{items[0], items[0]}, &restored, TrashOptions{}) }, 100, 100, 1)
	})
}

func TestUploadedNeverOverflows(t *testing.T) {
	for _, u := range []UploadItem{
		{Name: "a", Size: math.MaxInt64, Progress: 1},
		{Name: "b", Size: math.MaxInt64, Progress: 0.5},
		{Name: "c", Size: 10, Progress: 0.37},
	} {
		got := uploaded(u)
		if got < 0 || got > u.Size || u.Progress == 1 && got != u.Size {
			t.Errorf("uploaded(%+v) = %d", u, got)
		}
	}
	uploads := []UploadItem{{Name: "disk.img", Size: math.MaxInt64, Progress: 1}}
	render := func(c *ui.Context) { UploadList(c, "Uploads", &uploads, UploadListOptions{}) }
	ui.Render(render, 500, 300, 1)
}

func TestRecentFilesSelectionAnimates(t *testing.T) {
	opened := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { RecentFiles(c, "Recent", sampleFiles(), &opened, RecentFilesOptions{}) }), 500, 250)
	settles(t, tt)
	opened = 1
	tt.Frame()
	if !tt.FramePending() {
		t.Error("the selection color snaps")
	}
}
