package layout

import (
	"math"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestShelf(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(func(c *ui.Context) {
		_, seeAll := Shelf(c, "Recent", ShelfOptions{}, func() { ui.Text(c, "Album") })
		if seeAll.Clicked() {
			clicked++
		}
	}, 400, 200)
	r := find(t, tt, "See all")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "Shelf hover")
	taps := count(sound.Tap)
	if err := tt.Click("See all"); err != nil {
		t.Fatal(err)
	}
	if clicked != 1 || count(sound.Tap) != taps+1 {
		t.Errorf("clicked %d, cues %d", clicked, count(sound.Tap)-taps)
	}
	find(t, tt, "Album")
}

func TestIsland(t *testing.T) {
	two := false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			Island(c, func() {
				ui.Text(c, "One")
				if two {
					ui.Text(c, "Two")
				}
			})
		})
	}, 300, 200)
	two = true
	tt.Frame()
	settles(t, tt, "Island resize")
	if find(t, tt, "Two").X <= find(t, tt, "One").X {
		t.Error("new content did not join the pill")
	}
}

func TestBoard(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		columns := []func(){
			func() { ui.Text(c, "Backlog") },
			func() { ui.Text(c, "Doing") },
		}
		Board(c, columns)
	}, 400, 200)
	settles(t, tt, "Board columns")
	backlog, doing := find(t, tt, "Backlog"), find(t, tt, "Doing")
	if doing.X <= backlog.X || doing.Y != backlog.Y {
		t.Errorf("columns at %v %v", backlog, doing)
	}
}

func TestLadder(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Ladder(c, []func(){func() { ui.Text(c, "One") }, func() { ui.Text(c, "Two") }})
	}, 300, 200)
	settles(t, tt, "Ladder rows")
	one, two := find(t, tt, "One"), find(t, tt, "Two")
	if one.X+one.W >= two.X || two.Y <= one.Y {
		t.Errorf("rows did not alternate: %v %v", one, two)
	}
}

func TestGrouped(t *testing.T) {
	n := 2
	tt := ui.NewTester(func(c *ui.Context) {
		Grouped(c, []func(){
			func() { ui.Text(c, "First") },
			func() {
				if n >= 2 {
					ui.Text(c, "Second")
				}
			},
			func() {
				if n >= 3 {
					ui.Text(c, "Third")
				}
			},
		}[:n])
	}, 300, 200)
	n = 3
	tt.Frame()
	settles(t, tt, "Grouped enter")
	find(t, tt, "Third")
	n = 1
	tt.Frame()
	settles(t, tt, "Grouped exit")
	if tt.HasText("Second") || tt.HasText("Third") {
		t.Errorf("removed rows remain: %q", tt.Texts())
	}
}

func TestSheet(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Sheet(c, func() { ui.Text(c, "Letter") })
	}, 500, 300)
	settles(t, tt, "Sheet arrival")
	find(t, tt, "Letter")
}

func TestLayoutValidation(t *testing.T) {
	nan, inf := float32(math.NaN()), float32(math.Inf(1))
	none := func() {}
	items := []RailItem{{Icon: "house", Label: "Home"}}
	cases := map[string]func(c *ui.Context){
		"Gutter negative":     func(c *ui.Context) { Gutter(c, -1) },
		"Gutter NaN":          func(c *ui.Context) { Gutter(c, nan) },
		"Tiles zero":          func(c *ui.Context) { Tiles(c, 0, nil) },
		"Tiles Inf":           func(c *ui.Context) { Tiles(c, inf, nil) },
		"Ratio mismatch":      func(c *ui.Context) { Ratio(c, []float32{1}, nil) },
		"Ratio negative":      func(c *ui.Context) { Ratio(c, []float32{-1}, []func(){none}) },
		"Ratio NaN":           func(c *ui.Context) { Ratio(c, []float32{nan}, []func(){none}) },
		"Ratio zero sum":      func(c *ui.Context) { Ratio(c, []float32{0}, []func(){none}) },
		"Responsive zero":     func(c *ui.Context) { Responsive(c, 0, nil) },
		"Responsive Inf":      func(c *ui.Context) { Responsive(c, inf, nil) },
		"Stage zero":          func(c *ui.Context) { Stage(c, 0, 10, none) },
		"Stage NaN":           func(c *ui.Context) { Stage(c, 10, nan, none) },
		"Float NaN":           func(c *ui.Context) { Float(c, nan, 0, none) },
		"Float Inf":           func(c *ui.Context) { Float(c, 0, inf, none) },
		"Pane nil":            func(c *ui.Context) { Pane(c, "Pane", nil, none) },
		"Rail empty":          func(c *ui.Context) { Rail(c, nil, new(int)) },
		"Rail nil":            func(c *ui.Context) { Rail(c, items, nil) },
		"Rail out of range":   func(c *ui.Context) { i := 1; Rail(c, items, &i) },
		"Dock invalid side":   func(c *ui.Context) { Dock(c, DockSide(9), true, none) },
		"Corner invalid spot": func(c *ui.Context) { Corner(c, CornerSpot(9), none) },
	}
	for what, build := range cases {
		panics(t, what, func() { ui.Render(func(c *ui.Context) { build(c) }, 10, 10, 1) })
	}
}
