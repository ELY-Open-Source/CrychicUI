package data

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var fields, ops = []string{"Name", "Age"}, []string{"is", "contains"}

func TestFilterBuilder(t *testing.T) {
	var rs []Rule
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if FilterBuilder(c, "Filters", fields, ops, &rs, FilterBuilderOptions{}).Changed() {
			changed++
		}
	}, 600, 400)
	before := baseline()
	click(t, tt, "Add filter")
	if len(rs) != 1 || rs[0] != (Rule{Field: "Name", Op: "is"}) || !tt.FramePending() {
		t.Fatalf("add: %+v", rs)
	}
	settles(t, tt, "added")
	click(t, tt, "Value 1")
	tt.Type("Ada")
	tt.Frame()
	if rs[0].Value != "Ada" {
		t.Fatalf("value %q", rs[0].Value)
	}
	click(t, tt, "Add filter")
	click(t, tt, "Remove filter 1")
	if len(rs) != 1 || rs[0].Value != "" || !tt.FramePending() {
		t.Fatalf("remove first: %+v", rs)
	}
	cues(t, before, sound.Pop, sound.Pop, sound.Delete)
	settles(t, tt, "removed")
	panics(t, "unknown field", func() {
		render(func(c *ui.Context) {
			r := []Rule{{"Height", "is", ""}}
			FilterBuilder(c, "x", fields, ops, &r, FilterBuilderOptions{})
		})
	})
}

func TestQueryBuilder(t *testing.T) {
	q := Query{Rules: []Rule{{"Name", "is", "Ada"}}}
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if QueryBuilder(c, "Query", fields, ops, &q, QueryBuilderOptions{}).Changed() {
			changed++
		}
	}, 700, 500)
	before := baseline()
	click(t, tt, "Any")
	click(t, tt, "Add group")
	if !q.Any || len(q.Groups) != 1 || len(q.Groups[0].Rules) != 1 {
		t.Fatalf("any and group: %+v", q)
	}
	click(t, tt, "Remove group")
	click(t, tt, "Add rule")
	if len(q.Groups) != 0 || len(q.Rules) != 2 {
		t.Fatalf("remove group, add rule: %+v", q)
	}
	click(t, tt, "Add group")
	// The top rule is built before the group.
	click(t, tt, "Remove rule 1")
	if len(q.Rules) != 1 || q.Rules[0].Value != "" || len(q.Groups) != 1 || changed != 6 {
		t.Errorf("remove rule: %+v, changed %d", q, changed)
	}
	cues(t, before, sound.Toggle, sound.Pop, sound.Delete, sound.Pop, sound.Pop, sound.Delete)
}

func TestFacets(t *testing.T) {
	facets := []Facet{{"Color", []FacetOption{{"Red", 12}, {"Blue", 3}}}, {"Size", []FacetOption{{"Red", 1}}}}
	var chosen ui.Selection[string]
	tt := ui.NewTester(func(c *ui.Context) { Facets(c, "Filters", facets, &chosen) }, 400, 400)
	if !tt.HasText("12") {
		t.Fatalf("counts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Blue")
	if !chosen.Has("Color/Blue") || chosen.Len() != 1 {
		t.Errorf("chose %d", chosen.Len())
	}
	cues(t, before, sound.Toggle)
	panics(t, "negative", func() {
		render(func(c *ui.Context) {
			var s ui.Selection[string]
			Facets(c, "x", []Facet{{"a", []FacetOption{{"b", -1}}}}, &s)
		})
	})
}

func TestResults(t *testing.T) {
	results := []Result{{"Go tour", "Learn go fast"}, {"Rust book", "Ownership"}}
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) { Results(c, "Search", "go", results, &chosen, ResultsOptions{}) }, 500, 400)
	if !tt.FramePending() {
		t.Error("results do not stagger in")
	}
	settles(t, tt, "results")
	before := baseline()
	click(t, tt, "Rust book")
	if chosen != 1 {
		t.Errorf("chose %d", chosen)
	}
	cues(t, before, sound.Tap)
	empty := ui.NewTester(func(c *ui.Context) { ch := -1; Results(c, "Search", "zz", nil, &ch, ResultsOptions{}) }, 400, 200)
	if !empty.HasText("No results") {
		t.Errorf("empty %q", empty.Texts())
	}
}

func TestTableFooter(t *testing.T) {
	page, size := 0, 10
	tt := ui.NewTester(func(c *ui.Context) { TableFooter(c, 42, &page, &size, TableFooterOptions{}) }, 600, 300)
	if !tt.HasText("1–10 of 42") {
		t.Fatalf("range %q", tt.Texts())
	}
	before := baseline()
	for range 4 {
		click(t, tt, "Next page")
	}
	if page != 4 || !tt.HasText("41–42 of 42") {
		t.Fatalf("last page: %d, %q", page, tt.Texts())
	}
	click(t, tt, "Next page")
	if page != 4 {
		t.Errorf("Next past the end: %d", page)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	click(t, tt, "25")
	if size != 25 || page != 0 {
		t.Errorf("size: %d, page %d", size, page)
	}
	cues(t, before, sound.Tap, sound.Tap, sound.Tap, sound.Tap, sound.Select)
	panics(t, "page", func() {
		render(func(c *ui.Context) { p, s := 5, 10; TableFooter(c, 42, &p, &s, TableFooterOptions{}) })
	})
}

func TestBulkBar(t *testing.T) {
	count, chosen, clears := 0, -1, 0
	actions := []Action{{"Archive", "archive"}, {"Delete", "trash-2"}}
	tt := ui.NewTester(func(c *ui.Context) {
		if _, clear := BulkBar(c, count, actions, &chosen, BulkBarOptions{}); clear.Clicked() {
			clears++
			count = 0
		}
	}, 600, 300)
	if tt.HasText("Archive") {
		t.Fatal("bar shown with nothing selected")
	}
	count = 3
	tt.Frame()
	if !tt.HasText("3 selected") || !tt.FramePending() {
		t.Fatalf("bar %q", tt.Texts())
	}
	settles(t, tt, "bar")
	before := baseline()
	click(t, tt, "Archive")
	click(t, tt, "Clear selection")
	if chosen != 0 || clears != 1 {
		t.Errorf("chosen %d, clears %d", chosen, clears)
	}
	cues(t, before, sound.Tap, sound.Tap)
	settles(t, tt, "bar gone")
}

func TestRowActions(t *testing.T) {
	chosen, open := -1, false
	actions := []Action{{"Edit", "pencil"}, {"Delete", "trash-2"}}
	tt := ui.NewTester(func(c *ui.Context) { RowActions(c, "More actions", actions, &chosen, &open) }, 500, 400)
	before := baseline()
	click(t, tt, "More actions")
	if !open || !tt.HasText("Delete") {
		t.Fatalf("open %v, %q", open, tt.Texts())
	}
	settles(t, tt, "menu")
	click(t, tt, "Delete")
	if chosen != 1 || open {
		t.Errorf("chose %d, open %v", chosen, open)
	}
	cues(t, before, sound.Open, sound.Select)
}
