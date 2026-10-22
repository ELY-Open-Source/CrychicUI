package motion

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// darks counts pixels in r darker than mid gray.
func darks(tt *ui.Tester, r ui.Rect) int {
	img, n := tt.Image(), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); int(p.R)+int(p.G)+int(p.B) < 3*96 {
				n++
			}
		}
	}
	return n
}

func TestSweep(t *testing.T) {
	var run uint64
	tt := ui.NewTester(padded(func(c *ui.Context) { Sweep(c, run, "Mark this phrase") }), 400, 200)
	r := find(t, tt, "Mark this phrase")
	marked := func() int { return inked(tt, r) - darks(tt, r) }
	start := marked()
	settles(t, tt, "sweep")
	full := marked()
	if full <= start+100 {
		t.Errorf("marker did not grow: %d to %d", start, full)
	}
	if darks(tt, r) < 50 {
		t.Error("marker hides the text")
	}
	run = 1
	tt.Frame()
	tt.Advance(20 * time.Millisecond)
	if marked() >= full {
		t.Error("run change did not replay")
	}
}

func TestUnderlineGrow(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { UnderlineGrow(c, "Read more") }), 300, 200)
	th := ui.LightTheme()
	r := find(t, tt, "Read more")
	line := ui.Rect{X: r.X, Y: r.Y + r.H, W: r.W, H: 6}
	if pixels(tt, line, th.Accent) != 0 {
		t.Error("underline before hover")
	}
	tt.Move(r.X+5, r.Y+5)
	tt.Frame()
	settles(t, tt, "grow")
	if pixels(tt, line, th.Accent) < int(r.W) {
		t.Error("underline did not span the text")
	}
	tt.Move(290, 190)
	tt.Frame()
	settles(t, tt, "shrink")
	if pixels(tt, line, th.Accent) != 0 {
		t.Error("underline stayed after hover")
	}
}

func TestWipe(t *testing.T) {
	var run uint64
	tt := ui.NewTester(func(c *ui.Context) { Wipe(c, run, "Revealed") }, 300, 60)
	r := ui.Rect{W: 300, H: 60}
	if inked(tt, r) != 0 {
		t.Error("text shows before the wipe")
	}
	tt.Advance(300 * time.Millisecond)
	mid := inked(tt, r)
	settles(t, tt, "wipe")
	if full := inked(tt, r); !(0 < mid && mid < full) || !looks(tt, "Revealed") {
		t.Errorf("ink %d then %d", mid, full)
	}
	run = 1
	tt.Frame()
	tt.Advance(20 * time.Millisecond)
	if looks(tt, "Revealed") {
		t.Error("run change did not replay")
	}
}

func TestIconSwap(t *testing.T) {
	name := "sun"
	tt := ui.NewTester(func(c *ui.Context) { IconSwap(c, name) }, 100, 100)
	name = "moon"
	tt.Frame()
	settles(t, tt, "swap")
	panics(t, "unknown icon", func() { render(func(c *ui.Context) { IconSwap(c, "no-such-icon") }) })
}

func TestDrawCheck(t *testing.T) {
	var run uint64
	before := baseline()
	tt := ui.NewTester(padded(func(c *ui.Context) { DrawCheck(c, "Saved", run) }), 200, 200)
	th := ui.LightTheme()
	r := find(t, tt, "Saved")
	if tt.FramePending() || pixels(tt, r, th.Success) != 0 {
		t.Error("draws before its run changes")
	}
	run = 1
	tt.Frame()
	tt.Advance(150 * time.Millisecond)
	mid := pixels(tt, r, th.Success)
	settles(t, tt, "check")
	if full := pixels(tt, r, th.Success); !(0 < mid && mid < full) {
		t.Errorf("check ink %d then %d", mid, full)
	}
	cues(t, before, sound.Success)
}

func TestTextShimmer(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { TextShimmer(c, "Loading the feed") }), 400, 200)
	r := find(t, tt, "Loading the feed")
	moves(t, tt, "shimmer", func() float32 {
		img, sum := tt.Image(), 0
		for x := int(r.X); x < int(r.X+r.W); x++ {
			for y := int(r.Y); y < int(r.Y+r.H); y++ {
				if p := img.RGBAAt(x, y); p.R < 60 {
					sum += x
				}
			}
		}
		return float32(sum)
	})
}

func TestMagnify(t *testing.T) {
	picked := -1
	items := []MagnifyItem{{"house", "Home"}, {"mail", "Mail"}, {"camera", "Camera"}, {"music", "Music"}, {"settings", "Settings"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Magnify(c, items, &picked) }), 500, 200)
	home, last := find(t, tt, "Home"), find(t, tt, "Settings")
	tt.Move(home.X+home.W/2, home.Y+home.H/2)
	tt.Frame()
	settles(t, tt, "magnify")
	if w := find(t, tt, "Home").W; w < home.W*1.5 || find(t, tt, "Settings").W != last.W {
		t.Errorf("near %v, far %v", w, find(t, tt, "Settings").W)
	}
	before := baseline()
	click(t, tt, "Mail")
	if picked != 1 {
		t.Errorf("picked %d", picked)
	}
	cues(t, before, sound.Tap)
	panics(t, "bad icon", func() {
		render(func(c *ui.Context) { p := 0; Magnify(c, []MagnifyItem{{"no-such-icon", "x"}}, &p) })
	})
}
