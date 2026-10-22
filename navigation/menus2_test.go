package navigation

import (
	"math"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestIconToolbar(t *testing.T) {
	used := [2]int{-1, -1}
	groups := [][]NavItem{{{"Bold", "bold"}, {"Italic", "italic"}}, {{"Link", "link"}}}
	tt := ui.NewTester(func(c *ui.Context) { IconToolbar(c, "Format", groups, &used) }, 400, 100)
	before := baseline()
	click(t, tt, "Link")
	click(t, tt, "Italic")
	if used != [2]int{0, 1} {
		t.Errorf("used %v", used)
	}
	cues(t, before, sound.Tap, sound.Tap)
	panics(t, "empty group", func() { render(func(c *ui.Context) { u := [2]int{}; IconToolbar(c, "x", [][]NavItem{{}}, &u) }) })
}

func TestDockBar(t *testing.T) {
	opened, changed := -1, 0
	apps := []NavItem{{"Files", "folder"}, {"Mail", "mail"}}
	tt := ui.NewTester(func(c *ui.Context) {
		if DockBar(c, "Dock", apps, &opened).Changed() {
			changed++
		}
	}, 400, 120)
	before := baseline()
	click(t, tt, "Mail")
	if opened != 1 || changed != 1 {
		t.Errorf("opened %d, changed %d", opened, changed)
	}
	cues(t, before, sound.Tap)
	panics(t, "bad opened", func() { render(func(c *ui.Context) { o := 2; DockBar(c, "x", apps, &o) }) })
}

// menuCase drives one popover menu: open, then pick.
type menuCase struct {
	name, trigger, pick string
	view                func(c *ui.Context) ui.Element
	pickCue             sound.Cue
	picked              func() bool
}

func TestMenusPick(t *testing.T) {
	apps, appsOpen, app := []NavItem{{"Mail", "mail"}, {"Calendar", "calendar"}}, false, -1
	spaces, spacesOpen, space := []string{"Acme", "Globex"}, false, 0
	accountOpen, item := false, -1
	langs, langOpen, lang := []string{"English", "Deutsch"}, false, 0
	for _, m := range []menuCase{
		{"AppSwitcher", "Apps", "Calendar", func(c *ui.Context) ui.Element { return AppSwitcher(c, "Apps", apps, &app, &appsOpen) }, sound.Select, func() bool { return app == 1 && !appsOpen }},
		{"WorkspaceSwitcher", "Workspace", "Globex", func(c *ui.Context) ui.Element {
			return WorkspaceSwitcher(c, "Workspace", spaces, &space, &spacesOpen)
		}, sound.Select, func() bool { return space == 1 && !spacesOpen }},
		{"AccountNav", "Ada Lovelace", "Sign out", func(c *ui.Context) ui.Element {
			return AccountNav(c, "Ada Lovelace", nil, []string{"Profile", "Sign out"}, &item, &accountOpen)
		}, sound.Tap, func() bool { return item == 1 && !accountOpen }},
		{"LanguageMenu", "Language", "Deutsch", func(c *ui.Context) ui.Element { return LanguageMenu(c, "Language", langs, &lang, &langOpen) }, sound.Select, func() bool { return lang == 1 && !langOpen }},
	} {
		changed := 0
		tt := ui.NewTester(func(c *ui.Context) {
			if m.view(c).Changed() {
				changed++
			}
		}, 500, 400)
		before := baseline()
		click(t, tt, m.trigger)
		if changed != 1 || !tt.HasText(m.pick) {
			t.Fatalf("%s: open changed %d, %q", m.name, changed, tt.Texts())
		}
		click(t, tt, m.pick)
		settles(t, tt, m.name)
		if !m.picked() || changed != 2 {
			t.Errorf("%s: pick not reported, changed %d", m.name, changed)
		}
		cues(t, before, sound.Open, m.pickCue)
	}
	tw := ui.NewTester(func(c *ui.Context) { WorkspaceSwitcher(c, "W", spaces, &space, &spacesOpen) }, 300, 100)
	if !tw.HasText("Globex") {
		t.Error("the trigger does not show the picked workspace")
	}
}

func TestSortMenu(t *testing.T) {
	key, desc, open, changed := 0, false, false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SortMenu(c, "Sort", []string{"Name", "Date"}, &key, &desc, &open, SortMenuOptions{}).Changed() {
			changed++
		}
	}, 400, 400)
	before := baseline()
	click(t, tt, "Sort")
	click(t, tt, "Descending")
	if !desc || open || !tt.FramePending() {
		t.Fatalf("desc %v, open %v, turning %v", desc, open, tt.FramePending())
	}
	settles(t, tt, "turned")
	click(t, tt, "Sort")
	click(t, tt, "Date")
	if key != 1 || open || changed != 4 {
		t.Errorf("key %d, open %v, changed %d", key, open, changed)
	}
	cues(t, before, sound.Open, sound.Select, sound.Open, sound.Select)
}

func TestZoomControl(t *testing.T) {
	zoom := 1.0
	tt := ui.NewTester(func(c *ui.Context) { ZoomControl(c, "Zoom", &zoom, ZoomControlOptions{}) }, 400, 100)
	before := baseline()
	click(t, tt, "Fit")
	click(t, tt, "Zoom in")
	if zoom != 1.25 || !tt.FramePending() {
		t.Fatalf("zoom %v, rolling %v", zoom, tt.FramePending())
	}
	click(t, tt, "Zoom out")
	click(t, tt, "Zoom out")
	click(t, tt, "Fit")
	if zoom != 1 {
		t.Errorf("zoom %v", zoom)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	for name, z := range map[string]float64{"too far": 5, "NaN": math.NaN()} {
		panics(t, name, func() { render(func(c *ui.Context) { ZoomControl(c, "x", &z, ZoomControlOptions{}) }) })
	}
	fit, fitChanges := 1.0, 0
	tf := ui.NewTester(func(c *ui.Context) {
		if ZoomControl(c, "Zoom", &fit, ZoomControlOptions{FitTo: 0.9}).Changed() {
			fitChanges++
		}
	}, 400, 100)
	click(t, tf, "Fit")
	click(t, tf, "Fit")
	if fit != 0.9 || fitChanges != 1 {
		t.Errorf("fit %v, changes %d", fit, fitChanges)
	}
	panics(t, "infinite step", func() {
		render(func(c *ui.Context) { z := 1.0; ZoomControl(c, "x", &z, ZoomControlOptions{Step: math.Inf(1)}) })
	})
}

func TestFilterBar(t *testing.T) {
	var selected []string
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if FilterBar(c, "Status", []string{"Open", "Closed"}, &selected, FilterBarOptions{}).Changed() {
			changed++
		}
	}, 500, 100)
	if tt.HasText("Clear all") {
		t.Fatal("clear shows with nothing selected")
	}
	before := baseline()
	click(t, tt, "Open")
	if len(selected) != 1 || !tt.HasText("Clear all") {
		t.Fatalf("selected %v, %q", selected, tt.Texts())
	}
	click(t, tt, "Clear all")
	settles(t, tt, "cleared")
	if len(selected) != 0 || tt.HasText("Clear all") || changed != 2 {
		t.Errorf("selected %v, changed %d, %q", selected, changed, tt.Texts())
	}
	cues(t, before, sound.Toggle, sound.Tap)
}

func TestViewSwitch(t *testing.T) {
	view := 0
	tt := ui.NewTester(func(c *ui.Context) { ViewSwitch(c, "View", &view, ViewSwitchOptions{}) }, 300, 100)
	before := baseline()
	click(t, tt, "Board")
	if view != 2 {
		t.Errorf("view %d", view)
	}
	cues(t, before, sound.Select)
}
