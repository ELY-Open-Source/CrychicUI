package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDataDemosLog(t *testing.T) {
	type step func(*ui.Tester) error
	click := func(l string) step { return func(tt *ui.Tester) error { return tt.Click(l) } }
	key := func(mods ui.Modifiers, k ui.Key) step {
		return func(tt *ui.Tester) error { tt.Key(mods, k); return nil }
	}
	runs := map[string]struct {
		steps []step
		want  string
	}{
		"DataTable":     {[]step{click("Name")}, "sorted by Name, descending false"},
		"SelectTable":   {[]step{click("Select row 1")}, "1 chosen"},
		"ExpandTable":   {[]step{click("Details for row 1")}, "1 open"},
		"EditTable":     {[]step{click("Edit Name, row 1")}, "editing [0 0]"},
		"TreeTable":     {[]step{click("Expand src")}, "1 open"},
		"VirtualList":   {[]step{click("Invoice 00003")}, "chose 2"},
		"GroupedList":   {[]step{click("Grace")}, "chose [1 0]"},
		"ReorderList":   {[]step{click("Build"), key(ui.Alt, ui.KeyUp)}, "Plan, Build, Design, Test, Ship"},
		"SwipeList":     {[]step{key(0, ui.KeyTab), key(0, ui.KeyDelete)}, "3 left, 0 archived"},
		"CardList":      {[]step{click("Team")}, "chose Team"},
		"GridList":      {[]step{click("Music")}, "chose Music"},
		"CheckTree":     {[]step{click("Email")}, "1 checked"},
		"DirTree":       {[]step{click("Music")}, `chose "Music", 0 open`},
		"JSONTree":      {[]step{click("Expand user")}, "2 open"},
		"KVTable":       {[]step{click("Copy Host")}, "copied"},
		"PropertyGrid":  {[]step{key(0, ui.KeyTab), key(0, ui.KeyTab), key(0, ui.KeySpace)}, "visible false"},
		"Spreadsheet":   {[]step{key(0, ui.KeyTab), key(0, ui.KeyRight)}, "cell [1 2]"},
		"ColumnPicker":  {[]step{click("Role")}, "1 hidden, first Name"},
		"FilterBuilder": {[]step{click("Add filter")}, "2 rules"},
		"QueryBuilder":  {[]step{click("Any")}, "any true, 1 rules, 1 groups"},
		"Facets":        {[]step{click("Acme")}, "1 chosen"},
		"Results":       {[]step{click("Effective Go")}, "opened Effective Go"},
		"TableFooter":   {[]step{click("Next page")}, "page 2, 10 per page"},
		"BulkBar":       {[]step{click("Select row 1"), click("Archive")}, "Archive 1"},
		"RowActions":    {[]step{click("More actions")}, "open true, chose -1"},
	}
	for _, e := range Entries() {
		run, ok := runs[e.Name]
		if !ok || e.Category != "data" {
			continue
		}
		delete(runs, e.Name)
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		for _, s := range run.steps {
			if err := s(tt); err != nil {
				t.Fatalf("%s: %v", e.Name, err)
			}
			tt.Advance(500 * time.Millisecond)
		}
		if len(logs) == 0 || !strings.HasPrefix(logs[len(logs)-1], run.want) {
			t.Errorf("%s logged %q, want %q last", e.Name, logs, run.want)
		}
	}
	if len(runs) > 0 {
		t.Errorf("no data demo for %d runs", len(runs))
	}
}

func TestBulkBarDemoLogsSelection(t *testing.T) {
	var logs []string
	tt := ui.NewTester(func(c *ui.Context) { demoBulkBar(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
	tt.Advance(time.Second)
	if err := tt.Click("Select row 1"); err != nil {
		t.Fatal(err)
	}
	if len(logs) == 0 || logs[len(logs)-1] != "selected 1" {
		t.Errorf("logged %q, want selected 1", logs)
	}
}
