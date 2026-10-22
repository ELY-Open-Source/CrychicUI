package chart

import (
	"math"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestWheelsPick(t *testing.T) {
	names := []string{"iOS", "Web", "Mac"}
	vs := []float64{6, 3, 1}
	cases := []struct {
		name string
		view func(c *ui.Context, at *int)
		// x, y find part 0 from the box center.
		x, y float32
	}{
		{"Pie", func(c *ui.Context, at *int) { Pie(c, "Users", names, vs, at, PieOptions{}) }, 30, -40},
		{"Donut", func(c *ui.Context, at *int) { Donut(c, "Users", names, vs, at, DonutOptions{}) }, 30, -80},
		{"PolarArea", func(c *ui.Context, at *int) { PolarArea(c, "Users", names, vs, at, PolarAreaOptions{}) }, 20, -60},
		{"RadialBars", func(c *ui.Context, at *int) { RadialBars(c, "Users", names, vs, at, RadialBarsOptions{}) }, 20, -96},
		{"Sunburst", func(c *ui.Context, at *int) {
			Sunburst(c, "Users", []Node{{Name: "iOS", Children: []Node{{Name: "Phone", Value: 4}, {Name: "Pad", Value: 2}}}, {Name: "Web", Value: 3}}, at, SunburstOptions{})
		}, 20, -40},
		{"Radar", func(c *ui.Context, at *int) {
			Radar(c, "Skills", names, []Series{{Name: "Ada", Values: vs}}, at, RadarOptions{})
		}, 0, -50},
	}
	for _, k := range cases {
		t.Run(k.name, func(t *testing.T) {
			at := -1
			tt := ui.NewTester(func(c *ui.Context) { k.view(c, &at) }, 500, 500)
			if !tt.FramePending() {
				t.Error("nothing grows in")
			}
			settles(t, tt, "wheel")
			before := baseline()
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyRight)
			tt.Key(0, ui.KeyRight)
			tt.Key(0, ui.KeyEnter)
			tt.Frame()
			if at != 1 {
				t.Fatalf("keys picked %d, want 1", at)
			}
			r := find(t, tt, map[bool]string{true: "Skills", false: "Users"}[k.name == "Radar"])
			cx, cy := r.X+r.W/2, r.Y+r.W/2
			tt.Move(cx+k.x, cy+k.y)
			tt.Frame()
			tt.ClickAt(cx+k.x, cy+k.y)
			tt.Frame()
			if at != 0 {
				t.Fatalf("click picked %d, want 0", at)
			}
			cues(t, before, sound.Select, sound.Select)
		})
	}
}

func TestDrawnPick(t *testing.T) {
	two := []Series{{Name: "Web", Values: []float64{3, 5, 2}}, {Name: "App", Values: []float64{1, 2, 4}}}
	candles := []Candle{{"Jan", 10, 14, 9, 13}, {"Feb", 13, 15, 11, 12}, {"Mar", 12, 16, 12, 15}}
	cases := []struct {
		name string
		view func(c *ui.Context, at *int)
		tip  string
	}{
		{"Area", func(c *ui.Context, at *int) { Area(c, "x", months, two, at, AreaOptions{}) }, "Feb: Web 5, App 2"},
		{"StackedArea", func(c *ui.Context, at *int) { StackedArea(c, "x", months, two, at, StackedAreaOptions{}) }, "Feb: Web 5, App 2"},
		{"StepLine", func(c *ui.Context, at *int) { StepLine(c, "x", months, two, at, StepLineOptions{}) }, "Feb: Web 5, App 2"},
		{"OHLC", func(c *ui.Context, at *int) { OHLC(c, "x", candles, at, OHLCOptions{}) }, "Feb: open 13, high 15, low 11, close 12"},
		{"Violin", func(c *ui.Context, at *int) {
			Violin(c, "x", []Group{{"Jan", []float64{1, 2, 3}}, {"Feb", []float64{2, 4, 4, 6}}, {"Mar", []float64{5}}}, at, ViolinOptions{})
		}, "Feb: median 4, 4 values"},
	}
	for _, k := range cases {
		t.Run(k.name, func(t *testing.T) {
			at := -1
			tt := ui.NewTester(func(c *ui.Context) { k.view(c, &at) }, 500, 300)
			settles(t, tt, "chart")
			feb := find(t, tt, "Feb")
			tt.Move(feb.X+feb.W/2, 120)
			tt.Frame()
			if !tt.HasText(k.tip) {
				t.Fatalf("no tip %q in %q", k.tip, tt.Texts())
			}
			before := baseline()
			tt.ClickAt(feb.X+feb.W/2, 120)
			tt.Frame()
			if at != 1 {
				t.Errorf("at %d", at)
			}
			cues(t, before, sound.Select)
		})
	}
}

func TestSingleValueCharts(t *testing.T) {
	v := 40.0
	tt := ui.NewTester(func(c *ui.Context) {
		GaugeChart(c, "Load", v, 0, 100, GaugeChartOptions{Bands: []float64{60, 85}})
		Bullet(c, "Revenue", v, 50, BulletOptions{Bands: []float64{30, 60, 90}})
		Sparkline(c, "Trend", []float64{1, 4, 2, v}, SparklineOptions{})
		SparkBars(c, "Daily", []float64{1, 4, 2, v}, SparkBarsOptions{})
	}, 600, 600)
	if !tt.FramePending() {
		t.Error("nothing grows in")
	}
	settles(t, tt, "charts")
	if !tt.HasText("40") {
		t.Fatalf("gauge value missing: %q", tt.Texts())
	}
	find(t, tt, "Revenue: 40, target 50")
	v = 70
	tt.Frame()
	tt.Advance(50 * 1e6)
	if !tt.FramePending() {
		t.Error("needle did not spring")
	}
	settles(t, tt, "spring")
	panics(t, "gauge range", func() { render(func(c *ui.Context) { GaugeChart(c, "x", 120, 0, 100, GaugeChartOptions{}) }) })
	panics(t, "bands", func() { render(func(c *ui.Context) { Bullet(c, "x", 1, 2, BulletOptions{Bands: []float64{5, 3}}) }) })
	panics(t, "spark NaN", func() { render(func(c *ui.Context) { Sparkline(c, "x", []float64{math.NaN()}, SparklineOptions{}) }) })
	panics(t, "loop", func() {
		render(func(c *ui.Context) {
			n := -1
			Sankey(c, "x", []string{"a", "b"}, []Flow{{0, 1, 1}, {1, 0, 1}}, &n, SankeyOptions{})
		})
	})
}

func TestStats(t *testing.T) {
	if q := quartiles([]float64{5, 1, 3, 2, 4}); q != [5]float64{1, 2, 3, 4, 5} {
		t.Errorf("quartiles %v", q)
	}
	d := density([]float64{0, 0, 0, 10}, []float64{0, 5, 10})
	if !(d[0] > d[1] && d[1] > d[2]) {
		t.Errorf("density %v", d)
	}
	tiles := squarify([]float64{6, 3, 1})
	area := 0.0
	for _, tl := range tiles {
		w, h := float64(tl[0].Hi-tl[0].Lo), float64(tl[1].Hi-tl[1].Lo)
		if w <= 0 || h <= 0 {
			t.Fatalf("empty tile %v", tl)
		}
		area += w * h
	}
	if math.Abs(area-1) > 1e-4 || math.Abs(float64((tiles[0][0].Hi-tiles[0][0].Lo)*(tiles[0][1].Hi-tiles[0][1].Lo))-0.6) > 1e-4 {
		t.Errorf("tiles %v cover %v", tiles, area)
	}
}
