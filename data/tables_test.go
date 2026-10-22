package data

import (
	"math"
	"slices"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var people = [][]string{{"Ada", "36"}, {"Grace", "45"}, {"Linus", "28"}}

func TestDataTable(t *testing.T) {
	cols := []Column{{Title: "Name"}, {Title: "Age", Width: 60, End: true}}
	sort, changed := ui.SortOrder{}, 0
	tt := ui.NewTester(func(c *ui.Context) {
		rows := slices.Clone(people)
		slices.SortFunc(rows, func(a, b []string) int {
			if sort.Descending {
				a, b = b, a
			}
			if sort.Column == "Name" {
				return compare(a[0], b[0])
			}
			return 0
		})
		if DataTable(c, "People", cols, len(rows), &sort, func(r, col int) { ui.Text(c, rows[r][col]) }, DataTableOptions{}).Changed() {
			changed++
		}
	}, 500, 300)
	before := baseline()
	click(t, tt, "Name")
	if sort != (ui.SortOrder{Column: "Name"}) || changed != 1 {
		t.Fatalf("first click: %+v, changed %d", sort, changed)
	}
	click(t, tt, "Name")
	if !sort.Descending || !tt.FramePending() {
		t.Fatalf("second click: %+v, arrow turning %v", sort, tt.FramePending())
	}
	settles(t, tt, "sorted")
	if a, l := find(t, tt, "Ada"), find(t, tt, "Linus"); a.Y < l.Y {
		t.Errorf("descending: Ada above Linus")
	}
	cues(t, before, sound.Tick, sound.Tick)
	empty := ui.NewTester(func(c *ui.Context) {
		DataTable(c, "People", cols, 0, nil, func(int, int) {}, DataTableOptions{})
	}, 400, 200)
	if !empty.HasText("No rows") {
		t.Errorf("empty %q", empty.Texts())
	}
	panics(t, "no columns", func() {
		render(func(c *ui.Context) { DataTable(c, "x", nil, 0, nil, func(int, int) {}, DataTableOptions{}) })
	})
}

func compare(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func TestSelectTable(t *testing.T) {
	var chosen ui.Selection[int]
	changed := 0
	cols := []Column{{Title: "Name"}}
	tt := ui.NewTester(func(c *ui.Context) {
		if SelectTable(c, "People", cols, 3, &chosen, func(r, _ int) { ui.Text(c, people[r][0]) }, SelectTableOptions{}).Changed() {
			changed++
		}
	}, 500, 300)
	before := baseline()
	click(t, tt, "Select row 2")
	if !chosen.Has(1) || chosen.Len() != 1 {
		t.Fatalf("row 2: %d chosen", chosen.Len())
	}
	click(t, tt, "Select all")
	if chosen.Len() != 3 {
		t.Fatalf("all from mixed: %d chosen", chosen.Len())
	}
	click(t, tt, "Select all")
	if chosen.Len() != 0 || changed != 3 {
		t.Errorf("none: %d chosen, changed %d", chosen.Len(), changed)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle)
}

func TestExpandTable(t *testing.T) {
	var open ui.Selection[int]
	tt := ui.NewTester(func(c *ui.Context) {
		ExpandTable(c, "People", []Column{{Title: "Name"}}, 3, &open, func(r, _ int) { ui.Text(c, people[r][0]) },
			func(r int) { ui.Text(c, "Age "+people[r][1]) }, ExpandTableOptions{})
	}, 500, 400)
	before := baseline()
	click(t, tt, "Details for row 2")
	if !open.Has(1) || !tt.HasText("Age 45") {
		t.Fatalf("open: %v, %q", open.Has(1), tt.Texts())
	}
	// Find cuts rects by clips: the springing row starts short.
	early := find(t, tt, "Age 45").H
	settles(t, tt, "opened")
	if full := find(t, tt, "Age 45").H; early >= full {
		t.Errorf("detail did not spring open: %v then %v", early, full)
	}
	click(t, tt, "Details for row 2")
	settles(t, tt, "closed")
	if open.Has(1) || tt.HasText("Age 45") {
		t.Errorf("still open: %q", tt.Texts())
	}
	cues(t, before, sound.Toggle, sound.Toggle)
}

func TestEditTable(t *testing.T) {
	rows := [][]string{{"Ada", "36"}}
	at, changed := [2]int{-1, -1}, 0
	cols := []Column{{Title: "Name"}, {Title: "Age"}}
	tt := ui.NewTester(func(c *ui.Context) {
		if EditTable(c, "People", cols, rows, &at, EditTableOptions{}).Changed() {
			changed++
		}
	}, 500, 300)
	before := baseline()
	click(t, tt, "Edit Name, row 1")
	if at != [2]int{0, 0} || changed != 1 {
		t.Fatalf("edit: at %v, changed %d", at, changed)
	}
	tt.Command("selectAll")
	tt.Type("Grace")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if rows[0][0] != "Grace" || at != [2]int{-1, -1} {
		t.Fatalf("save: %q, at %v", rows[0][0], at)
	}
	cues(t, before, sound.Success)
	click(t, tt, "Edit Age, row 1")
	tt.Type("9")
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if rows[0][1] != "36" || at != [2]int{-1, -1} || changed != 4 {
		t.Errorf("cancel: %q, at %v, changed %d", rows[0][1], at, changed)
	}
	panics(t, "ragged", func() {
		render(func(c *ui.Context) {
			a := [2]int{-1, -1}
			EditTable(c, "x", cols, [][]string{{"a"}}, &a, EditTableOptions{})
		})
	})
	panics(t, "bad cell", func() {
		render(func(c *ui.Context) { a := [2]int{0, -1}; EditTable(c, "x", cols, rows, &a, EditTableOptions{}) })
	})
}

func TestTreeTable(t *testing.T) {
	rows := []TreeRow{{Cells: []string{"src", "2"}, Children: []TreeRow{{Cells: []string{"main.go", "1"}}}}, {Cells: []string{"go.mod", "1"}}}
	var open ui.Selection[string]
	tt := ui.NewTester(func(c *ui.Context) {
		TreeTable(c, "Files", []Column{{Title: "Name"}, {Title: "Files", Width: 60}}, rows, &open, TreeTableOptions{})
	}, 500, 300)
	if tt.HasText("main.go") {
		t.Fatal("child shown while closed")
	}
	before := baseline()
	click(t, tt, "Expand src")
	if !open.Has("src") || !tt.HasText("main.go") {
		t.Fatalf("open: %q", tt.Texts())
	}
	cues(t, before, sound.Toggle)
	panics(t, "twins", func() {
		render(func(c *ui.Context) {
			var o ui.Selection[string]
			TreeTable(c, "x", []Column{{Title: "Name"}}, []TreeRow{{Cells: []string{"a"}}, {Cells: []string{"a"}}}, &o, TreeTableOptions{})
		})
	})
}

func TestPivot(t *testing.T) {
	records := []Record{{"North", "Q1", 1}, {"North", "Q1", 2}, {"North", "Q2", 3}, {"South", "Q1", 4}}
	tt := ui.NewTester(func(c *ui.Context) { Pivot(c, "Sales", records, PivotOptions{}) }, 600, 300)
	if !tt.FramePending() {
		t.Error("cells do not fade in")
	}
	settles(t, tt, "pivot")
	for _, want := range []string{"6", "4", "7", "10"} {
		if !tt.HasText(want) {
			t.Errorf("no %s in %q", want, tt.Texts())
		}
	}
	panics(t, "NaN", func() {
		render(func(c *ui.Context) { Pivot(c, "x", []Record{{"a", "b", math.NaN()}}, PivotOptions{}) })
	})
}

func TestCompareTable(t *testing.T) {
	th := ui.LightTheme()
	tt := ui.NewTester(func(c *ui.Context) {
		CompareTable(c, "Plans", []string{"Free", "Pro"}, []Feature{{"Sync", []bool{false, true}}}, CompareTableOptions{})
	}, 600, 300)
	green := func() int {
		n, img := 0, tt.Image()
		for y := 0; y < 300; y++ {
			for x := 0; x < 600; x++ {
				if p := img.RGBAAt(x, y); near(p.R, th.Success.R) && near(p.G, th.Success.G) && near(p.B, th.Success.B) {
					n++
				}
			}
		}
		return n
	}
	if g := green(); g != 0 {
		t.Errorf("check drawn before it draws in: %d px", g)
	}
	find(t, tt, "Included")
	find(t, tt, "Not included")
	settles(t, tt, "checks")
	if g := green(); g < 10 {
		t.Errorf("check not drawn: %d px", g)
	}
}

func near(a, b uint8) bool { return a+8 >= b && b+8 >= a }

func TestDiffTable(t *testing.T) {
	cols := []Column{{Title: "Key"}, {Title: "Value"}}
	rows := []DiffRow{{Change: ChangeAdded, Cells: []string{"port", "80"}}, {Change: ChangeRemoved, Cells: []string{"debug", "on"}},
		{Change: ChangeEdited, Cells: []string{"host", "b.dev"}, Before: []string{"host", "a.dev"}}}
	tt := ui.NewTester(func(c *ui.Context) { DiffTable(c, "Config", cols, rows, DiffTableOptions{}) }, 600, 300)
	for _, want := range []string{"Added", "Removed", "Changed", "a.dev", "b.dev"} {
		if !tt.HasText(want) {
			t.Errorf("no %s in %q", want, tt.Texts())
		}
	}
	panics(t, "edited without before", func() {
		render(func(c *ui.Context) {
			DiffTable(c, "x", cols, []DiffRow{{Change: ChangeEdited, Cells: []string{"a", "b"}}}, DiffTableOptions{})
		})
	})
}

func TestSpreadsheet(t *testing.T) {
	cells := [][]string{{"1", "2"}, {"3", "4"}}
	at, changed := [2]int{0, 0}, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Spreadsheet(c, "Sheet", cells, &at).Changed() {
			changed++
		}
	}, 600, 300)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	if at != [2]int{1, 1} {
		t.Fatalf("arrows: at %v", at)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	tt.Command("selectAll")
	tt.Type("9")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if cells[1][1] != "9" {
		t.Fatalf("edit: %q", cells[1][1])
	}
	tt.Key(0, ui.KeyLeft)
	tt.Frame()
	if at != [2]int{1, 0} {
		t.Errorf("focus did not return to the sheet: at %v", at)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	if changed != 5 {
		t.Errorf("changed %d, want 5", changed)
	}
	panics(t, "outside", func() { render(func(c *ui.Context) { a := [2]int{2, 0}; Spreadsheet(c, "x", cells, &a) }) })
}
