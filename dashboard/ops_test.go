package dashboard

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestWidgetFrameRefreshReloadsBody(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { WidgetFrame(c, "Signups", WidgetFrameOptions{}, func() { ui.Text(c, "1,246") }) }), 500, 300)
	settles(t, tt)
	click(t, tt, "Widget actions")
	if err := tt.ChooseMenuItem("Refresh"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	r := rectOf(t, tt, "1,246")
	tt.Advance(40 * time.Millisecond)
	mid := tt.Image()
	settles(t, tt)
	end, moved := tt.Image(), 0
	for y := int(r.Y) - 4; y < int(r.Y+r.H)+4; y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if mid.RGBAAt(x, y) != end.RGBAAt(x, y) {
				moved++
			}
		}
	}
	if moved == 0 {
		t.Error("the body did not reload on refresh")
	}
}

func TestWidgetGridHighlightsDropTarget(t *testing.T) {
	titles := []string{"Revenue", "Churn"}
	tt := ui.NewTester(padded(func(c *ui.Context) { WidgetGrid(c, &titles, WidgetGridOptions{}) }), 500, 300)
	settles(t, tt)
	a, b := rectOf(t, tt, "Revenue"), rectOf(t, tt, "Churn")
	before := tt.Image().RGBAAt(int(b.X+b.W-6), int(b.Y+b.H/2))
	// The ghost then covers only the target's left side.
	tt.Press(a.X+a.W-8, a.Y+a.H/2)
	tt.Move(b.X+8, b.Y+b.H/2)
	tt.Advance(time.Second)
	over := tt.Image().RGBAAt(int(b.X+b.W-6), int(b.Y+b.H/2))
	tt.Release(b.X+8, b.Y+b.H/2)
	tt.Frame()
	if over == before {
		t.Error("the drop target does not light up under a drag")
	}
	if !slices.Equal(titles, []string{"Churn", "Revenue"}) {
		t.Errorf("titles %v", titles)
	}
	dup := []string{"A", "A"}
	panicsSaying(t, "WidgetGrid name", func(c *ui.Context) { WidgetGrid(c, &dup, WidgetGridOptions{}) })
}

func TestReportExportDoneState(t *testing.T) {
	format, exported := "PDF", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ReportExport(c, &format, []string{"PDF", "CSV"}, &exported, ReportExportOptions{})
	}), 500, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Export")
	settles(t, tt)
	if !exported || !tt.HasText("PDF report ready") {
		t.Fatalf("exported %v texts %q", exported, tt.Texts())
	}
	click(t, tt, "CSV")
	settles(t, tt)
	if exported || tt.HasText("report ready") {
		t.Errorf("a new format kept the old done state: %q", tt.Texts())
	}
	cues(t, before, sound.Success, sound.Select)
	odd := "DOCX"
	panicsSaying(t, `format "DOCX"`, func(c *ui.Context) { ReportExport(c, &odd, []string{"PDF"}, &exported, ReportExportOptions{}) })
}

func TestRangeBarStepsByItsLength(t *testing.T) {
	start := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 14)
	compare := false
	tt := ui.NewTester(padded(func(c *ui.Context) { RangeBar(c, &start, &end, &compare, RangeBarOptions{}) }), 600, 200)
	settles(t, tt)
	if !tt.HasText("14 days") || !tt.HasText("Oct 1, 2026 – Oct 15, 2026") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "Next")
	if start.Day() != 15 || end.Day() != 29 {
		t.Errorf("stepped to %v – %v", start, end)
	}
	far0, far1 := time.Date(1700, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC)
	want := fmt.Sprintf("%d days", (far1.Unix()-far0.Unix())/86400)
	one, day := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	ft := ui.NewTester(padded(func(c *ui.Context) {
		RangeBar(c, &far0, &far1, &compare, RangeBarOptions{})
		RangeBar(c, &one, &day, &compare, RangeBarOptions{Compare: "Versus"})
	}), 600, 300)
	settles(t, ft)
	if !ft.HasText(want) || !exact(ft, "1 day") {
		t.Errorf("want %q and 1 day in %q", want, ft.Texts())
	}
}

func TestChurnRiseReadsRed(t *testing.T) {
	points := []Point{{Label: "Sep", Value: 2}, {Label: "Oct", Value: 2.5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Churn(c, 2.5, points, ChurnOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("+0.5 pts") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestAnomalyFlags(t *testing.T) {
	points := []Point{{Label: "14:00", Value: 400}, {Label: "14:05", Value: 1180}, {Label: "14:10", Value: 410}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Anomaly(c, points, []int{1}, AnomalyOptions{}) }), 500, 400)
	tt.Advance(2 * time.Second)
	if !tt.HasText("1 flagged") || !tt.HasText("14:05") || !tt.HasText("1,180") {
		t.Errorf("texts %q", tt.Texts())
	}
	panicsSaying(t, "unique point index", func(c *ui.Context) { Anomaly(c, points, []int{1, 1}, AnomalyOptions{}) })
}

func TestScorecardAndHealthBands(t *testing.T) {
	goals := []ScoreGoal{{Label: "ARR", Current: 14, Target: 10}, {Label: "CSAT", Current: 91, Target: 95}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Scorecard(c, goals, ScorecardOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("1 of 2 met") || !exact(tt, "Met") || !exact(tt, "Behind") || !tt.HasText("91 of 95") {
		t.Errorf("texts %q", tt.Texts())
	}
	panicsSaying(t, "Scorecard name", func(c *ui.Context) {
		Scorecard(c, []ScoreGoal{{Label: "A", Target: 1}, {Label: "A", Target: 1}}, ScorecardOptions{})
	})
	for score, band := range map[int]string{80: "Healthy", 79: "Needs attention", 50: "Needs attention", 49: "Critical"} {
		ht := ui.NewTester(padded(func(c *ui.Context) { HealthRing(c, score, HealthRingOptions{}) }), 500, 300)
		ht.Frame()
		if !exact(ht, band) {
			t.Errorf("%d: want %q in %q", score, band, ht.Texts())
		}
	}
}

func TestQueueLatencyAndSignups(t *testing.T) {
	days := []SignupDay{{Label: "Mon", Count: 1200}, {Label: "Tue", Count: 34}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Queue(c, 1284, 36, QueueOptions{})
		Latency(c, 84, 212, 1480, LatencyOptions{})
		Signups(c, days, SignupsOptions{})
	}), 600, 1000)
	settles(t, tt)
	for _, w := range []string{"1,284", "36", "84 ms", "1.48 s", "1,234 total", "Tue"} {
		if !tt.HasText(w) {
			t.Errorf("want %q in %q", w, tt.Texts())
		}
	}
	mon := rectOf(t, tt, "Mon")
	hover(t, tt, mon.X+mon.W/2, mon.Y-30)
	if !exact(tt, "1,200") {
		t.Errorf("hovered Monday %q", tt.Texts())
	}
	panicsSaying(t, "Queue", func(c *ui.Context) { Queue(c, 1<<62, 1<<62, QueueOptions{}) })
	panicsSaying(t, "p50 ≤ p95 ≤ p99", func(c *ui.Context) { Latency(c, 10, 5, 20, LatencyOptions{}) })
}

func TestErrorRateMRRNPSAndBanner(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ErrorRate(c, 1, 2, ErrorRateOptions{History: []float64{1, 3, 1}})
		MRR(c, 1000, 100, 0, 1100, "USD", MRROptions{})
		NPSWidget(c, 50, 25, 25, NPSWidgetOptions{})
		WelcomeBanner(c, "Good morning", "All green.", WelcomeBannerOptions{Stats: []BannerStat{{Label: "Signups", Value: "241"}}})
	}), 600, 1400)
	settles(t, tt)
	for _, w := range []string{"Within limit", "Net new USD -1,000.00", "NPS 25", "50%", "Signups", "241"} {
		if !tt.HasText(w) {
			t.Errorf("want %q in %q", w, tt.Texts())
		}
	}
	panicsSaying(t, "at most four stats", func(c *ui.Context) {
		WelcomeBanner(c, "Hi", "Day", WelcomeBannerOptions{Stats: make([]BannerStat, 5)})
	})
}

func TestInsightsEmptyAndMetric(t *testing.T) {
	var cards []Insight
	tt := ui.NewTester(padded(func(c *ui.Context) { Insights(c, cards, InsightsOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("No insights yet") {
		t.Fatalf("empty %q", tt.Texts())
	}
	cards = []Insight{{Title: "Trials up", Body: "More trials.", Metric: "+18%"}}
	tt.Frame()
	settles(t, tt)
	if !tt.HasText("+18%") || tt.HasText("No insights yet") {
		t.Errorf("texts %q", tt.Texts())
	}
}
