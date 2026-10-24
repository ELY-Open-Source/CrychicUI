package files

import (
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestFileListPutsFoldersFirstAndSlidesOnSort(t *testing.T) {
	now := time.Now()
	list := []FileItem{
		{Name: "budget.xlsx", Size: 9000, Modified: now},
		{Name: "Archive", Dir: true, Modified: now},
		{Name: "agenda.md", Size: 120, Modified: now.Add(-time.Hour)},
	}
	selected, sort := -1, SortName
	tt := ui.NewTester(padded(func(c *ui.Context) { FileList(c, "Files", list, &selected, &sort, FileListOptions{}) }), 600, 300)
	settles(t, tt)
	texts := tt.Texts()
	if a, b := slices.Index(texts, "Archive"), slices.Index(texts, "agenda.md"); a < 0 || b < 0 || a > b {
		t.Fatalf("folders do not lead: %q", texts)
	}
	click(t, tt, "Size")
	if !tt.FramePending() {
		t.Error("rows jump to their new order")
	}
	settles(t, tt)
	if a, b := find(t, tt, "agenda.md"), find(t, tt, "budget.xlsx"); a.Y > b.Y {
		t.Errorf("size order: agenda at %v, budget at %v", a.Y, b.Y)
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			FileList(c, "x", []FileItem{{Name: "a"}, {Name: "a"}}, &selected, &sort, FileListOptions{})
		}, 100, 100, 1)
	})
}

func TestNewFolderRefusesSlashes(t *testing.T) {
	name := ""
	submitted := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, create := NewFolder(c, &name, NewFolderOptions{}); create.Changed() {
			submitted++
		}
	}), 500, 200)
	settles(t, tt)
	click(t, tt, "Folder name")
	tt.Type("Q3/plans")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if submitted != 0 || !tt.HasText(`Names cannot contain / or \`) {
		t.Errorf("submitted %d texts %q", submitted, tt.Texts())
	}
}

func TestRenameRowGlowsOnSave(t *testing.T) {
	name := "draft.md"
	tt := ui.NewTester(padded(func(c *ui.Context) { RenameRow(c, &name) }), 500, 200)
	settles(t, tt)
	r := find(t, tt, "draft.md")
	px := func() [4]uint8 {
		p := tt.Image().RGBAAt(int(r.X+r.W)+40, int(r.Y+r.H/2))
		return [4]uint8{p.R, p.G, p.B, p.A}
	}
	rest := px()
	before := baseline()
	tt.ClickAt(r.X+5, r.Y+5)
	tt.ClickAt(r.X+5, r.Y+5)
	tt.Type("final")
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	tt.Advance(300 * time.Millisecond)
	if name != "final.md" || px() == rest {
		t.Errorf("name %q, no glow after the rename", name)
	}
	cues(t, before, sound.Success)
	settles(t, tt)
	if px() != rest {
		t.Error("the glow stays")
	}
}

func TestDropzoneCountsFiles(t *testing.T) {
	paths := []string{"/Users/ana/q3.pdf"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Dropzone(c, "Upload", &paths, DropzoneOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("1 file added") {
		t.Errorf("texts %q", tt.Texts())
	}
	paths = append(paths, "/Users/ana/q4.pdf")
	tt.Frame()
	settles(t, tt)
	if !tt.HasText("2 files added") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMoveDialogNamesTarget(t *testing.T) {
	target := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MoveDialog(c, "Move file", []string{"Projects", "Projects/Billing"}, &target, MoveDialogOptions{})
	}), 500, 300)
	settles(t, tt)
	click(t, tt, "Projects/Billing")
	settles(t, tt)
	if !tt.HasText("Move to Projects/Billing") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			MoveDialog(c, "x", []string{"a", "a"}, &target, MoveDialogOptions{})
		}, 100, 100, 1)
	})
}

func TestPathBarKeepsIdentityAcrossPaths(t *testing.T) {
	parts := []string{"Home", "Projects", "Billing"}
	at := 2
	tt := ui.NewTester(padded(func(c *ui.Context) { PathBar(c, "Path", parts, &at) }), 500, 200)
	settles(t, tt)
	r := find(t, tt, "Projects")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	parts = []string{"Home", "Photos", "2026"}
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
	if at != 2 {
		t.Errorf("a new folder under the pointer fired: at %d", at)
	}
}

// find is the rect of s, else a fatal error.
func find(t *testing.T, tt *ui.Tester, s string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(s)
	if !ok {
		t.Fatalf("%q not shown in %q", s, tt.Texts())
	}
	return r
}
