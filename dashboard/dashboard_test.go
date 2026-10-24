package dashboard

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func padded(view func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) { ui.Column(c).Padding(10).Gap(8).Children(func() { view(c) }) }
}

func settles(t *testing.T, tt *ui.Tester) {
	t.Helper()
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Error("animation did not settle")
	}
}

func click(t *testing.T, tt *ui.Tester, name string) {
	t.Helper()
	if err := tt.Click(name); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
}

func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Warning)
	}
	return sound.Recent()
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	after := sound.Recent()
	var got []sound.Cue
	if len(before) == 32 {
		i := 0
		for i < len(after) && after[i] == sound.Warning {
			i++
		}
		got = after[i:]
	}
	if len(got) != len(want) {
		t.Errorf("cues %v, want %v", got, want)
		return
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("cues %v, want %v", got, want)
		}
	}
}

func TestKPIGrid(t *testing.T) {
	tiles := []KPI{{Label: "Revenue", Value: "$12k", Delta: "+8%", Points: []float64{1, 2}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { KPIGrid(c, tiles) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Revenue") || !tt.HasText("$12k") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMetricTile(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { MetricTile(c, "Users", "1,200", "+8%", []float64{1, 2}, MetricTileOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Users") || !tt.HasText("1,200") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestGoalProgress(t *testing.T) {
	goals := []Goal{{Quarter: "Q1", Current: 80, Target: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { GoalProgress(c, goals, GoalProgressOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Q1") || !tt.HasText("80 of 100") || !tt.HasText("80%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestRevenue(t *testing.T) {
	period := "Month"
	points := []Point{{Label: "Jan", Value: 10}, {Label: "Feb", Value: 20}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Revenue(c, points, "USD", &period, []string{"Month", "Year"}, RevenueOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Year")
	if period != "Year" {
		t.Errorf("period %q", period)
	}
	cues(t, before, sound.Select)
}

func TestUsers(t *testing.T) {
	points := []Point{{Label: "Jan", Value: 10}, {Label: "Feb", Value: 20}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Users(c, points, UsersOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("Active users") || !tt.HasText("20") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestUsersNegative(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Users(c, []Point{{Label: "A", Value: 1}, {Label: "B", Value: -1}}, UsersOptions{})
	}, 500, 200, 1)
}

func TestConversion(t *testing.T) {
	stages := []FunnelStage{{Label: "Visits", Count: 100}, {Label: "Signups", Count: 20}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Conversion(c, stages, ConversionOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Visits") || !tt.HasText("20%") || !tt.HasText("20% overall") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestConversionZeroStage(t *testing.T) {
	stages := []FunnelStage{{Label: "Visits", Count: 100}, {Label: "Signups", Count: 0}, {Label: "Active", Count: 0}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Conversion(c, stages, ConversionOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("0%") || strings.Contains(strings.Join(tt.Texts(), " "), "NaN") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTraffic(t *testing.T) {
	sources := []Share{{Name: "Search", Share: 0.6}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Traffic(c, sources, TrafficOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Search") || !tt.HasText("60%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTopPages(t *testing.T) {
	pages := []PageView{{Path: "/home", Views: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { TopPages(c, pages, TopPagesOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("/home") || !tt.HasText("100") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestRealtime(t *testing.T) {
	count := 10
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Realtime(c, count, RealtimeOptions{History: []int{4, 8, 10}})
	}), 500, 300)
	tt.Frame()
	if !tt.HasText("Visitors now") || !tt.HasText("10") {
		t.Fatalf("texts %q", tt.Texts())
	}
	count = 12
	tt.Advance(time.Second)
	if !tt.HasText("12") || !tt.FramePending() {
		t.Errorf("count %q, pulse pending %v", tt.Texts(), tt.FramePending())
	}
	panics := func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic")
			}
		}()
		ui.Render(func(c *ui.Context) { Realtime(c, -1, RealtimeOptions{}) }, 500, 180, 1)
	}
	panics()
}

func TestActivityFeed(t *testing.T) {
	events := []Activity{{ID: "a", Text: "Deployed", Time: "2m ago"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ActivityFeed(c, events, ActivityFeedOptions{}) }), 500, 180)
	settles(t, tt)
	events = append([]Activity{{ID: "n", Text: "New event", Time: "now"}}, events...)
	tt.Frame()
	if !tt.HasText("New event") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestLeaders(t *testing.T) {
	people := []Leader{{Name: "Ana", Score: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Leaders(c, people, LeadersOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("1") || !tt.HasText("Ana") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestUptime(t *testing.T) {
	days := []UptimeDay{{Label: "Mon", Status: Up}, {Label: "Tue", Status: Down}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Uptime(c, days, UptimeOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("50.00% uptime") || !tt.HasText("Mon") || !tt.HasText("Tue") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestStatusGrid(t *testing.T) {
	services := []Service{{Name: "API", Status: Up}}
	tt := ui.NewTester(padded(func(c *ui.Context) { StatusGrid(c, services, StatusGridOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("API") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestStatusGridMotion(t *testing.T) {
	services := []Service{{Name: "API", Status: Up}}
	tt := ui.NewTester(padded(func(c *ui.Context) { StatusGrid(c, services, StatusGridOptions{}) }), 500, 250)
	settles(t, tt)
	services[0].Status = Down
	tt.Frame()
	if !tt.FramePending() {
		t.Error("status color does not animate")
	}
}

func TestIncident(t *testing.T) {
	events := []IncidentEvent{{Text: "Investigating", Time: "10:00"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Incident(c, "API outage", events, IncidentOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("API outage") || !tt.HasText("Investigating") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSLA(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { SLA(c, 99.9, 99.5, SLAOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("99.90%") || !tt.HasText("Target 99.50%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSLAValidation(t *testing.T) {
	panics := func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic")
			}
		}()
		ui.Render(func(c *ui.Context) { SLA(c, 150, 99.5, SLAOptions{}) }, 500, 180, 1)
	}
	panics()
}

func TestRegions(t *testing.T) {
	regions := []RegionValue{{Name: "US", Value: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Regions(c, regions, RegionsOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("US") || !tt.HasText("100") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDevices(t *testing.T) {
	shares := []DeviceShare{{Kind: "Mobile", Share: 0.7}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Devices(c, shares, DevicesOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Mobile") || !tt.HasText("70%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCohortMotion(t *testing.T) {
	rows := []CohortRow{{Label: "Jan", Values: []float64{1, 0.8, 0.6, 0.4, 0.2}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Cohort(c, rows, CohortOptions{}) }), 500, 300)
	settles(t, tt)
	rows[0].Values[1] = 0.2
	tt.Frame()
	if !tt.FramePending() {
		t.Error("retention color does not animate")
	}
}

func TestCohort(t *testing.T) {
	rows := []CohortRow{{Label: "Jan", Values: []float64{1, 0.8, 0.6, 0.4, 0.2}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Cohort(c, rows, CohortOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Jan") || !tt.HasText("80%") {
		t.Errorf("texts %q", tt.Texts())
	}
}
