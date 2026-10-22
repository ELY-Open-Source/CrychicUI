package media

import (
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPiano(t *testing.T) {
	note, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Piano(c, "Keys", &note, PianoOptions{Keys: 7}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "piano")
	before := baseline()
	k := find(t, tt, "C4")
	x, y := k.X+k.W/2, k.Y+k.H*0.85
	tt.Press(x, y)
	tt.Frame()
	if note != 60 {
		t.Fatalf("note %d on pressing C4, want 60", note)
	}
	cues(t, before, sound.Note)
	tt.Advance(40 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("the key does not press down")
	}
	tt.Release(x, y)
	tt.Frame()
	cues(t, before, sound.Note)
	b := find(t, tt, "C#4")
	tt.Press(b.X+b.W/2, b.Y+b.H/2)
	tt.Frame()
	tt.Release(b.X+b.W/2, b.Y+b.H/2)
	tt.Frame()
	if note != 61 {
		t.Errorf("note %d on pressing C#4, want 61", note)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	cues(t, before, sound.Note, sound.Note, sound.Note)
	if changed != 3 {
		t.Errorf("changed %d times, want 3", changed)
	}
	panics(t, "black lowest key", func() {
		n := -1
		render(func(c *ui.Context) { Piano(c, "x", &n, PianoOptions{From: 61}) })
	})
	panics(t, "one key", func() {
		n := -1
		render(func(c *ui.Context) { Piano(c, "x", &n, PianoOptions{Keys: 1}) })
	})
	panics(t, "note below the keys", func() {
		n := 30
		render(func(c *ui.Context) { Piano(c, "x", &n, PianoOptions{}) })
	})
}

func TestPitchAndNames(t *testing.T) {
	if p := pitch(69); p != 440 {
		t.Errorf("A4 = %v Hz", p)
	}
	if p := pitch(60); p < 261.6 || p > 261.7 {
		t.Errorf("C4 = %v Hz", p)
	}
	for m, want := range map[int]string{60: "C4", 61: "C#4", 69: "A4", 72: "C5"} {
		if got := named(m); got != want {
			t.Errorf("named(%d) = %s, want %s", m, got, want)
		}
	}
}

func TestDrumPad(t *testing.T) {
	hit := -1
	pads := []sound.Piece{sound.Kick, sound.Snare, sound.Hat}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		DrumPad(c, "Pads", pads, &hit, DrumPadOptions{Names: []string{"Boom", "Tss", "Tick"}})
	}), 500, 300)
	settles(t, tt, "drum pad")
	r := find(t, tt, "Tss")
	x, y := r.X+r.W/2, r.Y+r.H*0.8
	rest := tt.Image().RGBAAt(int(x), int(y))
	before := baseline()
	tt.Press(x, y)
	tt.Frame()
	if hit != 1 {
		t.Fatalf("hit %d, want 1", hit)
	}
	tt.Advance(30 * time.Millisecond)
	lit := tt.Image().RGBAAt(int(x), int(y))
	tt.Release(x, y)
	settles(t, tt, "drum pad")
	if absDiff(lit.R, rest.R)+absDiff(lit.G, rest.G)+absDiff(lit.B, rest.B) < 30 {
		t.Errorf("pad did not flash: %v at rest, %v struck", rest, lit)
	}
	// The press focused Tss; Tab moves to Tick.
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if hit != 2 {
		t.Errorf("hit %d after Tab and Space, want 2", hit)
	}
	cues(t, before, sound.Drum, sound.Drum)
	tt2 := ui.NewTester(padded(func(c *ui.Context) { DrumPad(c, "Pads", pads, &hit, DrumPadOptions{}) }), 500, 300)
	tt2.Frame()
	if !slices.Contains(tt2.Texts(), "Snare") {
		t.Errorf("default names missing: %q", tt2.Texts())
	}
	panics(t, "unknown piece", func() {
		h := -1
		render(func(c *ui.Context) { DrumPad(c, "x", []sound.Piece{9}, &h, DrumPadOptions{}) })
	})
	panics(t, "one name for three", func() {
		h := -1
		render(func(c *ui.Context) { DrumPad(c, "x", pads, &h, DrumPadOptions{Names: []string{"a"}}) })
	})
}
