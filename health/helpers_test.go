package health

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func padded(view func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) { ui.Column(c).Padding(10).Gap(8).Children(func() { view(c) }) }
}

func settles(t *testing.T, tt *ui.Tester) {
	t.Helper()
	tt.Advance(2 * time.Second)
}

func click(t *testing.T, tt *ui.Tester, name string) {
	t.Helper()
	if err := tt.Click(name); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
}

// baseline fills history with Delete, which health never plays.
func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Delete)
	}
	return sound.Recent()
}

// heard lists cues played since baseline.
func heard(before []sound.Cue) []sound.Cue {
	after := sound.Recent()
	if len(before) != 32 || slices.Equal(before, after) {
		return nil
	}
	i := 0
	for i < len(after) && after[i] == sound.Delete {
		i++
	}
	return after[i:]
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	if got := heard(before); !slices.Equal(got, want) {
		t.Errorf("cues %v, want %v", got, want)
	}
}

func exact(tt *ui.Tester, s string) bool { return slices.Contains(tt.Texts(), s) }

func panics(t *testing.T, want string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want %q", r, want)
		}
	}()
	ui.Render(view, 600, 600, 1)
}

func rectOf(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("%q missing: %q", name, tt.Texts())
	}
	return r
}

// swapUnderPress presses name, reorders, then releases there.
func swapUnderPress(t *testing.T, tt *ui.Tester, name string, reorder func()) {
	t.Helper()
	r := rectOf(t, tt, name)
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	reorder()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func drag(tt *ui.Tester, from, to ui.Rect) {
	tt.Press(from.X+from.W/2, from.Y+from.H/2)
	tt.Move(to.X+to.W/2, to.Y+to.H/2)
	tt.Release(to.X+to.W/2, to.Y+to.H/2)
	tt.Frame()
}

// inside fails when r leaves left..right.
func inside(t *testing.T, name string, r ui.Rect, left, right float32) {
	t.Helper()
	if r.W <= 0 || r.X < left-0.5 || r.X+r.W > right+0.5 {
		t.Errorf("%s spans %.1f..%.1f, outside %.1f..%.1f", name, r.X, r.X+r.W, left, right)
	}
}
