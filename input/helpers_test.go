package input

import (
	"image"
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

func near(img *image.RGBA, x, y int, want ui.Color) bool {
	p := img.RGBAAt(x, y)
	d := func(a, b uint8) bool { return a-b < 4 || b-a < 4 }
	return d(p.R, want.R) && d(p.G, want.G) && d(p.B, want.B)
}

func rectChanged(a, b *image.RGBA, x0, y0, x1, y1 int) bool {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if a.RGBAAt(x, y) != b.RGBAAt(x, y) {
				return true
			}
		}
	}
	return false
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

func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Swipe)
	}
	return sound.Recent()
}

func appended(before []sound.Cue, want ...sound.Cue) bool {
	if len(before) != 32 {
		return false
	}
	filler := before[0]
	for _, cue := range before {
		if cue != filler {
			return false
		}
	}
	for _, cue := range want {
		if cue == filler {
			return false
		}
	}
	after := sound.Recent()
	if slices.Equal(before, after) {
		return len(want) == 0
	}
	start := 0
	for start < len(after) && after[start] == filler {
		start++
	}
	added := after[start:]
	if len(want) == 0 {
		return len(added) == 0
	}
	matched := 0
	for _, cue := range added {
		if matched < len(want) && cue == want[matched] {
			matched++
		}
	}
	return matched == len(want)
}
