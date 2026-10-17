package layout

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestStackGapsAndResizes(t *testing.T) {
	more := false
	var gap float32
	tt := ui.NewTester(func(c *ui.Context) {
		gap = c.Theme().Space(4)
		Stack(c, 4, func() {
			ui.Text(c, "Alpha")
			ui.Text(c, "Beta")
			if more {
				ui.Text(c, "Gamma")
			}
		})
	}, 300, 300)
	a, b := find(t, tt, "Alpha"), find(t, tt, "Beta")
	if got := b.Y - (a.Y + a.H); abs(got-gap) > 0.5 {
		t.Errorf("gap %v, want %v", got, gap)
	}
	more = true
	tt.Frame()
	settles(t, tt, "Stack size")
	find(t, tt, "Gamma")
}

func TestInlineWraps(t *testing.T) {
	var gap float32
	tt := ui.NewTester(func(c *ui.Context) {
		gap = c.Theme().Space(2)
		Inline(c, 2, func() {
			for _, s := range []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"} {
				ui.Text(c, s)
			}
		})
	}, 160, 300)
	a, b, e := find(t, tt, "Alpha"), find(t, tt, "Beta"), find(t, tt, "Epsilon")
	if got := b.X - (a.X + a.W); abs(got-gap) > 0.5 {
		t.Errorf("gap %v, want %v", got, gap)
	}
	if e.Y <= a.Y {
		t.Error("items did not wrap")
	}
}

func TestLayersOverlapAndFade(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Layers(c, []func(){func() { ui.Text(c, "Back") }, func() { ui.Text(c, "Front") }})
	}, 300, 200)
	settles(t, tt, "Layers fade")
	if b, f := find(t, tt, "Back"), find(t, tt, "Front"); b.X != f.X || b.Y != f.Y {
		t.Errorf("layers apart: %v %v", b, f)
	}
	panics(t, "Layers(nil)", func() { ui.Render(func(c *ui.Context) { Layers(c, nil) }, 10, 10, 1) })
}

func TestCentered(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Centered(c, func() { ui.Text(c, "Mid") })
	}, 300, 200)
	settles(t, tt, "Centered scale")
	r := find(t, tt, "Mid")
	if abs(r.X+r.W/2-150) > 1 || abs(r.Y+r.H/2-100) > 1 {
		t.Errorf("Mid at %v, not centered", r)
	}
}

func TestContainerWidth(t *testing.T) {
	var pad, rem float32
	width := float32(400)
	tt := ui.NewTester(func(c *ui.Context) {
		pad, rem = c.Theme().Space(6), c.Theme().Rem(64)
		Container(c, ContainerOptions{MaxWidth: width}, func() { ui.Text(c, "Edge") })
	}, 1200, 200)
	if x := find(t, tt, "Edge").X; abs(x-((1200-width)/2+pad)) > 0.5 {
		t.Errorf("content at %v", x)
	}
	width = 0
	tt.Frame()
	settles(t, tt, "Container width")
	if x := find(t, tt, "Edge").X; abs(x-((1200-rem)/2+pad)) > 0.5 {
		t.Errorf("default width puts content at %v", x)
	}
	panics(t, "negative MaxWidth", func() {
		ui.Render(func(c *ui.Context) { Container(c, ContainerOptions{MaxWidth: -1}, nil) }, 10, 10, 1)
	})
}

// inked finds ink in a column span of pixels.
func inked(tt *ui.Tester, x int, y0, y1 float32) bool {
	img := tt.Image()
	for y := max(int(y0), 0); y <= min(int(y1), img.Bounds().Dy()-1); y++ {
		if lum(img, x, y) < lum(img, 299, 199)-20000 {
			return true
		}
	}
	return false
}

func TestSectionDividerGrows(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Section(c, "Title", func() { ui.Text(c, "Body") })
	}, 300, 200)
	title, body := find(t, tt, "Title"), find(t, tt, "Body")
	if body.Y <= title.Y {
		t.Fatal("body above title")
	}
	if inked(tt, 280, title.Y+title.H, body.Y-1) {
		t.Error("divider drawn before it grew")
	}
	settles(t, tt, "Section divider")
	if !inked(tt, 280, title.Y+title.H, body.Y-1) {
		t.Error("divider never reached the right edge")
	}
}

func TestSurfaces(t *testing.T) {
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		ui.Column(c).Gap(10).Children(func() {
			Panel(c, func() { ui.Box(c).Size(100, 40) })
			Inset(c, func() { ui.Box(c).Size(100, 40) })
			Well(c, func() { ui.Box(c).Size(100, 40) })
		})
	}, 300, 300)
	tt.Advance(time.Second)
	img := tt.Image()
	for _, p := range [][2]int{{60, 30}, {60, 100}, {60, 170}} {
		if !near(img, p[0], p[1], th.Surface) {
			t.Errorf("pixel %v is not Surface", p)
		}
	}
	tt.SetDark(true)
	if !tt.FramePending() {
		t.Error("Inset color did not animate to the dark theme")
	}
	tt.Advance(2 * time.Second)
	if !near(tt.Image(), 60, 100, ui.DarkTheme().Surface) {
		t.Error("Inset did not settle on the dark surface")
	}
}

func TestCardLifts(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			Card(c, func() { ui.Box(c).Size(100, 60) })
		})
	}, 300, 200)
	tt.Move(70, 50)
	settles(t, tt, "Card lift")
}

func TestFrameClips(t *testing.T) {
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		ui.Column(c).Padding(20).Children(func() {
			Frame(c, func() { ui.Box(c).Size(100, 60).Background(th.Accent) })
		})
	}, 300, 200)
	settles(t, tt, "Frame fade")
	if img := tt.Image(); near(img, 20, 20, th.Accent) || !near(img, 70, 50, th.Accent) {
		t.Error("Frame did not clip its corner")
	}
}

func TestAspectBox(t *testing.T) {
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		ui.Column(c).Width(300).Children(func() {
			AspectBox(c, 2, func() { ui.Box(c).Fill().Background(th.Accent) })
		})
	}, 300, 300)
	if img := tt.Image(); !near(img, 150, 140, th.Accent) || near(img, 150, 160, th.Accent) {
		t.Error("AspectBox is not 300x150")
	}
	panics(t, "ratio 0", func() { ui.Render(func(c *ui.Context) { AspectBox(c, 0, nil) }, 10, 10, 1) })
}

func TestAppFrame(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		AppFrame(c, func() { ui.Text(c, "Top") }, func() { ui.Text(c, "Body") }, func() { ui.Text(c, "Status") })
	}, 400, 300)
	top, body, st := find(t, tt, "Top"), find(t, tt, "Body"), find(t, tt, "Status")
	if !(top.Y < body.Y && body.Y < st.Y) {
		t.Errorf("order %v %v %v", top, body, st)
	}
	mid := func(r ui.Rect) int { return int(r.X + r.W/2) }
	if inked(tt, mid(top), top.Y, top.Y+top.H) || inked(tt, mid(st), st.Y, st.Y+st.H) {
		t.Error("bars shown before they slid in")
	}
	settles(t, tt, "AppFrame bars")
	if !inked(tt, mid(top), top.Y-2, top.Y+top.H+2) || !inked(tt, mid(st), st.Y-2, st.Y+st.H+2) {
		t.Error("bars never arrived")
	}
}

func TestHolyGrail(t *testing.T) {
	aside := true
	tt := ui.NewTester(func(c *ui.Context) {
		s := HolyGrailSlots{
			Header: func() { ui.Text(c, "Head") }, Nav: func() { ui.Text(c, "Nav") },
			Main: func() { ui.Text(c, "Main") }, Footer: func() { ui.Text(c, "Foot") },
		}
		if aside {
			s.Aside = func() { ui.Text(c, "Aside") }
		}
		HolyGrail(c, s)
	}, 600, 400)
	nav, main, side := find(t, tt, "Nav"), find(t, tt, "Main"), find(t, tt, "Aside")
	if !(nav.X < main.X && main.X < side.X) || find(t, tt, "Head").Y >= main.Y || find(t, tt, "Foot").Y <= main.Y {
		t.Error("slots out of place")
	}
	aside = false
	tt.Frame()
	settles(t, tt, "HolyGrail aside exit")
	if tt.HasText("Aside") {
		t.Error("aside still shown")
	}
	panics(t, "no Main", func() { ui.Render(func(c *ui.Context) { HolyGrail(c, HolyGrailSlots{}) }, 10, 10, 1) })
}

func TestRejectsBadNumbers(t *testing.T) {
	nan, inf := float32(math.NaN()), float32(math.Inf(1))
	none := func() {}
	cases := map[string]func(c *ui.Context){
		"Stack gap -1":        func(c *ui.Context) { Stack(c, -1, none) },
		"Stack gap NaN":       func(c *ui.Context) { Stack(c, nan, none) },
		"Inline gap Inf":      func(c *ui.Context) { Inline(c, inf, none) },
		"Container NaN":       func(c *ui.Context) { Container(c, ContainerOptions{MaxWidth: nan}, none) },
		"Container Inf":       func(c *ui.Context) { Container(c, ContainerOptions{MaxWidth: inf}, none) },
		"AspectBox NaN":       func(c *ui.Context) { AspectBox(c, nan, none) },
		"AspectBox Inf":       func(c *ui.Context) { AspectBox(c, inf, none) },
		"AutoGrid NaN":        func(c *ui.Context) { AutoGrid(c, nan, nil) },
		"AutoGrid Inf":        func(c *ui.Context) { AutoGrid(c, inf, nil) },
		"SplitPane size NaN":  func(c *ui.Context) { SplitPane(c, &nan, none, none) },
		"SplitStack size Inf": func(c *ui.Context) { SplitStack(c, &inf, none, none) },
		"SplitPane size -1":   func(c *ui.Context) { SplitPane(c, new(float32(-1)), none, none) },
		"SplitStack size -1":  func(c *ui.Context) { SplitStack(c, new(float32(-1)), none, none) },
	}
	for what, build := range cases {
		panics(t, what, func() { ui.Render(func(c *ui.Context) { build(c) }, 10, 10, 1) })
	}
}
