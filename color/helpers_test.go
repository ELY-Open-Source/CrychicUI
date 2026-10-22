package color

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func find(t *testing.T, tt *ui.Tester, s string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(s)
	if !ok {
		t.Fatalf("%q not shown in %q", s, tt.Texts())
	}
	return r
}

func click(t *testing.T, tt *ui.Tester, s string) {
	t.Helper()
	if err := tt.Click(s); err != nil {
		t.Fatalf("click %q: %v", s, err)
	}
	tt.Frame()
}

func panics(t *testing.T, what string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", what)
		}
	}()
	fn()
}

func render(view func(c *ui.Context)) { ui.Render(view, 400, 600, 1) }

func settles(t *testing.T, tt *ui.Tester, what string) {
	t.Helper()
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Errorf("%s still moving", what)
	}
}

// baseline fills Recent with Swipe, which color never plays.
func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Swipe)
	}
	return sound.Recent()
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	after := sound.Recent()
	var got []sound.Cue
	if len(before) == 32 && !slices.Equal(before, after) {
		i := 0
		for i < len(after) && after[i] == sound.Swipe {
			i++
		}
		got = after[i:]
	}
	if !slices.Equal(got, want) {
		t.Errorf("cues %v, want %v", got, want)
	}
}

// close reports channels within one of each other.
func close(a, b ui.Color) bool {
	d := func(x, y uint8) bool { return max(x, y)-min(x, y) <= 1 }
	return d(a.R, b.R) && d(a.G, b.G) && d(a.B, b.B) && a.A == b.A
}

func TestColorMath(t *testing.T) {
	for _, c := range []ui.Color{ui.RGB(255, 0, 0), ui.RGB(18, 52, 86), ui.RGB(250, 200, 10), ui.RGB(128, 128, 128), ui.RGB(0, 0, 0)} {
		if h, s, v := hsv(c); !close(fromHSV(h, s, v, 255), c) {
			t.Errorf("hsv round trip of %v: %v", c, fromHSV(h, s, v, 255))
		}
		if h, s, l := hsl(c); !close(fromHSL(h, s, l, 255), c) {
			t.Errorf("hsl round trip of %v: %v", c, fromHSL(h, s, l, 255))
		}
		if l, ch, h := oklch(c); !inSRGB(l, ch, h) {
			t.Errorf("%v reads outside sRGB", c)
		}
	}
	if h, s, l := hsl(ui.RGB(0, 128, 255)); math.Abs(h-210) > 0.5 || math.Abs(s-1) > 0.01 || math.Abs(l-0.5) > 0.01 {
		t.Errorf("hsl of #0080FF: %v %v %v", h, s, l)
	}
	if l, ch, _ := oklch(ui.RGB(255, 255, 255)); math.Abs(l-1) > 1e-3 || ch > 1e-3 {
		t.Errorf("oklch of white: %v %v", l, ch)
	}
	if l, ch, h := oklch(ui.RGB(255, 0, 0)); math.Abs(l-0.628) > 0.002 || math.Abs(ch-0.2577) > 0.002 || math.Abs(h-29.23) > 0.2 {
		t.Errorf("oklch of red: %v %v %v", l, ch, h)
	}
	if inSRGB(0.7, 0.35, 150) {
		t.Error("oklch(0.7 0.35 150) reads inside sRGB")
	}
	if r := ratio(ui.RGB(0, 0, 0), ui.RGB(255, 255, 255)); math.Abs(r-21) > 1e-9 {
		t.Errorf("black on white: %v", r)
	}
	if r := ratio(ui.RGB(118, 118, 118), ui.RGB(255, 255, 255)); math.Abs(r-4.54) > 0.01 {
		t.Errorf("#767676 on white: %v", r)
	}
	if got := hexOf(ui.RGBA(18, 52, 86, 0.5)); got != "#12345680" {
		t.Errorf("hexOf: %s", got)
	}
	for s, ok := range map[string]bool{"#123": true, "abcd": true, "#A1B2C3": true, "#a1b2c3d4": true, "#12": false, "#12345": false, "#zzz": false, "": false} {
		if _, got := parseHex(s); got != ok {
			t.Errorf("parseHex(%q) = %v", s, got)
		}
	}
	if g := seen(ui.RGB(128, 128, 128), 0); !close(g, ui.RGB(128, 128, 128)) {
		t.Errorf("protan gray: %v", g)
	}
	if r, d := seen(ui.RGB(255, 0, 0), 0), seen(ui.RGB(0, 160, 0), 0); r.R > r.G+40 || d.R < 40 {
		t.Errorf("protan red %v, green %v", r, d)
	}
}
