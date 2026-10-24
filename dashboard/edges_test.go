package dashboard

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func narrow(view func(c *ui.Context)) func(c *ui.Context) {
	return padded(func(c *ui.Context) { ui.Column(c).Width(320).Children(func() { view(c) }) })
}

func TestCountsStayExactPast2To53(t *testing.T) {
	big := 1<<53 + 1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Realtime(c, big, RealtimeOptions{})
		TopPages(c, []PageView{{Path: "/a", Views: big}, {Path: "/b", Views: 1}}, TopPagesOptions{})
		Leaders(c, []Leader{{Name: "Ana", Score: big}}, LeadersOptions{})
	}), 600, 900)
	tt.Advance(2 * time.Second)
	if n := len(slices.DeleteFunc(tt.Texts(), func(s string) bool { return s != "9,007,199,254,740,993" })); n != 3 || !tt.HasText("9,007,199,254,740,994 views") {
		t.Errorf("texts %q", tt.Texts())
	}
	panicsSaying(t, "TopPages total overflows", func(c *ui.Context) {
		TopPages(c, []PageView{{Path: "/a", Views: math.MaxInt}, {Path: "/b", Views: 1}}, TopPagesOptions{})
	})
}

func TestGoalPercentOverflow(t *testing.T) {
	panicsSaying(t, `GoalProgress "Q1" percent overflows`, func(c *ui.Context) {
		GoalProgress(c, []Goal{{Quarter: "Q1", Current: 1, Target: 1e-308}}, GoalProgressOptions{})
	})
}

func TestRevenueSwitchFitsNarrowCard(t *testing.T) {
	periods := []string{"Day", "Week", "Month", "Quarter", "Half", "Year"}
	period := "Day"
	points := []Point{{Label: "a", Value: 1}, {Label: "b", Value: 2}}
	tt := ui.NewTester(narrow(func(c *ui.Context) { Revenue(c, points, "USD", &period, periods, RevenueOptions{}) }), 500, 500)
	settles(t, tt)
	last := rectOf(t, tt, "Year")
	if last.X+last.W > 10+320+0.5 || last.W <= 0 {
		t.Fatalf("Year at %v runs past the 320 card", last)
	}
	click(t, tt, "Year")
	if period != "Year" {
		t.Errorf("period %q", period)
	}
}

func TestCohortDenseAndLongLabel(t *testing.T) {
	values := make([]float64, 12)
	for i := range values {
		values[i] = 1 - float64(i)/100
	}
	long := "Customers who signed up during the August launch week"
	tt := ui.NewTester(narrow(func(c *ui.Context) {
		Cohort(c, []CohortRow{{Label: long, Size: 10, Values: values}}, CohortOptions{})
	}), 500, 400)
	settles(t, tt)
	a, b := rectOf(t, tt, "100%"), rectOf(t, tt, "99%")
	if a.X+a.W > b.X {
		t.Errorf("cells overlap: %v then %v", a, b)
	}
	hover(t, tt, b.X+b.W/2, b.Y+b.H/2)
	text := long + " · Week 1 · 99%"
	r := rectOf(t, tt, text)
	if r.X+r.W > 10+320+0.5 {
		t.Errorf("readout %v runs past the 320 card", r)
	}
}

func TestPanicsNameTheCaller(t *testing.T) {
	panicsSaying(t, `MetricTile "MRR"`, func(c *ui.Context) {
		MetricTile(c, "MRR", "$1", "+1%", []float64{1}, MetricTileOptions{})
	})
	panicsSaying(t, `MetricTile "MRR" value`, func(c *ui.Context) {
		MetricTile(c, "MRR", "$1", "+1%", []float64{1, math.NaN()}, MetricTileOptions{})
	})
	panicsSaying(t, `Incident "Checkout errors"`, func(c *ui.Context) {
		Incident(c, "Checkout errors", nil, IncidentOptions{})
	})
	panicsSaying(t, `KPIGrid tile "Revenue"`, func(c *ui.Context) {
		KPIGrid(c, []KPI{{Label: "Revenue", Value: "1", Points: []float64{1, math.Inf(1)}}})
	})
}
