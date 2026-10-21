package choice

import (
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestDropdown(t *testing.T) {
	plan, changes := "Free", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Dropdown(c, "Plan", []string{"Free", "Pro", "Team"}, &plan).Changed() {
			changes++
		}
	}, 300, 260)
	click(t, tt, "Plan")
	settles(t, tt, "Dropdown popup")
	before := baseline()
	click(t, tt, "Team")
	if plan != "Team" || changes != 1 || tt.HasText("Pro") {
		t.Fatalf("plan %q changes %d texts %q", plan, changes, tt.Texts())
	}
	cues(t, before, sound.Select)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if plan != "Pro" || changes != 2 {
		t.Errorf("keys chose %q, changes %d", plan, changes)
	}
}

func TestMultiDropdown(t *testing.T) {
	var sel []string
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if MultiDropdown(c, "Labels", []string{"Bug", "Docs", "UI"}, &sel).Changed() {
			changes++
		}
	}, 300, 260)
	click(t, tt, "Labels")
	before := baseline()
	click(t, tt, "Docs")
	click(t, tt, "UI")
	if !slices.Equal(sel, []string{"Docs", "UI"}) || changes != 2 {
		t.Fatalf("selection %v changes %d", sel, changes)
	}
	if !tt.HasText("2") || !tt.HasText("Bug") {
		t.Errorf("count or open popup missing: %q", tt.Texts())
	}
	cues(t, before, sound.Toggle, sound.Toggle)
	settles(t, tt, "MultiDropdown count")
}

func TestPopupFollowsTriggerWidth(t *testing.T) {
	var sel []string
	lang, w := "Go", float32(200)
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(w).Gap(8).Children(func() {
			MultiDropdown(c, "Labels", []string{"Bug", "Docs"}, &sel).FillWidth()
		})
	}, 600, 300)
	click(t, tt, "Labels")
	tt.Move(590, 290)
	tt.Advance(time.Second)
	w = 400
	tt.Frame()
	if r := find(t, tt, "Docs"); r.W < 380 {
		t.Errorf("MultiDropdown row %v after widening to 400", r)
	}
	click(t, tt, "Labels")
	tt2 := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(w).Children(func() {
			SearchDropdown(c, "Language", []string{"Go", "Rust"}, &lang, SearchDropdownOptions{}).FillWidth()
		})
	}, 600, 300)
	w = 200
	tt2.Frame()
	click(t, tt2, "Language")
	tt2.Move(590, 290)
	tt2.Advance(time.Second)
	w = 450
	tt2.Frame()
	if r := find(t, tt2, "Rust"); r.W < 430 {
		t.Errorf("SearchDropdown row %v after widening to 450", r)
	}
}

func TestGroupedDropdown(t *testing.T) {
	city := "Paris"
	tt := ui.NewTester(func(c *ui.Context) {
		GroupedDropdown(c, "City", []DropdownGroup{{"Europe", []string{"Paris", "Rome"}}, {"Asia", []string{"Tokyo"}}}, &city)
	}, 300, 300)
	click(t, tt, "City")
	find(t, tt, "Europe")
	find(t, tt, "Asia")
	before := baseline()
	click(t, tt, "Tokyo")
	if city != "Tokyo" {
		t.Fatalf("city %q", city)
	}
	cues(t, before, sound.Select)
	panics(t, "repeated option", func() {
		ui.Render(func(c *ui.Context) {
			GroupedDropdown(c, "x", []DropdownGroup{{"A", []string{"P"}}, {"B", []string{"P"}}}, &city)
		}, 50, 50, 1)
	})
}

func TestSearchDropdown(t *testing.T) {
	lang, changes := "Go", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SearchDropdown(c, "Language", []string{"Go", "Rust", "Ruby"}, &lang, SearchDropdownOptions{Empty: "Nothing"}).Changed() {
			changes++
		}
	}, 300, 300)
	click(t, tt, "Language")
	tt.Type("ru")
	find(t, tt, "Ruby")
	before := baseline()
	tt.Key(0, ui.KeyEnter)
	if lang != "Rust" || changes != 1 {
		t.Fatalf("lang %q changes %d", lang, changes)
	}
	cues(t, before, sound.Select)
	click(t, tt, "Language")
	tt.Type("zz")
	find(t, tt, "Nothing")
	settles(t, tt, "SearchDropdown empty note")
}

func TestTreeDropdown(t *testing.T) {
	item := "Docs"
	tree := []TreeItem{{Label: "Docs", Children: []TreeItem{{Label: "Guides"}, {Label: "API"}}}, {Label: "Blog"}}
	tt := ui.NewTester(func(c *ui.Context) { TreeDropdown(c, "Section", tree, &item, TreeDropdownOptions{}) }, 300, 300)
	click(t, tt, "Section")
	if tt.HasText("Guides") {
		t.Fatal("children shown before expanding")
	}
	before := baseline()
	click(t, tt, "Expand Docs")
	settles(t, tt, "TreeDropdown rows")
	find(t, tt, "Collapse Docs")
	click(t, tt, "API")
	if item != "API" {
		t.Fatalf("item %q", item)
	}
	cues(t, before, sound.Tick, sound.Select)
}

func TestCascade(t *testing.T) {
	tree := []TreeItem{{Label: "France", Children: []TreeItem{{Label: "Paris"}, {Label: "Lyon"}}}, {Label: "Japan", Children: []TreeItem{{Label: "Tokyo"}}}}
	var path []string
	tt := ui.NewTester(func(c *ui.Context) { Cascade(c, "Place", tree, &path) }, 500, 200)
	before := baseline()
	click(t, tt, "Japan")
	settles(t, tt, "Cascade column")
	click(t, tt, "Tokyo")
	if !slices.Equal(path, []string{"Japan", "Tokyo"}) {
		t.Fatalf("path %v", path)
	}
	click(t, tt, "France")
	if !slices.Equal(path, []string{"France"}) || !tt.HasText("Lyon") || tt.HasText("Tokyo") {
		t.Fatalf("path %v texts %q", path, tt.Texts())
	}
	cues(t, before, sound.Select, sound.Select, sound.Select)
	panics(t, "path off the tree", func() {
		bad := []string{"Mars"}
		ui.Render(func(c *ui.Context) { Cascade(c, "x", tree, &bad) }, 50, 50, 1)
	})
}

func TestListbox(t *testing.T) {
	fruit := "Apple"
	tt := ui.NewTester(func(c *ui.Context) { Listbox(c, "Fruit", []string{"Apple", "Pear", "Plum"}, &fruit) }, 240, 200)
	before := baseline()
	click(t, tt, "Plum")
	if fruit != "Plum" {
		t.Fatalf("fruit %q", fruit)
	}
	settles(t, tt, "Listbox highlight")
	tt.Key(0, ui.KeyUp)
	if fruit != "Pear" {
		t.Errorf("arrow up gave %q", fruit)
	}
	cues(t, before, sound.Select, sound.Select)
}

func TestTransfer(t *testing.T) {
	left, right := []string{"Alpha", "Beta"}, []string{"Gamma"}
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Transfer(c, &left, &right, TransferOptions{RightTitle: "Picked"}).Changed() {
			changes++
		}
	}, 500, 260)
	find(t, tt, "Available")
	find(t, tt, "Picked")
	before := baseline()
	a := find(t, tt, "Alpha")
	click(t, tt, "Alpha")
	if !slices.Equal(left, []string{"Beta"}) || !slices.Equal(right, []string{"Gamma", "Alpha"}) || changes != 1 {
		t.Fatalf("left %v right %v changes %d", left, right, changes)
	}
	if b := find(t, tt, "Alpha"); b.X != a.X {
		t.Errorf("flight starts at %v, want %v", b, a)
	}
	settles(t, tt, "Transfer flight")
	if b := find(t, tt, "Alpha"); b.X <= a.X+100 {
		t.Errorf("Alpha landed at %v", b)
	}
	click(t, tt, "Gamma")
	if !slices.Equal(left, []string{"Beta", "Gamma"}) {
		t.Errorf("left %v", left)
	}
	cues(t, before, sound.Drop, sound.Drop)
}
