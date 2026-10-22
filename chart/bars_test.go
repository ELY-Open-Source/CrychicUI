package chart

import (
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var months = []string{"Jan", "Feb", "Mar"}

func TestBars(t *testing.T) {
	th := ui.LightTheme()
	vs, at, changed := []float64{4, 8, 2}, -1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Bars(c, "Sales", months, vs, &at, BarsOptions{}).Changed() {
			changed++
		}
	}, 500, 300)
	whole := ui.Rect{W: 500, H: 300}
	tt.Advance(150 * time.Millisecond)
	early := inked(tt, whole, th.Accent)
	settles(t, tt, "bars")
	full := inked(tt, whole, th.Accent)
	if full == 0 || float64(early) > 0.8*float64(full) {
		t.Errorf("bars at 150 ms have %d of %d pixels; they should grow", early, full)
	}
	before := baseline()
	click(t, tt, "Feb: 8")
	if at != 1 || changed != 1 {
		t.Fatalf("at %d, changed %d", at, changed)
	}
	cues(t, before, sound.Select)
	settles(t, tt, "dim")
	vs = []float64{4, 8, 6}
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("bar did not spring to its new value")
	}
	settles(t, tt, "spring")
	panics(t, "NaN", func() {
		render(func(c *ui.Context) { n := -1; Bars(c, "x", months, []float64{1, math.NaN(), 2}, &n, BarsOptions{}) })
	})
	panics(t, "count", func() {
		render(func(c *ui.Context) { n := -1; Bars(c, "x", months, []float64{1}, &n, BarsOptions{}) })
	})
}

func TestLineHoverAndPick(t *testing.T) {
	at := -1
	tt := ui.NewTester(func(c *ui.Context) {
		Line(c, "Visits", months, []Series{{Name: "Web", Values: []float64{1, 5, 3}}}, &at, LineOptions{})
	}, 500, 300)
	settles(t, tt, "line")
	feb := find(t, tt, "Feb")
	tt.Move(feb.X+feb.W/2, 150)
	tt.Frame()
	if !tt.HasText("Feb: Web 5") {
		t.Fatalf("no tip: %q", tt.Texts())
	}
	before := baseline()
	tt.ClickAt(feb.X+feb.W/2, 150)
	tt.Frame()
	if at != 1 {
		t.Fatalf("at %d", at)
	}
	cues(t, before, sound.Select)
}
