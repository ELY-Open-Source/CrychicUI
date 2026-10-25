package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/devtools"
	"github.com/egoist/mygo/ui"
)

var devtoolsEntries2 = []Entry{
	{ID: 1284, Name: "ServiceMap", Category: "devtools", Summary: "Services and calls.", Demo: demoDevServiceMap},
	{ID: 1285, Name: "Health", Category: "devtools", Summary: "Endpoints with checks.", Demo: demoDevHealth},
	{ID: 1286, Name: "Tables", Category: "devtools", Summary: "Database tables with rows.", Demo: demoDevTables},
	{ID: 1287, Name: "Query", Category: "devtools", Summary: "SQL editor with run; cue Tap.", Demo: demoDevQuery},
	{ID: 1288, Name: "QueryResult", Category: "devtools", Summary: "Result grid with time.", Demo: demoDevQueryResult},
	{ID: 1289, Name: "Schema", Category: "devtools", Summary: "Columns and types.", Demo: demoDevSchema},
	{ID: 1290, Name: "Migrations", Category: "devtools", Summary: "Applied and pending.", Demo: demoDevMigrations},
	{ID: 1291, Name: "Queues", Category: "devtools", Summary: "Queues with depth.", Demo: demoDevQueues},
	{ID: 1292, Name: "Cron", Category: "devtools", Summary: "Jobs with next run.", Demo: demoDevCron},
	{ID: 1293, Name: "Webhooks", Category: "devtools", Summary: "Deliveries with status and retry; cue Send.", Demo: demoDevWebhooks},
	{ID: 1294, Name: "Secrets", Category: "devtools", Summary: "Masked secrets with reveal; cue Toggle.", Demo: demoDevSecrets},
	{ID: 1295, Name: "RateLimit", Category: "devtools", Summary: "Requests against limit.", Demo: demoDevRateLimit},
	{ID: 1296, Name: "Errors", Category: "devtools", Summary: "Grouped errors with counts.", Demo: demoDevErrors},
	{ID: 1297, Name: "Inspector", Category: "devtools", Summary: "Element tree and props.", Demo: demoDevInspector},
	{ID: 1298, Name: "Ping", Category: "devtools", Summary: "Latency dots updating.", Demo: demoDevPing},
	{ID: 1299, Name: "BuildLog", Category: "devtools", Summary: "Build steps with times.", Demo: demoDevBuildLog},
	{ID: 1300, Name: "Cache", Category: "devtools", Summary: "Hit rate and size.", Demo: demoDevCache},
}

func demoDevServiceMap(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected := textState(root, "service", "kiln-api")
	services := []devtools.Service{
		{Name: "storefront"}, {Name: "gateway"}, {Name: "auth"}, {Name: "kiln-api"},
		{Name: "postgres"}, {Name: "redis"}, {Name: "billing", State: devtools.HealthDegraded},
	}
	calls := []devtools.Call{
		{From: "storefront", To: "gateway", Rate: "420/s"}, {From: "gateway", To: "auth", Rate: "410/s"},
		{From: "gateway", To: "kiln-api", Rate: "395/s"}, {From: "auth", To: "postgres", Rate: "88/s"},
		{From: "kiln-api", To: "postgres", Rate: "1.2k/s"}, {From: "kiln-api", To: "redis", Rate: "3.4k/s"},
		{From: "kiln-api", To: "billing", Rate: "18/s", Failing: true},
	}
	root.Children(func() {
		if devtools.ServiceMap(c, services, calls, selected, devtools.ServiceMapOptions{}).Changed() {
			log("service " + *selected)
		}
	})
}

func demoDevHealth(c *ui.Context, log func(string)) {
	runs := func(fails ...int) []bool {
		h := make([]bool, 30)
		for i := range h {
			h[i] = true
		}
		for _, f := range fails {
			h[f] = false
		}
		return h
	}
	ms := time.Millisecond
	checks := []devtools.Check{
		{Endpoint: "GET /healthz", Latency: 12 * ms, History: runs()},
		{Endpoint: "GET /api/products", Latency: 48 * ms, History: runs(17)},
		{Endpoint: "POST /api/checkout", State: devtools.HealthDegraded, Latency: 812 * ms, History: runs(21, 25, 26, 29)},
		{Endpoint: "GET /api/search", State: devtools.HealthDown, Latency: 5 * time.Second, History: runs(24, 25, 26, 27, 28, 29)},
	}
	ui.Column(c).Width(500).Children(func() { devtools.Health(c, checks, devtools.HealthOptions{}) })
}

func demoDevTables(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	selected := textState(root, "table", "")
	tables := []devtools.Table{
		{Name: "orders", Rows: 48_213, Size: 18_400_000}, {Name: "order_items", Rows: 131_902, Size: 42_100_000},
		{Name: "products", Rows: 412, Size: 2_300_000}, {Name: "customers", Rows: 9_870, Size: 6_800_000},
		{Name: "inventory_events", Rows: 1_204_551, Size: 210_000_000}, {Name: "settings", Rows: 1, Size: 8_192},
	}
	root.Children(func() {
		if devtools.Tables(c, tables, selected, devtools.TablesOptions{}).Changed() {
			log("table " + *selected)
		}
	})
}

func demoDevQuery(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(500)
	sql := textState(root, "sql", "SELECT id, customer, total\nFROM orders\nWHERE created_at > now() - interval '1 day'\nORDER BY total DESC\nLIMIT 5;")
	root.Children(func() {
		q, run := devtools.Query(c, sql, devtools.QueryOptions{Database: "kiln_prod"})
		if q.Changed() {
			log(fmt.Sprintf("sqlRunes %d", len([]rune(*sql))))
		}
		if run.Clicked() {
			log(fmt.Sprintf("run lines %d", strings.Count(*sql, "\n")+1))
		}
	})
}

func demoDevQueryResult(c *ui.Context, log func(string)) {
	columns := []string{"id", "customer", "total", "status", "coupon"}
	rows := [][]string{
		{"10482", "Maya Chen", "148.00", "paid", "NULL"},
		{"10479", "Daniel Okafor", "96.50", "paid", "SPRING20"},
		{"10481", "Priya Nair", "72.00", "pending", "NULL"},
		{"10477", "Lucas Martin", "64.00", "refunded", "NULL"},
		{"10480", "Sofia Rossi", "48.00", "paid", "WELCOME10"},
	}
	ui.Column(c).Width(520).Children(func() {
		devtools.QueryResult(c, columns, rows, 12400*time.Microsecond, devtools.QueryResultOptions{})
	})
}

func demoDevSchema(c *ui.Context, log func(string)) {
	columns := []devtools.Column{
		{Name: "id", Type: "bigserial", Primary: true},
		{Name: "customer_id", Type: "bigint", Ref: "customers.id"},
		{Name: "total", Type: "numeric(10,2)", Default: "0"},
		{Name: "status", Type: "text", Default: "'pending'"},
		{Name: "coupon", Type: "text", Nullable: true},
		{Name: "created_at", Type: "timestamptz", Default: "now()"},
	}
	ui.Column(c).Width(480).Children(func() { devtools.Schema(c, "orders", columns, devtools.SchemaOptions{}) })
}

func demoDevMigrations(c *ui.Context, log func(string)) {
	migrations := []devtools.Migration{
		{ID: "20261002_create_orders", Name: "Create orders", Applied: true, At: "Oct 2"},
		{ID: "20261009_add_coupons", Name: "Add coupons to orders", Applied: true, At: "Oct 9"},
		{ID: "20261016_index_orders_created", Name: "Index orders by date", Applied: true, At: "Oct 16"},
		{ID: "20261023_add_gift_cards", Name: "Add gift cards"},
		{ID: "20261024_backfill_gift_card_totals", Name: "Backfill gift card totals"},
	}
	ui.Column(c).Width(460).Children(func() { devtools.Migrations(c, migrations, devtools.MigrationsOptions{}) })
}

func demoDevQueues(c *ui.Context, log func(string)) {
	s := time.Second
	queues := []devtools.Queue{
		{Name: "emails", Depth: 12, Workers: 4, Oldest: 3 * s},
		{Name: "image-resize", Depth: 1_284, Workers: 2, Oldest: 4 * time.Minute},
		{Name: "inventory-sync", Workers: 1},
		{Name: "webhooks-out", Depth: 37, Oldest: 12 * time.Minute},
	}
	ui.Column(c).Width(480).Children(func() { devtools.Queues(c, queues, devtools.QueuesOptions{}) })
}

func demoDevCron(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	born := ui.Local(root, "born", func() time.Time { return c.Now() })
	jobs := []devtools.CronJob{
		{Name: "sync-inventory", Schedule: "*/5 * * * *", Next: born.Add(2*time.Minute + 14*time.Second)},
		{Name: "reindex-search", Schedule: "0 */6 * * *", Next: born.Add(41 * time.Minute), Failed: true},
		{Name: "nightly-backup", Schedule: "0 3 * * *", Next: born.Add(5*time.Hour + 12*time.Minute)},
		{Name: "send-digest", Schedule: "0 8 * * MON", Next: born.Add(50 * time.Hour)},
		{Name: "prune-sessions", Schedule: "@hourly", Paused: true},
	}
	root.Children(func() { devtools.Cron(c, jobs, devtools.CronOptions{}) })
}

func demoDevWebhooks(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(500)
	retry := textState(root, "retry", "")
	deliveries := ui.Local(root, "deliveries", func() []devtools.Delivery {
		return []devtools.Delivery{
			{ID: "w1", Event: "order.created", URL: "https://hooks.slack.com/services/kiln/orders", When: "2 min ago", Status: 200, Attempts: 1},
			{ID: "w2", Event: "payment.succeeded", URL: "https://erp.kiln.co/webhooks/payments", When: "6 min ago", Status: 500, Attempts: 3},
			{ID: "w3", Event: "inventory.low", URL: "https://ops.kiln.co/alerts", When: "18 min ago", Attempts: 2},
			{ID: "w4", Event: "order.shipped", URL: "https://hooks.slack.com/services/kiln/shipping", When: "1 h ago", Status: 204, Attempts: 1},
		}
	})
	root.Children(func() {
		if devtools.Webhooks(c, *deliveries, retry, devtools.WebhooksOptions{}).Changed() {
			for i, d := range *deliveries {
				if d.ID == *retry {
					(*deliveries)[i].Attempts++
					(*deliveries)[i].Status = 200
				}
			}
			log("retried " + *retry)
			*retry = ""
		}
	})
}

func demoDevSecrets(c *ui.Context, log func(string)) {
	secrets := []devtools.Secret{
		{Name: "STRIPE_SECRET_KEY", Value: "demo-stripe-key-7f3a91c2", Updated: "2 days ago"},
		{Name: "DATABASE_URL", Value: "postgres://kiln:demo@db.internal:5432/kiln_prod", Updated: "Oct 12"},
		{Name: "MAIL_API_KEY", Value: "demo-mail-key-41b8e0", Updated: "Sep 30"},
		{Name: "SESSION_SECRET", Value: "demo-session-secret-0c55", Updated: "Aug 4"},
	}
	ui.Column(c).Width(480).Children(func() {
		if devtools.Secrets(c, secrets, devtools.SecretsOptions{}).Changed() {
			log("secret toggled")
		}
	})
}

func demoDevRateLimit(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(8)
	born := ui.Local(root, "born", func() time.Time { return c.Now() })
	used := intState(root, "used", 842)
	root.Children(func() {
		devtools.RateLimit(c, int64(*used), 1000, born.Add(42*time.Second), devtools.RateLimitOptions{Title: "Storefront API"})
		if ui.Button(c, "Send 40 requests").Clicked() {
			*used += 40
			log(fmt.Sprintf("used %d", *used))
		}
	})
}

func demoDevErrors(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected := textState(root, "error", "")
	groups := []devtools.ErrorGroup{
		{ID: "e1", Type: "TypeError", Message: "Cannot read properties of undefined (reading 'price')", Count: 1_284, Users: 312,
			Last: "2 min ago", Trend: []float64{12, 18, 15, 40, 96, 120, 88}},
		{ID: "e2", Type: "NetworkError", Message: "Failed to fetch /api/cart", Count: 438, Users: 96, Last: "14 min ago",
			Trend: []float64{30, 28, 41, 37, 52, 44, 39}},
		{ID: "e3", Type: "PaymentDeclined", Message: "card_declined: insufficient_funds", Count: 57, Users: 41, Last: "1 h ago",
			Trend: []float64{6, 9, 7, 8, 10, 9, 8}},
		{ID: "e4", Type: "ValidationError", Message: "Coupon SPRING20 has expired", Count: 12, Users: 1, Last: "Yesterday", Resolved: true},
	}
	root.Children(func() {
		if devtools.Errors(c, groups, selected, devtools.ErrorsOptions{}).Changed() {
			log("error " + *selected)
		}
	})
}

func demoDevInspector(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	selected := textState(root, "node", "cart-btn")
	attr := func(kv ...string) []devtools.Attr {
		var out []devtools.Attr
		for i := 0; i < len(kv); i += 2 {
			out = append(out, devtools.Attr{Key: kv[i], Value: kv[i+1]})
		}
		return out
	}
	tree := devtools.Node{ID: "app", Tag: "App", Props: attr("theme", `"system"`), Children: []devtools.Node{
		{ID: "header", Tag: "Header", Props: attr("sticky", "true"), Children: []devtools.Node{
			{ID: "logo", Tag: "Logo", Props: attr("href", `"/"`)},
			{ID: "cart-btn", Tag: "CartButton", Props: attr("count", "2", "label", `"Cart"`, "open", "false")},
		}},
		{ID: "grid", Tag: "ProductGrid", Props: attr("columns", "3"), Children: []devtools.Node{
			{ID: "card-1", Tag: "ProductCard", Props: attr("name", `"Speckled mug"`, "price", "24")},
			{ID: "card-2", Tag: "ProductCard", Props: attr("name", `"Linen apron"`, "price", "38")},
		}},
		{ID: "footer", Tag: "Footer"},
	}}
	root.Children(func() {
		if devtools.Inspector(c, tree, selected, devtools.InspectorOptions{}).Changed() {
			log("node " + *selected)
		}
	})
}

func demoDevPing(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	seq := intState(root, "seq", 40)
	tick := ui.Local(root, "tick", func() time.Time { return c.Now() })
	if c.Now().Sub(*tick) >= time.Second {
		*tick = c.Now()
		*seq++
	}
	c.After(time.Second - c.Now().Sub(*tick))
	base := []time.Duration{24, 3, 11, 96}
	names := []string{"api.kiln.co", "db.internal", "cdn.kiln.co", "eu.kiln.co"}
	hosts := make([]devtools.PingHost, len(names))
	for h, name := range names {
		samples := make([]devtools.PingSample, 40)
		for i := range samples {
			k := *seq - 40 + i
			rtt := (base[h] + time.Duration((k*7+h*13)%29)) * time.Millisecond
			if k%17 == 5 {
				rtt += 140 * time.Millisecond
			}
			samples[i] = devtools.PingSample{RTT: rtt, Lost: (k+h*5)%31 == 0}
		}
		hosts[h] = devtools.PingHost{Name: name, Seq: *seq, Samples: samples}
	}
	root.Children(func() { devtools.Ping(c, hosts, devtools.PingOptions{}) })
}

func demoDevBuildLog(c *ui.Context, log func(string)) {
	s := time.Second
	steps := []devtools.BuildStep{
		{Name: "Checkout", State: devtools.StepDone, Duration: 2 * s, Lines: []string{"Cloning kiln-co/kiln-api@a41f9c2", "HEAD is now at a41f9c2 Add gift card checkout"}},
		{Name: "Install dependencies", State: devtools.StepDone, Duration: 18 * s, Lines: []string{"go: downloading 84 modules", "go mod verify: all modules verified"}},
		{Name: "Lint", State: devtools.StepDone, Duration: 6 * s},
		{Name: "Test", State: devtools.StepFailed, Duration: 74 * s, Lines: []string{
			"ok  \tkiln/catalog\t1.204s", "--- FAIL: TestCheckoutTotals (0.02s)",
			"    cart_test.go:88: total = 47.99, want 48.00", "FAIL\tkiln/cart\t0.412s"}},
		{Name: "Build image", State: devtools.StepPending},
	}
	ui.Column(c).Width(500).Children(func() {
		if devtools.BuildLog(c, steps, devtools.BuildLogOptions{Title: "Build #1842"}).Changed() {
			log("step toggled")
		}
	})
}

func demoDevCache(c *ui.Context, log func(string)) {
	stats := devtools.CacheStats{Hits: 1_842_310, Misses: 96_412, Keys: 48_213, Evictions: 1_204,
		Size: 412_000_000, Capacity: 512_000_000, History: []float64{0.91, 0.92, 0.9, 0.93, 0.94, 0.95, 0.95}}
	ui.Column(c).Width(460).Children(func() { devtools.Cache(c, stats, devtools.CacheOptions{Title: "Redis"}) })
}
