package chart

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// spread is the height ink k covers down column x.
func spread(tt *ui.Tester, x int, r ui.Rect, k ui.Color) int {
	img := tt.Image()
	top, bot := -1, -1
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		if p := img.RGBAAt(x, y); near(p.R, k.R) && near(p.G, k.G) && near(p.B, k.B) {
			if top < 0 {
				top = y
			}
			bot = y
		}
	}
	return bot - top
}

// rides checks marks and line stay together while values change.
func rides(t *testing.T, tt *ui.Tester, x int, plot ui.Rect, ink ui.Color, dot float32, change func()) {
	t.Helper()
	most := int(dot) + 6
	if n := spread(tt, x, plot, ink); n < 0 || n > most {
		t.Fatalf("settled: ink spans %d px, want 0-%d", n, most)
	}
	change()
	tt.Frame()
	for _, ms := range []int{16, 34, 50, 50, 50, 100} {
		tt.Advance(time.Duration(ms) * time.Millisecond)
		if n := spread(tt, x, plot, ink); n < 0 || n > most {
			t.Fatalf("mid-spring: ink spans %d px, want the mark on the line (0-%d)", n, most)
		}
	}
}

func TestComboDotsRideLine(t *testing.T) {
	line := []float64{0, 5, 10}
	at := -1
	var ink ui.Color
	var dot float32
	tt := ui.NewTester(func(c *ui.Context) {
		ink, dot = tone(c.Theme(), 1), c.Theme().Space(3)
		Combo(c, "x", months, Series{Name: "Bars", Values: []float64{1, 1, 10}}, Series{Name: "Line", Values: line}, &at, ComboOptions{})
	}, 500, 340)
	settles(t, tt, "chart")
	bar := find(t, tt, "Feb: Bars 1, Line 5")
	rides(t, tt, int(bar.X+bar.W/2), find(t, tt, "x"), ink, dot, func() { line = []float64{0, 5, 40} })
}

func TestRangeBandRingRidesLine(t *testing.T) {
	mid, lo, hi := []float64{1, 5, 10}, []float64{0, 4, 9}, []float64{2, 6, 11}
	at := 1
	var ink ui.Color
	var dot float32
	tt := ui.NewTester(func(c *ui.Context) {
		ink, dot = c.Theme().Accent, c.Theme().Space(3)
		RangeBand(c, "x", months, mid, lo, hi, &at, RangeBandOptions{})
	}, 500, 340)
	settles(t, tt, "chart")
	feb := find(t, tt, "Feb")
	rides(t, tt, int(feb.X+feb.W/2), find(t, tt, "x"), ink, dot, func() {
		mid, lo, hi = []float64{1, 5, 40}, []float64{0, 4, 39}, []float64{2, 6, 41}
	})
}

func TestThresholdPointsRideLine(t *testing.T) {
	vs := []float64{0, 5, 10}
	at := -1
	var ink ui.Color
	var dot float32
	tt := ui.NewTester(func(c *ui.Context) {
		ink, dot = tone(c.Theme(), 0), c.Theme().Space(3)
		Threshold(c, "x", months, vs, 8, &at, ThresholdOptions{})
	}, 500, 340)
	settles(t, tt, "chart")
	feb := find(t, tt, "Feb")
	rides(t, tt, int(feb.X+feb.W/2), find(t, tt, "x"), ink, dot, func() { vs = []float64{0, 5, 40} })
}

func TestTreeLinksMeetBoxes(t *testing.T) {
	kids := []Node{{Name: "b", Value: 1}, {Name: "c", Value: 1}}
	at := -1
	var ink ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		ink = c.Theme().Border
		TreeDiagram(c, "x", []Node{{Name: "a", Children: kids}}, &at, TreeDiagramOptions{})
	}, 500, 340)
	settles(t, tt, "chart")
	kids = []Node{{Name: "b", Value: 1}, {Name: "c", Value: 1}, {Name: "d", Value: 1}, {Name: "e", Value: 1}}
	tt.Frame()
	for _, ms := range []int{16, 34, 50} {
		tt.Advance(time.Duration(ms) * time.Millisecond)
		b := find(t, tt, "b")
		p := tt.Image().RGBAAt(int(b.X+b.W/2), int(b.Y)-2)
		if !(near(p.R, ink.R) && near(p.G, ink.G) && near(p.B, ink.B)) {
			t.Fatalf("mid-move: no link just above box b (pixel %v, link %v)", p, ink)
		}
	}
}
