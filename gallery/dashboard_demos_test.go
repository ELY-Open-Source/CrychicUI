package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDashboardDemosLog(t *testing.T) {
	quiet := map[string]bool{"KPIGrid": true, "MetricTile": true, "GoalProgress": true, "Users": true,
		"Conversion": true, "Traffic": true, "TopPages": true, "Realtime": true, "ActivityFeed": true,
		"Leaders": true, "Uptime": true, "StatusGrid": true, "Incident": true, "SLA": true,
		"Regions": true, "Devices": true, "Cohort": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "dashboard" {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		if err := tt.Click("Year"); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), "period Year") {
			t.Errorf("%s logs %q", e.Name, logs)
		}
	}
	if seen != 18 {
		t.Errorf("checked %d dashboard demos", seen)
	}
}
