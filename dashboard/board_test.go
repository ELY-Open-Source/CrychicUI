package dashboard

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func rectOf(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("no %q in %q", name, tt.Texts())
	}
	return r
}

func hover(t *testing.T, tt *ui.Tester, x, y float32) {
	t.Helper()
	tt.Move(x, y)
	tt.Advance(time.Second)
}

func panicsWith(t *testing.T, view func(c *ui.Context)) { panicsSaying(t, "dashboard: ", view) }

// panicsSaying wants a panic whose message holds want.
func panicsSaying(t *testing.T, want string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want %q", r, want)
		}
	}()
	ui.Render(view, 500, 400, 1)
}

func exact(tt *ui.Tester, s string) bool { return slices.Contains(tt.Texts(), s) }

func TestRevenueKeepsPeriodAcrossReorder(t *testing.T) {
	period, periods := "Week", []string{"Week", "Month", "Year"}
	points := []Point{{Label: "Mon", Value: 100}, {Label: "Sun", Value: 150}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Revenue(c, points, "USD", &period, periods, RevenueOptions{}) }), 500, 500)
	settles(t, tt)
	if !tt.HasText("USD 250.00") || !tt.HasText("+50.0%") {
		t.Fatalf("total and change %q", tt.Texts())
	}
	r := rectOf(t, tt, "Year")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	periods[1], periods[2] = periods[2], periods[1]
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
	if period != "Week" {
		t.Errorf("released over a moved choice and picked %q", period)
	}
}

func TestRevenueRejectsBadPeriods(t *testing.T) {
	points := []Point{{Label: "a", Value: 1}, {Label: "b", Value: 2}}
	odd, same := "Day", "Week"
	panicsWith(t, func(c *ui.Context) { Revenue(c, points, "USD", &odd, []string{"Week", "Month"}, RevenueOptions{}) })
	panicsSaying(t, "Revenue name", func(c *ui.Context) { Revenue(c, points, "USD", &same, []string{"Week", "Week"}, RevenueOptions{}) })
}

func TestRevenueChartMorphs(t *testing.T) {
	period := "Week"
	sets := map[string][]Point{"Week": {{Label: "a", Value: 1}, {Label: "b", Value: 3}}, "Year": {{Label: "a", Value: 3}, {Label: "b", Value: 1}}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Revenue(c, sets[period], "USD", &period, []string{"Week", "Year"}, RevenueOptions{})
	}), 500, 500)
	settles(t, tt)
	top, axis := rectOf(t, tt, "USD 4.00"), rectOf(t, tt, "a")
	click(t, tt, "Year")
	tt.Advance(80 * time.Millisecond)
	mid := tt.Image()
	settles(t, tt)
	end, moved := tt.Image(), 0
	for y := int(top.Y + top.H + 4); y < int(axis.Y-4); y++ {
		for x := int(axis.X); x < int(axis.X+300); x++ {
			if mid.RGBAAt(x, y) != end.RGBAAt(x, y) {
				moved++
			}
		}
	}
	if period != "Year" || moved == 0 {
		t.Errorf("period %q; the line jumped instead of morphing", period)
	}
}

func TestUsersHoverShowsDay(t *testing.T) {
	points := []Point{{Label: "Mon", Value: 1200}, {Label: "Tue", Value: 3400}, {Label: "Wed", Value: 2000}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Users(c, points, UsersOptions{}) }), 500, 500)
	settles(t, tt)
	if !tt.HasText("2,000") || !tt.HasText("-41.2%") {
		t.Fatalf("latest %q", tt.Texts())
	}
	chart := rectOf(t, tt, "Mon")
	hover(t, tt, chart.X+5, chart.Y-40)
	if !tt.HasText("1,200") {
		t.Errorf("hovered first bar %q", tt.Texts())
	}
}

func TestUptimeHoverNamesDay(t *testing.T) {
	days := []UptimeDay{{Label: "Mon"}, {Label: "Tue", Status: Down}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Uptime(c, days, UptimeOptions{}) }), 500, 300)
	settles(t, tt)
	end := rectOf(t, tt, "Tue")
	hover(t, tt, end.X+end.W-4, end.Y-20)
	if !tt.HasText("Outage") {
		t.Errorf("hovered day %q", tt.Texts())
	}
}

func TestCohortHoverReadsCell(t *testing.T) {
	rows := []CohortRow{{Label: "Aug 4", Size: 1204, Values: []float64{1, 0.62}}, {Label: "Aug 11", Size: 90, Values: []float64{1}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Cohort(c, rows, CohortOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("1,204 users") {
		t.Fatalf("sizes %q", tt.Texts())
	}
	cell := rectOf(t, tt, "62%")
	hover(t, tt, cell.X+cell.W/2, cell.Y+cell.H/2)
	if !tt.HasText("Aug 4 · Week 1 · 62%") {
		t.Errorf("readout %q", tt.Texts())
	}
	panicsWith(t, func(c *ui.Context) { Cohort(c, []CohortRow{{Label: "a", Values: []float64{1.2}}}, CohortOptions{}) })
}

func TestStatusGridSummary(t *testing.T) {
	for want, services := range map[string][]Service{
		"All systems operational": {{Name: "API"}, {Name: "Web"}},
		"1 issue":                 {{Name: "API", Status: Degraded}, {Name: "Web"}},
		"2 issues":                {{Name: "API", Status: Down}, {Name: "Web", Status: Degraded}},
	} {
		tt := ui.NewTester(padded(func(c *ui.Context) { StatusGrid(c, services, StatusGridOptions{}) }), 500, 400)
		tt.Frame()
		if !exact(tt, want) {
			t.Errorf("want %q in %q", want, tt.Texts())
		}
	}
	panicsWith(t, func(c *ui.Context) { StatusGrid(c, []Service{{Name: "API"}, {Name: "API"}}, StatusGridOptions{}) })
}

func TestSLABudget(t *testing.T) {
	for _, tc := range []struct {
		actual, target float64
		want           []string
	}{
		{99.97, 99.9, []string{"Met", "70% of error budget left"}},
		{99.81, 99.9, []string{"Missed", "0% of error budget left"}},
		{100, 100, []string{"Met", "100% of error budget left"}},
		{99.99, 100, []string{"Missed", "0% of error budget left"}},
	} {
		tt := ui.NewTester(padded(func(c *ui.Context) { SLA(c, tc.actual, tc.target, SLAOptions{}) }), 500, 300)
		settles(t, tt)
		for _, w := range tc.want {
			if !tt.HasText(w) {
				t.Errorf("%v: want %q in %q", tc, w, tt.Texts())
			}
		}
	}
}

func TestRankedListsValidate(t *testing.T) {
	panicsSaying(t, "sorted", func(c *ui.Context) {
		Leaders(c, []Leader{{Name: "A", Score: 1}, {Name: "B", Score: 2}}, LeadersOptions{})
	})
	panicsWith(t, func(c *ui.Context) {
		Leaders(c, []Leader{{Name: "A", Score: 2}, {Name: "A", Score: 1}}, LeadersOptions{})
	})
	panicsWith(t, func(c *ui.Context) {
		Conversion(c, []FunnelStage{{Label: "a", Count: 10}, {Label: "b", Count: 20}}, ConversionOptions{})
	})
	panicsWith(t, func(c *ui.Context) {
		Traffic(c, []Share{{Name: "a", Share: 0.7}, {Name: "b", Share: 0.4}}, TrafficOptions{})
	})
	panicsWith(t, func(c *ui.Context) { KPIGrid(c, []KPI{{Label: "a", Value: "1"}, {Label: "a", Value: "2"}}) })
	panicsWith(t, func(c *ui.Context) { TopPages(c, []PageView{{Path: "/a"}, {Path: "/a"}}, TopPagesOptions{}) })
}

func TestActivityFeedEmptyAndArrival(t *testing.T) {
	var events []Activity
	tt := ui.NewTester(padded(func(c *ui.Context) { ActivityFeed(c, events, ActivityFeedOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("No activity yet") {
		t.Fatalf("empty %q", tt.Texts())
	}
	events = []Activity{{ID: "a", Actor: "Maya Chen", Text: "deployed api", Time: "2m"}}
	tt.Frame()
	if !tt.FramePending() || !tt.HasText("Maya Chen ") {
		t.Errorf("arrival %v texts %q", tt.FramePending(), tt.Texts())
	}
}

func TestGoalsMetAndIncidentResolved(t *testing.T) {
	goals := []Goal{{Quarter: "Q1", Current: 1240, Target: 1200}, {Quarter: "Q2", Current: 300, Target: 1200}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		GoalProgress(c, goals, GoalProgressOptions{})
		Incident(c, "Checkout errors", []IncidentEvent{{Stage: "Resolved", Text: "Fixed.", Time: "10:21"}}, IncidentOptions{Resolved: true})
	}), 500, 600)
	settles(t, tt)
	for _, w := range []string{"1 of 2 met", "1,240 of 1,200", "103%", "Resolved", "Fixed."} {
		if !tt.HasText(w) {
			t.Errorf("want %q in %q", w, tt.Texts())
		}
	}
	if tt.HasText("Ongoing") {
		t.Error("a resolved incident shows Ongoing")
	}
}
