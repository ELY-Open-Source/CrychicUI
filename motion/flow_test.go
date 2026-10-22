package motion

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestRings(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Rings(c, RingsOptions{}).Label("rings") }), 200, 200)
	r := find(t, tt, "rings")
	moves(t, tt, "rings", func() float32 { return float32(inked(tt, r)) })
	panics(t, "negative count", func() { render(func(c *ui.Context) { Rings(c, RingsOptions{Count: -1}) }) })
}

func TestMorph(t *testing.T) {
	open, changes := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Morph(c, "Panel", &open, func() { ui.Text(c, "Inside") }).Changed() {
			changes++
		}
	}), 400, 300)
	w0 := find(t, tt, "Panel").W
	before := baseline()
	click(t, tt, "Panel")
	if !open || changes != 1 {
		t.Fatalf("open %v, changes %d", open, changes)
	}
	cues(t, before, sound.Swipe)
	settles(t, tt, "grow")
	if w := find(t, tt, "Panel").W; w <= 4*w0 || !tt.HasText("Inside") {
		t.Errorf("open: width %v from %v, texts %q", w, w0, tt.Texts())
	}
	click(t, tt, "Panel")
	settles(t, tt, "shrink")
	if w := find(t, tt, "Panel").W; w != w0 || tt.HasText("Inside") {
		t.Errorf("closed: width %v, texts %q", w, tt.Texts())
	}
}

func TestSharedMove(t *testing.T) {
	open := -1
	tt := ui.NewTester(func(c *ui.Context) {
		SharedMove(c, []string{"Alpha", "Beta"}, &open, func(i int) { ui.Text(c, fmt.Sprint("Details ", i)) }, SharedMoveOptions{})
	}, 400, 300)
	row := find(t, tt, "Beta")
	before := baseline()
	click(t, tt, "Beta")
	if open != 1 {
		t.Fatalf("open %d", open)
	}
	tt.Advance(16 * time.Millisecond)
	flying := find(t, tt, "Beta").Y
	tt.Advance(2 * time.Second)
	if rest := find(t, tt, "Beta").Y; flying <= rest || !tt.HasText("Details 1") {
		t.Errorf("did not fly from the row at %v: %v then %v", row.Y, flying, rest)
	}
	click(t, tt, "Back")
	if open != -1 {
		t.Fatalf("back left open %d", open)
	}
	cues(t, before, sound.Open, sound.Close)
	panics(t, "open out of range", func() {
		render(func(c *ui.Context) { o := 2; SharedMove(c, []string{"a"}, &o, func(int) {}, SharedMoveOptions{}) })
	})
}

func TestShuffle(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	tt := ui.NewTester(func(c *ui.Context) { Shuffle(c, items, ShuffleOptions{}) }, 300, 400)
	ys := func() []float32 {
		var out []float32
		for _, s := range items {
			out = append(out, find(t, tt, s).Y)
		}
		return out
	}
	start := ys()
	before := baseline()
	click(t, tt, "Shuffle")
	cues(t, before, sound.Swipe)
	settles(t, tt, "rows slide")
	if slices.Equal(start, ys()) {
		t.Error("order did not change")
	}
	panics(t, "duplicate", func() { render(func(c *ui.Context) { Shuffle(c, []string{"a", "a"}, ShuffleOptions{}) }) })
}

func TestArriveList(t *testing.T) {
	items := []string{"b", "c"}
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { ArriveList(c, items) }, 300, 300)
	cues(t, before)
	y0 := find(t, tt, "b").Y
	items = []string{"a", "b", "c"}
	tt.Frame()
	cues(t, before, sound.Notify)
	settles(t, tt, "arrival")
	if y := find(t, tt, "b").Y; y <= y0 {
		t.Errorf("b did not move down: %v to %v", y0, y)
	}
	tt.Frame()
	cues(t, before, sound.Notify)
}

func TestPresence(t *testing.T) {
	shown := true
	tt := ui.NewTester(padded(func(c *ui.Context) { Presence(c, shown, func() { ui.Text(c, "Hi") }) }), 200, 200)
	find(t, tt, "Hi")
	shown = false
	tt.Frame()
	settles(t, tt, "exit")
	if tt.HasText("Hi") {
		t.Error("still shown")
	}
	shown = true
	tt.Frame()
	settles(t, tt, "enter")
	find(t, tt, "Hi")
}

func TestExpand(t *testing.T) {
	n := 1
	tt := ui.NewTester(func(c *ui.Context) {
		Expand(c, func() {
			for i := range n {
				ui.Text(c, fmt.Sprint("Line ", i))
			}
		}).Label("box")
	}, 300, 300)
	h1 := find(t, tt, "box").H
	n = 4
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	mid := find(t, tt, "box").H
	settles(t, tt, "grow")
	if h := find(t, tt, "box").H; !(h1 < mid && mid < h) {
		t.Errorf("heights %v, %v, %v", h1, mid, h)
	}
}

func TestCrossFade(t *testing.T) {
	b := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		CrossFade(c, b, func() { dot(c) }, func() { ui.Text(c, "Second") })
	}), 200, 200)
	th := ui.LightTheme()
	r := find(t, tt, "dot")
	b = true
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	p := tt.Image().RGBAAt(int(r.X+20), int(r.Y+20))
	if tinted(tt, int(r.X+20), int(r.Y+20), th.Accent) || tinted(tt, int(r.X+20), int(r.Y+20), th.Background) {
		t.Errorf("mid fade pixel %v is not between", p)
	}
	settles(t, tt, "fade")
	if tt.HasText("dot") || !tt.HasText("Second") {
		t.Errorf("after fade %q", tt.Texts())
	}
}
