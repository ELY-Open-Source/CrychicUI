package choice

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestValueSliderTicksOnRelease(t *testing.T) {
	v := 20.0
	var r ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() {
			r = ValueSlider(c, "Volume", &v, 0, 100, ValueSliderOptions{}).Width(220).LayoutBounds()
		})
	}, 320, 140)
	tt.Frame()
	before := baseline()
	y := r.Y + r.H/2
	tt.Press(r.X+r.W*0.2, y)
	tt.Move(r.X+r.W*0.5, y)
	tt.Move(r.X+r.W*0.8, y)
	if v < 70 || v > 90 {
		t.Fatalf("value %v after drag", v)
	}
	if !tt.HasText(fmt.Sprintf("%.0f", v)) {
		t.Errorf("bubble lacks the value: %q", tt.Texts())
	}
	cues(t, before)
	tt.Move(r.X+r.W+60, y+50)
	tt.Move(r.X+r.W*0.8, y)
	tt.Move(r.X+r.W+60, y+50)
	if v != 100 {
		t.Errorf("drag past the end gave %v", v)
	}
	cues(t, before)
	tt.Release(r.X+r.W+60, y+50)
	cues(t, before, sound.Tick)
	tt.Key(0, ui.KeyLeft)
	cues(t, before, sound.Tick, sound.Tick)
	settles(t, tt, "ValueSlider bubble")
}

func TestSliderArgs(t *testing.T) {
	v, w, half, fifth := 0.0, 1.0, 0.5, 0.2
	seven, four := 0.7, 0.4
	near := seven - four
	inf := math.Inf(1)
	for name, fn := range map[string]func(c *ui.Context){
		"ValueSlider step":    func(c *ui.Context) { ValueSlider(c, "x", &v, 0, 1, ValueSliderOptions{Step: inf}) },
		"ValueSlider bounds":  func(c *ui.Context) { ValueSlider(c, "x", &v, -inf, 1, ValueSliderOptions{}) },
		"VerticalSlider step": func(c *ui.Context) { VerticalSlider(c, "x", &v, 0, 1, VerticalSliderOptions{Step: math.NaN()}) },
		"VerticalSlider size": func(c *ui.Context) { VerticalSlider(c, "x", &v, 0, 1, VerticalSliderOptions{Height: -4}) },
		"RangeSlider step":    func(c *ui.Context) { RangeSlider(c, "x", &v, &w, 0, 1, inf, RangeSliderOptions{}) },
		"RangeSlider bounds":  func(c *ui.Context) { RangeSlider(c, "x", &v, &w, 0, inf, 1, RangeSliderOptions{}) },
		"RangeSlider grid":    func(c *ui.Context) { RangeSlider(c, "x", &v, &half, 0, 1, 0.3, RangeSliderOptions{}) },
		"RangeSlider float":   func(c *ui.Context) { RangeSlider(c, "x", &fifth, &near, 0, 1, 0.1, RangeSliderOptions{}) },
	} {
		panics(t, name, func() { ui.Render(fn, 100, 100, 1) })
	}
}

func TestVerticalSlider(t *testing.T) {
	v := 0.0
	var r ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			r = VerticalSlider(c, "Level", &v, 0, 10, VerticalSliderOptions{Step: 1}).LayoutBounds()
		})
	}, 200, 240)
	tt.Frame()
	if r.H != 160 {
		t.Fatalf("height %v, want 160", r.H)
	}
	th := *ui.LightTheme()
	th.Spacing = 5
	if h := heightWith(&th); h != 200 {
		t.Errorf("roomy height %v, want 200", h)
	}
	before := baseline()
	tt.Press(r.X+r.W/2, r.Y+r.H-10)
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	if v != 5 {
		t.Errorf("value %v, want 5 at the middle", v)
	}
	cues(t, before, sound.Tick)
	settles(t, tt, "VerticalSlider knob")
}

func TestTickSlider(t *testing.T) {
	v := 0.0
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			if TickSlider(c, "Size", &v, 0, 4, 1, TickSliderOptions{Labels: []string{"XS", "S", "M", "L", "XL"}}).Changed() {
				changes++
			}
		})
	}, 360, 160)
	for _, s := range []string{"XS", "S", "M", "L", "XL"} {
		find(t, tt, s)
	}
	before := baseline()
	m := find(t, tt, "M")
	tt.ClickAt(m.X+m.W/2, m.Y-22)
	if v != 2 || changes != 1 {
		t.Fatalf("value %v changes %d, want 2 and 1", v, changes)
	}
	tt.Key(0, ui.KeyRight)
	if v != 3 {
		t.Errorf("arrow gave %v, want 3", v)
	}
	cues(t, before, sound.Tick, sound.Tick)
	settles(t, tt, "TickSlider label")
	panics(t, "infinite step", func() {
		ui.Render(func(c *ui.Context) { TickSlider(c, "x", &v, 0, 1, math.Inf(1), TickSliderOptions{}) }, 100, 100, 1)
	})
	panics(t, "uneven step", func() {
		ui.Render(func(c *ui.Context) { TickSlider(c, "x", &v, 0, 1, 0.3, TickSliderOptions{}) }, 100, 100, 1)
	})
}

func TestRangeSlider(t *testing.T) {
	low, high := 20.0, 80.0
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			if RangeSlider(c, "Price", &low, &high, 0, 100, 10, RangeSliderOptions{Format: "$%.0f"}).Changed() {
				changes++
			}
		})
	}, 360, 120)
	find(t, tt, "$20")
	find(t, tt, "$80")
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if low != 30 || changes != 1 {
		t.Fatalf("low %v changes %d", low, changes)
	}
	cues(t, before, sound.Tick)
	settles(t, tt, "RangeSlider digits")
	find(t, tt, "$30")
	panics(t, "crossed values", func() {
		h := 10.0
		ui.Render(func(c *ui.Context) { RangeSlider(c, "x", &high, &h, 0, 100, 1, RangeSliderOptions{}) }, 100, 100, 1)
	})
}

func heightWith(th *ui.Theme) float32 {
	v, h := 0.0, float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(th)
		ui.Column(c).Children(func() {
			h = VerticalSlider(c, "Level", &v, 0, 1, VerticalSliderOptions{}).LayoutBounds().H
		})
	}, 200, 400)
	tt.Frame()
	return h
}

func TestRangeSliderKnobsMeet(t *testing.T) {
	low, high := 20.0, 30.0
	tt := ui.NewTester(func(c *ui.Context) {
		RangeSlider(c, "Price", &low, &high, 0, 95, 10, RangeSliderOptions{})
	}, 360, 120)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if low != 30 || high != 90 {
		t.Errorf("low %v high %v, want 30 and 90", low, high)
	}
}

func TestRangeSliderFineStep(t *testing.T) {
	low, high := 10000.0, 10001.0
	tt := ui.NewTester(func(c *ui.Context) {
		RangeSlider(c, "Fine", &low, &high, 10000, 10001, 0.0001, RangeSliderOptions{Format: "%.4f"})
	}, 360, 120)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if low != 10000.0002 {
		t.Errorf("low %v, want 10000.0002", low)
	}
}

func TestRangeSliderNeverCrosses(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	keys := []ui.Key{ui.KeyLeft, ui.KeyRight, ui.KeyHome, ui.KeyEnd, ui.KeyPageUp, ui.KeyPageDown, ui.KeyTab}
	for _, g := range [][3]float64{{0, 1, 0.1}, {0.3, 2.05, 0.25}, {-2.5, 40, 7.5}, {10000, 10001, 0.0001}} {
		lo, hi, step := g[0], g[1], g[2]
		low, high := lo, lo+step
		tt := ui.NewTester(func(c *ui.Context) { RangeSlider(c, "R", &low, &high, lo, hi, step, RangeSliderOptions{}) }, 360, 120)
		tt.Key(0, ui.KeyTab)
		for range 300 {
			tt.Key(0, keys[rng.IntN(len(keys))])
			if low > high || !onStep(low, lo, step) || !onStep(high, lo, step) {
				t.Fatalf("grid %v: low %v high %v", g, low, high)
			}
		}
	}
}
