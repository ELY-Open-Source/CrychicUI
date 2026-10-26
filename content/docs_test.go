package content

import (
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestArticleMetaFits(t *testing.T) {
	o := ArticleOptions{Author: "Aleksandra Papadopoulos", Role: "Senior ceramicist and kiln technician at Riverbend Clay",
		Date: time.Date(2026, 10, 14, 9, 0, 0, 0, time.UTC), Minutes: 7, Read: "Approximate reading time: about %d minutes"}
	tt := ui.NewTester(within(320, func(c *ui.Context) { Article(c, "Cone 6", []string{"A slow cool grows crystals."}, o) }), 320, 600)
	settles(t, tt)
	// Card padding is 24 DIPs.
	for _, s := range []string{o.Role, "Oct 14, 2026", "Approximate reading time: about 7 minutes"} {
		inside(t, s, rectOf(t, tt, s), 0, 320-24)
	}
	panics(t, "blank paragraph", func(c *ui.Context) { Article(c, "Cone 6", []string{" "}, ArticleOptions{Author: "Maya"}) })
	panics(t, "needs an author", func(c *ui.Context) { Article(c, "Cone 6", []string{"x"}, ArticleOptions{}) })
}

func outline() []Heading {
	return []Heading{{Title: "Overview"}, {Title: "Install"}, {Title: "Sign in", Depth: 1}, {Title: "Book a kiln"}, {Title: "Webhooks"}}
}

func TestDocTOCPicks(t *testing.T) {
	items, active, changes := outline(), 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if DocTOC(c, items, &active, DocTOCOptions{}).Changed() {
			changes++
		}
	}), 400, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Book a kiln")
	if active != 3 || changes != 1 || !tt.FramePending() {
		t.Errorf("active %d, changes %d, marker moving %v", active, changes, tt.FramePending())
	}
	settles(t, tt)
	click(t, tt, "Book a kiln")
	cues(t, before, sound.Select)
	if changes != 1 {
		t.Errorf("the active heading reported again: %d", changes)
	}
	panics(t, "depth 0 to 2", func(c *ui.Context) { DocTOC(c, []Heading{{Title: "A", Depth: 3}}, new(int), DocTOCOptions{}) })
	panics(t, "active one", func(c *ui.Context) { n := 9; DocTOC(c, outline(), &n, DocTOCOptions{}) })
	panics(t, "repeated", func(c *ui.Context) { DocTOC(c, []Heading{{Title: "A"}, {Title: "A"}}, new(int), DocTOCOptions{}) })
}

func TestDocTOCReorder(t *testing.T) {
	items, active := outline(), 0
	tt := ui.NewTester(padded(func(c *ui.Context) { DocTOC(c, items, &active, DocTOCOptions{}) }), 400, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Book a kiln", func() { items[3], items[4] = items[4], items[3] })
	if active != 0 {
		t.Errorf("a reorder under the pointer picked %q", items[active].Title)
	}
}

func sections() []DocSection {
	return []DocSection{
		{Title: "Getting started", Icon: "rocket", Pages: []string{"Introduction", "Quickstart"}},
		{Title: "Kilns", Icon: "flame", Pages: []string{"Book a kiln", "Firing logs"}},
	}
}

func TestDocSidebarOpensAndPicks(t *testing.T) {
	current, changes := "Quickstart", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if DocSidebar(c, sections(), &current).Changed() {
			changes++
		}
	}), 400, 500)
	settles(t, tt)
	if r, ok := tt.Find("Firing logs"); ok && r.H > 0 {
		t.Errorf("a closed section shows its pages: %v", r)
	}
	before := baseline()
	click(t, tt, "Kilns")
	settles(t, tt)
	click(t, tt, "Firing logs")
	click(t, tt, "Firing logs")
	cues(t, before, sound.Toggle, sound.Select)
	// The fold and the pick each report.
	if current != "Firing logs" || changes != 2 {
		t.Errorf("current %q, changes %d", current, changes)
	}
	panics(t, "current page", func(c *ui.Context) { s := "Nope"; DocSidebar(c, sections(), &s) })
	panics(t, "repeated", func(c *ui.Context) {
		s := "A"
		DocSidebar(c, []DocSection{{Title: "X", Icon: "flame", Pages: []string{"A"}}, {Title: "Y", Icon: "flame", Pages: []string{"A"}}}, &s)
	})
}

func TestAdmonitionTitles(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Admonition(c, Warning, "", "Wax the feet.", AdmonitionOptions{})
		Admonition(c, Tip, "", "Load by cone.", AdmonitionOptions{Kinds: [4]string{"", "Hint"}})
	}), 400, 400)
	settles(t, tt)
	if !exact(tt, "Warning") || !exact(tt, "Hint") {
		t.Errorf("titles: %q", tt.Texts())
	}
	panics(t, "known kind", func(c *ui.Context) { Admonition(c, 9, "", "x", AdmonitionOptions{}) })
}

func TestDocStepsNumbers(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { DocSteps(c, []DocStep{{Title: "Install"}, {Title: "Sign in"}, {Title: "Book"}}) }), 400, 400)
	settles(t, tt)
	for _, n := range []string{"1", "2", "3"} {
		if !exact(tt, n) {
			t.Errorf("step %s missing: %q", n, tt.Texts())
		}
	}
	panics(t, "1 to 20", func(c *ui.Context) { DocSteps(c, nil) })
	panics(t, "1 to 20", func(c *ui.Context) { DocSteps(c, slices.Repeat([]DocStep{{Title: "Wedge"}}, 21)) })
	panics(t, "needs a title", func(c *ui.Context) { DocSteps(c, []DocStep{{Title: " "}}) })
}

func TestDocCodeTabsSwitchAndCopy(t *testing.T) {
	samples := []CodeSample{{Name: "cURL", Lang: "bash", Code: "curl -X POST /v1/bookings"}, {Name: "Go", Lang: "go", Code: "client.Bookings.Create(ctx)"}}
	picked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { DocCodeTabs(c, samples, &picked, DocCodeTabsOptions{}) }), 520, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Go")
	click(t, tt, "Go")
	cues(t, before, sound.Select)
	settles(t, tt)
	click(t, tt, "Copy code")
	if picked != 1 || tt.Clipboard() != samples[1].Code {
		t.Errorf("picked %d, copied %q", picked, tt.Clipboard())
	}
	panics(t, "repeated", func(c *ui.Context) {
		DocCodeTabs(c, []CodeSample{samples[0], samples[0]}, new(int), DocCodeTabsOptions{})
	})
	panics(t, "picked one", func(c *ui.Context) { n := 2; DocCodeTabs(c, samples, &n, DocCodeTabsOptions{}) })
	panics(t, "needs code", func(c *ui.Context) { DocCodeTabs(c, []CodeSample{{Name: "Go"}}, new(int), DocCodeTabsOptions{}) })
}

func TestAPIRowAndPropTable(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		APIRow(c, "get", "/v1/shelves", "Old route.", APIRowOptions{Deprecated: true})
		PropTable(c, []Prop{{Name: "shelves", Type: "int", Default: "1"}, {Name: "kiln", Type: "string", Required: true}}, PropTableOptions{})
	}), 520, 600)
	settles(t, tt)
	for _, s := range []string{"GET", "Deprecated", "Default: 1", "Required", "2"} {
		if !exact(tt, s) {
			t.Errorf("%q missing: %q", s, tt.Texts())
		}
	}
	panics(t, "starting with /", func(c *ui.Context) { APIRow(c, "GET", "v1", "", APIRowOptions{}) })
	panics(t, "needs a type", func(c *ui.Context) { PropTable(c, []Prop{{Name: "a"}}, PropTableOptions{}) })
	panics(t, "repeated", func(c *ui.Context) {
		PropTable(c, []Prop{{Name: "a", Type: "int"}, {Name: "a", Type: "int"}}, PropTableOptions{})
	})
}

func TestMarkdownRenders(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Markdown(c, "## Glaze CLI\n\nBook kilns.", MarkdownOptions{Title: "README.md"}) }), 480, 400)
	settles(t, tt)
	if !tt.HasText("Glaze CLI") || !exact(tt, "README.md") {
		t.Errorf("texts: %q", tt.Texts())
	}
	panics(t, "needs text", func(c *ui.Context) { Markdown(c, "  ", MarkdownOptions{}) })
}

func TestFootnotesBackReportsEachClick(t *testing.T) {
	back, changes := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Footnotes(c, []string{"Cone 6 is 2232°F.", "Oxidation leaves speckles brown."}, &back, FootnotesOptions{}).Changed() {
			changes++
		}
	}), 480, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Back to reference 2")
	click(t, tt, "Back to reference 2")
	cues(t, before, sound.Tap, sound.Tap)
	if back != 1 || changes != 2 {
		t.Errorf("back %d, changes %d", back, changes)
	}
	panics(t, "is blank", func(c *ui.Context) { Footnotes(c, []string{"a", ""}, new(int), FootnotesOptions{}) })
}

func TestGlossaryGroupsByLetter(t *testing.T) {
	terms := []Term{{Word: "ırmak glaze", Meaning: "A Turkish river glaze."}, {Word: "Iron red", Meaning: "A red from iron."},
		{Word: "Irene's slip", Meaning: "A house slip."}, {Word: "3D printing", Meaning: "Clay printed in layers."}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Glossary(c, terms, GlossaryOptions{}) }), 480, 800)
	settles(t, tt)
	n := 0
	for _, s := range tt.Texts() {
		if s == "I" {
			n++
		}
	}
	// One letter chip and one heading.
	if n != 2 || !exact(tt, "#") {
		t.Errorf("letter I shows %d times, texts %q", n, tt.Texts())
	}
	// As text ırmak follows zinc; its group is I.
	pair := []Term{{Word: "Zinc white", Meaning: "A matte glaze."}, {Word: "ırmak glaze", Meaning: "A river glaze."}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Glossary(c, pair, GlossaryOptions{}) }), 480, 600)
	settles(t, tt)
	if i, z := rectOf(t, tt, "ırmak glaze"), rectOf(t, tt, "Zinc white"); i.Y > z.Y {
		t.Errorf("group I sits below group Z: %v, %v", i, z)
	}
	panics(t, "needs a meaning", func(c *ui.Context) { Glossary(c, []Term{{Word: "Slip"}}, GlossaryOptions{}) })
	panics(t, "repeated", func(c *ui.Context) { Glossary(c, []Term{{"Slip", "a"}, {"Slip", "b"}}, GlossaryOptions{}) })
}

func TestDocSearchFiltersAndPicks(t *testing.T) {
	pages := []DocPage{{Section: "Guides", Title: "Firing logs", Text: "Temperatures."}, {Section: "API", Title: "Webhooks", Text: "A POST when a firing cools."}}
	query, picked, changes := "", "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if DocSearch(c, &query, pages, &picked, DocSearchOptions{}).Changed() {
			changes++
		}
	}), 480, 500)
	settles(t, tt)
	click(t, tt, "Search docs")
	tt.Type("hook")
	settles(t, tt)
	if tt.HasText("Firing logs") || !tt.HasText("Webhooks") || changes != 1 {
		t.Errorf("query %q, changes %d, texts %q", query, changes, tt.Texts())
	}
	before := baseline()
	click(t, tt, "Webhooks")
	cues(t, before, sound.Select)
	if picked != "Webhooks" {
		t.Errorf("picked %q", picked)
	}
	click(t, tt, "Search docs")
	tt.Type("zzz")
	settles(t, tt)
	if !exact(tt, "No pages match") {
		t.Errorf("no empty state: %q", tt.Texts())
	}
	panics(t, "needs a section", func(c *ui.Context) {
		DocSearch(c, new(string), []DocPage{{Title: "A"}}, new(string), DocSearchOptions{})
	})
}
