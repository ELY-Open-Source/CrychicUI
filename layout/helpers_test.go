package layout

import (
	"image"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func near(img *image.RGBA, x, y int, want ui.Color) bool {
	p := img.RGBAAt(x, y)
	d := func(a, b uint8) bool { return a-b < 4 || b-a < 4 }
	return d(p.R, want.R) && d(p.G, want.G) && d(p.B, want.B)
}

func lum(img *image.RGBA, x, y int) int {
	p := img.RGBAAt(x, y)
	return int(p.R)*2126 + int(p.G)*7152 + int(p.B)*722
}

func count(cue sound.Cue) int {
	n := 0
	for _, c := range sound.Recent() {
		if c == cue {
			n++
		}
	}
	return n
}

func find(t *testing.T, tt *ui.Tester, s string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(s)
	if !ok {
		t.Fatalf("%q not shown in %q", s, tt.Texts())
	}
	return r
}

// settles checks an animation runs, then stops.
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

func panics(t *testing.T, what string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", what)
		}
	}()
	fn()
}

func abs(v float32) float32 { return max(v, -v) }
