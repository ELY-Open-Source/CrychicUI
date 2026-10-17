package kit

import (
	"bytes"
	"image"
	"log/slog"
	"strings"
	"testing"
	"time"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func near(img *image.RGBA, x, y int, want ui.Color) bool {
	p := img.RGBAAt(x, y)
	d := func(a, b uint8) bool { return a-b < 4 || b-a < 4 }
	return d(p.R, want.R) && d(p.G, want.G) && d(p.B, want.B)
}

func lum(r, g, b uint8) int { return int(r)*2126 + int(g)*7152 + int(b)*722 }

func count(cue sound.Cue) int {
	n := 0
	for _, c := range sound.Recent() {
		if c == cue {
			n++
		}
	}
	return n
}

// posed poses an accent box at (50,50).
func posed(pose func(ui.Element) ui.Element, bg, fill *ui.Color) func(*ui.Context) {
	return func(c *ui.Context) {
		t := c.Theme()
		*bg, *fill = t.Background, t.Accent
		ui.Box(c).Size(200, 200).Padding(50).Background(t.Background).Children(func() {
			pose(ui.Box(c).Size(100, 100).Background(t.Accent))
		})
	}
}

func TestPressShrinks(t *testing.T) {
	var bg, fill ui.Color
	tt := ui.NewTester(posed(Press, &bg, &fill), 200, 200)
	if !near(tt.Image(), 50, 50, fill) {
		t.Fatal("corner not filled at rest")
	}
	tt.Press(100, 100)
	tt.Advance(time.Second)
	if !near(tt.Image(), 50, 50, bg) {
		t.Error("corner still filled while pressed")
	}
	tt.Release(100, 100)
	tt.Advance(time.Second)
	if !near(tt.Image(), 50, 50, fill) {
		t.Error("corner not refilled after release")
	}
}

func TestLiftRaises(t *testing.T) {
	var bg, fill ui.Color
	tt := ui.NewTester(posed(Lift, &bg, &fill), 200, 200)
	if !near(tt.Image(), 100, 149, fill) || !near(tt.Image(), 100, 48, bg) {
		t.Fatal("box not at rest")
	}
	tt.Move(100, 100)
	tt.Advance(time.Second)
	if !near(tt.Image(), 100, 149, bg) || !near(tt.Image(), 100, 48, fill) {
		t.Error("box not lifted 2 DIP on hover")
	}
}

func TestTransitionsEnterAndExit(t *testing.T) {
	for name, tr := range map[string]ui.ElementTransition{"Enter": Enter, "Pop": Pop} {
		t.Run(name, func(t *testing.T) {
			var bg, fill ui.Color
			show := false
			tt := ui.NewTester(func(c *ui.Context) {
				th := c.Theme()
				bg, fill = th.Background, th.Accent
				ui.Column(c).Size(100, 100).Background(bg).Children(func() {
					if show {
						ui.Box(c).Key("b").Size(100, 100).Background(fill).Transition(tr)
					}
				})
			}, 100, 100)
			show = true
			tt.Frame()
			if !tt.FramePending() || !near(tt.Image(), 2, 2, bg) {
				t.Fatal("entry did not start hidden")
			}
			tt.Advance(time.Second)
			if tt.FramePending() || !near(tt.Image(), 2, 2, fill) {
				t.Fatal("entry did not settle")
			}
			show = false
			tt.Frame()
			if !tt.FramePending() {
				t.Fatal("exit did not animate")
			}
			tt.Advance(time.Second)
			if tt.FramePending() || !near(tt.Image(), 2, 2, bg) {
				t.Error("exit did not settle")
			}
		})
	}
}

func TestClickedPlaysOnce(t *testing.T) {
	n := 0
	before := count(sound.Tap)
	tt := ui.NewTester(func(c *ui.Context) {
		if Clicked(ui.Button(c, "Go"), sound.Tap) {
			n++
		}
	}, 200, 100)
	if count(sound.Tap) != before {
		t.Fatal("cue played without a click")
	}
	if err := tt.Click("Go"); err != nil {
		t.Fatal(err)
	}
	if n != 1 || count(sound.Tap) != before+1 {
		t.Errorf("clicks %d, taps %d, want 1 and %d", n, count(sound.Tap), before+1)
	}
}

func TestChangedPlaysOnce(t *testing.T) {
	n, on := 0, false
	before := count(sound.Toggle)
	tt := ui.NewTester(func(c *ui.Context) {
		if Changed(ui.Checkbox(c, &on, "Box"), sound.Toggle) {
			n++
		}
	}, 200, 100)
	if err := tt.Click("Box"); err != nil {
		t.Fatal(err)
	}
	if !on || n != 1 || count(sound.Toggle) != before+1 {
		t.Errorf("on %v, changes %d, toggles %d", on, n, count(sound.Toggle)-before)
	}
}

func TestCardSurface(t *testing.T) {
	var th ui.Theme
	img := ui.Render(func(c *ui.Context) {
		th = *c.Theme()
		ui.Box(c).Size(140, 140).Padding(20).Background(th.Background).Children(func() {
			Card(c).Size(100, 100)
		})
	}, 140, 140, 1)
	if !near(img, 70, 70, th.Background) {
		t.Error("card face is not Background")
	}
	if !near(img, 20, 70, th.Border) {
		t.Error("card edge is not Border")
	}
	p, q := img.RGBAAt(70, 120), img.RGBAAt(70, 135)
	if lum(p.R, p.G, p.B) >= lum(q.R, q.G, q.B) {
		t.Error("card casts no shadow")
	}
}

func TestMutedColor(t *testing.T) {
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		ui.Column(c).Grow(1).Background(th.Background).Children(func() {
			Muted(c, "Note")
		})
	}, 200, 100)
	r, ok := tt.Find("Note")
	if !ok {
		t.Fatal("Note not shown")
	}
	img, floor, inked := tt.Image(), lum(th.TextMuted.R, th.TextMuted.G, th.TextMuted.B), false
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			p := img.RGBAAt(x, y)
			if l := lum(p.R, p.G, p.B); l < floor-40000 {
				t.Fatalf("pixel at %d,%d darker than TextMuted", x, y)
			}
			inked = inked || !near(img, x, y, th.Background)
		}
	}
	if !inked {
		t.Error("Note drew nothing")
	}
}

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(old) })
	Log("button", "click", "count", 2)
	if s := buf.String(); !strings.Contains(s, `msg="crychic: button click" count=2`) {
		t.Errorf("log = %q", s)
	}
}
