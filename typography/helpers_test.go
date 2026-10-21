package typography

import (
	"image"
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

func count(cue sound.Cue) int {
	n := 0
	for _, c := range sound.Recent() {
		if c == cue {
			n++
		}
	}
	return n
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

func near(img *image.RGBA, x, y int, want ui.Color) bool {
	p := img.RGBAAt(x, y)
	d := func(a, b uint8) bool { return a-b < 4 || b-a < 4 }
	return d(p.R, want.R) && d(p.G, want.G) && d(p.B, want.B)
}

func textCount(tt *ui.Tester, want string) int {
	n := 0
	for _, s := range tt.Texts() {
		if s == want {
			n++
		}
	}
	return n
}
