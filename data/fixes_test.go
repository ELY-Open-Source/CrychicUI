package data

import (
	"encoding/json"
	"math"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestEditTableStartsFromCell(t *testing.T) {
	rows := [][]string{{"Ada", "36"}}
	at := [2]int{0, 0}
	tt := ui.NewTester(func(c *ui.Context) {
		EditTable(c, "People", []Column{{Title: "Name"}, {Title: "Age"}}, rows, &at, EditTableOptions{})
	}, 500, 300)
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if rows[0][0] != "Ada" || at != [2]int{-1, -1} {
		t.Errorf("saved %q, at %v", rows[0][0], at)
	}
}

func TestEditTableSwitchKeepsFocus(t *testing.T) {
	rows := [][]string{{"Ada", "36"}}
	at := [2]int{-1, -1}
	tt := ui.NewTester(func(c *ui.Context) {
		EditTable(c, "People", []Column{{Title: "Name"}, {Title: "Age"}}, rows, &at, EditTableOptions{})
	}, 500, 300)
	click(t, tt, "Edit Name, row 1")
	click(t, tt, "Edit Age, row 1")
	tt.Command("selectAll")
	tt.Type("40")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if rows[0][0] != "Ada" || rows[0][1] != "40" || at != [2]int{-1, -1} {
		t.Errorf("rows %q, at %v", rows[0], at)
	}
}

func TestJSONTreeSlashKeys(t *testing.T) {
	var v any
	if err := json.Unmarshal([]byte(`{"a":{"b":{"x":1}},"a/b":{"y":2}}`), &v); err != nil {
		t.Fatal(err)
	}
	var open ui.Selection[string]
	tt := ui.NewTester(func(c *ui.Context) { JSONTree(c, "Payload", v, &open, JSONTreeOptions{}) }, 400, 400)
	for _, k := range []string{"$", "a", "b"} {
		click(t, tt, "Expand "+k)
		settles(t, tt, k)
	}
	if !slices.Contains(tt.Texts(), "x") || slices.Contains(tt.Texts(), "y") {
		t.Errorf("shows %q", tt.Texts())
	}
}

func TestNumbersAndListsMustBeSound(t *testing.T) {
	n := 1.0
	for name, p := range map[string]Property{
		"NaN step": {Name: "n", Value: &n, Max: 2, Step: math.NaN()},
		"Inf step": {Name: "n", Value: &n, Max: 2, Step: math.Inf(1)},
	} {
		panics(t, name, func() { render(func(c *ui.Context) { PropertyGrid(c, "Props", []Property{p}) }) })
	}
	nan := math.NaN()
	panics(t, "NaN value", func() {
		render(func(c *ui.Context) { PropertyGrid(c, "Props", []Property{{Name: "n", Value: &nan, Max: 2}}) })
	})
	for _, sizes := range [][]int{{-1, 10}, {0, 10}} {
		page, size := 0, 10
		panics(t, "size", func() {
			render(func(c *ui.Context) { TableFooter(c, 42, &page, &size, TableFooterOptions{Sizes: sizes}) })
		})
	}
	var rs []Rule
	var q Query
	panics(t, "no fields", func() {
		render(func(c *ui.Context) { FilterBuilder(c, "Filters", nil, []string{"is"}, &rs, FilterBuilderOptions{}) })
	})
	panics(t, "no operators", func() {
		render(func(c *ui.Context) { QueryBuilder(c, "Query", []string{"Name"}, nil, &q, QueryBuilderOptions{}) })
	})
}

func TestReorderListRowsSlide(t *testing.T) {
	items := []string{"a", "b", "c", "d"}
	tt := ui.NewTester(func(c *ui.Context) { ReorderList(c, "Steps", &items, ReorderListOptions{}) }, 400, 400)
	settles(t, tt, "list")
	old, slot := find(t, tt, "d").Y, find(t, tt, "c").Y
	click(t, tt, "d")
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Frame()
	if y := find(t, tt, "d").Y; y < old-1 {
		t.Fatalf("d jumped to %v on the first frame, was %v", y, old)
	}
	tt.Advance(60 * time.Millisecond)
	mid := find(t, tt, "d").Y
	if !tt.FramePending() || mid <= slot+1 || mid > old+1 {
		t.Errorf("mid-slide d at %v, want between %v and %v", mid, slot, old)
	}
	settles(t, tt, "slide")
	if y := find(t, tt, "d").Y; y != slot {
		t.Errorf("d rests at %v, want %v", y, slot)
	}
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Advance(60 * time.Millisecond)
	before := find(t, tt, "d").Y
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Frame()
	if y := find(t, tt, "d").Y; y < before-6 || y > before+6 {
		t.Errorf("interrupted slide jumped from %v to %v", before, y)
	}
	settles(t, tt, "second slide")
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Frame()
	if a := find(t, tt, "a"); tt.FramePending() || find(t, tt, "d").Y != a.Y+old-slot {
		t.Errorf("reduced motion: pending %v, a %v, d %v", tt.FramePending(), a, find(t, tt, "d"))
	}
}

func TestColumnPickerRowsSlide(t *testing.T) {
	order, hidden := []string{"Name", "Age", "Role"}, ui.Selection[string]{}
	tt := ui.NewTester(func(c *ui.Context) { ColumnPicker(c, "Columns", &order, &hidden) }, 400, 400)
	settles(t, tt, "list")
	r, slot := find(t, tt, "Role"), find(t, tt, "Age").Y
	tt.ClickAt(r.X+r.W+40, r.Y+r.H/2)
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Advance(60 * time.Millisecond)
	mid := find(t, tt, "Role").Y
	if !slices.Equal(order, []string{"Name", "Role", "Age"}) || mid <= slot+1 || mid > r.Y+1 {
		t.Errorf("order %v, Role at %v, want between %v and %v", order, mid, slot, r.Y)
	}
	settles(t, tt, "slide")
}

func TestRowActionsOpenByCallerIsSilent(t *testing.T) {
	chosen, open := -1, true
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		RowActions(c, "More actions", []Action{{"Edit", "pencil"}}, &chosen, &open)
	}, 500, 400)
	tt.Frame()
	cues(t, before)
}
