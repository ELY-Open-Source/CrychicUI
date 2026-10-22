package navigation

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestTabsSlide(t *testing.T) {
	names := []string{"One", "Two", "Three"}
	for name, build := range map[string]func(*ui.Context, *int) ui.Element{
		"UnderlineTabs": func(c *ui.Context, a *int) ui.Element { return UnderlineTabs(c, "Tabs", names, a) },
		"PillTabs":      func(c *ui.Context, a *int) ui.Element { return PillTabs(c, "Tabs", names, a) },
		"VerticalTabs":  func(c *ui.Context, a *int) ui.Element { return VerticalTabs(c, "Tabs", names, a) },
	} {
		active := 0
		tt := ui.NewTester(func(c *ui.Context) { build(c, &active) }, 500, 300)
		before := baseline()
		click(t, tt, "Two")
		if active != 1 || !tt.FramePending() {
			t.Errorf("%s: active %d, sliding %v", name, active, tt.FramePending())
		}
		settles(t, tt, name)
		key := ui.KeyRight
		if name == "VerticalTabs" {
			key = ui.KeyDown
		}
		tt.Key(0, key)
		tt.Frame()
		if active != 2 {
			t.Errorf("%s: key moved to %d", name, active)
		}
		cues(t, before, sound.Select, sound.Select)
	}
}

func TestClosableTabs(t *testing.T) {
	names, active, changed := []string{"A", "B", "C"}, 1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ClosableTabs(c, "Files", &names, &active, ClosableTabsOptions{}).Changed() {
			changed++
		}
	}, 700, 100)
	click(t, tt, "A")
	if active != 0 || changed != 1 {
		t.Errorf("pick A: active %d, changed %d", active, changed)
	}
	active = 1
	tt.Frame()
	settles(t, tt, "picked")
	before := baseline()
	click(t, tt, "Close C")
	if !slices.Equal(names, []string{"A", "B"}) || active != 1 {
		t.Errorf("close C: %v, active %d", names, active)
	}
	click(t, tt, "Close B")
	if !slices.Equal(names, []string{"A"}) || active != 0 {
		t.Errorf("close the active last: %v, active %d", names, active)
	}
	click(t, tt, "New tab")
	if !slices.Equal(names, []string{"A", "Untitled 2"}) || active != 1 {
		t.Errorf("add: %v, active %d", names, active)
	}
	cues(t, before, sound.Close, sound.Close, sound.Pop)
	settles(t, tt, "added")
	click(t, tt, "Close A")
	settles(t, tt, "closed A")
	click(t, tt, "Close Untitled 2")
	if len(names) != 0 || active != -1 {
		t.Errorf("close all: %v, active %d", names, active)
	}
	settles(t, tt, "ClosableTabs")
	twins, a := []string{"A", "A"}, 0
	render(func(c *ui.Context) { ClosableTabs(c, "x", &twins, &a, ClosableTabsOptions{}) })
}

func TestScrollTabs(t *testing.T) {
	var names []string
	for i := range 20 {
		names = append(names, fmt.Sprintf("Tab %d", i+1))
	}
	active, changed, arrows := 0, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		root, back, forth := ScrollTabs(c, "Tabs", names, &active, ScrollTabsOptions{})
		if root.Changed() {
			changed++
		}
		if back.Clicked() || forth.Clicked() {
			arrows++
		}
	}, 400, 100)
	settles(t, tt, "ScrollTabs")
	click(t, tt, "Tab 2")
	if active != 1 || changed != 1 {
		t.Errorf("pick: active %d, changed %d", active, changed)
	}
	settles(t, tt, "picked")
	first := find(t, tt, "Tab 1").X
	before := baseline()
	click(t, tt, "Scroll forward")
	settles(t, tt, "scrolling")
	if find(t, tt, "Tab 1").X >= first {
		t.Error("the tabs did not scroll")
	}
	click(t, tt, "Scroll back")
	settles(t, tt, "scrolling back")
	if find(t, tt, "Tab 1").X != first {
		t.Error("the tabs did not scroll back")
	}
	if arrows != 2 {
		t.Errorf("arrows reported %d clicks", arrows)
	}
	cues(t, before, sound.Tap, sound.Tap)
	width := find(t, tt, "Tab 1").W
	click(t, tt, "Tab 1")
	tt.Key(0, ui.KeyEnd)
	settles(t, tt, "End")
	tt.Key(0, ui.KeyHome)
	settles(t, tt, "Home")
	if active != 0 || find(t, tt, "Tab 1").W != width {
		t.Errorf("Home: active %d, Tab 1 shows %v of %v", active, find(t, tt, "Tab 1").W, width)
	}
}
