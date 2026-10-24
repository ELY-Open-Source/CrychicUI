package agent

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func testSource() Source {
	return Source{Title: "Layout in MyGo", Host: "mygo.egoist.dev", URL: "https://mygo.egoist.dev/docs/ui", Snippet: "Rows, columns and grids."}
}

func TestCitation(t *testing.T) {
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Citation(c, 2, testSource()).Changed() {
			changed++
		}
	}), 400, 200)
	settles(t, tt, "citation")
	if !tt.HasText("2") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "2 Layout in MyGo")
	if changed != 1 || len(tt.OpenedURLs()) != 1 || tt.OpenedURLs()[0] != testSource().URL {
		t.Errorf("changed %d opened %q", changed, tt.OpenedURLs())
	}
	cues(t, before, sound.Open)
	panics(t, "zero number", func() { render(func(c *ui.Context) { Citation(c, 0, testSource()) }) })
	panics(t, "no URL", func() { render(func(c *ui.Context) { Citation(c, 1, Source{Title: "x", Host: "y"}) }) })
}

func TestSourceCard(t *testing.T) {
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if SourceCard(c, 1, testSource()).Changed() {
			changed++
		}
	}), 400, 300)
	settles(t, tt, "source")
	for _, s := range []string{"mygo.egoist.dev", "M", "Rows, columns and grids."} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	before := baseline()
	click(t, tt, "Layout in MyGo")
	if changed != 1 || len(tt.OpenedURLs()) != 1 || tt.OpenedURLs()[0] != testSource().URL {
		t.Errorf("changed %d opened %q", changed, tt.OpenedURLs())
	}
	cues(t, before, sound.Tap)
}

func TestModelPicker(t *testing.T) {
	selected, changed := 0, 0
	models := []Model{{Name: "Opus", Description: "Deep reasoning", Speed: "Slow", Cost: "$$$"}, {Name: "Haiku", Speed: "Fast", Cost: "$"}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ModelPicker(c, "Models", models, &selected).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "picker")
	if !tt.HasText("Deep reasoning") || !tt.HasText("$$$") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Haiku")
	click(t, tt, "Haiku")
	if selected != 1 || changed != 1 {
		t.Errorf("selected %d changed %d", selected, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "out of range", func() {
		n := 2
		render(func(c *ui.Context) { ModelPicker(c, "M", models, &n) })
	})
	panics(t, "duplicate names", func() {
		n := 0
		render(func(c *ui.Context) { ModelPicker(c, "M", []Model{{Name: "A"}, {Name: "A"}}, &n) })
	})
}

func TestModelPickerKeepsIdentityOnReorder(t *testing.T) {
	selected := -1
	models := []Model{{Name: "Opus"}, {Name: "Haiku"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ModelPicker(c, "Models", models, &selected) }), 500, 300)
	settles(t, tt, "picker")
	r := find(t, tt, "Opus")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	models[0], models[1] = models[1], models[0]
	tt.Frame()
	r = find(t, tt, "Opus")
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if selected != 1 {
		t.Errorf("selected %d, want Opus at 1", selected)
	}
}

func TestTokenCount(t *testing.T) {
	used := 1200
	tt := ui.NewTester(padded(func(c *ui.Context) { TokenCount(c, used, 4000, TokenCountOptions{}) }), 400, 200)
	settles(t, tt, "tokens")
	if !tt.HasText("2800 left") {
		t.Errorf("texts %q", tt.Texts())
	}
	used = 4500
	tt.Frame()
	settles(t, tt, "over")
	if !tt.HasText("500 over") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "zero limit", func() { render(func(c *ui.Context) { TokenCount(c, 1, 0, TokenCountOptions{}) }) })
	panics(t, "inexact count", func() { render(func(c *ui.Context) { TokenCount(c, exact+1, exact, TokenCountOptions{}) }) })
}

func TestContextMeter(t *testing.T) {
	used := 90_000
	tt := ui.NewTester(padded(func(c *ui.Context) { ContextMeter(c, used, 200_000, ContextMeterOptions{}) }), 400, 200)
	settles(t, tt, "meter")
	if !tt.HasText("45%") || !tt.HasText("90k of 200k tokens") || tt.HasText("Older messages will be summarized") {
		t.Errorf("texts %q", tt.Texts())
	}
	used = 170_000
	tt.Frame()
	settles(t, tt, "near full")
	if !tt.HasText("85%") || !tt.HasText("Older messages will be summarized") {
		t.Errorf("texts %q", tt.Texts())
	}
	if short(950) != "950" || short(12_400) != "12.4k" || short(1_000_000) != "1M" {
		t.Errorf("short %q %q %q", short(950), short(12_400), short(1_000_000))
	}
	panics(t, "over limit", func() { render(func(c *ui.Context) { ContextMeter(c, 3, 2, ContextMeterOptions{}) }) })
}

func TestCostMeter(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { CostMeter(c, 20.5, 50, "USD", CostMeterOptions{}) }), 400, 250)
	settles(t, tt, "cost")
	for _, s := range []string{"USD 20.50", "USD 50.00", "41% used", "Budget"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	panics(t, "zero budget", func() { render(func(c *ui.Context) { CostMeter(c, 1, 0, "USD", CostMeterOptions{}) }) })
}
