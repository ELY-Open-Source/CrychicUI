package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestChartDemosLog(t *testing.T) {
	buttons := map[string][2]string{
		"GaugeChart": {"Add load", "load 59"},
		"Bullet":     {"Book a sale", "revenue 250"},
		"Sparkline":  {"Next week", "week 1"},
		"SparkBars":  {"Next week", "week 1"},
	}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "chart" {
			continue
		}
		seen++
		press, want := "Shuffle data", "data 1"
		if b, ok := buttons[e.Name]; ok {
			press, want = b[0], b[1]
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		if err := tt.Click(press); err != nil {
			t.Fatalf("%s: %v", e.Name, err)
		}
		tt.Advance(500 * time.Millisecond)
		if len(logs) == 0 || !strings.HasPrefix(logs[len(logs)-1], want) {
			t.Errorf("%s logged %q, want %q", e.Name, logs, want)
		}
	}
	if seen != 30 {
		t.Errorf("%d chart demos, want 30", seen)
	}
}

func TestChartDemoLogsPick(t *testing.T) {
	var logs []string
	tt := ui.NewTester(func(c *ui.Context) { demoBars(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
	tt.Advance(time.Second)
	if err := tt.Click("Feb: 18"); err != nil {
		t.Fatal(err)
	}
	if len(logs) == 0 || logs[len(logs)-1] != "chose 1" {
		t.Errorf("logged %q", logs)
	}
}
