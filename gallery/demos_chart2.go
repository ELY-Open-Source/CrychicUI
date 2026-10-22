package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/chart"
	"github.com/egoist/mygo/ui"
)

func demoGaugeChart(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400).Gap(12)
	v := intState(root, "v", 42)
	root.Children(func() {
		chart.GaugeChart(c, "CPU load", float64(*v), 0, 100, chart.GaugeChartOptions{Bands: []float64{60, 85}})
		if button.Outline(c, "Add load").Clicked() {
			*v = (*v + 17) % 100
			log(fmt.Sprintf("load %d", *v))
		}
	})
}

func demoHistogram(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		var vs []float64
		for i := range 60 {
			vs = append(vs, float64((i*37+s*11)%50)+float64(i%7))
		}
		return chart.Histogram(c, "Response times", vs, at, chart.HistogramOptions{Bins: 8})
	})
}

func demoBoxPlot(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.BoxPlot(c, "Scores", chartGroups(s), at, chart.BoxPlotOptions{})
	})
}

func demoViolin(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Violin(c, "Scores", chartGroups(s), at, chart.ViolinOptions{})
	})
}

// chartGroups are three samples that move with shift.
func chartGroups(shift int) []chart.Group {
	return []chart.Group{{Name: "A", Values: rotated(chartWeb, shift)}, {Name: "B", Values: chartApp},
		{Name: "C", Values: []float64{9, 12, 14, 15, 16, 22, float64(10 + shift%5)}}}
}

func demoCandles(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Candles(c, "ACME", chartCandles(s), at, chart.CandlesOptions{})
	})
}

func demoOHLC(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.OHLC(c, "ACME", chartCandles(s), at, chart.OHLCOptions{})
	})
}

// chartCandles are six days of prices that move with shift.
func chartCandles(shift int) []chart.Candle {
	var out []chart.Candle
	for i, m := range chartMonths {
		o, cl := rotated(chartWeb, shift)[i], chartWeb[(i+1)%6]
		out = append(out, chart.Candle{Label: m, Open: o, High: max(o, cl) + 3, Low: min(o, cl) - 2, Close: cl})
	}
	return out
}

func demoHeatmap(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		days := []string{"Mon", "Tue", "Wed"}
		var vs [][]float64
		for r := range days {
			vs = append(vs, rotated(chartWeb, s+r))
		}
		return chart.Heatmap(c, "Visits by hour", days, chartMonths, vs, at, chart.HeatmapOptions{})
	})
}

func demoTreemap(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Treemap(c, "Disk use", []chart.Node{
			{Name: "Photos", Children: []chart.Node{{Name: "2024", Value: 40 + float64(s%3)*10}, {Name: "2025", Value: 25}}},
			{Name: "Apps", Children: []chart.Node{{Name: "Games", Value: 20}, {Name: "Tools", Value: 8}}},
			{Name: "Music", Value: 15},
		}, at, chart.TreemapOptions{})
	})
}

func demoSunburst(c *ui.Context, log func(string)) {
	chartBox(c, log, 400, func(s int, at *int) ui.Element {
		return chart.Sunburst(c, "Traffic", []chart.Node{
			{Name: "Search", Children: []chart.Node{{Name: "Google", Value: 30 + float64(s%3)*5}, {Name: "Bing", Value: 8}}},
			{Name: "Social", Children: []chart.Node{{Name: "X", Value: 12}, {Name: "Reddit", Value: 9}}},
			{Name: "Direct", Value: 20},
		}, at, chart.SunburstOptions{})
	})
}

func demoFunnel(c *ui.Context, log func(string)) {
	chartBox(c, log, 480, func(s int, at *int) ui.Element {
		return chart.Funnel(c, "Checkout", []string{"Visit", "Cart", "Pay", "Done"}, []float64{1000, 420 + float64(s%3)*40, 180, 120}, at, chart.FunnelOptions{})
	})
}

func demoSankey(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Sankey(c, "Budget", []string{"Salary", "Rent", "Food", "Savings", "Fun"},
			[]chart.Flow{{From: 0, To: 1, Value: 40}, {From: 0, To: 2, Value: 20 + float64(s%3)*5}, {From: 0, To: 3, Value: 25}, {From: 3, To: 4, Value: 10}},
			at, chart.SankeyOptions{})
	})
}

func demoRadar(c *ui.Context, log func(string)) {
	chartBox(c, log, 420, func(s int, at *int) ui.Element {
		return chart.Radar(c, "Skills", []string{"Speed", "Power", "Range", "Aim", "Guard"},
			[]chart.Series{{Name: "Ada", Values: rotated([]float64{8, 5, 7, 9, 4}, s)}, {Name: "Alan", Values: []float64{5, 8, 6, 4, 7}}}, at, chart.RadarOptions{})
	})
}

func demoWaterfall(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Waterfall(c, "Cash", []string{"Start", "Sales", "Costs", "Tax"}, []float64{40, 25 + float64(s%3)*5, -18, -7}, at, chart.WaterfallOptions{})
	})
}

func demoBullet(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520).Gap(12)
	v := intState(root, "v", 230)
	root.Children(func() {
		chart.Bullet(c, "Revenue", float64(*v), 260, chart.BulletOptions{Bands: []float64{150, 225, 300}})
		if button.Outline(c, "Book a sale").Clicked() {
			*v = (*v + 20) % 300
			log(fmt.Sprintf("revenue %d", *v))
		}
	})
}

func demoSparkline(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(12).AlignItems(ui.Center)
	s := intState(root, "s", 0)
	root.Children(func() {
		chart.Sparkline(c, "Visits this week", rotated(chartWeb, *s), chart.SparklineOptions{})
		if button.Outline(c, "Next week").Clicked() {
			*s++
			log(fmt.Sprintf("week %d", *s))
		}
	})
}

func demoSparkBars(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(12).AlignItems(ui.Center)
	s := intState(root, "s", 0)
	root.Children(func() {
		chart.SparkBars(c, "Orders this week", rotated(chartApp, *s), chart.SparkBarsOptions{})
		if button.Outline(c, "Next week").Clicked() {
			*s++
			log(fmt.Sprintf("week %d", *s))
		}
	})
}

func demoStepLine(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.StepLine(c, "Plan seats", chartMonths, chartTwo(s), at, chart.StepLineOptions{})
	})
}
