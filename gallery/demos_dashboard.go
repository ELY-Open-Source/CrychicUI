package gallery

import (
	"fmt"
	"time"

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
	tiles := []dashboard.KPI{
		{Label: "Revenue", Value: "$48.2k", Delta: "+12.4%", Icon: "dollar-sign", Points: []float64{31, 34, 33, 38, 41, 44, 48}},
		{Label: "Active users", Value: "8,941", Delta: "+5.1%", Icon: "users", Points: []float64{80, 82, 85, 84, 86, 88, 89}},
		{Label: "Churn", Value: "2.3%", Delta: "-0.4%", Icon: "user-minus", Points: []float64{3.1, 2.9, 2.8, 2.7, 2.5, 2.4, 2.3}},
		{Label: "Avg. session", Value: "6m 12s", Delta: "+0.8%", Icon: "timer", Points: []float64{5.8, 6, 5.9, 6.1, 6, 6.2, 6.2}},
	}
	ui.Column(c).Width(500).Children(func() { dashboard.KPIGrid(c, tiles) })
}

func demoMetricTile(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		dashboard.MetricTile(c, "Monthly recurring revenue", "$128,400", "+8.2%",
			[]float64{96, 101, 104, 109, 112, 118, 121, 128}, dashboard.MetricTileOptions{Icon: "trending-up", Note: "vs. $118,700 last month"})
	})
}

func demoGoalProgress(c *ui.Context, log func(string)) {
	goals := []dashboard.Goal{{Quarter: "Q1 2026", Current: 1240, Target: 1200}, {Quarter: "Q2 2026", Current: 1310, Target: 1400},
		{Quarter: "Q3 2026", Current: 1590, Target: 1600}, {Quarter: "Q4 2026", Current: 640, Target: 1800}}
	ui.Column(c).Width(400).Children(func() {
		dashboard.GoalProgress(c, goals, dashboard.GoalProgressOptions{Title: "New customers"})
	})
}

func demoRevenue(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	period := textState(root, "period", "Week")
	sets := map[string][]dashboard.Point{
		"Week": {{Label: "Mon", Value: 6120}, {Label: "Tue", Value: 7340}, {Label: "Wed", Value: 6980}, {Label: "Thu", Value: 8210},
			{Label: "Fri", Value: 9050}, {Label: "Sat", Value: 5430}, {Label: "Sun", Value: 4980}},
		"Month": {{Label: "Oct 1", Value: 41200}, {Label: "Oct 8", Value: 44800}, {Label: "Oct 15", Value: 43100}, {Label: "Oct 22", Value: 48110}},
		"Year": {{Label: "Nov 2025", Value: 98000}, {Label: "Feb", Value: 121000}, {Label: "May", Value: 134000},
			{Label: "Aug", Value: 152000}, {Label: "Oct 2026", Value: 176500}},
	}
	root.Children(func() {
		if dashboard.Revenue(c, sets[*period], "USD", period, []string{"Week", "Month", "Year"}, dashboard.RevenueOptions{}).Changed() {
			log("period " + *period)
		}
	})
}

func demoUsers(c *ui.Context, log func(string)) {
	counts := []float64{3120, 3380, 3290, 3510, 3720, 2410, 2280, 3460, 3640, 3590, 3810, 3920, 2560, 2490}
	points := make([]dashboard.Point, len(counts))
	for i, v := range counts {
		points[i] = dashboard.Point{Label: fmt.Sprintf("Oct %d", 11+i), Value: v}
	}
	ui.Column(c).Width(440).Children(func() { dashboard.Users(c, points, dashboard.UsersOptions{Title: "Daily active users"}) })
}

func demoConversion(c *ui.Context, log func(string)) {
	stages := []dashboard.FunnelStage{{Label: "Visited pricing", Count: 18420}, {Label: "Started trial", Count: 4210},
		{Label: "Invited a teammate", Count: 1960}, {Label: "Became paying", Count: 742}}
	ui.Column(c).Width(420).Children(func() {
		dashboard.Conversion(c, stages, dashboard.ConversionOptions{Title: "Trial funnel"})
	})
}

func demoTraffic(c *ui.Context, log func(string)) {
	sources := []dashboard.Share{{Name: "Organic search", Icon: "search", Share: 0.42}, {Name: "Direct", Icon: "mouse-pointer-click", Share: 0.24},
		{Name: "Referral", Icon: "link", Share: 0.16}, {Name: "Social", Icon: "share-2", Share: 0.11}, {Name: "Email", Icon: "mail", Share: 0.07}}
	ui.Column(c).Width(400).Children(func() { dashboard.Traffic(c, sources, dashboard.TrafficOptions{}) })
}

func demoTopPages(c *ui.Context, log func(string)) {
	pages := []dashboard.PageView{{Path: "/pricing", Views: 28410}, {Path: "/docs/getting-started", Views: 19230},
		{Path: "/blog/launch-week-2026", Views: 12880}, {Path: "/changelog", Views: 7640}, {Path: "/careers", Views: 3120}}
	ui.Column(c).Width(420).Children(func() { dashboard.TopPages(c, pages, dashboard.TopPagesOptions{}) })
}

func demoRealtime(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	start := ui.Local(root, "start", func() time.Time { return c.Now() })
	tick := int(c.Now().Sub(*start) / time.Second)
	c.After(time.Second - c.Now().Sub(*start)%time.Second)
	wave := []int{0, 3, 5, 2, -1, -4, -2, 1, 4, 6, 3, 0, -3, -5, -2}
	history := make([]int, 20)
	for i := range history {
		history[i] = 128 + wave[(tick+i)%len(wave)]
	}
	root.Children(func() {
		dashboard.Realtime(c, history[len(history)-1], dashboard.RealtimeOptions{Note: "on acme.dev in the last 30 minutes", History: history})
	})
}

func demoActivityFeed(c *ui.Context, log func(string)) {
	events := []dashboard.Activity{
		{ID: "e1", Actor: "Maya Chen", Text: "deployed api v2.14.0 to production", Time: "2m"},
		{ID: "e2", Icon: "credit-card", Text: "Northwind Traders upgraded to the Business plan", Time: "18m"},
		{ID: "e3", Actor: "Diego Alvarez", Text: "closed 4 issues in Billing", Time: "1h"},
		{ID: "e4", Icon: "triangle-alert", Text: "Error rate on /checkout crossed 2%", Time: "2h"},
		{ID: "e5", Actor: "Priya Raman", Text: "invited 3 people to Acme Design", Time: "3h"},
	}
	ui.Column(c).Width(440).Children(func() { dashboard.ActivityFeed(c, events, dashboard.ActivityFeedOptions{}) })
}

func demoLeaders(c *ui.Context, log func(string)) {
	people := []dashboard.Leader{{Name: "Sofia Okafor", Detail: "Enterprise, EMEA", Score: 412000}, {Name: "James Whitfield", Detail: "Mid-market", Score: 365500},
		{Name: "Aiko Tanaka", Detail: "Enterprise, APAC", Score: 318200}, {Name: "Lucas Moreau", Detail: "SMB", Score: 241900},
		{Name: "Hannah Becker", Detail: "Mid-market", Score: 198400}}
	ui.Column(c).Width(420).Children(func() {
		dashboard.Leaders(c, people, dashboard.LeadersOptions{Title: "Q3 bookings (USD)"})
	})
}

func demoUptime(c *ui.Context, log func(string)) {
	days := make([]dashboard.UptimeDay, 45)
	for i := range days {
		days[i] = dashboard.UptimeDay{Label: time.Date(2026, 9, 10+i, 0, 0, 0, 0, time.UTC).Format("Jan 2")}
	}
	days[12].Status, days[13].Status, days[31].Status = dashboard.Degraded, dashboard.Down, dashboard.Degraded
	ui.Column(c).Width(460).Children(func() { dashboard.Uptime(c, days, dashboard.UptimeOptions{Title: "api.acme.dev"}) })
}

func demoStatusGrid(c *ui.Context, log func(string)) {
	services := []dashboard.Service{{Name: "API", Detail: "42 ms"}, {Name: "Dashboard", Detail: "118 ms"},
		{Name: "Webhooks", Detail: "Delayed 4 min", Status: dashboard.Degraded}, {Name: "Search", Detail: "61 ms"},
		{Name: "Payments", Detail: "Stripe", Status: dashboard.Up}, {Name: "Email", Detail: "Bounces high", Status: dashboard.Down}}
	ui.Column(c).Width(460).Children(func() { dashboard.StatusGrid(c, services, dashboard.StatusGridOptions{}) })
}

func demoIncident(c *ui.Context, log func(string)) {
	events := []dashboard.IncidentEvent{
		{Stage: "Investigating", Text: "Checkout requests are failing for some customers in eu-west-1.", Time: "09:42"},
		{Stage: "Identified", Text: "A database failover left the payment worker pool on a read-only replica.", Time: "09:58"},
		{Stage: "Monitoring", Text: "Workers restarted on the primary; error rate is back under 0.1%.", Time: "10:21"},
	}
	ui.Column(c).Width(460).Children(func() {
		dashboard.Incident(c, "Elevated checkout errors in Europe", events, dashboard.IncidentOptions{Meta: "Major · Payments, Checkout"})
	})
}

func demoSLA(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Gap(12).Children(func() {
		dashboard.SLA(c, 99.97, 99.9, dashboard.SLAOptions{Title: "API availability · October"})
		dashboard.SLA(c, 99.81, 99.9, dashboard.SLAOptions{Title: "Webhook delivery · October"})
	})
}

func demoRegions(c *ui.Context, log func(string)) {
	regions := []dashboard.RegionValue{{Name: "United States", Code: "US", Value: 48210}, {Name: "Germany", Code: "DE", Value: 17430},
		{Name: "United Kingdom", Code: "GB", Value: 14980}, {Name: "Japan", Code: "JP", Value: 9120}, {Name: "Brazil", Code: "BR", Value: 6540}}
	ui.Column(c).Width(420).Children(func() { dashboard.Regions(c, regions, dashboard.RegionsOptions{Title: "Sessions by country"}) })
}

func demoDevices(c *ui.Context, log func(string)) {
	shares := []dashboard.DeviceShare{{Kind: "Desktop", Share: 0.58}, {Kind: "Mobile", Share: 0.36}, {Kind: "Tablet", Share: 0.06}}
	ui.Column(c).Width(400).Children(func() { dashboard.Devices(c, shares, dashboard.DevicesOptions{}) })
}

func demoCohort(c *ui.Context, log func(string)) {
	rows := []dashboard.CohortRow{
		{Label: "Aug 4", Size: 1204, Values: []float64{1, 0.62, 0.48, 0.41, 0.37, 0.35}},
		{Label: "Aug 11", Size: 1381, Values: []float64{1, 0.66, 0.51, 0.44, 0.40}},
		{Label: "Aug 18", Size: 1129, Values: []float64{1, 0.59, 0.46, 0.39}},
		{Label: "Aug 25", Size: 1460, Values: []float64{1, 0.68, 0.53}},
		{Label: "Sep 1", Size: 1512, Values: []float64{1, 0.71}},
	}
	ui.Column(c).Width(500).Children(func() { dashboard.Cohort(c, rows, dashboard.CohortOptions{Title: "Weekly retention"}) })
}
