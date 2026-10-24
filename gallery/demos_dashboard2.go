package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/dashboard"
	"github.com/egoist/mygo/ui"
)

var dashboardEntries2 = []Entry{
	{ID: 1074, Name: "Churn", Category: "dashboard", Summary: "Churn rate with trend.", Demo: demoChurn},
	{ID: 1075, Name: "MRR", Category: "dashboard", Summary: "MRR with new, expansion, churn.", Demo: demoMRR},
	{ID: 1076, Name: "RangeBar", Category: "dashboard", Summary: "Date range and compare toggle; cue Select.", Demo: demoRangeBar},
	{ID: 1077, Name: "WidgetFrame", Category: "dashboard", Summary: "Titled widget with menu and refresh; cue Tap.", Demo: demoWidgetFrame},
	{ID: 1078, Name: "WidgetGrid", Category: "dashboard", Summary: "Widgets reordered by drag; cue Drop.", Demo: demoWidgetGrid},
	{ID: 1079, Name: "ReportExport", Category: "dashboard", Summary: "Format choices and export; cue Success.", Demo: demoReportExport},
	{ID: 1080, Name: "Insights", Category: "dashboard", Summary: "Generated insight cards.", Demo: demoInsights},
	{ID: 1081, Name: "Anomaly", Category: "dashboard", Summary: "Chart with flagged points pulsing.", Demo: demoAnomaly},
	{ID: 1082, Name: "Scorecard", Category: "dashboard", Summary: "Goals with status per row.", Demo: demoScorecard},
	{ID: 1083, Name: "HealthRing", Category: "dashboard", Summary: "Overall health ring.", Demo: demoHealthRing},
	{ID: 1084, Name: "ServerLoad", Category: "dashboard", Summary: "CPU and memory bars that move.", Demo: demoServerLoad},
	{ID: 1085, Name: "Queue", Category: "dashboard", Summary: "Jobs waiting and running.", Demo: demoQueue},
	{ID: 1086, Name: "ErrorRate", Category: "dashboard", Summary: "Errors per minute with limit line.", Demo: demoErrorRate},
	{ID: 1087, Name: "Latency", Category: "dashboard", Summary: "p50, p95, p99 values.", Demo: demoLatency},
	{ID: 1088, Name: "Signups", Category: "dashboard", Summary: "Signups by day bars.", Demo: demoSignups},
	{ID: 1089, Name: "NPSWidget", Category: "dashboard", Summary: "Promoters, passives, detractors.", Demo: demoNPSWidget},
	{ID: 1090, Name: "WelcomeBanner", Category: "dashboard", Summary: "Greeting with the day's summary.", Demo: demoWelcomeBanner},
}

func demoChurn(c *ui.Context, log func(string)) {
	points := []dashboard.Point{{Label: "May", Value: 3.4}, {Label: "Jun", Value: 3.1}, {Label: "Jul", Value: 3.3},
		{Label: "Aug", Value: 2.9}, {Label: "Sep", Value: 2.6}, {Label: "Oct", Value: 2.3}}
	ui.Column(c).Width(380).Children(func() {
		dashboard.Churn(c, 2.3, points, dashboard.ChurnOptions{Title: "Monthly logo churn", Note: "41 of 1,782 accounts"})
	})
}

func demoMRR(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		dashboard.MRR(c, 128400, 9820, 4310, 3630, "USD", dashboard.MRROptions{Title: "MRR · October"})
	})
}

func demoRangeBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(500)
	start := ui.Local(root, "start", func() time.Time { return time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC) })
	end := ui.Local(root, "end", func() time.Time { return time.Date(2026, 10, 8, 0, 0, 0, 0, time.UTC) })
	compare := boolState(root, "compare", false)
	root.Children(func() {
		if dashboard.RangeBar(c, start, end, compare, dashboard.RangeBarOptions{}).Changed() {
			log(fmt.Sprintf("%v %v", start.Day(), *compare))
		}
	})
}

func demoWidgetFrame(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	root.Children(func() {
		frame := dashboard.WidgetFrame(c, "Signups this week", dashboard.WidgetFrameOptions{Icon: "user-plus", Note: "Updated 2m ago"}, func() {
			ui.Text(c, "1,246").Bold().FontSize(28)
			ui.Text(c, "Up 9% on last week, led by Thursday's launch post.").FontSize(13).Wrap()
		})
		if frame.Changed() {
			log("refreshed")
		}
	})
}

func demoWidgetGrid(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	titles := listState(root, "titles", "Revenue", "Active users", "Churn", "Top pages")
	root.Children(func() {
		if dashboard.WidgetGrid(c, titles, dashboard.WidgetGridOptions{}).Changed() {
			log(fmt.Sprint(*titles))
		}
	})
}

func demoReportExport(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	format := textState(root, "format", "PDF")
	exported := boolState(root, "exported", false)
	root.Children(func() {
		root, export := dashboard.ReportExport(c, format, []string{"PDF", "CSV", "XLSX", "JSON"}, exported,
			dashboard.ReportExportOptions{Title: "Export Q3 board report"})
		switch {
		case export.Clicked():
			log("exported " + *format)
		case root.Changed():
			log("format " + *format)
		}
	})
}

func demoInsights(c *ui.Context, log func(string)) {
	cards := []dashboard.Insight{
		{Title: "Trial conversion is up", Body: "Teams that invite a colleague in week one convert at twice the rate.", Icon: "trending-up", Metric: "+18%"},
		{Title: "Checkout errors cluster in Europe", Body: "71% of failed payments this week came from eu-west-1.", Icon: "triangle-alert", Metric: "-4.2%"},
		{Title: "Docs drive signups", Body: "Visitors from /docs/getting-started sign up 3.1× more than the site average.", Icon: "book-open", Metric: "3.1×"},
	}
	ui.Column(c).Width(440).Children(func() { dashboard.Insights(c, cards, dashboard.InsightsOptions{}) })
}

func demoAnomaly(c *ui.Context, log func(string)) {
	values := []float64{410, 432, 418, 445, 1180, 452, 437, 461, 448, 120, 470, 466}
	points := make([]dashboard.Point, len(values))
	for i, v := range values {
		points[i] = dashboard.Point{Label: fmt.Sprintf("Oct %d · 14:%02d", 24, i*5), Value: v}
	}
	ui.Column(c).Width(440).Children(func() {
		dashboard.Anomaly(c, points, []int{4, 9}, dashboard.AnomalyOptions{Title: "Requests per second"})
	})
}

func demoScorecard(c *ui.Context, log func(string)) {
	goals := []dashboard.ScoreGoal{{Label: "New ARR", Current: 1.42e6, Target: 1.3e6}, {Label: "Net revenue retention", Current: 112, Target: 110},
		{Label: "Enterprise logos", Current: 14, Target: 18}, {Label: "Support CSAT", Current: 91, Target: 95}}
	ui.Column(c).Width(440).Children(func() { dashboard.Scorecard(c, goals, dashboard.ScorecardOptions{Title: "Q3 scorecard"}) })
}

func demoHealthRing(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Children(func() {
		dashboard.HealthRing(c, 86, dashboard.HealthRingOptions{Label: "Platform health", Note: "2 of 14 checks degraded: webhooks and email."})
	})
}

func demoServerLoad(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	start := ui.Local(root, "start", func() time.Time { return c.Now() })
	tick := int(c.Now().Sub(*start) / time.Second)
	c.After(time.Second - c.Now().Sub(*start)%time.Second)
	cpu := []float64{42, 58, 71, 64, 49, 83, 92, 76}
	mem := []float64{61, 62, 64, 66, 65, 68, 71, 70}
	root.Children(func() {
		dashboard.ServerLoad(c, cpu[tick%len(cpu)], mem[tick%len(mem)], dashboard.ServerLoadOptions{Host: "api-prod-3"})
	})
}

func demoQueue(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Children(func() {
		dashboard.Queue(c, 1284, 36, dashboard.QueueOptions{Title: "Email delivery queue"})
	})
}

func demoErrorRate(c *ui.Context, log func(string)) {
	history := []float64{0.8, 1.1, 0.9, 1.4, 2.6, 3.4, 2.2, 1.6, 1.2, 1.0, 1.3, 2.4}
	ui.Column(c).Width(420).Children(func() {
		dashboard.ErrorRate(c, 2.4, 2, dashboard.ErrorRateOptions{Title: "Checkout errors", History: history})
	})
}

func demoLatency(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		dashboard.Latency(c, 84, 212, 1480, dashboard.LatencyOptions{Title: "API latency · last hour"})
	})
}

func demoSignups(c *ui.Context, log func(string)) {
	days := []dashboard.SignupDay{{Label: "Mon", Count: 184}, {Label: "Tue", Count: 212}, {Label: "Wed", Count: 198},
		{Label: "Thu", Count: 241}, {Label: "Fri", Count: 226}, {Label: "Sat", Count: 97}, {Label: "Sun", Count: 88}}
	ui.Column(c).Width(400).Children(func() { dashboard.Signups(c, days, dashboard.SignupsOptions{Title: "Signups this week"}) })
}

func demoNPSWidget(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() { dashboard.NPSWidget(c, 412, 168, 96, dashboard.NPSWidgetOptions{}) })
}

func demoWelcomeBanner(c *ui.Context, log func(string)) {
	stats := []dashboard.BannerStat{{Label: "New signups", Value: "241"}, {Label: "Open tickets", Value: "12"}, {Label: "MRR", Value: "$128k"}}
	ui.Column(c).Width(500).Children(func() {
		dashboard.WelcomeBanner(c, "Good morning, Maya", "Three deploys went out overnight and every check is green.",
			dashboard.WelcomeBannerOptions{Stats: stats})
	})
}
