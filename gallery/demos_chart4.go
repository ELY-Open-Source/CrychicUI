package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/chart"
	"github.com/egoist/mygo/ui"
)

func demoTreeDiagram(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		team := []chart.Node{{Name: "Design"}, {Name: "Research"}}
		if s%2 == 1 {
			team = append(team, chart.Node{Name: "Brand"})
		}
		return chart.TreeDiagram(c, "Org", []chart.Node{{Name: "CEO", Children: []chart.Node{
			{Name: "Product", Children: team},
			{Name: "Eng", Children: []chart.Node{{Name: "Web"}, {Name: "App"}}},
		}}}, at, chart.TreeDiagramOptions{})
	})
}

func demoDendrogram(c *ui.Context, log func(string)) {
	chartBox(c, log, 520, func(s int, at *int) ui.Element {
		lift := float64(s % 3)
		return chart.Dendrogram(c, "Species", []string{"Cat", "Lion", "Wolf", "Dog", "Fox"},
			[]chart.Join{{A: 0, B: 1, Height: 2 + lift}, {A: 2, B: 3, Height: 1.5}, {A: 6, B: 4, Height: 3}, {A: 5, B: 7, Height: 6}},
			at, chart.DendrogramOptions{})
	})
}

func demoParallel(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		axes := []string{"Price", "Speed", "Range", "Seats"}
		return chart.Parallel(c, "Cars", axes, []chart.Series{
			{Name: "Swift", Values: rotated([]float64{20, 210, 300, 4}, s)},
			{Name: "Atlas", Values: rotated([]float64{45, 150, 520, 2}, s)},
			{Name: "Pico", Values: rotated([]float64{14, 180, 400, 5}, s)},
		}, at, chart.ParallelOptions{})
	})
}

func demoVenn(c *ui.Context, log func(string)) {
	chartBox(c, log, 320, func(s int, at *int) ui.Element {
		counts := []float64{40, 32, 25, 12, 9, 7, float64(3 + s%4)}
		return chart.Venn(c, "Readers", []string{"News", "Sport", "Tech"}, counts, at, chart.VennOptions{Size: 260})
	})
}

func demoChartRings(c *ui.Context, log func(string)) {
	chartBox(c, log, 320, func(s int, at *int) ui.Element {
		return chart.Rings(c, "Goals", []string{"Move", "Exercise", "Stand"}, rotated([]float64{0.82, 0.55, 0.9}, s), at, chart.RingsOptions{Size: 220})
	})
}

func demoStreaming(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520).Gap(12)
	n := intState(root, "n", 12)
	root.Children(func() {
		vs := make([]float64, *n)
		for i := range vs {
			vs[i] = float64(40 + (i*37)%25)
		}
		chart.Streaming(c, "CPU %", vs, chart.StreamingOptions{Window: 12})
		if button.Outline(c, "Add a reading").Clicked() {
			*n++
			log(fmt.Sprintf("readings %d", *n))
		}
	})
}

func demoZoomChart(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(560)
	from, to := intState(root, "from", 1), intState(root, "to", 4)
	root.Children(func() {
		if chart.ZoomChart(c, "Visits", chartMonths, chartTwo(0), from, to, chart.ZoomChartOptions{}).Changed() {
			log(fmt.Sprintf("zoomed %s–%s", chartMonths[*from], chartMonths[*to]))
		}
	})
}

func demoCrosshair(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Crosshair(c, "Visits", chartMonths, chartTwo(s), at, chart.CrosshairOptions{})
	})
}

func demoToggleLegend(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(560)
	at := intState(root, "at", -1)
	hidden := ui.Local(root, "hidden", func() []bool { return []bool{false, false} })
	root.Children(func() {
		was := *at
		if chart.ToggleLegend(c, "Visits", chartMonths, chartTwo(0), hidden, at, chart.ToggleLegendOptions{}).Changed() {
			if *at != was {
				log(fmt.Sprintf("chose %d", *at))
			} else {
				log(fmt.Sprintf("hidden %v", *hidden))
			}
		}
	})
}

func demoAnnotated(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		notes := []chart.Note{{At: 2, Text: "Launch"}, {At: 4, Series: 1, Text: "Promo"}}
		return chart.Annotated(c, "Visits", chartMonths, chartTwo(s), notes, at, chart.AnnotatedOptions{})
	})
}

func demoThreshold(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.Threshold(c, "Latency ms", chartMonths, rotated(chartWeb, s), 25, at, chart.ThresholdOptions{})
	})
}

func demoPopulation(c *ui.Context, log func(string)) {
	chartBox(c, log, 520, func(s int, at *int) ui.Element {
		ages := []string{"0–19", "20–39", "40–59", "60+"}
		men := chart.Series{Name: "Men", Values: rotated([]float64{21, 28, 26, 18}, s)}
		women := chart.Series{Name: "Women", Values: rotated([]float64{20, 27, 27, 22}, s)}
		return chart.Population(c, "Age", ages, men, women, at, chart.PopulationOptions{})
	})
}

func demoMarimekko(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		regions := []string{"North", "South", "East"}
		return chart.Marimekko(c, "Market", regions, []chart.Series{
			{Name: "Ours", Values: rotated([]float64{30, 12, 20}, s)},
			{Name: "Theirs", Values: rotated([]float64{50, 18, 15}, s)},
		}, at, chart.MarimekkoOptions{})
	})
}

func demoWeekdayBars(c *ui.Context, log func(string)) {
	chartBox(c, log, 560, func(s int, at *int) ui.Element {
		return chart.WeekdayBars(c, "Orders", []chart.Series{
			{Name: "Last week", Values: rotated([]float64{12, 15, 14, 18, 22, 9, 7}, s)},
			{Name: "This week", Values: rotated([]float64{14, 17, 13, 21, 25, 11, 8}, s)},
		}, at, chart.WeekdayBarsOptions{})
	})
}

func demoChartCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	s := intState(root, "s", 0)
	root.Children(func() {
		vs := rotated(chartWeb, *s)
		delta := 100 * (vs[len(vs)-1] - vs[len(vs)-2]) / vs[len(vs)-2]
		chart.ChartCard(c, "Visits", fmt.Sprintf("%.0fk", vs[len(vs)-1]), delta, chart.ChartCardOptions{}, func() {
			chart.Sparkline(c, "Visits trend", vs, chart.SparklineOptions{Width: 300})
		})
		if button.Outline(c, "Next day").Clicked() {
			*s++
			log(fmt.Sprintf("day %d", *s))
		}
	})
}
