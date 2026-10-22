package overlay

import (
	"fmt"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// lit reports whether label's center shows the page background.
func lit(t *testing.T, tt *ui.Tester, label string) bool {
	t.Helper()
	r, bg := find(t, tt, label), ui.LightTheme().Background
	p := tt.Image().RGBAAt(int(r.X+r.W/2), int(r.Y+r.H/2))
	return p.R == bg.R && p.G == bg.G && p.B == bg.B
}

func TestTour(t *testing.T) {
	step, open, changed := 0, true, 0
	steps := []TourStep{{"Name it", "Type a name."}, {"Save it", "Keep your work."}}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Gap(200).Children(func() {
			a := ui.ButtonBase(c).Label("Name field").Size(120, 32)
			b := ui.ButtonBase(c).Label("Save button").Size(120, 32)
			if Tour(c, "Intro", []ui.Element{a, b}, steps, &step, &open, TourOptions{}).Changed() {
				changed++
			}
		})
	}, 600, 600)
	settles(t, tt, "Tour")
	if !lit(t, tt, "Name field") || lit(t, tt, "Save button") {
		t.Fatalf("step 1: the hole is not on the first anchor")
	}
	before := baseline()
	click(t, tt, "Next")
	if step != 1 || !tt.FramePending() {
		t.Fatalf("next: step %d, moving %v", step, tt.FramePending())
	}
	settles(t, tt, "next")
	if lit(t, tt, "Name field") || !lit(t, tt, "Save button") {
		t.Errorf("step 2: the hole did not move to the second anchor")
	}
	click(t, tt, "Done")
	if open || changed != 2 {
		t.Errorf("done: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Tick, sound.Tick)
	panics(t, "anchor per step", func() {
		render(func(c *ui.Context) { s, o := 0, true; Tour(c, "x", nil, steps, &s, &o, TourOptions{}) })
	})
}

func TestPreviewPane(t *testing.T) {
	key := -1
	tt := ui.NewTester(func(c *ui.Context) {
		PreviewPane(c, "Preview", key, PreviewPaneOptions{}, func() { ui.Text(c, fmt.Sprint("Item ", key)) })
	}, 400, 300)
	if !tt.HasText("Select an item to preview") {
		t.Fatalf("empty %q", tt.Texts())
	}
	key = 1
	tt.Frame()
	settles(t, tt, "first item")
	key = 2
	tt.Frame()
	// Texts leaves out the page fading away.
	if !tt.HasText("Item 2") || !tt.FramePending() {
		t.Fatalf("no cross-fade: %q", tt.Texts())
	}
	settles(t, tt, "cross-fade")
	if tt.HasText("Item 1") {
		t.Errorf("old item stayed: %q", tt.Texts())
	}
	panics(t, "bad key", func() { render(func(c *ui.Context) { PreviewPane(c, "x", -2, PreviewPaneOptions{}, func() {}) }) })
}
