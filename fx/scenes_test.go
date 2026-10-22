package fx

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestScenesMove(t *testing.T) {
	hi := func(c *ui.Context) func() { return func() { ui.Text(c, "Hi") } }
	scenes := []struct {
		name string
		view func(c *ui.Context) ui.Element
	}{
		{"GridBackdrop", func(c *ui.Context) ui.Element { return GridBackdrop(c, hi(c)) }},
		{"DotBackdrop", func(c *ui.Context) ui.Element { return DotBackdrop(c, hi(c)) }},
		{"GradientBackdrop", func(c *ui.Context) ui.Element { return GradientBackdrop(c, hi(c)) }},
		{"Aurora", func(c *ui.Context) ui.Element { return Aurora(c, hi(c)) }},
		{"Starfield", func(c *ui.Context) ui.Element { return Starfield(c, hi(c)) }},
		{"Particles", func(c *ui.Context) ui.Element { return Particles(c, hi(c)) }},
		{"Snowfall", func(c *ui.Context) ui.Element { return Snowfall(c, hi(c)) }},
		{"Rainfall", func(c *ui.Context) ui.Element { return Rainfall(c, hi(c)) }},
		{"Bubbles", func(c *ui.Context) ui.Element { return Bubbles(c, hi(c)) }},
		{"Meteors", func(c *ui.Context) ui.Element { return Meteors(c, hi(c)) }},
		{"Sparkles", func(c *ui.Context) ui.Element { return Sparkles(c, hi(c)) }},
		{"BorderBeam", func(c *ui.Context) ui.Element { return BorderBeam(c, hi(c)) }},
		{"ConicBorder", func(c *ui.Context) ui.Element { return ConicBorder(c, hi(c)) }},
		{"ShineBorder", func(c *ui.Context) ui.Element { return ShineBorder(c, hi(c)) }},
		{"RetroGrid", func(c *ui.Context) ui.Element { return RetroGrid(c, hi(c)) }},
		{"RippleBackdrop", func(c *ui.Context) ui.Element { return RippleBackdrop(c, hi(c)) }},
		{"FlowingText", func(c *ui.Context) ui.Element { return FlowingText(c, "Flowing") }},
		{"ShinyText", func(c *ui.Context) ui.Element { return ShinyText(c, "Shiny") }},
	}
	for _, s := range scenes {
		text := s.name == "FlowingText" || s.name == "ShinyText"
		box := "scene"
		if text {
			box = strings.TrimSuffix(s.name, "Text")
		}
		t.Run(s.name, func(t *testing.T) { checkScene(t, s.name, box, s.view) })
	}
}

// checkScene checks the loop, the box and reduced motion.
func checkScene(t *testing.T, name, box string, view func(c *ui.Context) ui.Element) {
	t.Helper()
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			if e := view(c); box == "scene" {
				e.Label("scene")
			}
		})
	}, 300, 220)
	if box == "scene" {
		find(t, tt, "Hi")
	}
	moves(t, tt, name)
	r := find(t, tt, box)
	// One pixel more holds the antialiased edge.
	r = ui.Rect{X: r.X - 1, Y: r.Y - 1, W: r.W + 2, H: r.H + 2}
	if out := ink(tt, ui.Rect{W: 300, H: 220}) - ink(tt, r); out != 0 {
		t.Errorf("%s paints %d pixels outside its box", name, out)
	}
	holds(t, tt, name)
}

func TestGridBackdropDrawsHairlines(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { GridBackdrop(c, func() {}) }, 300, 200)
	if n := count(tt, ui.Rect{W: 300, H: 160}, ui.LightTheme().Border); n < 300 {
		t.Errorf("grid lines: %d border pixels", n)
	}
}

func TestGradientBackdropBlends(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { GradientBackdrop(c, func() {}) }, 300, 200)
	img, bg := tt.Image(), ui.LightTheme().Background
	a, b := img.RGBAAt(4, 4), img.RGBAAt(296, 156)
	if a == b || (a.R == bg.R && a.G == bg.G && a.B == bg.B) {
		t.Errorf("corners %v and %v", a, b)
	}
}

func TestRetroGridBelowHorizon(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { RetroGrid(c, func() {}) }, 300, 200)
	bg := ui.LightTheme().Background
	if count(tt, ui.Rect{W: 300, H: 60}, bg) != 300*60 || count(tt, ui.Rect{Y: 80, W: 300, H: 80}, bg) == 300*80 {
		t.Error("grid not confined below the horizon")
	}
}

func TestBackdropsNeedContent(t *testing.T) {
	panics(t, "nil content", func() { render(func(c *ui.Context) { Starfield(c, nil) }) })
	panics(t, "empty text", func() { render(func(c *ui.Context) { ShinyText(c, "") }) })
}
