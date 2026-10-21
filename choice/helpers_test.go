package choice

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

// baseline fills Recent with Swipe so later cues show.
func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Swipe)
	}
	return sound.Recent()
}

// played lists the cues since baseline.
func played(before []sound.Cue) []sound.Cue {
	after := sound.Recent()
	i := 0
	for i < len(after) && after[i] == sound.Swipe {
		i++
	}
	if len(before) != 32 || slices.Equal(before, after) {
		return nil
	}
	return after[i:]
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	if got := played(before); !slices.Equal(got, want) {
		t.Errorf("cues %v, want %v", got, want)
	}
}
