package gallery

import (
	"github.com/ELY-Open-Source/CrychicUI/dashboard"
	"github.com/egoist/mygo/ui"
)

var dashboardEntries = []Entry{
	{ID: 1056, Name: "KPIGrid", Category: "dashboard", Summary: "Grid of KPI tiles; values roll.", Demo: demoKPIGrid},
	{ID: 1057, Name: "MetricTile", Category: "dashboard", Summary: "Value, delta and sparkline.", Demo: demoMetricTile},
	{ID: 1058, Name: "GoalProgress", Category: "dashboard", Summary: "Goal bars by quarter.", Demo: demoGoalProgress},
	{ID: 1059, Name: "Revenue", Category: "dashboard", Summary: "Revenue card with chart and period select; cue Select.", Demo: demoRevenue},
	{ID: 1060, Name: "Users", Category: "dashboard", Summary: "Active users with chart.", Demo: demoUsers},
	{ID: 1061, Name: "Conversion", Category: "dashboard", Summary: "Funnel with rates.", Demo: demoConversion},
	{ID: 1062, Name: "Traffic", Category: "dashboard", Summary: "Sources with shares.", Demo: demoTraffic},
	{ID: 1063, Name: "TopPages", Category: "dashboard", Summary: "Pages with views bars.", Demo: demoTopPages},
	{ID: 1064, Name: "Realtime", Category: "dashboard", Summary: "Live visitor count that ticks.", Demo: demoRealtime},
	{ID: 1065, Name: "ActivityFeed", Category: "dashboard", Summary: "Recent events stream; slide in.", Demo: demoActivityFeed},
	{ID: 1066, Name: "Leaders", Category: "dashboard", Summary: "Top performers with ranks.", Demo: demoLeaders},
	{ID: 1067, Name: "Uptime", Category: "dashboard", Summary: "Day bars colored by status.", Demo: demoUptime},
	{ID: 1068, Name: "StatusGrid", Category: "dashboard", Summary: "Services as colored cells.", Demo: demoStatusGrid},
	{ID: 1069, Name: "Incident", Category: "dashboard", Summary: "Incident card with timeline.", Demo: demoIncident},
	{ID: 1070, Name: "SLA", Category: "dashboard", Summary: "Met percentage against target.", Demo: demoSLA},
	{ID: 1071, Name: "Regions", Category: "dashboard", Summary: "Regions with values bars.", Demo: demoRegions},
	{ID: 1072, Name: "Devices", Category: "dashboard", Summary: "Desktop, mobile, tablet shares.", Demo: demoDevices},
	{ID: 1073, Name: "Cohort", Category: "dashboard", Summary: "Retention grid shaded.", Demo: demoCohort},
}

func demoKPIGrid(c *ui.Context, log func(string)) {
	tiles := []dashboard.KPI{{Label: "Revenue", Value: "$12k", Delta: "+8%", Points: []float64{1, 2}}}
	dashboard.KPIGrid(c, tiles)
}

func demoMetricTile(c *ui.Context, log func(string)) {
	dashboard.MetricTile(c, "Users", "1,200", "+8%", []float64{1, 2})
}

func demoGoalProgress(c *ui.Context, log func(string)) {
	goals := []dashboard.Goal{{Quarter: "Q1", Current: 80, Target: 100}}
	dashboard.GoalProgress(c, goals)
}

func demoRevenue(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	period := textState(root, "period", "Month")
	points := []dashboard.Point{{Label: "Jan", Value: 10}, {Label: "Feb", Value: 20}}
	root.Children(func() {
		if dashboard.Revenue(c, points, period, []string{"Month", "Year"}, dashboard.RevenueOptions{}).Changed() {
			log("period " + *period)
		}
	})
}

func demoUsers(c *ui.Context, log func(string)) {
	points := []dashboard.Point{{Label: "Jan", Value: 10}, {Label: "Feb", Value: 20}}
	dashboard.Users(c, points, dashboard.UsersOptions{})
}

func demoConversion(c *ui.Context, log func(string)) {
	stages := []dashboard.FunnelStage{{Label: "Visits", Count: 100}, {Label: "Signups", Count: 20}}
	dashboard.Conversion(c, stages)
}

func demoTraffic(c *ui.Context, log func(string)) {
	sources := []dashboard.Share{{Name: "Search", Share: 0.6}}
	dashboard.Traffic(c, sources)
}

func demoTopPages(c *ui.Context, log func(string)) {
	pages := []dashboard.PageView{{Path: "/home", Views: 100}}
	dashboard.TopPages(c, pages)
}

func demoRealtime(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	count := intState(root, "count", 10)
	root.Children(func() {
		dashboard.Realtime(c, count, dashboard.RealtimeOptions{})
	})
}

func demoActivityFeed(c *ui.Context, log func(string)) {
	events := []dashboard.Activity{{ID: "deploy", Text: "Deployed", Time: "2m ago"}}
	dashboard.ActivityFeed(c, events)
}

func demoLeaders(c *ui.Context, log func(string)) {
	people := []dashboard.Leader{{Name: "Ana", Score: 100}}
	dashboard.Leaders(c, people)
}

func demoUptime(c *ui.Context, log func(string)) {
	days := []dashboard.UptimeDay{{Label: "Mon", Status: dashboard.Up}, {Label: "Tue", Status: dashboard.Down}}
	dashboard.Uptime(c, days)
}

func demoStatusGrid(c *ui.Context, log func(string)) {
	services := []dashboard.Service{{Name: "API", Status: dashboard.Up}}
	dashboard.StatusGrid(c, services)
}

func demoIncident(c *ui.Context, log func(string)) {
	events := []dashboard.IncidentEvent{{Text: "Investigating", Time: "10:00"}}
	dashboard.Incident(c, "API outage", events, dashboard.IncidentOptions{})
}

func demoSLA(c *ui.Context, log func(string)) {
	dashboard.SLA(c, 99.9, 99.5, dashboard.SLAOptions{})
}

func demoRegions(c *ui.Context, log func(string)) {
	regions := []dashboard.RegionValue{{Name: "US", Value: 100}}
	dashboard.Regions(c, regions)
}

func demoDevices(c *ui.Context, log func(string)) {
	shares := []dashboard.DeviceShare{{Kind: "Mobile", Share: 0.7}}
	dashboard.Devices(c, shares)
}

func demoCohort(c *ui.Context, log func(string)) {
	rows := []dashboard.CohortRow{{Label: "Jan", Values: []float64{1, 0.8, 0.6, 0.4, 0.2}}}
	dashboard.Cohort(c, rows)
}
