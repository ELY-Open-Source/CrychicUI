package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDevtools2DemosLog(t *testing.T) {
	press := func(names ...string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, name := range names {
				if err := tt.Click(name); err != nil {
					return err
				}
				tt.Frame()
			}
			return nil
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"ServiceMap": {press("redis"), "service redis"},
		"Tables":     {press("customers"), "table customers"},
		"Query":      {press("Run"), "run lines 5"},
		"Webhooks":   {press("Retry payment.succeeded"), "retried w2"},
		"RateLimit":  {press("Send 40 requests"), "used 882"},
		"Errors":     {press("NetworkError: Failed to fetch /api/cart"), "error e2"},
		"Inspector":  {press("ProductCard card-2"), "node card-2"},
		"Secrets":    {press("Reveal MAIL_API_KEY"), "secret toggled"},
		"BuildLog":   {press("Checkout"), "step toggled"},
	}
	quiet := map[string]bool{"Health": true, "QueryResult": true, "Schema": true, "Migrations": true, "Queues": true,
		"Cron": true, "Ping": true, "Cache": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "devtools" || e.ID < 1284 {
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
		tt.Advance(2 * time.Second)
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
		t.Errorf("checked %d devtools demos", seen)
	}
}
