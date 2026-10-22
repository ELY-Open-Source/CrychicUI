package chart

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestMarksPick(t *testing.T) {
	two := []Series{{Name: "Web", Values: []float64{3, 5, 2}}, {Name: "App", Values: []float64{1, 2, 4}}}
	candles := []Candle{{"Mon", 10, 14, 9, 13}, {"Tue", 13, 15, 11, 12}, {"Wed", 12, 16, 12, 15}}
	cases := []struct {
		name string
		view func(c *ui.Context, at *int) ui.Element
		tip  string
		want int
	}{
		{"HBars", func(c *ui.Context, at *int) ui.Element {
			return HBars(c, "x", months, []float64{3, 6, 1}, at, HBarsOptions{})
		}, "Mar: 1", 2},
		{"GroupedBars", func(c *ui.Context, at *int) ui.Element {
			return GroupedBars(c, "x", months, two, at, GroupedBarsOptions{})
		}, "Feb, App: 2", 1},
		{"StackedBars", func(c *ui.Context, at *int) ui.Element {
			return StackedBars(c, "x", months, two, at, StackedBarsOptions{})
		}, "Mar, Web: 2", 2},
		{"Histogram", func(c *ui.Context, at *int) ui.Element {
			return Histogram(c, "x", []float64{0, 1, 1, 2, 3, 3, 3, 4}, at, HistogramOptions{Bins: 4})
		}, "3–4: 4", 3},
		{"Waterfall", func(c *ui.Context, at *int) ui.Element {
			return Waterfall(c, "x", []string{"Start", "Loss"}, []float64{10, -4}, at, WaterfallOptions{})
		}, "Total: +6", 2},
		{"Heatmap", func(c *ui.Context, at *int) ui.Element {
			return Heatmap(c, "x", []string{"A", "B"}, []string{"X", "Y"}, [][]float64{{1, 2}, {3, 4}}, at, HeatmapOptions{})
		}, "B, X: 3", 2},
		{"Treemap", func(c *ui.Context, at *int) ui.Element {
			return Treemap(c, "x", []Node{{Name: "Big", Value: 6}, {Name: "Mid", Value: 3}, {Name: "Low", Value: 1}}, at, TreemapOptions{})
		}, "Mid: 3", 1},
		{"Scatter", func(c *ui.Context, at *int) ui.Element {
			return Scatter(c, "x", []Point{{Name: "a", X: 1, Y: 2}, {Name: "b", X: 3, Y: 1}}, at, ScatterOptions{})
		}, "b: 3, 1", 1},
		{"Bubble", func(c *ui.Context, at *int) ui.Element {
			return Bubble(c, "x", []Point{{Name: "a", X: 1, Y: 2, Size: 9}, {Name: "b", X: 3, Y: 1, Size: 1}}, at, BubbleOptions{})
		}, "a: 1, 2", 0},
		{"BoxPlot", func(c *ui.Context, at *int) ui.Element {
			return BoxPlot(c, "x", []Group{{"A", []float64{1, 2, 3, 4, 5}}, {"B", []float64{2, 4, 6}}}, at, BoxPlotOptions{})
		}, "B: median 4, 3–5", 1},
		{"Candles", func(c *ui.Context, at *int) ui.Element { return Candles(c, "x", candles, at, CandlesOptions{}) },
			"Tue: open 13, high 15, low 11, close 12", 1},
		{"Funnel", func(c *ui.Context, at *int) ui.Element {
			return Funnel(c, "x", []string{"Visit", "Cart", "Buy"}, []float64{100, 40, 10}, at, FunnelOptions{})
		}, "Cart: 40", 1},
		{"Sankey", func(c *ui.Context, at *int) ui.Element {
			return Sankey(c, "x", []string{"Ads", "Site", "Buy"}, []Flow{{0, 1, 5}, {1, 2, 3}}, at, SankeyOptions{})
		}, "Site: 5", 1},
	}
	for _, k := range cases {
		t.Run(k.name, func(t *testing.T) {
			at, changed := -1, 0
			tt := ui.NewTester(func(c *ui.Context) {
				if k.view(c, &at).Changed() {
					changed++
				}
			}, 600, 400)
			if !tt.FramePending() {
				t.Error("nothing grows in")
			}
			settles(t, tt, "chart")
			before := baseline()
			click(t, tt, k.tip)
			if at != k.want || changed != 1 {
				t.Errorf("at %d, want %d; changed %d", at, k.want, changed)
			}
			cues(t, before, sound.Select)
			settles(t, tt, "pick")
		})
	}
}
