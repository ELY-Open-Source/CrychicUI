package fx

import (
	"hash/fnv"
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

// shot hashes the last frame.
func shot(tt *ui.Tester) uint64 {
	h := fnv.New64a()
	h.Write(tt.Image().Pix)
	return h.Sum64()
}

// moves checks the frame changes soon and keeps changing.
func moves(t *testing.T, tt *ui.Tester, what string) {
	t.Helper()
	a, changed := shot(tt), false
	for range 6 {
		tt.Advance(100 * time.Millisecond)
		changed = changed || shot(tt) != a
	}
	tt.Advance(3 * time.Second)
	if !changed || !tt.FramePending() {
		t.Errorf("%s: changed %v, still drawing %v", what, changed, tt.FramePending())
	}
}

// holds checks reduced motion stops the frames and the picture.
func holds(t *testing.T, tt *ui.Tester, what string) {
	t.Helper()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	tt.Advance(time.Second)
	a := shot(tt)
	tt.Advance(time.Second)
	if tt.FramePending() || shot(tt) != a {
		t.Errorf("%s: moves under reduced motion", what)
	}
}

// count is how many pixels in r are exactly k.
func count(tt *ui.Tester, r ui.Rect, k ui.Color) int {
	img, n := tt.Image(), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); p.R == k.R && p.G == k.G && p.B == k.B {
				n++
			}
		}
	}
	return n
}

// ink counts pixels in r off the background.
func ink(tt *ui.Tester, r ui.Rect) int {
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

// baseline fills Recent with Error, which fx never plays.
func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Error)
	}
	return sound.Recent()
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	after := sound.Recent()
	var got []sound.Cue
	if len(before) == 32 && !slices.Equal(before, after) {
		i := 0
		for i < len(after) && after[i] == sound.Error {
			i++
		}
		got = after[i:]
	}
	if !slices.Equal(got, want) {
		t.Errorf("cues %v, want %v", got, want)
	}
}
