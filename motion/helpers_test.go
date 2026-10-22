package motion

import (
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
		t.Fatal(err)
	}
}

// settles checks that something moves, then rests.
func settles(t *testing.T, tt *ui.Tester, what string) {
	t.Helper()
	if !tt.FramePending() {
		t.Errorf("%s: no animation", what)
	}
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Errorf("%s: animation never settled", what)
	}
}

// loops checks that frames keep coming.
func loops(t *testing.T, tt *ui.Tester, what string) {
	t.Helper()
	tt.Advance(3 * time.Second)
	if !tt.FramePending() {
		t.Errorf("%s: loop stopped", what)
	}
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

func render(view func(c *ui.Context)) { ui.Render(view, 200, 200, 1) }

// dot is a labeled accent square.
func dot(c *ui.Context) { ui.Box(c).Size(40, 40).Background(c.Theme().Accent).Label("dot") }

// tinted reports whether the frame pixel is exactly k.
func tinted(tt *ui.Tester, x, y int, k ui.Color) bool {
	p := tt.Image().RGBAAt(x, y)
	return p.R == k.R && p.G == k.G && p.B == k.B
}

// inked counts pixels in r that differ from the background.
func inked(tt *ui.Tester, r ui.Rect) int {
	bg, img, n := ui.LightTheme().Background, tt.Image(), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); p.R != bg.R || p.G != bg.G || p.B != bg.B {
				n++
			}
		}
	}
	return n
}

// baseline fills Recent with Tap so later cues show.
func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Tap)
	}
	return sound.Recent()
}

func played(before []sound.Cue) []sound.Cue {
	after := sound.Recent()
	if len(before) != 32 || slices.Equal(before, after) {
		return nil
	}
	i := 0
	for i < len(after) && after[i] == sound.Tap {
		i++
	}
	return after[i:]
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	if got := played(before); !slices.Equal(got, want) {
		t.Errorf("cues %v, want %v", got, want)
	}
}
