package color

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestPaletteGrowsLightToDark(t *testing.T) {
	base := ui.RGB(40, 110, 220)
	var shades []string
	tt := ui.NewTester(func(c *ui.Context) {
		_, ch, h := oklch(base)
		shades = shades[:0]
		for i := range 5 {
			shades = append(shades, hexOf(ui.Oklch(float32(0.95-0.7*float64(i)/4), float32(ch), float32(h)).SRGB()))
		}
		Palette(c, base, PaletteOptions{Steps: 5})
	}, 600, 120)
	if !tt.FramePending() {
		t.Error("steps did not grow in")
	}
	settles(t, tt, "Palette")
	first, last := find(t, tt, shades[0]), find(t, tt, shades[4])
	if first.X >= last.X {
		t.Error("light is not first")
	}
	if a, _ := parseHex(shades[0]); luminance(a) < 0.7 {
		t.Errorf("lightest %s is dark", shades[0])
	}
	panics(t, "one step", func() { render(func(c *ui.Context) { Palette(c, base, PaletteOptions{Steps: 1}) }) })
}

func TestContrastMarks(t *testing.T) {
	fg, bg := ui.RGB(118, 118, 118), ui.RGB(255, 255, 255)
	tt := ui.NewTester(func(c *ui.Context) { Contrast(c, fg, bg, ContrastOptions{}) }, 500, 200)
	if !tt.HasText("4.54:1") {
		t.Fatalf("texts %q", tt.Texts())
	}
	find(t, tt, "AA passes")
	find(t, tt, "AAA fails")
	fg = ui.RGB(0, 0, 0)
	tt.Frame()
	settles(t, tt, "Contrast")
	if !tt.HasText("21.00:1") {
		t.Errorf("black on white %q", tt.Texts())
	}
	find(t, tt, "AAA passes")
}

func TestColorNameNearest(t *testing.T) {
	names := []NamedColor{{"Red", ui.RGB(255, 0, 0)}, {"Green", ui.RGB(0, 128, 0)}, {"Blue", ui.RGB(0, 0, 255)}}
	value := ui.RGB(240, 20, 30)
	tt := ui.NewTester(func(c *ui.Context) { ColorName(c, value, names) }, 400, 100)
	if !tt.HasText("Red") {
		t.Fatalf("texts %q", tt.Texts())
	}
	value = ui.RGB(20, 150, 40)
	tt.Frame()
	if !tt.HasText("Green") || !tt.FramePending() {
		t.Errorf("after change %q, fading %v", tt.Texts(), tt.FramePending())
	}
	panics(t, "no names", func() { render(func(c *ui.Context) { ColorName(c, value, nil) }) })
}

func TestBeforeAfter(t *testing.T) {
	after := ui.RGB(0, 0, 255)
	tt := ui.NewTester(func(c *ui.Context) { BeforeAfter(c, ui.RGB(255, 0, 0), after, BeforeAfterOptions{}) }, 400, 200)
	settles(t, tt, "BeforeAfter")
	if !tt.HasText("Before") || !tt.HasText("#FF0000") || !tt.HasText("#0000FF") {
		t.Fatalf("texts %q", tt.Texts())
	}
	after = ui.RGB(0, 200, 0)
	tt.Frame()
	if !tt.FramePending() || !tt.HasText("#00C800") {
		t.Errorf("after did not animate: %q", tt.Texts())
	}
}

func TestThemePreviewFollowsTheme(t *testing.T) {
	th := ui.LightTheme()
	tt := ui.NewTester(func(c *ui.Context) { ThemePreview(c, "Preview", th) }, 400, 300)
	r := find(t, tt, "Preview")
	px := func() ui.Color {
		p := tt.Image().RGBAAt(int(r.X+r.W-8), int(r.Y+r.H-8))
		return ui.RGB(p.R, p.G, p.B)
	}
	if !close(px(), th.Background) {
		t.Fatalf("light corner %v, want %v", px(), th.Background)
	}
	th = ui.DarkTheme()
	tt.Frame()
	if !tt.FramePending() {
		t.Error("theme change did not animate")
	}
	settles(t, tt, "ThemePreview")
	if !close(px(), th.Background) {
		t.Errorf("dark corner %v, want %v", px(), th.Background)
	}
}

func TestColorBlindRows(t *testing.T) {
	red := ui.RGB(255, 0, 0)
	tt := ui.NewTester(func(c *ui.Context) { ColorBlind(c, []ui.Color{red}, ColorBlindOptions{}) }, 500, 300)
	settles(t, tt, "ColorBlind")
	for k, name := range []string{"Protan", "Deutan", "Tritan"} {
		find(t, tt, name)
		find(t, tt, hexOf(seen(red, k)))
	}
	if hexOf(seen(red, 0)) == hexOf(red) {
		t.Error("protan red looks unchanged")
	}
	panics(t, "no colors", func() { render(func(c *ui.Context) { ColorBlind(c, nil, ColorBlindOptions{}) }) })
}
