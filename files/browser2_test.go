package files

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDownloadRowTimeLeft(t *testing.T) {
	opened := false
	item := FileItem{Name: "dataset.zip", Size: 10_000_000}
	tt := ui.NewTester(padded(func(c *ui.Context) { DownloadRow(c, item, 0.5, 1_000_000, &opened, DownloadRowOptions{}) }), 500, 200)
	settles(t, tt)
	if !tt.HasText("5s left") || tt.HasText("Open") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestColumnBrowserHidesUnchosenDepths(t *testing.T) {
	columns := [][]string{{"Finance", "Marketing"}, {"Reports", "Receipts"}, {"Q3.pdf"}}
	selected := []int{0, -1, -1}
	tt := ui.NewTester(padded(func(c *ui.Context) { ColumnBrowser(c, "Browse", columns, &selected) }), 600, 300)
	settles(t, tt)
	if !tt.HasText("Reports") || tt.HasText("Q3.pdf") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "Reports")
	settles(t, tt)
	if selected[1] != 0 || !tt.HasText("Q3.pdf") {
		t.Errorf("selected %v texts %q", selected, tt.Texts())
	}
	click(t, tt, "Marketing")
	settles(t, tt)
	if selected[1] != -1 || tt.HasText("Q3.pdf") {
		t.Errorf("deeper columns stayed: %v %q", selected, tt.Texts())
	}
}

func TestConflictShowsBothCopies(t *testing.T) {
	at := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	mine := FileItem{Name: "budget.xlsx", Size: 120_000, Modified: at}
	theirs := FileItem{Name: "budget.xlsx", Size: 98_000, Modified: at.Add(-time.Hour)}
	choice := ConflictNone
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Conflict(c, "budget.xlsx changed twice", mine, theirs, &choice, ConflictOptions{})
	}), 500, 300)
	settles(t, tt)
	if !tt.HasText("Your copy") || !tt.HasText("Their copy") || !tt.HasText("Oct 24, 08:00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestFileTagsAndThumbsMove(t *testing.T) {
	flags := []bool{false, false}
	tags := ui.NewTester(padded(func(c *ui.Context) { FileTags(c, "Tags", []string{"Invoices", "Urgent"}, &flags) }), 500, 150)
	settles(t, tags)
	click(t, tags, "Urgent")
	if !tags.FramePending() {
		t.Error("the tag check does not pop")
	}
	zoom, chosen := 1.0, -1
	thumbs := ui.NewTester(padded(func(c *ui.Context) { Thumbs(c, "Thumbs", sampleFiles(), &zoom, &chosen, ThumbsOptions{}) }), 500, 400)
	settles(t, thumbs)
	before := find(t, thumbs, "a.go")
	zoom = 2
	thumbs.Frame()
	if !thumbs.FramePending() {
		t.Error("tiles jump to the new zoom")
	}
	settles(t, thumbs)
	if after := find(t, thumbs, "a.go"); after.W <= before.W {
		t.Errorf("tile did not grow: %v -> %v", before.W, after.W)
	}
}
