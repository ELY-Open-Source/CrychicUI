package files

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestFileCard(t *testing.T) {
	item := FileItem{Name: "note.txt", Size: 12, Modified: time.Now()}
	changed := 0
	action := FileActionNone
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _ := FileCard(c, item, &action, FileCardOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "File actions")
	if err := tt.ChooseMenuItem("Rename"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if changed != 1 || action != FileActionRename {
		t.Errorf("changed %d action %v", changed, action)
	}
	cues(t, before, sound.Tap)
}

func TestFolderCard(t *testing.T) {
	opened := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		FolderCard(c, "Documents", 8, &opened, FolderCardOptions{})
	}), 500, 200)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Documents")
	if !opened {
		t.Error("folder did not open")
	}
	cues(t, before, sound.Open)
	one := false
	single := ui.NewTester(padded(func(c *ui.Context) {
		FolderCard(c, "One", 1, &one, FolderCardOptions{})
	}), 500, 200)
	if !single.HasText("1 item") || single.HasText("1 items") {
		t.Errorf("singular count %q", single.Texts())
	}
}

func TestPathBar(t *testing.T) {
	at := 2
	tt := ui.NewTester(padded(func(c *ui.Context) { PathBar(c, []string{"home", "work", "docs"}, &at) }), 500, 200)
	settles(t, tt)
	before := baseline()
	click(t, tt, "work")
	if at != 1 {
		t.Errorf("at %d", at)
	}
	cues(t, before, sound.Tap)
}

func TestNewFolder(t *testing.T) {
	name := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		NewFolder(c, &name, NewFolderOptions{})
	}), 500, 200)
	settles(t, tt)
	click(t, tt, "Folder name")
	tt.Type("Assets")
	before := baseline()
	click(t, tt, "Create")
	if name != "Assets" {
		t.Errorf("name %q", name)
	}
	cues(t, before, sound.Pop)
}

func TestNewFolderEnter(t *testing.T) {
	name := ""
	submitted := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, create := NewFolder(c, &name, NewFolderOptions{})
		if create.Changed() {
			submitted++
		}
	}), 500, 200)
	settles(t, tt)
	click(t, tt, "Folder name")
	tt.Type("Assets")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if submitted != 1 {
		t.Errorf("Enter submissions %d", submitted)
	}
}

func TestRenameRow(t *testing.T) {
	name := "old.txt"
	tt := ui.NewTester(padded(func(c *ui.Context) { RenameRow(c, &name) }), 500, 200)
	settles(t, tt)
	r, ok := tt.Find("old.txt")
	if !ok {
		t.Fatal("old name missing")
	}
	before := baseline()
	tt.ClickAt(r.X+5, r.Y+5)
	tt.ClickAt(r.X+5, r.Y+5)
	tt.Type("new")
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if name != "new.txt" {
		t.Errorf("name %q", name)
	}
	cues(t, before, sound.Success)
}

func TestRenameRowRejectsEmpty(t *testing.T) {
	name := "old.txt"
	tt := ui.NewTester(padded(func(c *ui.Context) { RenameRow(c, &name) }), 500, 200)
	settles(t, tt)
	r, ok := tt.Find("old.txt")
	if !ok {
		t.Fatal("old name missing")
	}
	tt.ClickAt(r.X+5, r.Y+5)
	tt.ClickAt(r.X+5, r.Y+5)
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if name != "old.txt" {
		t.Errorf("empty rename kept %q", name)
	}
}

func TestRenameRowExternalUpdate(t *testing.T) {
	name := "old.txt"
	tt := ui.NewTester(padded(func(c *ui.Context) { RenameRow(c, &name) }), 500, 200)
	settles(t, tt)
	name = "new.txt"
	tt.Frame()
	r, ok := tt.Find("new.txt")
	if !ok {
		t.Fatal("new name missing")
	}
	tt.ClickAt(r.X+5, r.Y+5)
	tt.ClickAt(r.X+5, r.Y+5)
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if name != "new.txt" {
		t.Errorf("rollback restored %q", name)
	}
}

func TestMoveDialog(t *testing.T) {
	target := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MoveDialog(c, "Move file", []string{"home", "work/docs"}, &target, MoveDialogOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "work/docs")
	click(t, tt, "Move")
	if target != 1 {
		t.Errorf("target %d", target)
	}
	cues(t, before, sound.Select, sound.Drop)
}
