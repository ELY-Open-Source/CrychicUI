package motion

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestFadeIn(t *testing.T) {
	var run uint64
	tt := ui.NewTester(func(c *ui.Context) { FadeIn(c, run, func() { dot(c) }) }, 200, 100)
	th := ui.LightTheme()
	r := find(t, tt, "dot")
	cx, cy := int(r.X+20), int(r.Y+20)
	if tinted(tt, cx, cy, th.Accent) {
		t.Error("starts opaque")
	}
	settles(t, tt, "fade")
	if !tinted(tt, cx, cy, th.Accent) {
		t.Error("never opaque")
	}
	run = 1
	tt.Frame()
	tt.Advance(20 * time.Millisecond)
	if tinted(tt, cx, cy, th.Accent) {
		t.Error("run change did not replay")
	}
	settles(t, tt, "replay")
}

func TestSlideIn(t *testing.T) {
	side := Bottom
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { SlideIn(c.Key(side), SlideInOptions{From: side}, func() { dot(c) }) })
	}, 200, 200)
	y0 := find(t, tt, "dot").Y
	settles(t, tt, "slide")
	if y1 := find(t, tt, "dot").Y; y0-y1 < 20 {
		t.Errorf("from bottom: %v to %v", y0, y1)
	}
	side = Start
	tt.SetDirection(ui.RightToLeft)
	tt.Frame()
	x0 := find(t, tt, "dot").X
	settles(t, tt, "slide start")
	if x1 := find(t, tt, "dot").X; x0-x1 < 20 {
		t.Errorf("rtl start should come from the right: %v to %v", x0, x1)
	}
	panics(t, "bad side", func() { render(func(c *ui.Context) { SlideIn(c, SlideInOptions{From: 9}, func() {}) }) })
}

func TestScaleIn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { ScaleIn(c, func() { dot(c) }) }, 200, 100)
	if w := find(t, tt, "dot").W; w > 37 {
		t.Errorf("starts at %v wide, want 36", w)
	}
	settles(t, tt, "scale")
	if w := find(t, tt, "dot").W; w != 40 {
		t.Errorf("rests at %v wide", w)
	}
}

func TestBlurIn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { BlurIn(c, func() { dot(c) }) })
	}, 200, 200)
	bg := ui.LightTheme().Background
	r := find(t, tt, "dot")
	x, y := int(r.X)-3, int(r.Y+20)
	if tinted(tt, x, y, bg) {
		t.Error("edge starts sharp")
	}
	settles(t, tt, "blur")
	if !tinted(tt, x, y, bg) {
		t.Error("edge never sharpens")
	}
}

func TestFlipIn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { FlipIn(c, func() { dot(c) }) })
	}, 200, 200)
	if h := find(t, tt, "dot").H; h > 10 {
		t.Errorf("starts %v tall, want edge-on", h)
	}
	settles(t, tt, "flip")
	if h := find(t, tt, "dot").H; h != 40 {
		t.Errorf("rests %v tall", h)
	}
}

func TestStagger(t *testing.T) {
	n := 3
	tt := ui.NewTester(func(c *ui.Context) {
		Stagger(c, n, func(i int) { ui.Text(c, fmt.Sprint("Item ", i)) }, StaggerOptions{})
	}, 300, 300)
	tt.Advance(2 * time.Second)
	rest := []float32{find(t, tt, "Item 0").Y, find(t, tt, "Item 2").Y}
	tt2 := ui.NewTester(func(c *ui.Context) {
		Stagger(c, 3, func(i int) { ui.Text(c, fmt.Sprint("Item ", i)) }, StaggerOptions{})
	}, 300, 300)
	tt2.Advance(120 * time.Millisecond)
	if d0, d2 := find(t, tt2, "Item 0").Y-rest[0], find(t, tt2, "Item 2").Y-rest[1]; d0 >= d2 {
		t.Errorf("item 0 off by %v, item 2 by %v: no stagger", d0, d2)
	}
	n = 4
	tt.Frame()
	tt.Advance(250 * time.Millisecond)
	if tt.FramePending() {
		t.Error("an item added later waited for its index")
	}
	panics(t, "negative n", func() { render(func(c *ui.Context) { Stagger(c, -1, func(int) {}, StaggerOptions{}) }) })
}
