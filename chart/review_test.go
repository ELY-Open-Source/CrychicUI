package chart

import (
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

func TestWheelShrinks(t *testing.T) {
	names, vs := []string{"a", "b", "c"}, []float64{1, 1, 1}
	at := -1
	tt := ui.NewTester(func(c *ui.Context) { Pie(c, "Users", names, vs, &at, PieOptions{}) }, 500, 500)
	settles(t, tt, "pie")
	tt.Key(0, ui.KeyTab)
	for range 3 {
		tt.Key(0, ui.KeyRight)
	}
	tt.Frame()
	names, vs = names[:2], vs[:2]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if at != -1 {
		t.Errorf("a gone part was picked: %d", at)
	}
}

func TestRestingPointerKeepsKeys(t *testing.T) {
	at := -1
	tt := ui.NewTester(func(c *ui.Context) { Pie(c, "Users", []string{"a", "b", "c"}, []float64{1, 1, 1}, &at, PieOptions{}) }, 500, 500)
	settles(t, tt, "pie")
	r := find(t, tt, "Users")
	tt.Move(r.X+r.W/2+30, r.Y+r.W/2-40)
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	for range 3 {
		tt.Frame()
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if at != 1 {
		t.Errorf("Enter picked %d, want the key's 1", at)
	}
}

func TestWheelClickHitsPointer(t *testing.T) {
	for _, key := range []ui.Key{ui.KeyRight, ui.KeyEscape} {
		at := -1
		tt := ui.NewTester(func(c *ui.Context) { Pie(c, "Users", []string{"a", "b", "c"}, []float64{1, 1, 1}, &at, PieOptions{}) }, 500, 500)
		settles(t, tt, "pie")
		r := find(t, tt, "Users")
		x, y := r.X+r.W/2+30, r.Y+r.W/2-40
		tt.Move(x, y)
		tt.Frame()
		tt.Key(0, ui.KeyTab)
		tt.Key(0, key)
		tt.Frame()
		before := baseline()
		for range 2 {
			tt.ClickAt(x, y)
			tt.Frame()
		}
		if at != 0 {
			t.Errorf("after %v a click picked %d, want 0", key, at)
		}
		cues(t, before, sound.Select, sound.Select)
	}
}

func TestSunburstHitsShownAngles(t *testing.T) {
	vs := []float64{9, 1}
	at := -1
	tt := ui.NewTester(func(c *ui.Context) {
		Sunburst(c, "Users", []Node{{Name: "a", Value: vs[0]}, {Name: "b", Value: vs[1]}}, &at, SunburstOptions{})
	}, 500, 500)
	settles(t, tt, "sunburst")
	r := find(t, tt, "Users")
	vs = []float64{1, 9}
	tt.Frame()
	tt.Advance(16 * time.Millisecond)
	tt.ClickAt(r.X+r.W/2, r.Y+r.W/2+60)
	tt.Frame()
	if at != 0 {
		t.Errorf("click on the shown a picked %d", at)
	}
}

func TestWheelHitsWhileGrowing(t *testing.T) {
	const ms = 150
	g := ui.EaseOut(float32(ms) / 700)
	two, nodes := []string{"a", "b"}, []Node{{Name: "a", Value: 1}, {Name: "b", Value: 1}}
	cases := []struct {
		name   string
		view   func(c *ui.Context, at *int)
		turn   float64
		reach  func(w float32) float32
		picked int
	}{
		{"Pie", func(c *ui.Context, at *int) { Pie(c, "x", two, []float64{1, 1}, at, PieOptions{}) },
			1.5 * math.Pi * float64(g), func(w float32) float32 { return w / 4 }, 1},
		{"Sunburst", func(c *ui.Context, at *int) { Sunburst(c, "x", nodes, at, SunburstOptions{}) },
			1.5 * math.Pi * float64(g), func(w float32) float32 { return w / 4 }, 1},
		{"PolarArea", func(c *ui.Context, at *int) { PolarArea(c, "x", two, []float64{1, 1}, at, PolarAreaOptions{}) },
			0.5 * math.Pi, func(w float32) float32 { return w * 0.4 }, -1},
		{"RadialBars", func(c *ui.Context, at *int) { RadialBars(c, "x", two, []float64{1, 1}, at, RadialBarsOptions{}) },
			1.3 * math.Pi, func(w float32) float32 { o := w/2 - 4; return o - o*0.75/2*0.4 }, -1},
	}
	if 1.5*float64(g) >= 1 {
		t.Fatalf("grow-in %v too far along for these spots", g)
	}
	for _, k := range cases {
		t.Run(k.name, func(t *testing.T) {
			at := -1
			tt := ui.NewTester(func(c *ui.Context) { k.view(c, &at) }, 500, 500)
			r := find(t, tt, "x")
			tt.Advance(ms * time.Millisecond)
			d := k.reach(r.W)
			x := r.X + r.W/2 + d*float32(math.Sin(k.turn))
			y := r.Y + r.W/2 - d*float32(math.Cos(k.turn))
			tt.Move(x, y)
			tt.ClickAt(x, y)
			tt.Frame()
			if at != k.picked {
				t.Errorf("click on the drawn shape picked %d, want %d", at, k.picked)
			}
		})
	}
}

func TestLiftedSliceHits(t *testing.T) {
	for _, donut := range []bool{false, true} {
		at := -1
		var size, outer, lift float32
		tt := ui.NewTester(func(c *ui.Context) {
			th := c.Theme()
			size, lift = th.Space(56), th.Space(1.5)
			outer = size/2 - th.Space(2)
			if donut {
				Donut(c, "x", []string{"a", "b"}, []float64{1, 1}, &at, DonutOptions{})
			} else {
				Pie(c, "x", []string{"a", "b"}, []float64{1, 1}, &at, PieOptions{})
			}
		}, 500, 500)
		at = 0
		tt.Frame()
		settles(t, tt, "lift")
		r := find(t, tt, "x")
		at = -1
		tt.Frame()
		before := baseline()
		mid := mc.Pie([]float64{1, 1}, 0.02)[0].Mid()
		x, y := mc.Polar(r.X+r.W/2, r.Y+r.W/2, outer+lift/2, mid)
		tt.ClickAt(x, y)
		tt.Frame()
		if at != 0 {
			t.Errorf("donut %v: click on the lifted strip picked %d", donut, at)
		}
		cues(t, before, sound.Select)
	}
}

func TestValuesSpring(t *testing.T) {
	cases := []struct {
		name string
		view func(c *ui.Context, flip bool)
	}{
		{"Sunburst", func(c *ui.Context, flip bool) {
			v := map[bool]float64{false: 2, true: 5}[flip]
			Sunburst(c, "x", []Node{{Name: "a", Value: 3}, {Name: "b", Value: v}}, new(-1), SunburstOptions{})
		}},
		{"Violin", func(c *ui.Context, flip bool) {
			mid := map[bool]float64{false: 2, true: 2.5}[flip]
			Violin(c, "x", []Group{{"a", []float64{1, mid, 3}}, {"b", []float64{0, 10}}}, new(-1), ViolinOptions{})
		}},
		{"Candles", func(c *ui.Context, flip bool) {
			low := map[bool]float64{false: 11, true: 10}[flip]
			Candles(c, "x", []Candle{{"a", 10, 14, 9, 13}, {"b", 13, 15, low, 12}, {"c", 12, 16, 12, 15}}, new(-1), CandlesOptions{})
		}},
		{"BoxPlot", func(c *ui.Context, flip bool) {
			lo := map[bool]float64{false: 1, true: 1.5}[flip]
			BoxPlot(c, "x", []Group{{"a", []float64{lo, 2, 3, 4, 5}}, {"b", []float64{0, 5, 10}}}, new(-1), BoxPlotOptions{})
		}},
	}
	for _, k := range cases {
		t.Run(k.name, func(t *testing.T) {
			flip := false
			tt := ui.NewTester(func(c *ui.Context) { k.view(c, flip) }, 500, 400)
			settles(t, tt, "chart")
			flip = true
			tt.Frame()
			tt.Advance(30 * time.Millisecond)
			if !tt.FramePending() {
				t.Error("new values snapped")
			}
			settles(t, tt, "spring")
		})
	}
}

func TestSankeyNodesGrow(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Sankey(c, "x", []string{"Ads", "Site", "Buy"}, []Flow{{0, 1, 5}, {1, 2, 3}}, new(-1), SankeyOptions{})
	}, 500, 400)
	tt.Advance(100 * time.Millisecond)
	early := find(t, tt, "Ads: 5")
	settles(t, tt, "sankey")
	if late := find(t, tt, "Ads: 5"); early.H >= late.H-1 {
		t.Errorf("node did not grow: %v then %v", early, late)
	}
}

func TestFlatViolin(t *testing.T) {
	for _, d := range density([]float64{0, 0, 0}, []float64{-1, 0.2, 1}) {
		if math.IsNaN(d) || math.IsInf(d, 0) {
			t.Fatalf("density %v", d)
		}
	}
	if d := density([]float64{0, 0, 0}, []float64{-1, 0.2, 1}); d[1] <= 0 {
		t.Errorf("no sample reached the grid: %v", d)
	}
}

func TestLinePaintStaysInPlot(t *testing.T) {
	vs := []float64{8, 2, 10}
	var accent ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		accent = c.Theme().Accent
		ui.Column(c).Padding(60).Children(func() {
			Line(c, "Visits", months, []Series{{Name: "Web", Values: vs}}, new(-1), LineOptions{})
		})
	}, 500, 400)
	settles(t, tt, "line")
	vs = []float64{1, 2, 1}
	tt.Frame()
	r := find(t, tt, "Visits")
	for range 10 {
		tt.Advance(16 * time.Millisecond)
		if n := inked(tt, ui.Rect{X: r.X, Y: 0, W: r.W, H: r.Y}, accent); n > 0 {
			t.Fatalf("%d series pixels above the plot", n)
		}
	}
}

func TestTreemapNests(t *testing.T) {
	at := -1
	tt := ui.NewTester(func(c *ui.Context) {
		Treemap(c, "Disk", []Node{
			{Name: "Photos", Children: []Node{{Name: "Old", Value: 4}, {Name: "New", Value: 2}}},
			{Name: "Music", Value: 3},
		}, &at, TreemapOptions{})
	}, 500, 400)
	settles(t, tt, "treemap")
	parent, kid := find(t, tt, "Photos: 6"), find(t, tt, "Old: 4")
	if kid.X < parent.X || kid.Y <= parent.Y || kid.X+kid.W > parent.X+parent.W || kid.Y+kid.H > parent.Y+parent.H {
		t.Fatalf("child %v not inside parent %v", kid, parent)
	}
	click(t, tt, "Old: 4")
	tt.Frame()
	if at != 1 {
		t.Errorf("at %d, want the child's 1", at)
	}
	panics(t, "zero leaf", func() {
		render(func(c *ui.Context) { Treemap(c, "x", []Node{{Name: "a"}}, new(-1), TreemapOptions{}) })
	})
}

func TestTipsTakeWords(t *testing.T) {
	candles := []Candle{{"Jan", 10, 14, 9, 13}, {"Feb", 13, 15, 11, 12}}
	tt := ui.NewTester(func(c *ui.Context) {
		Candles(c, "k", candles, new(-1), CandlesOptions{Tip: "%s: o %s h %s b %s c %s"})
		BoxPlot(c, "b", []Group{{"Jan", []float64{1, 2, 3}}}, new(-1), BoxPlotOptions{Tip: "%s: mediana %s (%s–%s)"})
	}, 500, 700)
	settles(t, tt, "charts")
	find(t, tt, "Jan: o 10 h 14 b 9 c 13")
	find(t, tt, "Jan: mediana 2 (1.5–2.5)")
	tt = ui.NewTester(func(c *ui.Context) {
		Violin(c, "v", []Group{{"Jan", []float64{1, 2, 3}}, {"Feb", []float64{2, 4, 6}}}, new(-1), ViolinOptions{Tip: "%s: mediana %s, %d valores"})
	}, 500, 300)
	settles(t, tt, "violin")
	feb := find(t, tt, "Feb")
	tt.Move(feb.X+feb.W/2, 120)
	tt.Frame()
	if !tt.HasText("Feb: mediana 4, 3 valores") {
		t.Errorf("no tip in %q", tt.Texts())
	}
}
