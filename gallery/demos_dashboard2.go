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
	points := []dashboard.Point{{Label: "Jan", Value: 2}, {Label: "Feb", Value: 3}}
	dashboard.Churn(c, 3, points, dashboard.ChurnOptions{})
}

func demoMRR(c *ui.Context, log func(string)) {
	dashboard.MRR(c, 1000, 100, 50, 25, "USD", dashboard.MRROptions{})
}

func demoRangeBar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
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
	root := ui.Column(c)
	root.Children(func() {
		if dashboard.WidgetFrame(c, "Revenue", dashboard.WidgetFrameOptions{}, func() { ui.Text(c, "Body") }).Changed() {
			log("refreshed")
		}
	})
}

func demoWidgetGrid(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	titles := listState(root, "titles", "Alpha", "Beta")
	root.Children(func() {
		if dashboard.WidgetGrid(c, titles).Changed() {
			log(fmt.Sprint(*titles))
		}
	})
}

func demoReportExport(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	format := textState(root, "format", "PDF")
	exported := boolState(root, "exported", false)
	root.Children(func() {
		root, export := dashboard.ReportExport(c, format, []string{"PDF", "CSV"}, exported, dashboard.ReportExportOptions{})
		if root.Changed() {
			log("format " + *format)
		}
		if export.Clicked() {
			log("exported " + *format)
			*exported = false
		}
	})
}

func demoInsights(c *ui.Context, log func(string)) {
	cards := []dashboard.Insight{{Title: "Growth", Body: "Revenue rose"}}
	dashboard.Insights(c, cards)
}

func demoAnomaly(c *ui.Context, log func(string)) {
	points := []dashboard.Point{{Label: "A", Value: 1}, {Label: "B", Value: 3}}
	dashboard.Anomaly(c, points, []int{1})
}

func demoScorecard(c *ui.Context, log func(string)) {
	goals := []dashboard.ScoreGoal{{Label: "Revenue", Current: 100, Target: 100}}
	dashboard.Scorecard(c, goals, dashboard.ScorecardOptions{})
}

func demoHealthRing(c *ui.Context, log func(string)) {
	dashboard.HealthRing(c, 90, dashboard.HealthRingOptions{})
}

func demoServerLoad(c *ui.Context, log func(string)) {
	dashboard.ServerLoad(c, 60, 70, dashboard.ServerLoadOptions{})
}

func demoQueue(c *ui.Context, log func(string)) {
	dashboard.Queue(c, 3, 2, dashboard.QueueOptions{})
}

func demoErrorRate(c *ui.Context, log func(string)) {
	dashboard.ErrorRate(c, 3, 2, dashboard.ErrorRateOptions{})
}

func demoLatency(c *ui.Context, log func(string)) {
	dashboard.Latency(c, 10, 20, 30, dashboard.LatencyOptions{})
}

func demoSignups(c *ui.Context, log func(string)) {
	days := []dashboard.SignupDay{{Label: "Mon", Count: 10}}
	dashboard.Signups(c, days)
}

func demoNPSWidget(c *ui.Context, log func(string)) {
	dashboard.NPSWidget(c, 80, 0, 20, dashboard.NPSWidgetOptions{})
}

func demoWelcomeBanner(c *ui.Context, log func(string)) {
	dashboard.WelcomeBanner(c, "Good morning", "3 tasks today")
}
