package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/chart"
	"github.com/egoist/mygo/ui"
)

var chartEntries = []Entry{
	{711, "Line", "chart", "Line series over axes.", demoLine},
	{712, "Area", "chart", "Filled line series.", demoArea},
	{713, "StackedArea", "chart", "Areas stacked.", demoStackedArea},
	{714, "Bars", "chart", "Vertical bars.", demoBars},
	{715, "HBars", "chart", "Horizontal bars.", demoHBars},
	{716, "GroupedBars", "chart", "Bars side by side.", demoGroupedBars},
	{717, "StackedBars", "chart", "Bars stacked.", demoStackedBars},
	{718, "Pie", "chart", "Slices; the chosen one lifts.", demoPie},
	{719, "Donut", "chart", "A ring with its total.", demoDonut},
	{720, "RadialBars", "chart", "Concentric arcs.", demoRadialBars},
	{721, "GaugeChart", "chart", "Bands and a needle.", demoGaugeChart},
	{722, "Scatter", "chart", "Points on two axes.", demoScatter},
	{723, "Bubble", "chart", "Points sized by a third value.", demoBubble},
	{724, "Histogram", "chart", "Binned counts.", demoHistogram},
	{725, "BoxPlot", "chart", "Quartiles and whiskers.", demoBoxPlot},
	{726, "Violin", "chart", "Mirrored density.", demoViolin},
	{727, "Candles", "chart", "Open, high, low, close.", demoCandles},
	{728, "OHLC", "chart", "Ticks on range lines.", demoOHLC},
	{729, "Heatmap", "chart", "Cells shaded by value.", demoHeatmap},
	{730, "Treemap", "chart", "Rectangles sized by value.", demoTreemap},
	{731, "Sunburst", "chart", "Rings of nested shares.", demoSunburst},
	{732, "Funnel", "chart", "Narrowing stages with rates.", demoFunnel},
	{733, "Sankey", "chart", "Flows between nodes.", demoSankey},
	{734, "Radar", "chart", "Values on spokes.", demoRadar},
	{735, "PolarArea", "chart", "Wedges sized by value.", demoPolarArea},
	{736, "Waterfall", "chart", "Steps to a total.", demoWaterfall},
	{737, "Bullet", "chart", "A bar against a target.", demoBullet},
	{738, "Sparkline", "chart", "A tiny line.", demoSparkline},
	{739, "SparkBars", "chart", "Tiny bars.", demoSparkBars},
	{740, "StepLine", "chart", "A line in steps.", demoStepLine},
}

var (
	chartMonths = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	chartWeb    = []float64{12, 18, 15, 24, 28, 31}
	chartApp    = []float64{6, 9, 14, 13, 19, 22}
	chartShares = []string{"iOS", "Android", "Web", "Desktop"}
	chartUsers  = []float64{42, 31, 18, 9}
)

// chartBox holds a chart and a shuffle button.
func chartBox(c *ui.Context, log func(string), w float32, build func(shift int, at *int) ui.Element) {
	root := ui.Column(c).Width(w).Gap(12)
	at := intState(root, "at", -1)
	shift := intState(root, "shift", 0)
	root.Children(func() {
		if build(*shift, at).Changed() {
			log(fmt.Sprintf("chose %d", *at))
		}
		if button.Outline(c, "Shuffle data").Clicked() {
			*shift, *at = *shift+1, -1
			log(fmt.Sprintf("data %d", *shift))
		}
	})
}

// rotated turns vs by n places.
func rotated(vs []float64, n int) []float64 {
	out := make([]float64, len(vs))
	for i := range vs {
		out[i] = vs[(i+n)%len(vs)]
	}
	return out
}

func chartTwo(shift int) []chart.Series {
	return []chart.Series{{Name: "Web", Values: rotated(chartWeb, shift)}, {Name: "App", Values: rotated(chartApp, shift)}}
}

func demoLine(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Line(c, "Visits", chartMonths, chartTwo(s), at, chart.LineOptions{})
	})
}

func demoArea(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Area(c, "Visits", chartMonths, chartTwo(s), at, chart.AreaOptions{})
	})
}

func demoStackedArea(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.StackedArea(c, "Visits", chartMonths, chartTwo(s), at, chart.StackedAreaOptions{})
	})
}

func demoBars(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Bars(c, "Sales", chartMonths, rotated(chartWeb, s), at, chart.BarsOptions{})
	})
}

func demoHBars(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.HBars(c, "Users", chartShares, rotated(chartUsers, s), at, chart.HBarsOptions{})
	})
}

func demoGroupedBars(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.GroupedBars(c, "Visits", chartMonths, chartTwo(s), at, chart.GroupedBarsOptions{})
	})
}

func demoStackedBars(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.StackedBars(c, "Visits", chartMonths, chartTwo(s), at, chart.StackedBarsOptions{})
	})
}

func demoPie(c *ui.Context, log func(string)) {
	chartBox(c, log, 400, func(s int, at *int) ui.Element {
		return chart.Pie(c, "Users by platform", chartShares, rotated(chartUsers, s), at, chart.PieOptions{})
	})
}

func demoDonut(c *ui.Context, log func(string)) {
	chartBox(c, log, 400, func(s int, at *int) ui.Element {
		return chart.Donut(c, "Users by platform", chartShares, rotated(chartUsers, s), at, chart.DonutOptions{})
	})
}

func demoRadialBars(c *ui.Context, log func(string)) {
	chartBox(c, log, 400, func(s int, at *int) ui.Element {
		return chart.RadialBars(c, "Goals", chartShares, rotated(chartUsers, s), at, chart.RadialBarsOptions{})
	})
}

func demoPolarArea(c *ui.Context, log func(string)) {
	chartBox(c, log, 400, func(s int, at *int) ui.Element {
		return chart.PolarArea(c, "Users by platform", chartShares, rotated(chartUsers, s), at, chart.PolarAreaOptions{})
	})
}

func demoScatter(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Scatter(c, "Height and weight", chartPoints(s), at, chart.ScatterOptions{})
	})
}

func demoBubble(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Bubble(c, "Cities", chartPoints(s), at, chart.BubbleOptions{})
	})
}

// chartPoints are six named points that move with shift.
func chartPoints(shift int) []chart.Point {
	var pts []chart.Point
	for i, m := range chartMonths {
		pts = append(pts, chart.Point{Name: m, X: chartWeb[i], Y: rotated(chartApp, shift)[i], Size: chartUsers[i%4]})
	}
	return pts
}
