package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDashboard2DemosLog(t *testing.T) {
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"RangeBar": {func(tt *ui.Tester) error { return tt.Click("Previous") }, "24"},
		"WidgetFrame": {func(tt *ui.Tester) error {
			if err := tt.Click("Widget actions"); err != nil {
				return err
			}
			return tt.ChooseMenuItem("Refresh")
		}, "refreshed"},
		"WidgetGrid": {func(tt *ui.Tester) error {
			a, _ := tt.Find("Revenue")
			b, _ := tt.Find("Active users")
			tt.Press(a.X+a.W/2, a.Y+a.H/2)
			tt.Move(b.X+b.W/2, b.Y+b.H/2)
			tt.Release(b.X+b.W/2, b.Y+b.H/2)
			return nil
		}, "[Active users Revenue Churn Top pages]"},
		"ReportExport": {func(tt *ui.Tester) error {
			return tt.Click("CSV")
		}, "format CSV"},
	}
	quiet := map[string]bool{"Churn": true, "MRR": true, "Insights": true, "Anomaly": true, "Scorecard": true,
		"HealthRing": true, "ServerLoad": true, "Queue": true, "ErrorRate": true, "Latency": true,
		"Signups": true, "NPSWidget": true, "WelcomeBanner": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "dashboard" || e.ID < 1074 {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 17 {
		t.Errorf("checked %d second dashboard demos", seen)
	}
}
