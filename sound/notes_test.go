package sound

import (
	"fmt"
	"math"
	"slices"
	"testing"
)

func mustPanic(t *testing.T, what string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", what)
		}
	}()
	fn()
}

func TestNotesAndDrumsRecord(t *testing.T) {
	resetRecent(t)
	PlayNote(261.63)
	PlayDrum(Snare)
	Play(Note)
	if got, want := Recent(), []Cue{Note, Drum, Note}; !slices.Equal(got, want) {
		t.Errorf("recent = %v, want %v", got, want)
	}
}

func TestNotesAndDrumsReject(t *testing.T) {
	for _, hz := range []float64{0, 19.9, 20001, math.NaN(), math.Inf(1)} {
		mustPanic(t, fmt.Sprintf("PlayNote(%v)", hz), func() { PlayNote(hz) })
	}
	mustPanic(t, "PlayDrum(pieceCount)", func() { PlayDrum(pieceCount) })
	if got := Piece(9).String(); got != "piece(9)" {
		t.Errorf("unknown piece name = %q", got)
	}
}

func TestTonesBounded(t *testing.T) {
	for _, hz := range []float64{27.5, 261.63, 4186, 15000} {
		bounded(t, fmt.Sprintf("tone %v", hz), render(tone(hz)))
	}
	if slices.Equal(render(tone(440)), render(tone(466.16))) {
		t.Error("A4 and A#4 sound the same")
	}
}

func TestDrumsBoundedAndDistinct(t *testing.T) {
	ps := Pieces()
	if len(ps) != 5 {
		t.Errorf("pieces = %d, want 5", len(ps))
	}
	for i, a := range ps {
		bounded(t, a.String(), render(drumScores[a]))
		for _, b := range ps[i+1:] {
			if slices.Equal(render(drumScores[a]), render(drumScores[b])) {
				t.Errorf("%v and %v sound the same", a, b)
			}
		}
	}
}
