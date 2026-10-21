package layout

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestDock(t *testing.T) {
	for _, tc := range []struct {
		name string
		side DockSide
	}{{"right", DockRight}, {"bottom", DockBottom}} {
		t.Run(tc.name, func(t *testing.T) {
			open := true
			clicked := 0
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ui.Box(c).Size(400, 200).Children(func() {
					if ui.Button(c, "Canvas").Fill().Clicked() {
						clicked++
					}
					bounds = Dock(c, tc.side, open, func() { ui.Text(c, "Inspector") }).LayoutBounds()
				})
			}, 400, 200)
			tt.Frame()
			initial := bounds
			find(t, tt, "Inspector")
			if tc.side == DockRight && abs(initial.X+initial.W-400) > 0.5 {
				t.Errorf("right root at %v", initial)
			}
			if tc.side == DockBottom && abs(initial.Y+initial.H-200) > 0.5 {
				t.Errorf("bottom root at %v", initial)
			}
			open = false
			tt.Frame()
			tt.Advance(50 * time.Millisecond)
			if bounds != initial {
				t.Errorf("exit moved root from %v to %v", initial, bounds)
			}
			settles(t, tt, "Dock exit")
			if tt.HasText("Inspector") {
				t.Error("closed dock remains")
			}
			if err := tt.Click("Canvas"); err != nil {
				t.Fatal(err)
			}
			if clicked != 1 {
				t.Error("closed dock blocked the canvas")
			}
			open = true
			tt.Frame()
			settles(t, tt, "Dock enter")
			find(t, tt, "Inspector")
		})
	}
}

func TestRail(t *testing.T) {
	active, changed := 0, 0
	pad := float32(0)
	items := []RailItem{{Icon: "house", Label: "Home"}, {Icon: "search", Label: "Search"}}
	tt := ui.NewTester(func(c *ui.Context) {
		pad = c.Theme().Space(2)
		if Rail(c, items, &active).Changed() {
			changed++
		}
	}, 300, 200)
	if home := find(t, tt, "Home"); abs(home.X-pad) > 0.5 {
		t.Errorf("rail stretched, first slot at %v", home)
	}
	taps := count(sound.Tap)
	if err := tt.Click("Search"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "Rail selection")
	if active != 1 || changed != 1 || count(sound.Tap) != taps+1 {
		t.Errorf("active %d, changed %d, cues %d", active, changed, count(sound.Tap)-taps)
	}
}

func TestTiles(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			Tiles(c, 80, []func(){
				func() { ui.Text(c, "One") },
				func() { ui.Text(c, "Two") },
			})
		})
	}, 300, 200)
	one, two := find(t, tt, "One"), find(t, tt, "Two")
	if one.Y != two.Y || two.X <= one.X {
		t.Errorf("tile positions %v %v", one, two)
	}
	tt.Move(one.X+one.W/2, one.Y+one.H/2)
	settles(t, tt, "Tile lift")
}

func TestShowcase(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Showcase(c, func() { ui.Text(c, "Featured") })
	}, 300, 200)
	settles(t, tt, "Showcase")
	r := find(t, tt, "Featured")
	if got := r.X + r.W/2; abs(got-150) > 1 {
		t.Errorf("content center %v", got)
	}
}

func TestRatio(t *testing.T) {
	weights := []float32{1, 3}
	tt := ui.NewTester(func(c *ui.Context) {
		Ratio(c, weights, []func(){func() { ui.Text(c, "A") }, func() { ui.Text(c, "B") }})
	}, 400, 200)
	before := find(t, tt, "B").X
	weights = []float32{3, 1}
	tt.Frame()
	settles(t, tt, "Ratio weights")
	if after := find(t, tt, "B").X; after <= before {
		t.Errorf("second column moved from %v to %v", before, after)
	}
}

func TestResponsive(t *testing.T) {
	width := float32(500)
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Width(width).Children(func() {
			Responsive(c, 400, []func(){func() { ui.Text(c, "A") }, func() { ui.Text(c, "B") }})
		})
	}, 600, 300)
	tt.Advance(time2s)
	a, b := find(t, tt, "A"), find(t, tt, "B")
	if a.Y != b.Y {
		t.Errorf("wide items stacked: %v %v", a, b)
	}
	width = 250
	tt.Frame()
	settles(t, tt, "Responsive stack")
	a, b = find(t, tt, "A"), find(t, tt, "B")
	if b.Y <= a.Y {
		t.Errorf("narrow items remained inline: %v %v", a, b)
	}
}

func TestBreakpoint(t *testing.T) {
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		Breakpoint(c,
			func() { ui.Box(c).Label("Compact").Size(100, 100).Background(th.Accent) },
			func() { ui.Box(c).Label("Regular").Size(100, 100).Background(th.Danger) },
			func() { ui.Box(c).Label("Wide").Size(100, 100).Background(th.Success) },
		)
	}, 500, 200)
	tt.Advance(time2s)
	if img := tt.Image(); !near(img, 50, 50, th.Accent) {
		t.Error("compact slot missing")
	}
	tt.SetSize(800, 200)
	if !tt.FramePending() {
		t.Error("regular slot did not cross-fade")
	}
	if img := tt.Image(); !near(img, 50, 50, th.Accent) {
		t.Error("regular slot arrived opaque")
	}
	tt.Advance(time2s)
	if !tt.HasText("Regular") || tt.HasText("Compact") || !near(tt.Image(), 50, 50, th.Danger) {
		t.Errorf("regular texts %q", tt.Texts())
	}
	tt.SetSize(1100, 200)
	if !tt.FramePending() {
		t.Error("wide slot did not cross-fade")
	}
	if img := tt.Image(); !near(img, 50, 50, th.Danger) {
		t.Error("wide slot arrived opaque")
	}
	tt.Advance(time2s)
	if !tt.HasText("Wide") || tt.HasText("Regular") || !near(tt.Image(), 50, 50, th.Success) {
		t.Errorf("wide texts %q", tt.Texts())
	}
}

func TestStage(t *testing.T) {
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		ui.Box(c).Size(400, 300).Children(func() {
			Stage(c, 200, 100, func() { ui.Box(c).Fill().Background(th.Accent) })
		})
	}, 400, 300)
	settles(t, tt, "Stage scale")
	img := tt.Image()
	if !near(img, 390, 150, th.Accent) || !near(img, 10, 10, th.Surface) {
		t.Error("content was not scaled and letterboxed")
	}
}

func TestCorner(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Size(300, 200).Children(func() {
			Corner(c, CornerTopRight, func() { ui.Text(c, "Badge") })
		})
	}, 300, 200)
	settles(t, tt, "Corner pop")
	r := find(t, tt, "Badge")
	if r.X < 150 || r.Y > 100 {
		t.Errorf("badge at %v", r)
	}
}

func TestFloat(t *testing.T) {
	x, y := float32(40), float32(30)
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Size(300, 200).Children(func() {
			Float(c, x, y, func() { ui.Text(c, "Pin") })
		})
	}, 300, 200)
	r := find(t, tt, "Pin")
	if abs(r.X-x) > 0.5 || abs(r.Y-y) > 0.5 {
		t.Errorf("pin at %v", r)
	}
	x, y = 160, 80
	tt.Frame()
	settles(t, tt, "Float offset")
	r = find(t, tt, "Pin")
	if abs(r.X-x) > 0.5 || abs(r.Y-y) > 0.5 {
		t.Errorf("moved pin at %v", r)
	}
}
