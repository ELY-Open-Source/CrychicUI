package kit

import (
	"bytes"
	"image"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
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

func TestNumeric(t *testing.T) {
	text := "12"
	tt := ui.NewTester(func(c *ui.Context) { Numeric(c, text) }, 200, 100)
	text = "13"
	tt.Frame()
	if !tt.FramePending() {
		t.Error("digit change did not animate")
	}
	tt.Advance(time.Second)
	if !tt.HasText("13") || tt.FramePending() {
		t.Errorf("texts %q, pending %v", tt.Texts(), tt.FramePending())
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

func TestArriveFadesIn(t *testing.T) {
	for _, delay := range []time.Duration{0, 200 * time.Millisecond} {
		var bg, fill ui.Color
		tt := ui.NewTester(func(c *ui.Context) {
			th := c.Theme()
			bg, fill = th.Background, th.Accent
			ui.Column(c).Size(100, 100).Background(bg).Children(func() {
				Arrive(ui.Box(c).Size(100, 100).Background(fill), ui.Pose{Fade: 1}, delay)
			})
		}, 100, 100)
		if !tt.FramePending() || !near(tt.Image(), 50, 50, bg) {
			t.Fatalf("delay %v: box shown before it arrived", delay)
		}
		tt.Advance(delay + 100*time.Millisecond)
		if near(tt.Image(), 50, 50, bg) || near(tt.Image(), 50, 50, fill) {
			t.Errorf("delay %v: box not mid-fade", delay)
		}
		tt.Advance(time.Second)
		if tt.FramePending() || !near(tt.Image(), 50, 50, fill) {
			t.Errorf("delay %v: box did not settle", delay)
		}
	}
}

func TestShade(t *testing.T) {
	if s := Shade(0.5); s.R != 0 || s.G != 0 || s.B != 0 || s.A != 128 {
		t.Errorf("Shade(0.5) = %+v", s)
	}
}

func TestSettleRebuildsWhenLayoutMoves(t *testing.T) {
	var w, prev, used float32 = 40, 0, 0
	builds := 0
	tt := ui.NewTester(func(c *ui.Context) {
		builds++
		b := ui.Box(c).FillWidth().MaxWidth(w).Height(10)
		used = prev
		prev = b.LayoutBounds().W
		Settle(c, b, []float32{used}, func() []float32 { return []float32{b.LayoutBounds().W} })
	}, 100, 100)
	if builds != 3 || used != 40 {
		t.Fatalf("first frame: %d builds used %v, want 3 using 40", builds, used)
	}
	builds = 0
	tt.Frame()
	if builds != 1 {
		t.Errorf("steady layout ran %d builds, want 1", builds)
	}
	builds, w = 0, 60
	tt.Frame()
	if builds != 3 || used != 60 {
		t.Errorf("changed layout: %d builds used %v, want 3 using 60", builds, used)
	}
	builds = 0
	tt.SetSize(50, 100)
	if builds != 3 || used != 50 {
		t.Errorf("resize: %d builds used %v, want 3 using 50", builds, used)
	}
}

func TestInsideClipsToRoundedBox(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Margin(20, 0, 0, 20).Size(100, 60).Draw(Inside(12, func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X - 20, Y: r.Y - 20, W: r.W + 40, H: r.H + 40}, c.Theme().Danger, 0)
		}))
	}, 200, 120)
	th, img := ui.LightTheme(), tt.Image()
	is := func(x, y int, k ui.Color) bool { p := img.RGBAAt(x, y); return p.R == k.R && p.G == k.G && p.B == k.B }
	if !is(70, 50, th.Danger) || !is(10, 50, th.Background) || !is(20, 20, th.Background) || !is(150, 50, th.Background) {
		t.Errorf("center %v, left %v, corner %v, right %v", img.RGBAAt(70, 50), img.RGBAAt(10, 50), img.RGBAAt(20, 20), img.RGBAAt(150, 50))
	}
}
