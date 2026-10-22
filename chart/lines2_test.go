package chart

import (
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var (
	web = Series{Name: "Web", Values: []float64{3, 5.5, 2}}
	app = Series{Name: "App", Values: []float64{1, 2.5, 4}}
)

func TestLineVariantsPick(t *testing.T) {
	cases := []struct {
		name string
		view func(c *ui.Context, at *int)
		tip  string
	}{
		{"Spline", func(c *ui.Context, at *int) { Spline(c, "x", months, []Series{web, app}, at, SplineOptions{}) }, "Feb: Web 5.5, App 2.5"},
		{"Stream", func(c *ui.Context, at *int) { Stream(c, "x", months, []Series{web, app}, at, StreamOptions{}) }, "Feb: Web 5.5, App 2.5"},
		{"Crosshair", func(c *ui.Context, at *int) { Crosshair(c, "x", months, []Series{web}, at, CrosshairOptions{}) }, "Feb: Web 5.5"},
		{"Annotated", func(c *ui.Context, at *int) {
			Annotated(c, "x", months, []Series{web}, []Note{{At: 0, Text: "Launch"}}, at, AnnotatedOptions{})
		}, "Feb: Web 5.5"},
		{"Threshold", func(c *ui.Context, at *int) { Threshold(c, "x", months, web.Values, 4, at, ThresholdOptions{}) }, "Feb: x 5.5"},
		{"ToggleLegend", func(c *ui.Context, at *int) {
			ToggleLegend(c, "x", months, []Series{web, app}, &[]bool{false, false}, at, ToggleLegendOptions{})
		}, "Feb: Web 5.5, App 2.5"},
		{"RangeBand", func(c *ui.Context, at *int) {
			RangeBand(c, "x", months, web.Values, []float64{2, 4, 1}, []float64{4, 7, 3}, at, RangeBandOptions{})
		}, "Feb: 5.5 (4–7)"},
	}
	for _, k := range cases {
		t.Run(k.name, func(t *testing.T) {
			at := -1
			tt := ui.NewTester(func(c *ui.Context) { k.view(c, &at) }, 500, 340)
			if !tt.FramePending() {
				t.Error("nothing grows in")
			}
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

func TestCrosshairReads(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Crosshair(c, "x", months, []Series{web, app}, new(-1), CrosshairOptions{}) }, 500, 340)
	settles(t, tt, "chart")
	if slices.Contains(tt.Texts(), "5.5") {
		t.Fatal("readout shows before hover")
	}
	feb := find(t, tt, "Feb")
	tt.Move(feb.X+feb.W/2, 120)
	tt.Frame()
	if !slices.Contains(tt.Texts(), "5.5") || !slices.Contains(tt.Texts(), "2.5") {
		t.Errorf("no readouts in %q", tt.Texts())
	}
}

func TestAnnotatedNote(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Annotated(c, "x", months, []Series{web}, []Note{{At: 1, Text: "Launch"}}, new(-1), AnnotatedOptions{})
	}, 500, 340)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("note did not arrive")
	}
	settles(t, tt, "note")
	find(t, tt, "Launch")
}

func TestThresholdTurnsDanger(t *testing.T) {
	// A second 9 holds the axis still while Feb drops.
	vs := []float64{3, 9, 2, 9}
	var danger ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		danger = c.Theme().Danger
		Threshold(c, "x", []string{"Jan", "Feb", "Mar", "Apr"}, vs, 6, new(-1), ThresholdOptions{})
	}, 500, 340)
	settles(t, tt, "chart")
	all := ui.Rect{W: 500, H: 340}
	past := inked(tt, all, danger)
	vs = []float64{3, 4, 2, 9}
	tt.Frame()
	settles(t, tt, "spring")
	if under := inked(tt, all, danger); past-under < 10 {
		t.Errorf("danger pixels %d past, %d under", past, under)
	}
	find(t, tt, "Limit 6")
}

func TestToggleLegendHides(t *testing.T) {
	hidden := []bool{false, false}
	changed := 0
	var green ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		green = tone(c.Theme(), 1)
		if ToggleLegend(c, "x", months, []Series{web, app}, &hidden, new(-1), ToggleLegendOptions{}).Changed() {
			changed++
		}
	}, 500, 380)
	settles(t, tt, "chart")
	all := ui.Rect{W: 500, H: 380}
	shown := inked(tt, all, green)
	before := baseline()
	click(t, tt, "App")
	tt.Frame()
	if !hidden[1] || changed != 1 {
		t.Fatalf("hidden %v, changed %d", hidden, changed)
	}
	cues(t, before, sound.Toggle)
	tt.Advance(50 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("series did not fade")
	}
	settles(t, tt, "fade")
	if gone := inked(tt, all, green); gone*4 > shown {
		t.Errorf("green pixels %d then %d", shown, gone)
	}
}

func TestStreamingScrolls(t *testing.T) {
	vs := []float64{1, 3, 2}
	tt := ui.NewTester(func(c *ui.Context) { Streaming(c, "Load", vs, StreamingOptions{Window: 5}) }, 500, 340)
	settles(t, tt, "chart")
	vs = append(vs, 7)
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("line did not scroll")
	}
	settles(t, tt, "scroll")
	if !slices.Contains(tt.Texts(), "7") {
		t.Errorf("readout missing in %q", tt.Texts())
	}
	// The same value keeps the readout still.
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	vs = append(vs, 7)
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	if tt.FramePending() {
		t.Error("scrolled under Reduce Motion")
	}
}

func TestZoomChart(t *testing.T) {
	labels := []string{"a", "b", "c", "d", "e"}
	s := Series{Name: "s", Values: []float64{1, 4, 2, 5, 3}}
	from, to, changed := 0, 2, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ZoomChart(c, "x", labels, []Series{s}, &from, &to, ZoomChartOptions{}).Changed() {
			changed++
		}
	}, 500, 400)
	settles(t, tt, "chart")
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if from != 1 || to != 3 {
		t.Fatalf("Right gave %d–%d", from, to)
	}
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if from != 1 || to != 2 {
		t.Fatalf("Up gave %d–%d", from, to)
	}
	cues(t, before, sound.Tick, sound.Tick)
	settles(t, tt, "keys")
	r := find(t, tt, "Zoom range")
	tt.Press(r.X+2, r.Y+r.H/2)
	for i := 1; i <= 8; i++ {
		tt.Advance(16 * time.Millisecond)
		tt.Move(r.X+2+(r.W-4)*float32(i)/8, r.Y+r.H/2)
	}
	tt.Release(r.X+r.W-2, r.Y+r.H/2)
	tt.Frame()
	if from != 0 || to != 4 || changed < 3 {
		t.Errorf("drag gave %d–%d, changed %d", from, to, changed)
	}
	if !tt.HasText("e") {
		t.Errorf("zoomed chart misses e: %q", tt.Texts())
	}
}

func TestNetworkGraphDrags(t *testing.T) {
	at, changed := -1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if NetworkGraph(c, "x", []string{"Ann", "Bo", "Cy"}, []Edge{{0, 1}, {1, 2}}, &at, NetworkGraphOptions{}).Changed() {
			changed++
		}
	}, 500, 340)
	settles(t, tt, "chart")
	r := find(t, tt, "Bo")
	before := baseline()
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	for i := 1; i <= 8; i++ {
		tt.Advance(16 * time.Millisecond)
		tt.Move(x+float32(i)*5, y)
	}
	tt.Release(x+40, y)
	tt.Frame()
	if moved := find(t, tt, "Bo"); moved.X-r.X < 30 {
		t.Errorf("node at %v after drag from %v", moved, r)
	}
	if at != -1 || changed != 1 {
		t.Errorf("at %d, changed %d", at, changed)
	}
	cues(t, before, sound.Drop)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if changed != 2 {
		t.Errorf("arrow key changed %d", changed)
	}
}

func TestChartCard(t *testing.T) {
	var danger ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		danger = c.Theme().Danger
		ChartCard(c, "Revenue", "$12k", -3.5, ChartCardOptions{}, func() {
			Sparkline(c, "Trend", []float64{1, 3, 2}, SparklineOptions{})
		})
	}, 400, 300)
	if !tt.FramePending() {
		t.Error("card did not arrive")
	}
	settles(t, tt, "card")
	for _, s := range []string{"Revenue", "$12k", "-3.5%"} {
		if !slices.Contains(tt.Texts(), s) {
			t.Errorf("no %q in %q", s, tt.Texts())
		}
	}
	if inked(tt, ui.Rect{W: 400, H: 300}, danger) == 0 {
		t.Error("a fall is not danger")
	}
}

func TestComboLinePicks(t *testing.T) {
	at := -1
	tt := ui.NewTester(func(c *ui.Context) {
		Combo(c, "x", months, Series{Name: "Bars", Values: []float64{1, 1, 10}}, Series{Name: "Line", Values: []float64{0, 5, 10}}, &at, ComboOptions{})
	}, 500, 340)
	settles(t, tt, "chart")
	plot, feb := find(t, tt, "x"), find(t, tt, "Feb")
	before := baseline()
	hit := float32(-1)
	for y := plot.Y; y < feb.Y && at < 0; y += 2 {
		tt.ClickAt(feb.X+feb.W/2, y)
		tt.Frame()
		hit = y
	}
	// Feb's bar fills only the bottom tenth.
	if at != 1 || hit > plot.Y+(feb.Y-plot.Y)*0.7 {
		t.Errorf("at %d, first hit at y %v of %v–%v", at, hit, plot.Y, feb.Y)
	}
	cues(t, before, sound.Select)
}

func TestZoomRangeClearsPick(t *testing.T) {
	labels := []string{"a", "b", "c", "d", "e"}
	from, to := 0, 4
	tt := ui.NewTester(func(c *ui.Context) {
		ZoomChart(c, "x", labels, []Series{{Name: "s", Values: []float64{1, 4, 2, 5, 3}}}, &from, &to, ZoomChartOptions{})
	}, 500, 400)
	settles(t, tt, "chart")
	tt.Key(0, ui.KeyTab)
	for range 4 {
		tt.Key(0, ui.KeyRight)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	to = 1
	tt.Frame()
	settles(t, tt, "zoom")
}

func TestMarimekkoThinColumn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Marimekko(c, "x", []string{"a", "b"}, []Series{{Name: "s", Values: []float64{1, 999}}}, new(-1), MarimekkoOptions{})
	}, 500, 340)
	settles(t, tt, "chart")
	if r := find(t, tt, "a, s: 1 (100%)"); r.W > 2 {
		t.Errorf("a 0.1%% column is %v wide", r.W)
	}
}

func TestC31Springs(t *testing.T) {
	cases := []struct {
		name string
		view func(c *ui.Context, flip bool)
	}{
		{"Waffle", func(c *ui.Context, flip bool) {
			Waffle(c, "x", []string{"a", "b"}, []float64{3, map[bool]float64{false: 1, true: 3}[flip]}, new(-1), WaffleOptions{})
		}},
		{"Rings", func(c *ui.Context, flip bool) {
			Rings(c, "x", []string{"a"}, []float64{map[bool]float64{false: 0.2, true: 0.8}[flip]}, new(-1), RingsOptions{})
		}},
		{"Lollipop", func(c *ui.Context, flip bool) {
			Lollipop(c, "x", months, []float64{3, map[bool]float64{false: 1, true: 3}[flip], 2}, new(-1), LollipopOptions{})
		}},
		{"Dendrogram", func(c *ui.Context, flip bool) {
			Dendrogram(c, "x", []string{"a", "b", "c"}, []Join{{0, 1, map[bool]float64{false: 1, true: 2}[flip]}, {3, 2, 3}}, new(-1), DendrogramOptions{})
		}},
		{"Bump", func(c *ui.Context, flip bool) {
			Bump(c, "x", months, []Series{web, {Name: "App", Values: []float64{1, map[bool]float64{false: 2, true: 9}[flip], 4}}}, new(-1), BumpOptions{})
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
