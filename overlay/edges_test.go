package overlay

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTourKeepsFocusInside(t *testing.T) {
	step, open, hits := 0, true, 0
	steps := []TourStep{{"Name it", ""}, {"Save it", ""}}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Gap(200).Children(func() {
			a := ui.ButtonBase(c).Label("Name field").Size(120, 32)
			if a.Clicked() {
				hits++
			}
			b := ui.ButtonBase(c).Label("Save button").Size(120, 32)
			Tour(c, "Intro", []ui.Element{a, b}, steps, &step, &open, TourOptions{})
		})
	}, 600, 600)
	settles(t, tt, "Tour")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if step != 1 || hits != 0 {
		t.Errorf("Tab, Enter: step %d, anchor pressed %d times", step, hits)
	}
}

func TestPromptDialogTakesValueOnOpen(t *testing.T) {
	value, open := "old", false
	tt := ui.NewTester(func(c *ui.Context) {
		PromptDialog(c, "Rename", &value, &open, PromptDialogOptions{})
	}, 600, 400)
	value, open = "new", true
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if value != "new" || open {
		t.Errorf("Enter wrote %q back, open %v", value, open)
	}
}

func TestContextMenuStaysInWindow(t *testing.T) {
	list := []MenuItem{{"Rename", "", ""}, {"Delete", "", ""}}
	chosen, open := -1, false
	tt := ui.NewTester(func(c *ui.Context) {
		target := ui.ButtonBase(c).Label("Canvas").Size(600, 400)
		ContextMenu(c, "Canvas actions", target, list, &chosen, &open)
	}, 600, 400)
	tt.RightClickAt(590, 390)
	settles(t, tt, "corner menu")
	r := find(t, tt, "Canvas actions")
	if r.X < 0 || r.Y < 0 || r.X+r.W > 600 || r.Y+r.H > 400 || r.W < 60 || r.H < 60 {
		t.Fatalf("menu at %+v, not whole inside the window", r)
	}
	tt.RightClickAt(100, 100)
	settles(t, tt, "moved menu")
	if r := find(t, tt, "Canvas actions"); !open || r.X < 95 || r.X > 110 || r.Y < 95 || r.Y > 110 {
		t.Fatalf("right click again: open %v, menu at %+v", open, r)
	}
	click(t, tt, "Delete")
	if chosen != 1 || open {
		t.Errorf("pick: chosen %d, open %v", chosen, open)
	}
}

func TestSearchReachesEveryMatch(t *testing.T) {
	var commands []MenuItem
	for i := range 12 {
		commands = append(commands, MenuItem{Label: fmt.Sprint("Open file ", i)})
	}
	commands = append(commands, MenuItem{Label: "Open"})
	chosen, open := -1, false
	tt := ui.NewTester(func(c *ui.Context) {
		CommandPalette(c, "Commands", commands, &chosen, &open, CommandPaletteOptions{})
	}, 600, 500)
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Frame()
	tt.Type("open")
	for range 12 {
		tt.Key(0, ui.KeyDown)
	}
	settles(t, tt, "scrolled list")
	if r, ok := tt.Find("Open"); !ok || r.H < 20 || r.Y+r.H > 500 {
		t.Fatalf("last match not in view: %+v, found %v", r, ok)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 12 {
		t.Errorf("chose %d, want the last match", chosen)
	}
}

func TestSearchRevealsAfterFilter(t *testing.T) {
	var commands []MenuItem
	for i := range 40 {
		name := "Open file "
		if i >= 30 {
			name = "Save file "
		}
		commands = append(commands, MenuItem{Label: fmt.Sprint(name, i)})
	}
	chosen, open := -1, false
	tt := ui.NewTester(func(c *ui.Context) {
		CommandPalette(c, "Commands", commands, &chosen, &open, CommandPaletteOptions{})
	}, 600, 500)
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Frame()
	for range 30 {
		tt.Key(0, ui.KeyDown)
	}
	tt.Type("open")
	settles(t, tt, "filtered list")
	if r, ok := tt.Find("Open file 0"); !ok || r.H < 20 {
		t.Fatalf("highlighted row hidden: %+v, found %v", r, ok)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 0 {
		t.Errorf("chose %d", chosen)
	}
}

func TestPaletteForgetsHighlight(t *testing.T) {
	commands := []MenuItem{{"Open file", "", ""}, {"Save file", "", ""}}
	chosen, open := -1, false
	tt := ui.NewTester(func(c *ui.Context) {
		CommandPalette(c, "Commands", commands, &chosen, &open, CommandPaletteOptions{})
	}, 600, 500)
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Frame()
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 0 {
		t.Errorf("reopened palette kept its highlight: chose %d", chosen)
	}
}
