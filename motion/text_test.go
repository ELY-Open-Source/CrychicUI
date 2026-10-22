package motion

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestTypewriter(t *testing.T) {
	text := "Hi there you"
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { Typewriter(c, text, TypewriterOptions{}) }, 300, 80)
	tt.Advance(160 * time.Millisecond)
	part := find(t, tt, text).W
	tt.Advance(time.Second)
	full := find(t, tt, text).W
	if part <= 0 || part >= full/2 {
		t.Errorf("typed widths %v then %v", part, full)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	text = "New"
	tt.Frame()
	if w := find(t, tt, "New").W; w > 0 {
		t.Errorf("new text did not restart: %v wide", w)
	}
	loops(t, tt, "caret")
	if find(t, tt, "New").W == 0 {
		t.Error("new text never typed")
	}
	before = baseline()
	burst := ui.NewTester(func(c *ui.Context) { Typewriter(c, "a b c d", TypewriterOptions{Every: time.Millisecond}) }, 300, 80)
	burst.Advance(20 * time.Millisecond)
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	panics(t, "empty", func() { render(func(c *ui.Context) { Typewriter(c, "", TypewriterOptions{}) }) })
}

func TestScramble(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Scramble(c, "CRYCHIC", ScrambleOptions{}) }, 300, 80)
	if looks(tt, "CRYCHIC") {
		t.Error("starts settled")
	}
	settles(t, tt, "scramble")
	if !looks(tt, "CRYCHIC") {
		t.Error("never settled")
	}
	tt2 := ui.NewTester(func(c *ui.Context) { Scramble(c, "AB C", ScrambleOptions{Glyphs: "#"}) }, 300, 80)
	if !looks(tt2, "## #") {
		t.Error("custom glyphs not drawn")
	}
}

func TestWordRotate(t *testing.T) {
	turns := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if WordRotate(c, []string{"fast", "calm"}, WordRotateOptions{Every: time.Second}).Changed() {
			turns++
		}
	}, 300, 80)
	find(t, tt, "fast")
	tt.Advance(1100 * time.Millisecond)
	find(t, tt, "calm")
	tt.Advance(time.Second)
	find(t, tt, "fast")
	if turns != 2 || !tt.FramePending() {
		t.Errorf("turns %d, pending %v", turns, tt.FramePending())
	}
	panics(t, "empty word", func() { render(func(c *ui.Context) { WordRotate(c, []string{""}, WordRotateOptions{}) }) })
}

func TestWordFade(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { WordFade(c, "one two three") }, 300, 80)
	y0 := find(t, tt, "three").Y
	settles(t, tt, "words")
	if y1 := find(t, tt, "three").Y; y0 <= y1 {
		t.Errorf("last word did not rise: %v to %v", y0, y1)
	}
	panics(t, "blank", func() { render(func(c *ui.Context) { WordFade(c, "  ") }) })
}

func TestLetterPull(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(30).Children(func() { LetterPull(c, "Crychic") })
	}, 300, 120)
	r := find(t, tt, "Crychic")
	if n := inked(tt, r); n != 0 {
		t.Errorf("starts with %d pixels inked", n)
	}
	settles(t, tt, "letters")
	if n := inked(tt, r); n < 50 {
		t.Errorf("letters never landed: %d pixels", n)
	}
}

func TestWaveText(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { WaveText(c, "Wave") }, 300, 80)
	r := find(t, tt, "Wave")
	a := inkRows(tt, r)
	tt.Advance(300 * time.Millisecond)
	if slices.Equal(a, inkRows(tt, r)) {
		t.Error("letters did not move")
	}
	loops(t, tt, "wave")
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	if tt.FramePending() {
		t.Error("reduced motion still waves")
	}
}

// inkRows counts inked pixels per row of r.
func inkRows(tt *ui.Tester, r ui.Rect) []int {
	bg, img := ui.LightTheme().Background, tt.Image()
	var rows []int
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		n := 0
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); p.R != bg.R || p.G != bg.G || p.B != bg.B {
				n++
			}
		}
		rows = append(rows, n)
	}
	return rows
}

// looks compares tt's frame with a plain Text of s.
func looks(tt *ui.Tester, s string, features ...string) bool {
	b := tt.Image().Bounds()
	ref := ui.NewTester(func(c *ui.Context) { ui.Text(c, s).FontFeatures(features...) }, b.Dx(), b.Dy())
	r := ui.Rect{W: float32(b.Dx()), H: float32(b.Dy())}
	return slices.Equal(inkRows(tt, r), inkRows(ref, r))
}

func TestNumberTicker(t *testing.T) {
	v := 1000.0
	tt := ui.NewTester(func(c *ui.Context) { NumberTicker(c, v, NumberTickerOptions{}) }, 300, 80)
	if !looks(tt, "0", "tnum") {
		t.Error("does not start at 0")
	}
	tt.Advance(300 * time.Millisecond)
	if looks(tt, "0", "tnum") || looks(tt, "1000", "tnum") {
		t.Error("no count between 0 and 1000")
	}
	tt.Advance(time.Second)
	if !looks(tt, "1000", "tnum") {
		t.Error("does not end at 1000")
	}
	v = 500
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	if looks(tt, "500", "tnum") || looks(tt, "1000", "tnum") {
		t.Error("no count between 1000 and 500")
	}
	settles(t, tt, "count down")
	if !looks(tt, "500", "tnum") {
		t.Error("does not end at 500")
	}
	tt2 := ui.NewTester(func(c *ui.Context) { NumberTicker(c, 2.5, NumberTickerOptions{From: 2.5, Format: "%.1f"}) }, 300, 80)
	if !looks(tt2, "2.5", "tnum") || tt2.FramePending() {
		t.Error("format or from ignored")
	}
	panics(t, "NaN", func() { render(func(c *ui.Context) { NumberTicker(c, math.NaN(), NumberTickerOptions{}) }) })
}

func TestRollingDigits(t *testing.T) {
	v := int64(41)
	tt := ui.NewTester(func(c *ui.Context) { RollingDigits(c, v) }, 300, 80)
	find(t, tt, "41")
	v = 42
	tt.Frame()
	settles(t, tt, "roll")
	find(t, tt, "42")
}
