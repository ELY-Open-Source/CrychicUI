package button

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

func appended(before []sound.Cue, want ...sound.Cue) bool {
	after := sound.Recent()
	maxDrop := 0
	if len(before) == 32 {
		maxDrop = len(before)
	}
	for drop := 0; drop <= maxDrop; drop++ {
		kept := len(before) - drop
		if kept > len(after) || !slices.Equal(before[drop:], after[:kept]) {
			continue
		}
		added := after[kept:]
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
	return false
}

func TestAppendedRequiresNewCue(t *testing.T) {
	sound.Play(sound.Tap)
	before := sound.Recent()
	if appended(before, sound.Tap) {
		t.Error("unchanged history passed")
	}
	sound.Play(sound.Tap)
	if !appended(before, sound.Tap) {
		t.Error("new cue was not found")
	}
}

func TestAppendedRejectsHistoricalCue(t *testing.T) {
	sound.Play(sound.Tap)
	before := sound.Recent()
	sound.Play(sound.Open)
	if appended(before, sound.Tap) {
		t.Error("old Tap counted as new")
	}
	if appended(before) {
		t.Error("unexpected Open was accepted")
	}
	if !appended(before, sound.Open) {
		t.Error("new Open was not found")
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
