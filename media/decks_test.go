package media

import (
	"bytes"
	"image/color"
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// moves reports whether the frame changes over d.
func moves(tt *ui.Tester, d time.Duration) bool {
	a := bytes.Clone(tt.Image().Pix)
	tt.Advance(d)
	return !bytes.Equal(a, tt.Image().Pix)
}

func TestMetronome(t *testing.T) {
	bpm, playing := 120.0, false
	tt := ui.NewTester(padded(func(c *ui.Context) { Metronome(c, "Click", &bpm, &playing, MetronomeOptions{}) }), 500, 300)
	tt.Frame()
	before := baseline()
	click(t, tt, "Play")
	tt.Advance(1100 * time.Millisecond)
	cues(t, before, sound.Toggle, sound.Tick, sound.Tick)
	if !moves(tt, 120*time.Millisecond) {
		t.Error("arm does not swing")
	}
	playing = false
	tt.Frame()
	settles(t, tt, "stopped metronome")
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	still := baseline()
	playing = true
	tt.Frame()
	tt.Advance(650 * time.Millisecond)
	if tt.FramePending() {
		t.Error("arm swings under reduced motion")
	}
	tt.Advance(450 * time.Millisecond)
	cues(t, still, sound.Tick, sound.Tick)
	panics(t, "tempo out of range", func() {
		b, on := 10.0, false
		render(func(c *ui.Context) { Metronome(c, "x", &b, &on, MetronomeOptions{}) })
	})
	panics(t, "min above max", func() {
		b, on := 60.0, false
		render(func(c *ui.Context) { Metronome(c, "x", &b, &on, MetronomeOptions{Min: 90, Max: 50}) })
	})
}

func TestEqualizer(t *testing.T) {
	var acc color.RGBA
	var curve float32
	gains, changed := []float64{0, 0, 0}, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		acc, curve = rgba(c.Theme().Accent), c.Theme().Space(16)
		if Equalizer(c, "EQ", []string{"60 Hz", "1 kHz", "8 kHz"}, gains, EqualizerOptions{}).Changed() {
			changed++
		}
	}), 420, 400)
	tt.Frame()
	settles(t, tt, "equalizer")
	r := find(t, tt, "EQ")
	x0, mid := r.X+r.W/6, r.Y+curve/2
	if !inkNear(tt, x0, mid, acc) {
		t.Error("flat curve not on the middle line")
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	for range 6 {
		tt.Key(0, ui.KeyUp)
	}
	tt.Frame()
	if gains[0] != 6 || changed != 6 {
		t.Fatalf("gain %v changed %d after six Ups", gains[0], changed)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	if !tt.FramePending() {
		t.Error("curve does not move")
	}
	settles(t, tt, "curve")
	if lift := float32(6) / 12 * curve * 0.42; !inkNear(tt, x0, mid-lift, acc) {
		t.Error("curve did not rise at the first band")
	}
	panics(t, "one gain short", func() {
		render(func(c *ui.Context) { Equalizer(c, "x", []string{"a", "b"}, []float64{0}, EqualizerOptions{}) })
	})
	panics(t, "gain past range", func() {
		render(func(c *ui.Context) { Equalizer(c, "x", []string{"a"}, []float64{20}, EqualizerOptions{}) })
	})
}

func TestTuner(t *testing.T) {
	f, changed := 98.0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Tuner(c, "Radio", &f, 87.5, 108, TunerOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	tt.Frame()
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if math.Abs(f-98.1) > 1e-9 || changed != 1 {
		t.Fatalf("freq %v changed %d after Right", f, changed)
	}
	cues(t, before, sound.Tick)
	if !slices.Contains(tt.Texts(), "98.1 MHz") {
		t.Errorf("texts %q, want the readout", tt.Texts())
	}
	if !tt.FramePending() {
		t.Error("band does not slide")
	}
	settles(t, tt, "band")
	panics(t, "freq past band", func() {
		v := 120.0
		render(func(c *ui.Context) { Tuner(c, "x", &v, 87.5, 108, TunerOptions{}) })
	})
}

// spins checks a deck turns, then rests.
func spins(t *testing.T, name string, deck func(c *ui.Context, playing bool)) {
	t.Helper()
	playing := false
	tt := ui.NewTester(padded(func(c *ui.Context) { deck(c, playing) }), 500, 420)
	tt.Frame()
	settles(t, tt, name)
	playing = true
	tt.Frame()
	tt.Advance(300 * time.Millisecond)
	if !moves(tt, 200*time.Millisecond) {
		t.Errorf("%s does not turn while playing", name)
	}
	playing = false
	tt.Frame()
	settles(t, tt, name)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	playing = true
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	if tt.FramePending() {
		t.Errorf("%s turns under reduced motion", name)
	}
}

func TestTurntable(t *testing.T) {
	spins(t, "turntable", func(c *ui.Context, on bool) { Turntable(c, "Deck", red, on, TurntableOptions{}) })
	panics(t, "no label", func() { render(func(c *ui.Context) { Turntable(c, "", nil, false, TurntableOptions{}) }) })
}

func TestCassette(t *testing.T) {
	spins(t, "cassette", func(c *ui.Context, on bool) { Cassette(c, "Tape", "Side A", on, CassetteOptions{}) })
	panics(t, "negative width", func() {
		render(func(c *ui.Context) { Cassette(c, "x", "", false, CassetteOptions{Width: -1}) })
	})
}

func TestVisualizer(t *testing.T) {
	spins(t, "visualizer", func(c *ui.Context, on bool) { Visualizer(c, "Bars", blue, on, VisualizerOptions{}) })
	panics(t, "negative bars", func() {
		render(func(c *ui.Context) { Visualizer(c, "x", nil, false, VisualizerOptions{Bars: -1}) })
	})
}
