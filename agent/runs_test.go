package agent

import (
	"image/color"
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCitation(t *testing.T) {
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Citation(c, 1, "https://example.com/docs").Changed() {
			changed++
		}
	}), 400, 200)
	settles(t, tt, "citation")
	before := baseline()
	click(t, tt, "1")
	if changed != 1 {
		t.Errorf("changed %d", changed)
	}
	if got := tt.OpenedURLs(); len(got) != 1 || got[0] != "https://example.com/docs" {
		t.Errorf("opened %v", got)
	}
	cues(t, before, sound.Open)
	panics(t, "zero number", func() { render(func(c *ui.Context) { Citation(c, 0, "https://x.co") }) })
	panics(t, "empty url", func() { render(func(c *ui.Context) { Citation(c, 1, "") }) })
}

func TestSourceCard(t *testing.T) {
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if SourceCard(c, "MyGo docs", "mygo.egoist.dev", "Native UI for Go.").Changed() {
			changed++
		}
	}), 400, 300)
	settles(t, tt, "card")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("MyGo docs") {
		t.Fatal("card does not take the focus")
	}
	before := baseline()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if changed != 1 {
		t.Errorf("changed %d after Enter", changed)
	}
	click(t, tt, "MyGo docs")
	if changed != 2 {
		t.Errorf("changed %d after pointer", changed)
	}
	cues(t, before, sound.Tap, sound.Tap)
	panics(t, "empty host", func() { render(func(c *ui.Context) { SourceCard(c, "t", "", "") }) })
}

func TestModelPicker(t *testing.T) {
	models := []Model{{Name: "GPT-6.1", Speed: "Fast", Cost: "$$$"}, {Name: "Local 8B", Speed: "Slow", Cost: "$"}}
	selected, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ModelPicker(c, "Models", models, &selected).Changed() {
			changed++
		}
	}), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("rows do not arrive")
	}
	settles(t, tt, "models")
	if !tt.HasText("Fast") || !tt.HasText("$$$") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Local 8B")
	if selected != 1 || changed != 1 {
		t.Fatalf("selected %d changed %d", selected, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "empty name", func() {
		render(func(c *ui.Context) { ModelPicker(c, "x", []Model{{Name: ""}}, &selected) })
	})
}

func TestTokenCount(t *testing.T) {
	used := 1248
	tt := ui.NewTester(padded(func(c *ui.Context) { TokenCount(c, used, 4096, TokenCountOptions{}) }), 400, 200)
	settles(t, tt, "count")
	if !tt.HasText("1,248") || !tt.HasText("4,096") || !tt.HasText("tokens") {
		t.Errorf("texts %q", tt.Texts())
	}
	used = 1500
	tt.Frame()
	if !tt.FramePending() {
		t.Error("digits do not roll")
	}
	settles(t, tt, "count")
	if !tt.HasText("1,500") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "negative", func() { render(func(c *ui.Context) { TokenCount(c, -1, 5, TokenCountOptions{}) }) })
	panics(t, "zero limit", func() { render(func(c *ui.Context) { TokenCount(c, 0, 0, TokenCountOptions{}) }) })
}

func TestContextMeter(t *testing.T) {
	used := 0.5
	var acc, warn color.RGBA
	tt := ui.NewTester(padded(func(c *ui.Context) {
		acc, warn = rgba(c.Theme().Accent), rgba(c.Theme().Warning)
		ContextMeter(c, used, ContextMeterOptions{})
	}), 400, 200)
	settles(t, tt, "ring")
	if !tt.HasText("50%") {
		t.Errorf("texts %q", tt.Texts())
	}
	r := find(t, tt, "Context used")
	rad := r.W/2 - 6
	if !at(tt, r.X+r.W/2+rad, r.Y+r.H/2, acc) {
		t.Error("no accent arc at half full")
	}
	used = 0.9
	tt.Frame()
	settles(t, tt, "ring")
	if !at(tt, r.X+r.W/2+rad, r.Y+r.H/2, warn) {
		t.Error("ring does not warn near full")
	}
	short := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(50).Children(func() { ContextMeter(c, 1, ContextMeterOptions{}).Height(20) })
	}, 400, 300)
	short.Frame()
	short.Advance(2 * time.Second)
	if bg := short.Image().RGBAAt(82, 34); bg.R < 240 && bg.G < 240 {
		t.Errorf("ring paints outside its box: %v", bg)
	}
	panics(t, "over one", func() { render(func(c *ui.Context) { ContextMeter(c, 1.5, ContextMeterOptions{}) }) })
}

func TestCostMeter(t *testing.T) {
	amount := 12.4
	tt := ui.NewTester(padded(func(c *ui.Context) { CostMeter(c, amount, "USD", CostMeterOptions{}) }), 400, 200)
	settles(t, tt, "cost")
	if !tt.HasText("USD 12.40") || !tt.HasText("Spend") {
		t.Errorf("texts %q", tt.Texts())
	}
	amount = 13
	tt.Frame()
	if !tt.FramePending() {
		t.Error("digits do not roll")
	}
	settles(t, tt, "cost")
	if !tt.HasText("USD 13.00") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "NaN", func() { render(func(c *ui.Context) { CostMeter(c, math.NaN(), "USD", CostMeterOptions{}) }) })
	panics(t, "empty currency", func() { render(func(c *ui.Context) { CostMeter(c, 1, "", CostMeterOptions{}) }) })
}
