package devtools

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// swapPress presses a, swaps the data under it, and releases.
func swapPress(t *testing.T, tt *ui.Tester, name string, swap func()) {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("no %q", name)
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	swap()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func TestServiceMap(t *testing.T) {
	services := []Service{{Name: "web"}, {Name: "api"}, {Name: "db", State: HealthDown}}
	calls := []Call{{From: "web", To: "api", Rate: "40/s"}, {From: "api", To: "db", Rate: "90/s", Failing: true}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { ServiceMap(c, services, calls, &selected, ServiceMapOptions{}) }), 520, 400)
	settles(t, tt)
	if !exact(tt, "Pick a service to see its calls") {
		t.Errorf("no hint before a pick: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "api")
	cues(t, before, sound.Select)
	settles(t, tt)
	if selected != "api" || !exact(tt, "40/s") || !exact(tt, "90/s") || !exact(tt, "Healthy") {
		t.Errorf("api's calls: %q", tt.Texts())
	}
	panics(t, "form a loop", func(c *ui.Context) {
		ServiceMap(c, services, append(calls, Call{From: "db", To: "web"}), &selected, ServiceMapOptions{})
	})
	panics(t, "two known services", func(c *ui.Context) {
		ServiceMap(c, services, []Call{{From: "web", To: "cache"}}, &selected, ServiceMapOptions{})
	})
}

func TestHealthCounts(t *testing.T) {
	runs := make([]bool, 60)
	for i := range runs {
		runs[i] = i != 7
	}
	checks := []Check{{Endpoint: "GET /a", History: runs}, {Endpoint: "GET /b", State: HealthDegraded, History: []bool{true}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Health(c, checks, HealthOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "1 of 2 healthy") || !exact(tt, "98.3% passed") || !exact(tt, "100.0% passed") {
		t.Errorf("health: %q", tt.Texts())
	}
	checks = checks[:1]
	tt = ui.NewTester(padded(func(c *ui.Context) { Health(c, checks, HealthOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "All healthy") {
		t.Errorf("one endpoint up: %q", tt.Texts())
	}
	panics(t, "at most 60 checks", func(c *ui.Context) {
		Health(c, []Check{{Endpoint: "x", History: make([]bool, 61)}}, HealthOptions{})
	})
}

func TestTablesPickAndPlurals(t *testing.T) {
	tables := []Table{{Name: "orders", Rows: 48213, Size: 100}, {Name: "settings", Rows: 1, Size: 10}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Tables(c, tables, &selected, TablesOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "48,213 rows") || !exact(tt, "1 row") {
		t.Errorf("rows: %q", tt.Texts())
	}
	swapPress(t, tt, "orders", func() { tables[0], tables[1] = tables[1], tables[0] })
	if selected == "settings" {
		t.Error("a press on orders picked settings after the swap")
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { Tables(c, nil, &selected, TablesOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "No tables") {
		t.Errorf("empty: %q", tt.Texts())
	}
}

func TestQueryRun(t *testing.T) {
	sql := "  "
	var runs int
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, run := Query(c, &sql, QueryOptions{}); run.Clicked() {
			runs++
		}
	}), 500, 400)
	settles(t, tt)
	before := baseline()
	if err := tt.Click("Run"); err == nil {
		tt.Frame()
	}
	if runs != 0 || !exact(tt, "1 line") {
		t.Errorf("blank SQL ran %d times: %q", runs, tt.Texts())
	}
	sql = "SELECT 1\nFROM orders;"
	tt.Frame()
	click(t, tt, "Run")
	cues(t, before, sound.Tap)
	if runs != 1 || !exact(tt, "2 lines") {
		t.Errorf("runs %d: %q", runs, tt.Texts())
	}
}

func TestQueryResultGrid(t *testing.T) {
	cols := []string{"id", "name"}
	rows := [][]string{{"1", "Maya Chen"}, {"NULL", "Priya Nair"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { QueryResult(c, cols, rows, 1500*time.Microsecond, QueryResultOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "2 rows") || !exact(tt, "1.5 ms") {
		t.Errorf("grid: %q", tt.Texts())
	}
	head, _ := tt.Find("id")
	cell, _ := tt.Find("1")
	if math.Abs(float64(head.X+head.W-(cell.X+cell.W))) > 1 {
		t.Errorf("numeric header %v and cell %v do not share a right edge", head, cell)
	}
	panics(t, "1 cells for 2 columns", func(c *ui.Context) {
		QueryResult(c, cols, [][]string{{"1"}}, 0, QueryResultOptions{})
	})
	tt = ui.NewTester(padded(func(c *ui.Context) { QueryResult(c, cols, nil, 0, QueryResultOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "No rows") || !exact(tt, "0 rows") {
		t.Errorf("empty grid: %q", tt.Texts())
	}
}

func TestSchemaAndMigrations(t *testing.T) {
	cols := []Column{{Name: "id", Type: "bigint", Primary: true}, {Name: "coupon", Type: "text", Nullable: true, Default: "''"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Schema(c, "orders", cols, SchemaOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "2 columns") || !exact(tt, "NULL") || !exact(tt, "default ''") {
		t.Errorf("schema: %q", tt.Texts())
	}
	panics(t, `"coupon" needs a type`, func(c *ui.Context) { Schema(c, "orders", []Column{{Name: "coupon"}}, SchemaOptions{}) })
	ms := []Migration{{ID: "1", Name: "Create", Applied: true, At: "Oct 2"}, {ID: "2", Name: "Add"}, {ID: "3", Name: "Backfill"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Migrations(c, ms, MigrationsOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "1 applied") || !exact(tt, "2 pending") || strings.Count(strings.Join(tt.Texts(), "|"), "Next") != 1 {
		t.Errorf("migrations: %q", tt.Texts())
	}
	panics(t, `"3" is applied after a pending one`, func(c *ui.Context) {
		Migrations(c, []Migration{ms[0], ms[1], {ID: "3", Name: "Late", Applied: true, At: "Oct 9"}}, MigrationsOptions{})
	})
}

func TestQueuesDepth(t *testing.T) {
	qs := []Queue{{Name: "emails", Depth: 1500, Workers: 1, Oldest: 4 * time.Minute}, {Name: "hooks", Depth: 3}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Queues(c, qs, QueuesOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "1,503 waiting") || !exact(tt, "1 worker") || !exact(tt, "0 workers") || !exact(tt, "oldest 4m 0s") {
		t.Errorf("queues: %q", tt.Texts())
	}
	panics(t, "total depth overflows", func(c *ui.Context) {
		Queues(c, []Queue{{Name: "a", Depth: math.MaxInt64}, {Name: "b", Depth: 1}}, QueuesOptions{})
	})
}

func TestCronCountsDown(t *testing.T) {
	var start time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if start.IsZero() {
			start = c.Now()
		}
		jobs := []CronJob{{Name: "sync", Schedule: "*/5 * * * *", Next: start.Add(2*time.Minute + 14*time.Second)},
			{Name: "prune", Schedule: "@hourly", Paused: true},
			{Name: "far", Schedule: "@yearly", Next: time.Date(2500, 1, 1, 0, 0, 0, 0, time.UTC)}}
		Cron(c, jobs, CronOptions{})
	}), 500, 400)
	tt.Frame()
	if !exact(tt, "in 2m 14s") || !exact(tt, "Paused") {
		t.Fatalf("cron: %q", tt.Texts())
	}
	tt.Advance(1500 * time.Millisecond)
	if !exact(tt, "in 2m 13s") {
		t.Errorf("after 1.5 s: %q", tt.Texts())
	}
	tt.Advance(3 * time.Minute)
	if !exact(tt, "Due now") {
		t.Errorf("past due: %q", tt.Texts())
	}
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, "in -") {
			t.Errorf("far job reads %q", s)
		}
	}
}

func TestWebhooksRetry(t *testing.T) {
	ds := []Delivery{{ID: "a", Event: "order.created", URL: "https://x", Status: 200, Attempts: 1},
		{ID: "b", Event: "payment.failed", URL: "https://y", Attempts: 3}}
	retry := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Webhooks(c, ds, &retry, WebhooksOptions{}) }), 500, 300)
	settles(t, tt)
	if !exact(tt, "1 failing") || !exact(tt, "Timeout") || !exact(tt, "1 attempt") || !exact(tt, "3 attempts") {
		t.Errorf("webhooks: %q", tt.Texts())
	}
	if _, ok := tt.Find("Retry order.created"); ok {
		t.Error("a delivered hook offers retry")
	}
	before := baseline()
	click(t, tt, "Retry payment.failed")
	cues(t, before, sound.Send)
	if retry != "b" {
		t.Errorf("retry = %q", retry)
	}
	panics(t, "a status of 0 or 100-599", func(c *ui.Context) {
		Webhooks(c, []Delivery{{ID: "c", Event: "e", URL: "u", Status: 42, Attempts: 1}}, &retry, WebhooksOptions{})
	})
}

func TestSecretsReveal(t *testing.T) {
	secrets := []Secret{{Name: "API_KEY", Value: "abc123"}, {Name: "DB_URL", Value: "postgres://x"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Secrets(c, secrets, SecretsOptions{}) }), 500, 300)
	settles(t, tt)
	for _, s := range tt.Texts() {
		if strings.Contains(s, "abc123") || strings.Contains(s, "postgres") {
			t.Fatalf("hidden secret shown: %q", s)
		}
	}
	before := baseline()
	click(t, tt, "Reveal API_KEY")
	cues(t, before, sound.Toggle)
	settles(t, tt)
	if !exact(tt, "abc123") || strings.Contains(strings.Join(tt.Texts(), "|"), "postgres://x") {
		t.Errorf("reveal one: %q", tt.Texts())
	}
	secrets[0], secrets[1] = secrets[1], secrets[0]
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "abc123") || strings.Contains(strings.Join(tt.Texts(), "|"), "postgres://x") {
		t.Errorf("reveal followed position, not name: %q", tt.Texts())
	}
	click(t, tt, "Hide API_KEY")
	settles(t, tt)
	if strings.Contains(strings.Join(tt.Texts(), "|"), "abc123") {
		t.Error("hide left the value shown")
	}
}

func TestRateLimitRounding(t *testing.T) {
	cases := []struct {
		used, limit int64
		want        []string
	}{
		{999, 1000, []string{"99%", "1 left"}},
		{1000, 1000, []string{"100%", "Limited"}},
		{math.MaxInt64, 1, []string{"100%", "Limited"}},
		{0, 10, []string{"0%", "10 left"}},
	}
	for _, k := range cases {
		tt := ui.NewTester(padded(func(c *ui.Context) { RateLimit(c, k.used, k.limit, time.Time{}, RateLimitOptions{}) }), 420, 300)
		settles(t, tt)
		for _, w := range k.want {
			if !exact(tt, w) {
				t.Errorf("%d of %d: missing %q in %q", k.used, k.limit, w, tt.Texts())
			}
		}
	}
	var start time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if start.IsZero() {
			start = c.Now()
		}
		RateLimit(c, 5, 10, start.Add(3*time.Second), RateLimitOptions{})
	}), 420, 300)
	tt.Frame()
	if !exact(tt, "Resets in 3s") {
		t.Errorf("reset: %q", tt.Texts())
	}
	tt.Advance(4 * time.Second)
	if !exact(tt, "Window reset") {
		t.Errorf("after reset: %q", tt.Texts())
	}
}

func TestErrorsGroups(t *testing.T) {
	groups := []ErrorGroup{{ID: "a", Type: "TypeError", Message: "x", Count: 1, Users: 1, Trend: []float64{1, 3}},
		{ID: "b", Type: "RangeError", Message: "y", Count: 1200, Users: 40, Resolved: true}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Errors(c, groups, &selected, ErrorsOptions{}) }), 520, 300)
	settles(t, tt)
	if !exact(tt, "1,201 events") || !exact(tt, "1 user") || !exact(tt, "40 users") || !exact(tt, "Resolved") {
		t.Errorf("errors: %q", tt.Texts())
	}
	swapPress(t, tt, "TypeError: x", func() { groups[0], groups[1] = groups[1], groups[0] })
	if selected == "b" {
		t.Error("a press on a picked b after the swap")
	}
	panics(t, "total count overflows", func(c *ui.Context) {
		Errors(c, []ErrorGroup{{ID: "a", Type: "T", Count: math.MaxInt64}, {ID: "b", Type: "T", Count: 1}}, &selected, ErrorsOptions{})
	})
}

func TestInspectorProps(t *testing.T) {
	tree := Node{ID: "app", Tag: "App", Children: []Node{{ID: "btn", Tag: "Button", Props: []Attr{{Key: "label", Value: `"Pay"`}}}, {ID: "box", Tag: "Box"}}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Inspector(c, tree, &selected, InspectorOptions{}) }), 480, 400)
	settles(t, tt)
	if !exact(tt, "Pick an element to see its props") {
		t.Errorf("no hint: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Button btn")
	cues(t, before, sound.Select)
	settles(t, tt)
	if selected != "btn" || !exact(tt, `"Pay"`) {
		t.Errorf("props: %q", tt.Texts())
	}
	click(t, tt, "Box box")
	settles(t, tt)
	if !exact(tt, "No props") {
		t.Errorf("no props: %q", tt.Texts())
	}
	panics(t, "repeated", func(c *ui.Context) {
		Inspector(c, Node{ID: "a", Tag: "A", Children: []Node{{ID: "a", Tag: "B"}}}, &selected, InspectorOptions{})
	})
}

func TestPingSummary(t *testing.T) {
	ms := time.Millisecond
	hosts := []PingHost{
		{Name: "api", Seq: 4, Samples: []PingSample{{RTT: 10 * ms}, {RTT: 30 * ms}, {Lost: true}, {RTT: 20 * ms}}},
		{Name: "down", Seq: 2, Samples: []PingSample{{Lost: true}, {Lost: true}}},
		{Name: "huge", Seq: 2, Samples: []PingSample{{RTT: math.MaxInt64}, {RTT: math.MaxInt64 - 1}}},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { Ping(c, hosts, PingOptions{}) }), 480, 400)
	settles(t, tt)
	if !exact(tt, "20 ms avg") || !exact(tt, "25% loss") || !exact(tt, "No replies") {
		t.Errorf("ping: %q", tt.Texts())
	}
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, "-") {
			t.Errorf("huge average reads %q", s)
		}
	}
	panics(t, "a count at least that", func(c *ui.Context) {
		Ping(c, []PingHost{{Name: "x", Seq: 0, Samples: []PingSample{{}}}}, PingOptions{})
	})
}

func TestBuildLogFolds(t *testing.T) {
	steps := []BuildStep{{Name: "Checkout", State: StepDone, Duration: time.Second, Lines: []string{"cloned"}},
		{Name: "Lint", State: StepDone, Duration: time.Second},
		{Name: "Test", State: StepFailed, Duration: 2 * time.Second, Lines: []string{"--- FAIL: TestTotals"}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { BuildLog(c, steps, BuildLogOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "--- FAIL: TestTotals") || exact(tt, "cloned") || !exact(tt, "Failed") || !exact(tt, "4 s") {
		t.Fatalf("failed step not open first: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Checkout")
	cues(t, before, sound.Toggle)
	settles(t, tt)
	if !exact(tt, "cloned") {
		t.Errorf("Checkout did not open: %q", tt.Texts())
	}
	// Lint has no output, so clicking it toggles nothing.
	_ = tt.Click("Lint")
	tt.Frame()
	cues(t, before, sound.Toggle)
	panics(t, "total time overflows", func(c *ui.Context) {
		BuildLog(c, []BuildStep{{Name: "a", Duration: math.MaxInt64}, {Name: "b", Duration: 1}}, BuildLogOptions{})
	})
}

func TestCacheRate(t *testing.T) {
	cases := []struct {
		s    CacheStats
		want string
	}{
		{CacheStats{Hits: 99999, Misses: 1, Capacity: 1}, "99.9%"},
		{CacheStats{Hits: 1, Misses: 99999, Capacity: 1}, "0.1%"},
		{CacheStats{Hits: 3, Misses: 1, Capacity: 1}, "75.0%"},
		{CacheStats{Capacity: 1}, "No requests"},
	}
	for _, k := range cases {
		tt := ui.NewTester(padded(func(c *ui.Context) { Cache(c, k.s, CacheOptions{}) }), 460, 400)
		settles(t, tt)
		if !exact(tt, k.want) {
			t.Errorf("%+v: want %q in %q", k.s, k.want, tt.Texts())
		}
	}
	panics(t, "size within a positive capacity", func(c *ui.Context) { Cache(c, CacheStats{Size: 2, Capacity: 1}, CacheOptions{}) })
}

func TestTogglesReportAndFit(t *testing.T) {
	stamp := "2026-10-25T10:35:16.123456789-04:00 (America/New_York)"
	secrets := []Secret{{Name: "API_KEY", Value: "abc123", Updated: stamp}}
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Secrets(c, secrets, SecretsOptions{}).Changed() {
			changes++
		}
	}), 320, 300)
	settles(t, tt)
	name, _ := tt.Find("API_KEY")
	eye, _ := tt.Find("Reveal API_KEY")
	when, _ := tt.Find(stamp)
	if name.W <= 0 || name.X+name.W > eye.X || when.X+when.W > eye.X {
		t.Errorf("name %v and time %v crowd the reveal button %v", name, when, eye)
	}
	click(t, tt, "Reveal API_KEY")
	if changes != 1 {
		t.Errorf("reveal reported %d changes", changes)
	}
	steps := []BuildStep{{Name: "Checkout", State: StepDone, Lines: []string{"cloned"}}}
	changes = 0
	tt = ui.NewTester(padded(func(c *ui.Context) {
		if BuildLog(c, steps, BuildLogOptions{}).Changed() {
			changes++
		}
	}), 500, 300)
	settles(t, tt)
	click(t, tt, "Checkout")
	if changes != 1 {
		t.Errorf("a fold reported %d changes", changes)
	}
	selected := ""
	panics(t, `"a" trend has a negative count`, func(c *ui.Context) {
		Errors(c, []ErrorGroup{{ID: "a", Type: "T", Count: 1, Trend: []float64{-1, 2}}}, &selected, ErrorsOptions{})
	})
}
