package dashboard

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestChurn(t *testing.T) {
	points := []Point{{Label: "Jan", Value: 2}, {Label: "Feb", Value: 3}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Churn(c, 3, points, ChurnOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("Churn rate") || !tt.HasText("3.0%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestChurnHistoryValidation(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Churn(c, 3, []Point{{Label: "A", Value: -1}, {Label: "B", Value: 120}}, ChurnOptions{})
	}, 500, 200, 1)
}

func TestMRR(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { MRR(c, 1000, 100, 50, 25, "USD", MRROptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("USD 1,000.00") || !tt.HasText("Expansion") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMRRCopy(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MRR(c, 1000, 100, 50, 25, "USD", MRROptions{New: "New money", Expansion: "Expanded", Churn: "Lost"})
	}), 500, 250)
	if !tt.HasText("New money") || !tt.HasText("Expanded") || !tt.HasText("Lost") {
		t.Errorf("custom labels %q", tt.Texts())
	}
}

func TestRangeBar(t *testing.T) {
	start := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(7 * 24 * time.Hour)
	compare := false
	tt := ui.NewTester(padded(func(c *ui.Context) { RangeBar(c, &start, &end, &compare, RangeBarOptions{}) }), 600, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Previous")
	if start.Day() != 24 {
		t.Errorf("start %v", start)
	}
	click(t, tt, "Compare")
	if !compare {
		t.Error("compare not enabled")
	}
	cues(t, before, sound.Select, sound.Select)
}

func TestRangeBarDST(t *testing.T) {
	zone, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(2024, 11, 2, 0, 0, 0, 0, zone)
	end := start.AddDate(0, 0, 7)
	compare := false
	tt := ui.NewTester(padded(func(c *ui.Context) { RangeBar(c, &start, &end, &compare, RangeBarOptions{}) }), 600, 180)
	settles(t, tt)
	click(t, tt, "Next")
	if start.Day() != 9 {
		t.Errorf("start %v", start)
	}
}

func TestWidgetFrame(t *testing.T) {
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if WidgetFrame(c, "Revenue", WidgetFrameOptions{}, func() { ui.Text(c, "Body") }).Changed() {
			changed++
		}
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Widget actions")
	if err := tt.ChooseMenuItem("Refresh"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if changed != 1 {
		t.Errorf("changed %d", changed)
	}
	cues(t, before, sound.Tap)
}

func TestWidgetGrid(t *testing.T) {
	titles := []string{"Alpha", "Beta"}
	tt := ui.NewTester(padded(func(c *ui.Context) { WidgetGrid(c, &titles, WidgetGridOptions{}) }), 500, 300)
	settles(t, tt)
	a, _ := tt.Find("Alpha")
	b, _ := tt.Find("Beta")
	before := baseline()
	tt.Press(a.X+a.W/2, a.Y+a.H/2)
	tt.Move(b.X+b.W/2, b.Y+b.H/2)
	tt.Release(b.X+b.W/2, b.Y+b.H/2)
	tt.Frame()
	if !slices.Equal(titles, []string{"Beta", "Alpha"}) {
		t.Errorf("titles %v", titles)
	}
	cues(t, before, sound.Drop)
}

func TestWidgetGridUpward(t *testing.T) {
	titles := []string{"Alpha", "Beta"}
	tt := ui.NewTester(padded(func(c *ui.Context) { WidgetGrid(c, &titles, WidgetGridOptions{}) }), 500, 300)
	settles(t, tt)
	a, _ := tt.Find("Alpha")
	b, _ := tt.Find("Beta")
	tt.Press(b.X+b.W/2, b.Y+b.H/2)
	tt.Move(a.X+a.W/2, a.Y+a.H/2)
	tt.Release(a.X+a.W/2, a.Y+a.H/2)
	tt.Frame()
	if !slices.Equal(titles, []string{"Beta", "Alpha"}) {
		t.Errorf("titles %v", titles)
	}
}

func TestWidgetGridIsolation(t *testing.T) {
	aTitles := []string{"Alpha", "Beta"}
	bTitles := []string{"Gamma", "Alpha"}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		c.Key("grid-a")
		WidgetGrid(c, &aTitles, WidgetGridOptions{})
		c.Key("grid-b")
		WidgetGrid(c, &bTitles, WidgetGridOptions{})
	}), 500, 450)
	settles(t, tt)
	a, _ := tt.Find("Alpha")
	g, _ := tt.Find("Gamma")
	tt.Press(a.X+a.W/2, a.Y+a.H/2)
	tt.Move(g.X+g.W/2, g.Y+g.H/2)
	tt.Release(g.X+g.W/2, g.Y+g.H/2)
	tt.Frame()
	if !slices.Equal(aTitles, []string{"Alpha", "Beta"}) || !slices.Equal(bTitles, []string{"Gamma", "Alpha"}) {
		t.Errorf("cross-grid drop: A %v B %v", aTitles, bTitles)
	}
}

func TestReportExport(t *testing.T) {
	format, exported := "PDF", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ReportExport(c, &format, []string{"PDF", "CSV"}, &exported, ReportExportOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "CSV")
	click(t, tt, "Export")
	if format != "CSV" || !exported {
		t.Errorf("format %q exported %v", format, exported)
	}
	cues(t, before, sound.Select, sound.Success)
}

func TestInsights(t *testing.T) {
	cards := []Insight{{Title: "Growth", Body: "Revenue rose"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Insights(c, cards, InsightsOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Growth") || !tt.HasText("Revenue rose") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestAnomaly(t *testing.T) {
	points := []Point{{Label: "A", Value: 1}, {Label: "B", Value: 3}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Anomaly(c, points, []int{1}, AnomalyOptions{}) }), 500, 250)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("anomaly point does not pulse")
	}
}

func TestAnomalyPosition(t *testing.T) {
	points := []Point{{Label: "A", Value: 1}, {Label: "B", Value: 3}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Anomaly(c, points, []int{1}, AnomalyOptions{}) }), 500, 400)
	tt.Advance(2 * time.Second)
	head, row := rectOf(t, tt, "1 flagged"), rectOf(t, tt, "B")
	ink := ui.LightTheme().Danger
	img := tt.Image()
	red := func(x0, x1 int) int {
		n := 0
		for y := int(head.Y + head.H); y < int(row.Y); y++ {
			for x := x0; x < x1; x++ {
				if got := img.RGBAAt(x, y); got.R == ink.R && got.G == ink.G && got.B == ink.B {
					n++
				}
			}
		}
		return n
	}
	if right, left := red(int(head.X+head.W)-40, int(head.X+head.W)), red(10, 60); right == 0 || left != 0 {
		t.Errorf("flag should sit at the last point: right %d, left %d", right, left)
	}
}

func TestAnomalyOffscreenValidation(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Height(1000)
		Anomaly(c, []Point{{Label: "A", Value: 1}, {Label: "B", Value: 2}}, []int{-1}, AnomalyOptions{})
	}, 500, 100)
}

func TestScorecard(t *testing.T) {
	goals := []ScoreGoal{{Label: "Revenue", Current: 100, Target: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Scorecard(c, goals, ScorecardOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Revenue") || !tt.HasText("Met") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestHealthRing(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { HealthRing(c, 90, HealthRingOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Health") || !exact(tt, "90") || !tt.HasText("Healthy") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestServerLoad(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ServerLoad(c, 60, 70, ServerLoadOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("CPU") || !tt.HasText("60%") || !tt.HasText("70%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestQueue(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Queue(c, 3, 2, QueueOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Waiting") || !tt.HasText("Running") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestErrorRate(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ErrorRate(c, 3, 2, ErrorRateOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("errors per minute") || !tt.HasText("Limit 2") || !tt.HasText("Over limit") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestErrorRateLimitLine(t *testing.T) {
	for _, limit := range []float64{2, 1e9} {
		tt := ui.NewTester(padded(func(c *ui.Context) { ErrorRate(c, 3, limit, ErrorRateOptions{}) }), 500, 300)
		settles(t, tt)
		r := rectOf(t, tt, "3.0")
		ink, img, n := ui.LightTheme().Text, tt.Image(), 0
		for y := int(r.Y + r.H + 8); y < int(r.Y+r.H+100); y++ {
			for x := int(r.X); x < int(r.X)+20; x++ {
				if got := img.RGBAAt(x, y); got.R == ink.R && got.G == ink.G && got.B == ink.B {
					n++
				}
			}
		}
		if limit == 2 && n == 0 {
			t.Error("limit line missing")
		}
		if limit == 1e9 && n == 0 {
			t.Error("a far limit should still draw at the top of the chart")
		}
	}
}

func TestLatency(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Latency(c, 10, 20, 30, LatencyOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("p50") || !tt.HasText("20 ms") || !tt.HasText("30 ms") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSignups(t *testing.T) {
	days := []SignupDay{{Label: "Mon", Count: 10}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Signups(c, days, SignupsOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Mon") || !tt.HasText("10") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestNPSWidget(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { NPSWidget(c, 80, 0, 20, NPSWidgetOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("NPS 60") || !tt.HasText("Promoters") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestNPSWidgetOverflow(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { NPSWidget(c, math.MaxInt, 0, 0, NPSWidgetOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("NPS 100") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestNPSWidgetRounding(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { NPSWidget(c, 100000000000314187, 0, 0, NPSWidgetOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("NPS 100") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestWelcomeBanner(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { WelcomeBanner(c, "Good morning", "3 tasks today", WelcomeBannerOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Good morning") || !tt.HasText("3 tasks today") {
		t.Errorf("texts %q", tt.Texts())
	}
}
