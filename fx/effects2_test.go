package fx

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestTrail(t *testing.T) {
	tt := still(Trail)
	rest := shot(tt)
	bg := ui.LightTheme().Background
	for x := float32(40); x <= 200; x += 20 {
		tt.Move(x, 80)
		tt.Advance(16 * time.Millisecond)
	}
	if !lit(tt, 200, 80, bg) || !lit(tt, 160, 80, bg) || lit(tt, 260, 80, bg) {
		t.Error("no trail behind the pointer")
	}
	if weight(tt, ui.Rect{X: 195, Y: 75, W: 10, H: 10}, bg) <= weight(tt, ui.Rect{X: 75, Y: 75, W: 10, H: 10}, bg) {
		t.Error("older dots did not fade")
	}
	tt.Advance(time.Second)
	if tt.FramePending() || shot(tt) != rest {
		t.Error("trail stayed after the pointer stopped")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Move(120, 60)
	tt.Frame()
	if tt.FramePending() || shot(tt) != rest {
		t.Error("trail drawn under reduced motion")
	}
}

func TestTrailStopsWhenMotionTurnsOff(t *testing.T) {
	tt := still(Trail)
	rest := shot(tt)
	tt.Move(100, 80)
	tt.Advance(16 * time.Millisecond)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	if tt.FramePending() || shot(tt) != rest {
		t.Error("dots kept fading under reduced motion")
	}
}

func TestTapsStopWhenMotionTurnsOff(t *testing.T) {
	for name, view := range map[string]func(c *ui.Context, content func()) ui.Element{"ClickBurst": ClickBurst,
		"WaterRipple": WaterRipple} {
		tt := still(view)
		rest := shot(tt)
		tt.ClickAt(80, 60)
		tt.Advance(100 * time.Millisecond)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		tt.Frame()
		if tt.FramePending() || shot(tt) != rest {
			t.Errorf("%s kept drawing under reduced motion", name)
		}
	}
}

// tapped clicks view at (80, 60) and samples it mid-life.
func tapped(t *testing.T, view func(c *ui.Context, content func()) ui.Element, cue sound.Cue, mid time.Duration) {
	t.Helper()
	before := baseline()
	tt := still(view)
	rest := shot(tt)
	tt.ClickAt(80, 60)
	tt.Frame()
	cues(t, before, cue)
	tt.Advance(mid)
	if n := ink(tt, ui.Rect{X: 30, Y: 10, W: 100, H: 100}); n == 0 || ink(tt, ui.Rect{X: 200, W: 100, H: 160}) != 0 {
		t.Errorf("%d pixels near the click, or some far off", n)
	}
	tt.Advance(time.Second)
	if tt.FramePending() || shot(tt) != rest {
		t.Error("drawing stayed after its life")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	before = baseline()
	tt.ClickAt(80, 60)
	tt.Frame()
	cues(t, before, cue)
	if tt.FramePending() || shot(tt) != rest {
		t.Error("drawn under reduced motion")
	}
}

func TestClickBurst(t *testing.T) { tapped(t, ClickBurst, sound.Pop, 200*time.Millisecond) }

func TestWaterRipple(t *testing.T) { tapped(t, WaterRipple, sound.Drop, 300*time.Millisecond) }

func TestCaretBlinks(t *testing.T) {
	th := ui.LightTheme()
	tt := ui.NewTester(func(c *ui.Context) { ui.Column(c).Padding(10).Children(func() { Caret(c) }) }, 60, 60)
	if !near(tt.Image().RGBAAt(11, 20), th.Accent) {
		t.Fatalf("caret %v not shown", tt.Image().RGBAAt(11, 20))
	}
	tt.Advance(800 * time.Millisecond)
	if !near(tt.Image().RGBAAt(11, 20), th.Background) {
		t.Error("caret did not blink off")
	}
	holds(t, tt, "Caret")
	if !near(tt.Image().RGBAAt(11, 20), th.Accent) {
		t.Error("caret hidden under reduced motion")
	}
}

func TestFlickerDims(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Flicker(c, "OPEN LATE") }, 300, 60)
	r := find(t, tt, "OPEN LATE")
	bg := ui.LightTheme().Background
	full := weight(tt, r, bg)
	tt.Advance(440 * time.Millisecond)
	if dim := weight(tt, r, bg); dim*2 >= full {
		t.Errorf("lit %d, dipped %d", full, dim)
	}
	tt.Advance(560 * time.Millisecond)
	if weight(tt, r, bg) != full {
		t.Error("did not relight")
	}
	holds(t, tt, "Flicker")
	panics(t, "empty text", func() { render(func(c *ui.Context) { Flicker(c, "") }) })
}

func TestFrostedBlursWhatIsBehind(t *testing.T) {
	stripes := func(p *ui.Painter, r ui.Rect) {
		for x := float32(0); x < r.W; x += 8 {
			p.Fill(ui.Rect{X: r.X + x, Y: r.Y, W: 4, H: r.H}, ui.LightTheme().Danger, 0)
		}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Grid(c).Columns(1).Size(300, 160).Children(func() {
			ui.Box(c).ColumnStart(1).RowStart(1).Draw(stripes)
			ui.Box(c).ColumnStart(1).RowStart(1).Padding(40).Children(func() {
				Frosted(c, func() { ui.Box(c).Size(120, 40).Label("pane") })
			})
		})
	}, 300, 160)
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Error("still moving after it arrived")
	}
	r := find(t, tt, "pane")
	img := tt.Image()
	y := int(r.Y + 20)
	in := abs(int(img.RGBAAt(int(r.X)+41, y).G) - int(img.RGBAAt(int(r.X)+45, y).G))
	out := abs(int(img.RGBAAt(1, 5).G) - int(img.RGBAAt(5, 5).G))
	if in*4 > out {
		t.Errorf("stripe contrast %d under the pane, %d outside", in, out)
	}
	panics(t, "nil content", func() { render(func(c *ui.Context) { Frosted(c, nil) }) })
}
