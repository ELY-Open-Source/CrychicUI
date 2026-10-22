package color

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// retype replaces the focused field's text.
func retype(tt *ui.Tester, s string) {
	tt.Command("selectAll")
	tt.Type(s)
	tt.Frame()
}

func TestPickerCuesOnClose(t *testing.T) {
	value, changed := ui.RGB(200, 40, 40), 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Picker(c, "Brand", &value).Changed() {
			changed++
		}
	}, 500, 600)
	before := baseline()
	click(t, tt, "Brand")
	click(t, tt, "Hex")
	retype(tt, "00FF00")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if value != ui.RGB(0, 255, 0) || changed == 0 {
		t.Fatalf("value %v, changed %d", value, changed)
	}
	cues(t, before)
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	cues(t, before, sound.Select)
	click(t, tt, "Brand")
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	cues(t, before, sound.Select)
}

func TestWellOpens(t *testing.T) {
	value := ui.RGB(10, 120, 200)
	tt := ui.NewTester(func(c *ui.Context) { Well(c, "Accent", &value) }, 500, 600)
	before := baseline()
	click(t, tt, "Accent")
	find(t, tt, "Hex")
	cues(t, before, sound.Open)
	settles(t, tt, "Well")
}

func TestHexField(t *testing.T) {
	value := ui.RGB(255, 0, 0)
	tt := ui.NewTester(func(c *ui.Context) { HexField(c, "Hex color", &value, HexFieldOptions{}) }, 500, 200)
	before := baseline()
	click(t, tt, "Hex color")
	retype(tt, "#00ff0080")
	if value != ui.RGBA(0, 255, 0, 128.0/255) {
		t.Errorf("value %v", value)
	}
	cues(t, before, sound.Tick)
	retype(tt, "#zz")
	if !tt.HasText("Use #RGB, #RRGGBB or #RRGGBBAA") || value.G != 255 {
		t.Errorf("bad hex: value %v, texts %q", value, tt.Texts())
	}
}

func TestRGBFields(t *testing.T) {
	value := ui.RGB(10, 20, 30)
	tt := ui.NewTester(func(c *ui.Context) { RGBFields(c, "Color", &value, RGBFieldsOptions{}) }, 600, 200)
	before := baseline()
	click(t, tt, "R")
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if value != ui.RGB(11, 20, 30) {
		t.Errorf("value %v", value)
	}
	cues(t, before, sound.Tick)
}

func TestHSLFieldsKeepHue(t *testing.T) {
	value := ui.RGB(0, 128, 255)
	tt := ui.NewTester(func(c *ui.Context) { HSLFields(c, "Color", &value, HSLFieldsOptions{}) }, 600, 200)
	click(t, tt, "S")
	retype(tt, "0")
	if value.R != value.G || value.G != value.B {
		t.Fatalf("saturation 0 is not gray: %v", value)
	}
	retype(tt, "100")
	if h, _, _ := hsl(value); math.Abs(h-210) > 2 {
		t.Errorf("hue lost through gray: %v (hue %v)", value, h)
	}
}

func TestOklchFieldsGamut(t *testing.T) {
	value := ui.RGB(0, 128, 255)
	tt := ui.NewTester(func(c *ui.Context) { OklchFields(c, "Color", &value, OklchFieldsOptions{}) }, 600, 200)
	note := "Outside sRGB; most screens show the nearest color"
	if tt.HasText(note) {
		t.Fatal("an sRGB color shows the gamut note")
	}
	click(t, tt, "C")
	retype(tt, "0.37")
	if !tt.HasText(note) {
		t.Errorf("chroma 0.37 shows no note: %q", tt.Texts())
	}
}

func TestHueRing(t *testing.T) {
	hue, changed := 0.0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if HueRing(c, "Hue", &hue).Changed() {
			changed++
		}
	}, 300, 300)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	for range 16 {
		tt.Key(0, ui.KeyRight)
	}
	tt.Frame()
	if hue != 16 || changed == 0 {
		t.Errorf("keys: hue %v, changed %d", hue, changed)
	}
	cues(t, before, slices.Repeat([]sound.Cue{sound.Tick}, 16)...)
	r := find(t, tt, "Hue")
	tt.Press(r.X+r.W-6, r.Y+r.H/2)
	tt.Frame()
	tt.Release(r.X+r.W-6, r.Y+r.H/2)
	tt.Frame()
	if math.Abs(hue-90) > 3 {
		t.Errorf("press at the right: hue %v", hue)
	}
	panics(t, "hue 360", func() { render(func(c *ui.Context) { h := 360.0; HueRing(c, "x", &h) }) })
}

func TestShadePadKeepsHue(t *testing.T) {
	value := fromHSV(200, 0.5, 0.5, 255)
	tt := ui.NewTester(func(c *ui.Context) { ShadePad(c, "Shade", &value) }, 400, 300)
	r := find(t, tt, "Shade")
	at := func(fx, fy float32) {
		tt.Press(r.X+r.W*fx, r.Y+r.H*fy)
		tt.Frame()
		tt.Release(r.X+r.W*fx, r.Y+r.H*fy)
		tt.Frame()
	}
	baseline()
	at(0.999, 0.001)
	if _, s, v := hsv(value); s < 0.98 || v < 0.98 {
		t.Errorf("top right: %v", value)
	}
	at(0.001, 0.999)
	at(0.5, 0.5)
	if h, _, _ := hsv(value); math.Abs(h-200) > 3 {
		t.Errorf("hue after black: %v (%v)", h, value)
	}
	ticks := 0
	for _, q := range sound.Recent() {
		if q == sound.Tick {
			ticks++
		} else if q != sound.Swipe {
			t.Errorf("ShadePad played %v", q)
		}
	}
	if ticks < 2 {
		t.Errorf("%d Ticks for three moves", ticks)
	}
}

func TestAlphaSlider(t *testing.T) {
	value := ui.RGB(20, 40, 60)
	tt := ui.NewTester(func(c *ui.Context) { AlphaSlider(c, "Opacity", &value) }, 400, 100)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	for range 10 {
		tt.Key(0, ui.KeyLeft)
	}
	tt.Frame()
	if value.A < 229 || value.A > 231 || value.R != 20 {
		t.Errorf("value %v", value)
	}
	cues(t, before, slices.Repeat([]sound.Cue{sound.Tick}, 10)...)
}

func TestSwatchGrid(t *testing.T) {
	colors := []ui.Color{ui.RGB(255, 0, 0), ui.RGB(0, 255, 0), ui.RGB(0, 0, 255)}
	value, changed := colors[0], 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SwatchGrid(c, "Swatches", colors, &value, SwatchGridOptions{}).Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	click(t, tt, "#00FF00")
	if value != colors[1] || changed != 1 {
		t.Errorf("value %v, changed %d", value, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "short names", func() {
		render(func(c *ui.Context) { SwatchGrid(c, "x", colors, &value, SwatchGridOptions{Names: []string{"a"}}) })
	})
}

func TestGradientEditor(t *testing.T) {
	stops := []ui.GradientStop{{Color: ui.RGB(255, 0, 0), At: 0}, {Color: ui.RGB(0, 0, 255), At: 1}}
	tt := ui.NewTester(func(c *ui.Context) { GradientEditor(c, "Fill", &stops, GradientEditorOptions{}) }, 500, 200)
	before := baseline()
	bar := find(t, tt, "Fill")
	tt.ClickAt(bar.X+bar.W/2, bar.Y+bar.H/2)
	tt.Frame()
	if len(stops) != 3 || math.Abs(float64(stops[2].At-0.5)) > 0.02 {
		t.Fatalf("stops %v", stops)
	}
	cues(t, before, sound.Pop)
	h := find(t, tt, "Stop at 0%")
	tt.Press(h.X+h.W/2, h.Y+h.H/2)
	for i := 1; i <= 8; i++ {
		tt.Move(h.X+h.W/2+bar.W/4*float32(i)/8, h.Y+h.H/2)
		tt.Advance(16 * time.Millisecond)
	}
	tt.Release(h.X+h.W/2+bar.W/4, h.Y+h.H/2)
	tt.Frame()
	if math.Abs(float64(stops[0].At-0.25)) > 0.03 {
		t.Errorf("dragged stop at %v", stops[0].At)
	}
	cues(t, before, sound.Pop, sound.Drop)
	panics(t, "one stop", func() {
		render(func(c *ui.Context) { GradientEditor(c, "x", &[]ui.GradientStop{{}}, GradientEditorOptions{}) })
	})
}

func TestRecentColors(t *testing.T) {
	a, b := ui.RGB(255, 0, 0), ui.RGB(0, 0, 255)
	colors := []ui.Color{a, b, a}
	value := a
	tt := ui.NewTester(func(c *ui.Context) { RecentColors(c, "Recent", colors, &value, RecentColorsOptions{}) }, 500, 100)
	before := baseline()
	click(t, tt, "#0000FF")
	if value != b {
		t.Errorf("value %v", value)
	}
	cues(t, before, sound.Select)
	colors = append([]ui.Color{ui.RGB(0, 200, 0)}, colors...)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("the new color did not slide in")
	}
	settles(t, tt, "RecentColors")
	colors = nil
	tt.Frame()
	if !tt.HasText("No recent colors") {
		t.Errorf("empty %q", tt.Texts())
	}
}

func TestHarmony(t *testing.T) {
	base := ui.RGB(220, 40, 40)
	value := base
	tt := ui.NewTester(func(c *ui.Context) { Harmony(c, "Harmony", base, &value, HarmonyOptions{}) }, 500, 300)
	l, ch, h := oklch(base)
	comp := ui.Oklch(float32(l), float32(ch), float32(h+180)).SRGB()
	before := baseline()
	click(t, tt, hexOf(comp))
	if value != comp {
		t.Errorf("value %v, want %v", value, comp)
	}
	if ch2, h2 := func() (float64, float64) { _, c2, hh := oklch(comp); return c2, hh }(); math.Abs(math.Mod(h2-h+360, 360)-180) > 6 || ch2 < 0.05 {
		t.Errorf("complement hue %v from %v", h2, h)
	}
	cues(t, before, sound.Select)
}

func TestEyedropper(t *testing.T) {
	colors := []ui.Color{ui.RGB(255, 0, 0), ui.RGB(0, 255, 0), ui.RGB(0, 0, 255), ui.RGB(0, 0, 0)}
	value := colors[0]
	tt := ui.NewTester(func(c *ui.Context) { Eyedropper(c, "Strip", colors, &value) }, 600, 200)
	r := find(t, tt, "Strip")
	before := baseline()
	tt.ClickAt(r.X+r.W*0.6, r.Y+r.H/2)
	tt.Frame()
	if value != colors[2] {
		t.Errorf("click at 60%%: %v", value)
	}
	cues(t, before, sound.Pop)
}

// outside fails when anything is painted outside s's box.
func outside(t *testing.T, tt *ui.Tester, s string) {
	t.Helper()
	r := find(t, tt, s)
	img := tt.Image()
	bg := img.RGBAAt(0, 0)
	for y := int(r.Y) - 3; y < int(r.Y+r.H)+3; y++ {
		for _, x := range []int{int(r.X) - 2, int(r.X) - 1, int(r.X + r.W), int(r.X+r.W) + 1} {
			if img.RGBAAt(x, y) != bg {
				t.Errorf("%s paints outside its box at %d,%d", s, x, y)
				return
			}
		}
	}
}

func TestPaintStaysInside(t *testing.T) {
	hue := 90.0
	tt := ui.NewTester(func(c *ui.Context) { ui.Box(c).Padding(20).Children(func() { HueRing(c, "Hue", &hue) }) }, 300, 300)
	settles(t, tt, "HueRing")
	outside(t, tt, "Hue")
	value := ui.RGB(20, 40, 60)
	ta := ui.NewTester(func(c *ui.Context) { ui.Box(c).Padding(20).Children(func() { AlphaSlider(c, "Opacity", &value) }) }, 400, 100)
	settles(t, ta, "AlphaSlider")
	outside(t, ta, "Opacity")
	stops := []ui.GradientStop{{Color: ui.RGB(255, 0, 0), At: 0}, {Color: ui.RGB(0, 0, 255), At: 1}}
	tg := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Padding(20).Children(func() { GradientEditor(c, "Fill", &stops, GradientEditorOptions{}) })
	}, 500, 200)
	settles(t, tg, "GradientEditor")
	root, first, last := find(t, tg, "Fill"), find(t, tg, "Stop at 0%"), find(t, tg, "Stop at 100%")
	if first.X < root.X-1 || last.X+last.W > root.X+root.W+first.W+1 {
		t.Errorf("handles %v and %v hang outside the bar %v", first, last, root)
	}
}

func TestEyedropperPaletteShrinks(t *testing.T) {
	colors := []ui.Color{ui.RGB(255, 0, 0), ui.RGB(0, 255, 0), ui.RGB(0, 0, 255), ui.RGB(0, 0, 0)}
	value := colors[0]
	tt := ui.NewTester(func(c *ui.Context) { Eyedropper(c, "Strip", colors, &value) }, 600, 200)
	r := find(t, tt, "Strip")
	tt.Move(r.X+r.W-4, r.Y+r.H/2)
	tt.Frame()
	tt.Move(r.X, r.Y+r.H+60)
	tt.Frame()
	colors = colors[:1]
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if value != colors[0] {
		t.Errorf("after shrink picked %v", value)
	}
}

func TestShadePadHueHolds(t *testing.T) {
	value := ui.RGB(0, 170, 255)
	tt := ui.NewTester(func(c *ui.Context) { ShadePad(c, "Shade", &value) }, 400, 300)
	r := find(t, tt, "Shade")
	at := func(s, v float32) {
		tt.Press(r.X+r.W*s, r.Y+r.H*(1-v))
		tt.Frame()
		tt.Release(r.X+r.W*s, r.Y+r.H*(1-v))
		tt.Frame()
	}
	h0, _, _ := hsv(value)
	at(0.01, 0.5)
	at(0.8, 0.5)
	if h, _, _ := hsv(value); math.Abs(h-h0) > 1.5 {
		t.Errorf("hue drifted from %v to %v", h0, h)
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	cues(t, before, sound.Tick)
}
