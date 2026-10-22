package chart

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestShapePicks(t *testing.T) {
	day := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name   string
		view   func(c *ui.Context, at *int)
		x, y   float32
		picked int
	}{
		// Top left cell belongs to the first name.
		{"Waffle", func(c *ui.Context, at *int) { Waffle(c, "x", []string{"a", "b"}, []float64{3, 1}, at, WaffleOptions{}) }, 0.03, 0.03, 0},
		// The outer ring's arc starts at twelve.
		{"Rings", func(c *ui.Context, at *int) {
			Rings(c, "x", []string{"a", "b"}, []float64{0.5, 0.9}, at, RingsOptions{})
		}, 0.55, 0.06, 0},
		// The middle of three circles is ABC.
		{"Venn", func(c *ui.Context, at *int) {
			Venn(c, "x", []string{"A", "B", "C"}, []float64{1, 2, 3, 4, 5, 6, 7}, at, VennOptions{})
		}, 0.5, 0.5, 6},
		// March 1 2026 is a Sunday: column 0, row 0.
		{"YearHeat", func(c *ui.Context, at *int) {
			YearHeat(c, "x", day, []float64{1, 2, 3, 4, 5, 6, 7, 8}, at, YearHeatOptions{})
		}, 0.02, 0.05, 0},
	}
	for _, k := range cases {
		t.Run(k.name, func(t *testing.T) {
			at := -1
			tt := ui.NewTester(func(c *ui.Context) { k.view(c, &at) }, 500, 500)
			if !tt.FramePending() {
				t.Error("nothing grows in")
			}
			settles(t, tt, "chart")
			before := baseline()
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyRight)
			tt.Key(0, ui.KeyRight)
			tt.Key(0, ui.KeyEnter)
			tt.Frame()
			if at != 1 {
				t.Fatalf("keys picked %d, want 1", at)
			}
			r := find(t, tt, "x")
			x, y := r.X+r.W*k.x, r.Y+r.H*k.y
			tt.Move(x, y)
			tt.ClickAt(x, y)
			tt.Frame()
			if at != k.picked {
				t.Errorf("click picked %d, want %d", at, k.picked)
			}
			cues(t, before, sound.Select, sound.Select)
		})
	}
}

func TestSizesValidate(t *testing.T) {
	n, v := []string{"a", "b"}, []float64{1, 2}
	for _, size := range []float32{-1, float32(math.NaN())} {
		for name, fn := range map[string]func(c *ui.Context){
			"Pie":        func(c *ui.Context) { Pie(c, "x", n, v, new(-1), PieOptions{Size: size}) },
			"Waffle":     func(c *ui.Context) { Waffle(c, "x", n, v, new(-1), WaffleOptions{Size: size}) },
			"Rings":      func(c *ui.Context) { Rings(c, "x", n, []float64{0.1, 0.2}, new(-1), RingsOptions{Size: size}) },
			"Venn":       func(c *ui.Context) { Venn(c, "x", n, []float64{1, 2, 3}, new(-1), VennOptions{Size: size}) },
			"PolarArea":  func(c *ui.Context) { PolarArea(c, "x", n, v, new(-1), PolarAreaOptions{Size: size}) },
			"RadialBars": func(c *ui.Context) { RadialBars(c, "x", n, v, new(-1), RadialBarsOptions{Size: size}) },
			"Sunburst": func(c *ui.Context) {
				Sunburst(c, "x", []Node{{Name: "a", Value: 1}}, new(-1), SunburstOptions{Size: size})
			},
			"Radar": func(c *ui.Context) {
				Radar(c, "x", []string{"a", "b", "c"}, []Series{{Name: "s", Values: []float64{1, 2, 3}}}, new(-1), RadarOptions{Size: size})
			},
			"GaugeChart": func(c *ui.Context) { GaugeChart(c, "x", 1, 0, 2, GaugeChartOptions{Size: size}) },
			"YearHeat":   func(c *ui.Context) { YearHeat(c, "x", time.Now(), v, new(-1), YearHeatOptions{Cell: size}) },
		} {
			panics(t, name, func() { render(fn) })
		}
	}
}

func TestShares(t *testing.T) {
	if got := shares([]float64{1, 1, 1}, 100); !slices.Equal(got, []int{34, 33, 33}) {
		t.Errorf("thirds %v", got)
	}
	if got := shares([]float64{3, 1}, 100); !slices.Equal(got, []int{75, 25}) {
		t.Errorf("quarters %v", got)
	}
}

func TestC31Panics(t *testing.T) {
	day := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := Series{Name: "s", Values: []float64{1, 2, 3}}
	for name, fn := range map[string]func(c *ui.Context){
		"venn counts": func(c *ui.Context) { Venn(c, "x", []string{"a", "b"}, []float64{1, 2}, new(-1), VennOptions{}) },
		"venn sets":   func(c *ui.Context) { Venn(c, "x", []string{"a"}, []float64{1}, new(-1), VennOptions{}) },
		"rings range": func(c *ui.Context) { Rings(c, "x", []string{"a"}, []float64{1.5}, new(-1), RingsOptions{}) },
		"waffle zero": func(c *ui.Context) { Waffle(c, "x", []string{"a"}, []float64{0}, new(-1), WaffleOptions{}) },
		"year days":   func(c *ui.Context) { YearHeat(c, "x", day, make([]float64, 372), new(-1), YearHeatOptions{}) },
		"gantt order": func(c *ui.Context) { Gantt(c, "x", []Task{{"a", day, day}}, new(-1), GanttOptions{}) },
		"dendro join": func(c *ui.Context) {
			Dendrogram(c, "x", []string{"a", "b"}, []Join{{0, 0, 1}}, new(-1), DendrogramOptions{})
		},
		"dendro low": func(c *ui.Context) {
			Dendrogram(c, "x", []string{"a", "b", "c"}, []Join{{0, 1, 2}, {3, 2, 1}}, new(-1), DendrogramOptions{})
		},
		"graph edge": func(c *ui.Context) {
			NetworkGraph(c, "x", []string{"a", "b"}, []Edge{{0, 0}}, new(-1), NetworkGraphOptions{})
		},
		"note point": func(c *ui.Context) {
			Annotated(c, "x", months, []Series{s}, []Note{{At: 5, Text: "n"}}, new(-1), AnnotatedOptions{})
		},
		"band order": func(c *ui.Context) {
			RangeBand(c, "x", months, s.Values, []float64{2, 2, 2}, []float64{4, 4, 4}, new(-1), RangeBandOptions{})
		},
		"error sign": func(c *ui.Context) {
			ErrorBars(c, "x", months, s.Values, []float64{1, -1, 1}, new(-1), ErrorBarsOptions{})
		},
		"pict icons": func(c *ui.Context) {
			Pictogram(c, "x", []string{"a"}, []float64{101}, "user", new(-1), PictogramOptions{})
		},
		"pict glyph": func(c *ui.Context) {
			Pictogram(c, "x", []string{"a"}, []float64{1}, "no-such-icon", new(-1), PictogramOptions{})
		},
		"marimekko 0": func(c *ui.Context) {
			Marimekko(c, "x", []string{"a"}, []Series{{Name: "s", Values: []float64{0}}}, new(-1), MarimekkoOptions{})
		},
		"population": func(c *ui.Context) {
			Population(c, "x", []string{"a"}, Series{"l", []float64{-1}}, Series{"r", []float64{1}}, new(-1), PopulationOptions{})
		},
		"stream window": func(c *ui.Context) { Streaming(c, "x", []float64{1}, StreamingOptions{Window: 1}) },
		"zoom range":    func(c *ui.Context) { ZoomChart(c, "x", months, []Series{s}, new(2), new(1), ZoomChartOptions{}) },
		"toggle flags": func(c *ui.Context) {
			ToggleLegend(c, "x", months, []Series{s}, &[]bool{}, new(-1), ToggleLegendOptions{})
		},
		"card title": func(c *ui.Context) { ChartCard(c, "", "1", 0, ChartCardOptions{}, func() {}) },
	} {
		panics(t, name, func() { render(fn) })
	}
}
