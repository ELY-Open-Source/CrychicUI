package navigation

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var site = []Node{{"Guides", []Node{{"Install", nil}, {"Deploy", nil}}}, {"Changelog", nil}}

func TestNavTreeAndFileNav(t *testing.T) {
	for name, build := range map[string]func(*ui.Context, *string, map[string]bool) ui.Element{
		"NavTree": func(c *ui.Context, s *string, o map[string]bool) ui.Element { return NavTree(c, "Docs", site, s, o) },
		"FileNav": func(c *ui.Context, s *string, o map[string]bool) ui.Element { return FileNav(c, "Docs", site, s, o) },
	} {
		chosen, open := "", map[string]bool{}
		tt := ui.NewTester(func(c *ui.Context) { build(c, &chosen, open) }, 400, 400)
		if tt.HasText("Install") {
			t.Fatalf("%s: closed folder shows children", name)
		}
		before := baseline()
		click(t, tt, "Guides")
		if !tt.HasText("Install") || !tt.FramePending() || !open["Guides"] {
			t.Fatalf("%s: did not open: %v, %q", name, open, tt.Texts())
		}
		settles(t, tt, name)
		click(t, tt, "Deploy")
		if chosen != "Guides/Deploy" {
			t.Errorf("%s: chose %q", name, chosen)
		}
		leaf := sound.Tap
		if name == "FileNav" {
			leaf = sound.Select
		}
		tt.Key(0, ui.KeyUp)
		tt.Key(0, ui.KeyUp)
		tt.Key(0, ui.KeyLeft)
		tt.Frame()
		settles(t, tt, "closing")
		if tt.HasText("Install") || open["Guides"] {
			t.Errorf("%s: Left did not close the folder: %v", name, open)
		}
		cues(t, before, sound.Toggle, leaf, sound.Toggle)
	}
	panics(t, "twin labels", func() {
		render(func(c *ui.Context) { s := ""; NavTree(c, "x", []Node{{"A", nil}, {"A", nil}}, &s, map[string]bool{}) })
	})
	for label, want := range map[string]string{"main.go": "file-code", "notes.md": "file-text", "logo.svg": "file-image", "LICENSE": "file"} {
		if got := fileIcon(Node{Label: label}, false); got != want {
			t.Errorf("fileIcon(%s) = %s", label, got)
		}
	}
	if fileIcon(Node{Label: "src", Children: []Node{}}, true) != "folder-open" {
		t.Error("an open folder is not folder-open")
	}
}

func TestAccordionNav(t *testing.T) {
	gs := []LinkGroup{{"Account", []string{"Profile", "Security"}}, {"Billing", []string{"Plans", "Invoices"}}}
	picked, opened, changed := [2]int{-1, -1}, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if AccordionNav(c, "Settings", gs, &picked, &opened).Changed() {
			changed++
		}
	}, 400, 400)
	if !tt.HasText("Profile") || tt.HasText("Plans") {
		t.Fatalf("first group should open: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Billing")
	settles(t, tt, "AccordionNav")
	if !tt.HasText("Plans") || tt.HasText("Profile") || opened != 1 || changed != 1 {
		t.Errorf("one at a time: open %d, changed %d, %q", opened, changed, tt.Texts())
	}
	click(t, tt, "Invoices")
	if picked != [2]int{1, 1} || changed != 2 {
		t.Errorf("picked %v, changed %d", picked, changed)
	}
	cues(t, before, sound.Toggle, sound.Tap)
}

func TestDrilldown(t *testing.T) {
	menu := []Node{{"Settings", []Node{{"Privacy", nil}, {"Sound", nil}}}, {"Help", nil}}
	chosen, changed := "", 0
	var trail []string
	tt := ui.NewTester(func(c *ui.Context) {
		if Drilldown(c, "Menu", menu, &chosen, &trail, DrilldownOptions{}).Changed() {
			changed++
		}
	}, 400, 400)
	before := baseline()
	click(t, tt, "Settings")
	if !tt.FramePending() || len(trail) != 1 || changed != 1 {
		t.Errorf("down: slid %v, trail %v, changed %d", tt.FramePending(), trail, changed)
	}
	settles(t, tt, "down")
	if !tt.HasText("Privacy") || tt.HasText("Help") {
		t.Fatalf("submenu %q", tt.Texts())
	}
	click(t, tt, "Privacy")
	if chosen != "Settings/Privacy" {
		t.Errorf("chose %q", chosen)
	}
	click(t, tt, "Back")
	settles(t, tt, "up")
	if !tt.HasText("Help") || len(trail) != 0 || changed != 3 {
		t.Errorf("back: trail %v, changed %d, %q", trail, changed, tt.Texts())
	}
	cues(t, before, sound.Tap, sound.Tap, sound.Tap)
	panics(t, "bad trail", func() {
		render(func(c *ui.Context) {
			s, tr := "", []string{"Help"}
			Drilldown(c, "x", menu, &s, &tr, DrilldownOptions{})
		})
	})
}
