package choice

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDecimalTopReachable(t *testing.T) {
	for _, s := range []struct {
		name string
		view func(c *ui.Context, v *float64)
		top  bool
	}{
		{"ValueSlider", func(c *ui.Context, v *float64) { ValueSlider(c, "Level", v, 0, 0.3, ValueSliderOptions{Step: 0.1}) }, false},
		{"TickSlider", func(c *ui.Context, v *float64) { TickSlider(c, "Level", v, 0, 0.3, 0.1, TickSliderOptions{}) }, false},
		{"VerticalSlider", func(c *ui.Context, v *float64) {
			VerticalSlider(c, "Level", v, 0, 0.3, VerticalSliderOptions{Step: 0.1})
		}, true},
	} {
		v := 0.0
		tt := ui.NewTester(func(c *ui.Context) { ui.Column(c).Padding(10).Children(func() { s.view(c, &v) }) }, 320, 300)
		r := find(t, tt, "Level")
		x, y := r.X+r.W-1, r.Y+r.H/2
		if s.top {
			x, y = r.X+r.W/2, r.Y+1
		}
		tt.Press(x, y)
		tt.Release(x, y)
		if v != 0.3 {
			t.Errorf("%s: the far end gave %v, want 0.3", s.name, v)
		}
	}
	low, high := 0.0, 0.1
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() { RangeSlider(c, "Level", &low, &high, 0, 0.3, 0.1, RangeSliderOptions{}) })
	}, 320, 300)
	r := find(t, tt, "Level")
	tt.Press(r.X+r.W-1, r.Y+r.H/2)
	tt.Release(r.X+r.W-1, r.Y+r.H/2)
	if high != 0.3 {
		t.Errorf("RangeSlider: the far end gave %v, want 0.3", high)
	}
}

// pressSwap presses name, reorders, then releases there.
func pressSwap(t *testing.T, tt *ui.Tester, name string, reorder func()) {
	t.Helper()
	r := find(t, tt, name)
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	reorder()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func TestOptionsKeepIdentity(t *testing.T) {
	labels, picked := []string{"Alpha", "Beta", "Gamma"}, 2
	tt := ui.NewTester(func(c *ui.Context) { Segments(c, "Greek", labels, &picked) }, 400, 200)
	pressSwap(t, tt, "Alpha", func() { labels = []string{"Beta", "Alpha", "Gamma"} })
	if picked != 2 {
		t.Errorf("Segments: a reorder under the press picked %q", labels[picked])
	}

	options, chosen := []string{"Alpha", "Beta", "Gamma"}, "Gamma"
	tt = ui.NewTester(func(c *ui.Context) { Dropdown(c, "Greek", options, &chosen) }, 400, 400)
	click(t, tt, "Greek")
	pressSwap(t, tt, "Alpha", func() { options = []string{"Beta", "Alpha", "Gamma"} })
	if chosen != "Gamma" {
		t.Errorf("Dropdown: a reorder under the press chose %q", chosen)
	}

	groups, chosen := []DropdownGroup{{"Greek", []string{"Alpha", "Beta"}}, {"Latin", []string{"Gamma"}}}, "Gamma"
	tt = ui.NewTester(func(c *ui.Context) { GroupedDropdown(c, "Letter", groups, &chosen) }, 400, 400)
	click(t, tt, "Letter")
	pressSwap(t, tt, "Alpha", func() { groups = []DropdownGroup{{"Greek", []string{"Beta", "Alpha"}}, {"Latin", []string{"Gamma"}}} })
	if chosen != "Gamma" {
		t.Errorf("GroupedDropdown: a reorder under the press chose %q", chosen)
	}
}

func TestCascadeColumnsFollowPath(t *testing.T) {
	items := []TreeItem{
		{Label: "North", Children: []TreeItem{{Label: "Warehouse"}, {Label: "Office"}}},
		{Label: "South", Children: []TreeItem{{Label: "Warehouse"}, {Label: "Depot"}}},
	}
	path := []string{"North"}
	tt := ui.NewTester(func(c *ui.Context) { Cascade(c, "Site", items, &path) }, 500, 300)
	pressSwap(t, tt, "Warehouse", func() { path = []string{"South"} })
	if len(path) != 1 || path[0] != "South" {
		t.Errorf("a press begun under North chose %q", path)
	}
}

func TestLongPopupsScroll(t *testing.T) {
	var many []string
	var leaves []TreeItem
	for i := 1; i <= 40; i++ {
		many = append(many, fmt.Sprintf("Option %02d", i))
		leaves = append(leaves, TreeItem{Label: many[i-1]})
	}
	chosen, picks := many[0], []string{}
	for _, p := range []struct {
		name string
		view func(c *ui.Context)
	}{
		{"Dropdown", func(c *ui.Context) { Dropdown(c, "Pick", many, &chosen) }},
		{"GroupedDropdown", func(c *ui.Context) { GroupedDropdown(c, "Pick", []DropdownGroup{{"All", many}}, &chosen) }},
		{"TreeDropdown", func(c *ui.Context) { TreeDropdown(c, "Pick", leaves, &chosen, TreeDropdownOptions{}) }},
		{"SearchDropdown", func(c *ui.Context) { SearchDropdown(c, "Pick", many, &chosen, SearchDropdownOptions{}) }},
		{"MultiDropdown", func(c *ui.Context) { MultiDropdown(c, "Pick", many, &picks) }},
	} {
		tt := ui.NewTester(func(c *ui.Context) { ui.Column(c).Padding(10).Children(func() { p.view(c) }) }, 320, 400)
		click(t, tt, "Pick")
		tt.Frame()
		at := find(t, tt, "Option 05")
		for range 20 {
			tt.Scroll(at.X+at.W/2, at.Y+at.H/2, 0, 200)
			tt.Frame()
		}
		r, ok := tt.Find("Option 40")
		if !ok || r.H <= 0 || r.Y < 0 || r.Y+r.H > 400 {
			t.Errorf("%s: the last option stays out of reach at %v", p.name, r)
		}
	}
}

func TestPanicsNameLabel(t *testing.T) {
	for name, view := range map[string]func(c *ui.Context){
		`"Sync"`:   func(c *ui.Context) { Check(c, "Sync", nil) },
		`"Volume"`: func(c *ui.Context) { v := 5.0; ValueSlider(c, "Volume", &v, 0, 1, ValueSliderOptions{}) },
		`"Fruit"`:  func(c *ui.Context) { s := "Kiwi"; Dropdown(c, "Fruit", []string{"Pear"}, &s) },
		`"Greek"`:  func(c *ui.Context) { i := 0; Segments(c, "Greek", []string{"A", "A"}, &i) },
	} {
		func() {
			defer func() {
				if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), name) {
					t.Errorf("panic %v, want the label %s", r, name)
				}
			}()
			ui.Render(view, 300, 200, 1)
		}()
	}
}
