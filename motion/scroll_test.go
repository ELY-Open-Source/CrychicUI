package motion

import (
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestScrollProgress(t *testing.T) {
	st := ui.ScrollState{MaxY: 100}
	tt := ui.NewTester(func(c *ui.Context) { ScrollProgress(c, "Read", &st).Width(200) }, 200, 40)
	th := ui.LightTheme()
	st.Y = 50
	tt.Frame()
	settles(t, tt, "bar")
	if !tinted(tt, 90, 1, th.Accent) || tinted(tt, 110, 1, th.Accent) {
		t.Error("bar does not end at half")
	}
	st.Y = math.MaxFloat32
	tt.Frame()
	tt.Advance(time.Second)
	if !tinted(tt, 198, 1, th.Accent) {
		t.Error("bar not full at the end")
	}
}

func TestScrollProgressFollowsTheWheel(t *testing.T) {
	st := ui.ScrollState{}
	tt := ui.NewTester(func(c *ui.Context) {
		ScrollProgress(c, "Read", &st).Width(200)
		ui.Scroll(c).TrackScroll(&st).Height(100).Children(func() {
			for i := range 20 {
				ui.Text(c, fmt.Sprint("Line ", i))
			}
		})
	}, 200, 140)
	tt.Scroll(100, 60, 0, 150)
	tt.Advance(time.Second)
	if st.Y <= 0 || !tinted(tt, 5, 1, ui.LightTheme().Accent) {
		t.Errorf("wheel scroll: Y %v, bar empty", st.Y)
	}
}

func TestScrollProgressFollowsContent(t *testing.T) {
	st, n := ui.ScrollState{}, 20
	tt := ui.NewTester(func(c *ui.Context) {
		ScrollProgress(c, "Read", &st).Width(200)
		ui.Scroll(c).TrackScroll(&st).Height(100).Children(func() {
			for i := range n {
				ui.Text(c, fmt.Sprint("Line ", i))
			}
		})
	}, 200, 140)
	tt.Scroll(100, 60, 0, 100)
	tt.Advance(time.Second)
	half := st.Y / st.MaxY
	n = 80
	tt.Frame()
	tt.Advance(time.Second)
	want := st.Y / st.MaxY
	x := int(200 * (half + want) / 2)
	if want >= half || tinted(tt, x, 1, ui.LightTheme().Accent) {
		t.Errorf("bar kept %v of the old content; want %v", half, want)
	}
}

func TestScrollFade(t *testing.T) {
	st := ui.ScrollState{MaxY: 500}
	tt := ui.NewTester(padded(func(c *ui.Context) { ScrollFade(c, &st, ScrollFadeOptions{}, func() { dot(c) }) }), 200, 200)
	th := ui.LightTheme()
	r := find(t, tt, "dot")
	if !tinted(tt, int(r.X+20), int(r.Y+20), th.Accent) {
		t.Error("not opaque at the top")
	}
	st.Y = 80
	tt.Frame()
	if !tinted(tt, int(r.X+20), int(r.Y+10), th.Background) {
		t.Error("not gone after the distance")
	}
}

func TestScrollReveal(t *testing.T) {
	st := ui.ScrollState{}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Scroll(c).TrackScroll(&st).Height(120).Children(func() {
			ScrollReveal(c, 10, func(i int) { ui.Box(c).Size(100, 30).Background(c.Theme().Accent).Label(fmt.Sprint("Row ", i)) })
		})
	}, 200, 120)
	th := ui.LightTheme()
	settles(t, tt, "first rows")
	if r := find(t, tt, "Row 0"); !tinted(tt, int(r.X+50), int(r.Y+15), th.Accent) {
		t.Error("row in view not revealed")
	}
	st.Y = math.MaxFloat32
	tt.Frame()
	tt.Frame()
	r := find(t, tt, "Row 9")
	if tinted(tt, int(r.X+50), int(r.Y+r.H/2), th.Accent) {
		t.Error("row 9 shown before it arrived")
	}
	settles(t, tt, "row 9")
	if r := find(t, tt, "Row 9"); !tinted(tt, int(r.X+50), int(r.Y+r.H/2), th.Accent) {
		t.Error("row 9 never revealed")
	}
}

func TestSequence(t *testing.T) {
	var run uint64
	clicked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, action := Sequence(c, run, "Welcome", "Three steps to start.", SequenceOptions{}); action.Clicked() {
			clicked++
		}
	}), 400, 300)
	title, act := find(t, tt, "Welcome").Y, find(t, tt, "Continue").Y
	tt.Advance(200 * time.Millisecond)
	if dt, da := title-find(t, tt, "Welcome").Y, act-find(t, tt, "Continue").Y; dt <= da {
		t.Errorf("title moved %v, action %v: not in order", dt, da)
	}
	settles(t, tt, "sequence")
	rest := find(t, tt, "Welcome").Y
	run = 1
	tt.Frame()
	tt.Advance(20 * time.Millisecond)
	if find(t, tt, "Welcome").Y <= rest {
		t.Error("run change did not replay")
	}
	tt.Advance(2 * time.Second)
	before := baseline()
	click(t, tt, "Continue")
	if clicked != 1 {
		t.Errorf("clicked %d", clicked)
	}
	cues(t, before, sound.Tap)
}

func TestKeyPath(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { KeyPath(c) }), 300, 200)
	th := ui.LightTheme()
	moves(t, tt, "key path", func() float32 { x, y := spot(tt, th.Accent); return x + 1000*y })
	stillUnderReduceMotion(t, tt, "key path")
}

func TestArcCounter(t *testing.T) {
	v := 75.0
	tt := ui.NewTester(padded(func(c *ui.Context) { ArcCounter(c, "Score", v, 100, ArcCounterOptions{}) }), 300, 300)
	th := ui.LightTheme()
	r := find(t, tt, "Score")
	arc := func() int { return pixels(tt, r, th.Accent) }
	number := func() float64 {
		for _, s := range tt.Texts() {
			if n, err := strconv.ParseFloat(s, 64); err == nil {
				return n
			}
		}
		t.Fatalf("no number in %q", tt.Texts())
		return 0
	}
	if arc() != 0 || number() != 0 {
		t.Errorf("starts with arc %d, number %v", arc(), number())
	}
	tt.Advance(200 * time.Millisecond)
	mid, n := arc(), number()
	settles(t, tt, "arc")
	if !(0 < mid && mid < arc()) || !(0 < n && n < 75) || number() != 75 {
		t.Errorf("arc %d then %d, number %v then %v", mid, arc(), n, number())
	}
	panics(t, "over total", func() { render(func(c *ui.Context) { ArcCounter(c, "x", 2, 1, ArcCounterOptions{}) }) })
}

// pixels counts the pixels in r that are exactly k.
func pixels(tt *ui.Tester, r ui.Rect, k ui.Color) int {
	img, n := tt.Image(), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); p.R == k.R && p.G == k.G && p.B == k.B {
				n++
			}
		}
	}
	return n
}
