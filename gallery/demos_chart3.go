package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/chart"
	"github.com/egoist/mygo/ui"
)

var chartEntries2 = []Entry{
	{741, "Spline", "chart", "A smooth curve through points.", demoSpline},
	{742, "Combo", "chart", "Bars and a line on two axes.", demoCombo},
	{743, "RangeBand", "chart", "A line in a shaded band.", demoRangeBand},
	{744, "ErrorBars", "chart", "Points with error whiskers.", demoErrorBars},
	{745, "DotPlot", "chart", "A dot per category.", demoDotPlot},
	{746, "Lollipop", "chart", "Stems with dots.", demoLollipop},
	{747, "Dumbbell", "chart", "Two dots joined per category.", demoDumbbell},
	{748, "Slope", "chart", "Lines between two columns.", demoSlope},
	{749, "Bump", "chart", "Ranks over time.", demoBump},
	{750, "Stream", "chart", "Areas stacked around a center.", demoStream},
	{751, "Waffle", "chart", "A grid of cells by share.", demoWaffle},
	{752, "Pictogram", "chart", "Rows of icons by count.", demoPictogram},
	{753, "YearHeat", "chart", "A year of day cells.", demoYearHeat},
	{754, "Gantt", "chart", "Task bars on a time axis.", demoGantt},
	{755, "NetworkGraph", "chart", "Nodes and edges; nodes drag.", demoNetworkGraph},
	{756, "TreeDiagram", "chart", "A hierarchy of linked boxes.", demoTreeDiagram},
	{757, "Dendrogram", "chart", "A cluster tree of joins.", demoDendrogram},
	{758, "Parallel", "chart", "Lines across parallel axes.", demoParallel},
	{759, "Venn", "chart", "Overlapping sets with counts.", demoVenn},
	{760, "Rings", "chart", "Concentric progress rings.", demoChartRings},
	{761, "Streaming", "chart", "A line that scrolls as values arrive.", demoStreaming},
	{762, "ZoomChart", "chart", "A line with a brush to zoom.", demoZoomChart},
	{763, "Crosshair", "chart", "A crosshair with value readouts.", demoCrosshair},
	{764, "ToggleLegend", "chart", "A legend that hides series.", demoToggleLegend},
	{765, "Annotated", "chart", "Notes on points.", demoAnnotated},
	{766, "Threshold", "chart", "Points past a limit turn red.", demoThreshold},
	{767, "Population", "chart", "Back-to-back bars by group.", demoPopulation},
	{768, "Marimekko", "chart", "Stacked bars of variable width.", demoMarimekko},
	{769, "WeekdayBars", "chart", "Bars by weekday per week.", demoWeekdayBars},
	{770, "ChartCard", "chart", "Title, value, change and a chart.", demoChartCard},
}

func demoSpline(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Spline(c, "Visits", chartMonths, chartTwo(s), at, chart.SplineOptions{})
	})
}

func demoCombo(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		margin := chart.Series{Name: "Margin %", Values: rotated([]float64{22, 25, 21, 30, 33, 29}, s)}
		return chart.Combo(c, "Sales and margin", chartMonths, chart.Series{Name: "Sales", Values: rotated(chartWeb, s)}, margin, at, chart.ComboOptions{})
	})
}

func demoRangeBand(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		mid := rotated(chartWeb, s)
		lo, hi := make([]float64, len(mid)), make([]float64, len(mid))
		for i, v := range mid {
			lo[i], hi[i] = v-3, v+4
		}
		return chart.RangeBand(c, "Forecast", chartMonths, mid, lo, hi, at, chart.RangeBandOptions{})
	})
}

func demoErrorBars(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.ErrorBars(c, "Load time", chartMonths, rotated(chartWeb, s), []float64{2, 3, 1.5, 4, 2, 3}, at, chart.ErrorBarsOptions{})
	})
}

func demoDotPlot(c *ui.Context, log func(string)) {
	chartBox(c, log, 520, func(s int, at *int) ui.Element {
		return chart.DotPlot(c, "Users", chartShares, rotated(chartUsers, s), at, chart.DotPlotOptions{})
	})
}

func demoLollipop(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Lollipop(c, "Signups", chartMonths, rotated(chartApp, s), at, chart.LollipopOptions{})
	})
}

func demoDumbbell(c *ui.Context, log func(string)) {
	chartBox(c, log, 520, func(s int, at *int) ui.Element {
		before := chart.Series{Name: "2025", Values: rotated(chartUsers, s)}
		after := chart.Series{Name: "2026", Values: rotated([]float64{48, 26, 25, 12}, s)}
		return chart.Dumbbell(c, "Users", chartShares, before, after, at, chart.DumbbellOptions{})
	})
}

func demoSlope(c *ui.Context, log func(string)) {
	chartBox(c, log, 480, func(s int, at *int) ui.Element {
		before := chart.Series{Name: "2025", Values: rotated(chartUsers, s)}
		after := chart.Series{Name: "2026", Values: rotated([]float64{48, 33, 21, 12}, s)}
		return chart.Slope(c, "Share by platform", chartShares, before, after, at, chart.SlopeOptions{})
	})
}

func demoBump(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		third := chart.Series{Name: "Mac", Values: rotated([]float64{14, 12, 20, 18, 22, 25}, s)}
		return chart.Bump(c, "Rank by visits", chartMonths, append(chartTwo(s), third), at, chart.BumpOptions{})
	})
}

func demoStream(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		third := chart.Series{Name: "Mac", Values: rotated([]float64{4, 8, 6, 9, 7, 10}, s)}
		return chart.Stream(c, "Listening", chartMonths, append(chartTwo(s), third), at, chart.StreamOptions{})
	})
}

func demoWaffle(c *ui.Context, log func(string)) {
	chartBox(c, log, 320, func(s int, at *int) ui.Element {
		return chart.Waffle(c, "Users", chartShares, rotated(chartUsers, s), at, chart.WaffleOptions{Size: 240})
	})
}

func demoPictogram(c *ui.Context, log func(string)) {
	chartBox(c, log, 520, func(s int, at *int) ui.Element {
		return chart.Pictogram(c, "Team size", chartShares, rotated([]float64{8, 5, 3, 2}, s), "user", at, chart.PictogramOptions{})
	})
}

func demoYearHeat(c *ui.Context, log func(string)) {
	chartBox(c, log, 600, func(s int, at *int) ui.Element {
		counts := make([]float64, 140)
		for i := range counts {
			counts[i] = float64((i*7 + s*13) % 11)
		}
		start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
		return chart.YearHeat(c, "Commits", start, counts, at, chart.YearHeatOptions{})
	})
}

func demoGantt(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		day := time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC)
		stretch := s % 3
		return chart.Gantt(c, "Launch plan", []chart.Task{
			{Name: "Design", Start: day, End: day.AddDate(0, 0, 5)},
			{Name: "Build", Start: day.AddDate(0, 0, 4), End: day.AddDate(0, 0, 14+stretch)},
			{Name: "Test", Start: day.AddDate(0, 0, 12), End: day.AddDate(0, 0, 18+stretch)},
			{Name: "Ship", Start: day.AddDate(0, 0, 18+stretch), End: day.AddDate(0, 0, 20+stretch)},
		}, at, chart.GanttOptions{})
	})
}

func demoNetworkGraph(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	at := intState(root, "at", -1)
	root.Children(func() {
		names := []string{"Ada", "Bo", "Cy", "Di", "Ed", "Flo"}
		edges := []chart.Edge{{From: 0, To: 1}, {From: 1, To: 2}, {From: 2, To: 0}, {From: 2, To: 3}, {From: 3, To: 4}, {From: 4, To: 5}}
		was := *at
		if chart.NetworkGraph(c, "Friends", names, edges, at, chart.NetworkGraphOptions{}).Changed() {
			if *at != was {
				log(fmt.Sprintf("chose %d", *at))
			} else {
				log("moved a node")
			}
		}
	})
}
