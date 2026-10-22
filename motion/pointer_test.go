package motion

import (
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func nan() float64 { return math.NaN() }

func TestFlipCard(t *testing.T) {
	flipped, changes := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if FlipCard(c, "Card", &flipped, func() { ui.Text(c, "Front") }, func() { ui.Text(c, "Back") }).Changed() {
			changes++
		}
	}), 300, 200)
	before := baseline()
	click(t, tt, "Card")
	if !flipped || changes != 1 {
		t.Fatalf("flipped %v, changes %d", flipped, changes)
	}
	cues(t, before, sound.Swipe)
	w0 := find(t, tt, "Front").W
	tt.Advance(50 * time.Millisecond)
	if find(t, tt, "Front").W == w0 {
		t.Error("card did not turn")
	}
	tt.Advance(2 * time.Second)
	if tt.FramePending() || tt.HasText("Front") || !tt.HasText("Back") {
		t.Errorf("after the turn: pending %v, texts %q", tt.FramePending(), tt.Texts())
	}
	tall := false
	tt2 := ui.NewTester(padded(func(c *ui.Context) {
		FlipCard(c, "Tall", &tall, func() { ui.Text(c, "One") }, func() {
			ui.Text(c, "A")
			ui.Text(c, "B")
			ui.Text(c, "C")
		})
	}), 300, 300)
	card, line := find(t, tt2, "Tall"), find(t, tt2, "One").H
	if card.H < 3*line {
		t.Errorf("card %v tall cannot hold a three-line back", card.H)
	}
	tt2.ClickAt(card.X+5, card.Y+card.H-3)
	if !tall {
		t.Fatal("a click near the bottom did not flip")
	}
	tt2.Advance(2 * time.Second)
	if last := find(t, tt2, "C"); last.Y+last.H > card.Y+card.H {
		t.Errorf("back face spills: card %v, last line %v", card, last)
	}
	panics(t, "nil flipped", func() { render(func(c *ui.Context) { FlipCard(c, "x", nil, func() {}, func() {}) }) })
}

func TestTiltCard(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		TiltCard(c, func() { ui.Box(c).Size(120, 80).Label("face") })
	}), 300, 220)
	flat := find(t, tt, "face")
	tt.Move(flat.X+5, flat.Y+5)
	tt.Advance(time.Second)
	if find(t, tt, "face") == flat {
		t.Error("card did not tilt toward the pointer")
	}
	tt.Move(2, 2)
	tt.Frame()
	settles(t, tt, "spring flat")
	if r := find(t, tt, "face"); math.Abs(float64(r.W-flat.W)) > 0.5 || math.Abs(float64(r.H-flat.H)) > 0.5 {
		t.Errorf("rests at %v, flat was %v", r, flat)
	}
}

func TestHoverZoom(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		HoverZoom(c, func() {
			ui.Box(c).Size(100, 60).Center().Children(func() { ui.Box(c).Size(50, 30).Label("pic") })
		})
	}), 300, 200)
	r := find(t, tt, "pic")
	tt.Move(r.X+25, r.Y+15)
	tt.Frame()
	settles(t, tt, "zoom in")
	if w := find(t, tt, "pic").W; w < 52.9 {
		t.Errorf("hovered width %v, want 53", w)
	}
	tt.Move(2, 2)
	tt.Frame()
	settles(t, tt, "zoom out")
	if w := find(t, tt, "pic").W; w != 50 {
		t.Errorf("rest width %v", w)
	}
}

func TestParallax(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Parallax(c, []float32{0, 1}, func(i int) {
			ui.Box(c).Fill().Center().Children(func() { ui.Box(c).Size(20, 20).Label([]string{"far", "near"}[i]) })
		}, ParallaxOptions{})
	}, 300, 200)
	x0 := find(t, tt, "near").X
	tt.Move(5, 80)
	tt.Advance(time.Second)
	left := find(t, tt, "near").X
	tt.Move(295, 80)
	tt.Advance(time.Second)
	right := find(t, tt, "near").X
	if left <= x0 || right >= x0 || find(t, tt, "far").X != x0 {
		t.Errorf("near %v, left %v, right %v, far %v", x0, left, right, find(t, tt, "far").X)
	}
	panics(t, "NaN depth", func() {
		render(func(c *ui.Context) { Parallax(c, []float32{float32(nan())}, func(int) {}, ParallaxOptions{}) })
	})
}

func TestMarquee(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Marquee(c, MarqueeOptions{}, func() { ui.Text(c, "News") }).Width(200)
	}, 300, 80)
	// The clip cuts the first copy as it leaves.
	shown := func() float32 { return find(t, tt, "News").W }
	tt.Advance(50 * time.Millisecond)
	w0 := shown()
	tt.Advance(200 * time.Millisecond)
	if w1 := shown(); w1 >= w0 {
		t.Fatalf("did not scroll: %v to %v", w0, w1)
	}
	tt.Move(100, 10)
	tt.Frame()
	paused := shown()
	tt.Advance(500 * time.Millisecond)
	if w := shown(); w != paused || tt.FramePending() {
		t.Errorf("hover did not pause: %v to %v, pending %v", paused, w, tt.FramePending())
	}
	tt.Move(290, 70)
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	if shown() == paused {
		t.Error("did not resume after hover")
	}
	wide := ui.NewTester(func(c *ui.Context) {
		Marquee(c, MarqueeOptions{}, func() { ui.Text(c, "Hi") }).Width(300)
	}, 300, 40)
	for range 12 {
		wide.Advance(170 * time.Millisecond)
		if n := inked(wide, ui.Rect{X: 220, W: 80, H: 30}); n == 0 {
			t.Fatal("short content left the right of a wide marquee empty")
		}
	}
	for _, rm := range []bool{false, true} {
		odd := ui.NewTester(func(c *ui.Context) {
			Marquee(c, MarqueeOptions{}, func() { ui.Box(c).Size(24.4, 24.4).Background(c.Theme().Accent) }).Width(112.8)
		}, 200, 60)
		odd.SetPreferences(ui.Preferences{ReduceMotion: rm})
		odd.Move(10, 10)
		odd.Frame()
		odd.Advance(time.Second)
		if odd.FramePending() {
			t.Errorf("fractional sizes, paused, reduce motion %v: keeps rebuilding", rm)
		}
	}
	panics(t, "negative speed", func() { render(func(c *ui.Context) { Marquee(c, MarqueeOptions{Speed: -1}, func() {}) }) })
}

func TestVerticalMarquee(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		VerticalMarquee(c, VerticalMarqueeOptions{Height: 60}, func() { ui.Text(c, "Line") })
	}, 300, 120)
	tt.Advance(50 * time.Millisecond)
	h0 := find(t, tt, "Line").H
	tt.Advance(500 * time.Millisecond)
	if h1 := find(t, tt, "Line").H; h1 >= h0 {
		t.Errorf("did not scroll up: visible height %v to %v", h0, h1)
	}
	loops(t, tt, "vertical marquee")
	panics(t, "NaN height", func() {
		render(func(c *ui.Context) { VerticalMarquee(c, VerticalMarqueeOptions{Height: float32(nan())}, func() {}) })
	})
}
