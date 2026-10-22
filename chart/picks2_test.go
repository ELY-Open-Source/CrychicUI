package chart

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestMarkPicks(t *testing.T) {
	web := Series{Name: "Web", Values: []float64{3, 5, 2}}
	app := Series{Name: "App", Values: []float64{1, 2, 4}}
	day := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	tree := []Node{{Name: "Root", Children: []Node{{Name: "Left"}, {Name: "Right"}}}}
	var week []Series
	for _, n := range []string{"W1", "W2"} {
		week = append(week, Series{Name: n, Values: []float64{1, 2, 3, 4, 5, 6, 7}})
	}
	cases := []struct {
		name string
		view func(c *ui.Context, at *int) ui.Element
		tip  string
		want int
	}{
		{"ErrorBars", func(c *ui.Context, at *int) ui.Element {
			return ErrorBars(c, "x", months, web.Values, []float64{1, 1, 1}, at, ErrorBarsOptions{})
		}, "Feb: 5 ± 1", 1},
		{"DotPlot", func(c *ui.Context, at *int) ui.Element {
			return DotPlot(c, "x", months, web.Values, at, DotPlotOptions{})
		}, "Feb: 5", 1},
		{"Lollipop", func(c *ui.Context, at *int) ui.Element {
			return Lollipop(c, "x", months, web.Values, at, LollipopOptions{})
		}, "Mar: 2", 2},
		{"Dumbbell", func(c *ui.Context, at *int) ui.Element {
			return Dumbbell(c, "x", months, web, app, at, DumbbellOptions{})
		}, "Feb: Web 5, App 2", 1},
		{"Slope", func(c *ui.Context, at *int) ui.Element { return Slope(c, "x", months, web, app, at, SlopeOptions{}) }, "Mar: Web 2, App 4", 2},
		{"Bump", func(c *ui.Context, at *int) ui.Element {
			return Bump(c, "x", months, []Series{web, app}, at, BumpOptions{})
		}, "App: #1 in Mar", 1},
		{"Parallel", func(c *ui.Context, at *int) ui.Element {
			return Parallel(c, "x", months, []Series{web, app}, at, ParallelOptions{})
		}, "App, Feb: 2", 1},
		{"Population", func(c *ui.Context, at *int) ui.Element {
			return Population(c, "x", months, web, app, at, PopulationOptions{})
		}, "Feb: Web 5, App 2", 1},
		{"Marimekko", func(c *ui.Context, at *int) ui.Element {
			return Marimekko(c, "x", months, []Series{web, app}, at, MarimekkoOptions{})
		}, "Mar, App: 4 (67%)", 2},
		{"WeekdayBars", func(c *ui.Context, at *int) ui.Element { return WeekdayBars(c, "x", week, at, WeekdayBarsOptions{}) }, "Wed, W2: 3", 2},
		{"Gantt", func(c *ui.Context, at *int) ui.Element {
			return Gantt(c, "x", []Task{{"Plan", day, day.AddDate(0, 0, 3)}, {"Build", day.AddDate(0, 0, 2), day.AddDate(0, 0, 9)}}, at, GanttOptions{})
		}, "Build: Jan 7 – Jan 14", 1},
		{"Combo", func(c *ui.Context, at *int) ui.Element { return Combo(c, "x", months, web, app, at, ComboOptions{}) }, "Feb: Web 5, App 2", 1},
		{"Pictogram", func(c *ui.Context, at *int) ui.Element {
			return Pictogram(c, "x", months, web.Values, "user", at, PictogramOptions{})
		}, "Feb: 5", 1},
		{"TreeDiagram", func(c *ui.Context, at *int) ui.Element { return TreeDiagram(c, "x", tree, at, TreeDiagramOptions{}) }, "Right", 2},
		{"Dendrogram", func(c *ui.Context, at *int) ui.Element {
			return Dendrogram(c, "x", []string{"a", "b", "c"}, []Join{{0, 1, 1}, {3, 2, 3}}, at, DendrogramOptions{})
		}, "c", 2},
		{"NetworkGraph", func(c *ui.Context, at *int) ui.Element {
			return NetworkGraph(c, "x", []string{"Ann", "Bo", "Cy"}, []Edge{{0, 1}, {1, 2}}, at, NetworkGraphOptions{})
		}, "Bo", 1},
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
			tt.Frame()
			if at != k.want || changed != 1 {
				t.Errorf("at %d, want %d; changed %d", at, k.want, changed)
			}
			cues(t, before, sound.Select)
			settles(t, tt, "pick")
		})
	}
}
