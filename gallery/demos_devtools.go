package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/devtools"
	"github.com/egoist/mygo/ui"
)

var devtoolsEntries = []Entry{
	{ID: 1266, Name: "LogView", Category: "devtools", Summary: "Leveled logs with search and follow; cue Toggle.", Demo: demoDevLogView},
	{ID: 1267, Name: "LevelFilter", Category: "devtools", Summary: "Debug to error chips; cue Toggle.", Demo: demoDevLevelFilter},
	{ID: 1268, Name: "Metrics", Category: "devtools", Summary: "Metric charts in a grid.", Demo: demoDevMetrics},
	{ID: 1269, Name: "Traces", Category: "devtools", Summary: "Span waterfall bars.", Demo: demoDevTraces},
	{ID: 1270, Name: "SpanDetail", Category: "devtools", Summary: "Attributes of a span.", Demo: demoDevSpanDetail},
	{ID: 1271, Name: "Requests", Category: "devtools", Summary: "Requests with method, status, time; cue Select.", Demo: demoDevRequests},
	{ID: 1272, Name: "RequestDetail", Category: "devtools", Summary: "Headers and body tabs.", Demo: demoDevRequestDetail},
	{ID: 1273, Name: "FlameGraph", Category: "devtools", Summary: "Stacked frames by width; cue Select.", Demo: demoDevFlameGraph},
	{ID: 1274, Name: "Memory", Category: "devtools", Summary: "Heap over time.", Demo: demoDevMemory},
	{ID: 1275, Name: "CPU", Category: "devtools", Summary: "Cores with load bars.", Demo: demoDevCPU},
	{ID: 1276, Name: "FPS", Category: "devtools", Summary: "Frame rate meter.", Demo: demoDevFPS},
	{ID: 1277, Name: "Events", Category: "devtools", Summary: "Event stream with filter.", Demo: demoDevEvents},
	{ID: 1278, Name: "StateTree", Category: "devtools", Summary: "App state tree; changes flash.", Demo: demoDevStateTree},
	{ID: 1279, Name: "Flags", Category: "devtools", Summary: "Feature flags with toggles; cue Toggle.", Demo: demoDevFlags},
	{ID: 1280, Name: "EnvBadge", Category: "devtools", Summary: "Env pill colored by stage.", Demo: demoDevEnvBadge},
	{ID: 1281, Name: "Deploy", Category: "devtools", Summary: "Deploy card with steps.", Demo: demoDevDeploy},
	{ID: 1282, Name: "Deploys", Category: "devtools", Summary: "History with rollback; cue Warning.", Demo: demoDevDeploys},
	{ID: 1283, Name: "Pods", Category: "devtools", Summary: "Containers with status.", Demo: demoDevPods},
}

func demoDevLogView(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	query, follow := textState(root, "query", ""), boolState(root, "follow", true)
	lines := []devtools.LogLine{
		{ID: "1", Time: "09:12:01", Level: devtools.LevelInfo, Message: "server listening on :8080", Source: "main.go:42"},
		{ID: "2", Time: "09:12:03", Level: devtools.LevelDebug, Message: "loaded 128 products from cache", Source: "catalog.go:77"},
		{ID: "3", Time: "09:12:09", Level: devtools.LevelInfo, Message: "POST /checkout 201 in 84 ms", Source: "http.go:19"},
		{ID: "4", Time: "09:12:15", Level: devtools.LevelWarn, Message: "stripe webhook retried (attempt 2)", Source: "billing.go:131"},
		{ID: "5", Time: "09:12:21", Level: devtools.LevelError, Message: "inventory sync failed: timeout after 5 s", Source: "sync.go:58"},
		{ID: "6", Time: "09:12:24", Level: devtools.LevelInfo, Message: "GET /products 200 in 12 ms", Source: "http.go:19"},
	}
	root.Children(func() {
		if devtools.LogView(c, lines, query, follow, devtools.LogViewOptions{Height: 200}).Changed() {
			log(fmt.Sprintf("follow %v queryRunes %d", *follow, len([]rune(*query))))
		}
	})
}

func demoDevLevelFilter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	on := ui.Local(root, "on", func() [4]bool { return [4]bool{false, true, true, true} })
	root.Children(func() {
		if devtools.LevelFilter(c, on, [4]int{214, 1380, 37, 4}, devtools.LevelFilterOptions{}).Changed() {
			log(fmt.Sprint("levels ", *on))
		}
	})
}

func demoDevMetrics(c *ui.Context, log func(string)) {
	metrics := []devtools.Metric{
		{Name: "Requests per second", Unit: "req/s", Points: []float64{182, 196, 240, 228, 265, 301, 288}},
		{Name: "p95 latency", Unit: "ms", Points: []float64{142, 138, 151, 129, 118, 121, 112}, Lower: true},
		{Name: "Error rate", Unit: "%", Points: []float64{0.4, 0.3, 0.5, 0.9, 0.6, 0.4, 0.3}, Lower: true},
		{Name: "Checkout conversion", Unit: "%", Points: []float64{2.1, 2.3, 2.2, 2.6, 2.8, 2.7, 3.1}},
	}
	ui.Column(c).Width(500).Children(func() { devtools.Metrics(c, metrics, devtools.MetricsOptions{}) })
}

func demoTraceSpans() []devtools.Span {
	ms := time.Millisecond
	return []devtools.Span{
		{ID: "a", Name: "POST /checkout", Service: "gateway", Duration: 248 * ms},
		{ID: "b", Parent: "a", Name: "auth.verify", Service: "auth", Start: 4 * ms, Duration: 18 * ms},
		{ID: "c", Parent: "a", Name: "cart.load", Service: "kiln-api", Start: 24 * ms, Duration: 41 * ms},
		{ID: "d", Parent: "c", Name: "SELECT cart_items", Service: "postgres", Start: 28 * ms, Duration: 22 * ms},
		{ID: "e", Parent: "a", Name: "stripe.charge", Service: "billing", Start: 70 * ms, Duration: 151 * ms},
		{ID: "f", Parent: "a", Name: "email.receipt", Service: "kiln-api", Start: 224 * ms, Duration: 20 * ms, Error: true},
	}
}

func demoDevTraces(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected := textState(root, "span", "e")
	root.Children(func() {
		if devtools.Traces(c, demoTraceSpans(), selected, devtools.TracesOptions{Title: "Trace 4f2a9c"}).Changed() {
			log("span " + *selected)
		}
	})
}

func demoDevSpanDetail(c *ui.Context, log func(string)) {
	span := demoTraceSpans()[4]
	attrs := []devtools.Attr{{Key: "http.method", Value: "POST"}, {Key: "http.url", Value: "https://api.stripe.com/v1/charges"},
		{Key: "stripe.amount", Value: "4800"}, {Key: "stripe.currency", Value: "usd"}, {Key: "peer.service", Value: "stripe"}}
	ui.Column(c).Width(460).Children(func() { devtools.SpanDetail(c, span, attrs, devtools.SpanDetailOptions{}) })
}

func demoDevRequests(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected := textState(root, "req", "")
	ms := time.Millisecond
	reqs := []devtools.Request{
		{ID: "1", Method: "GET", Path: "/api/products?page=2", Status: 200, Time: 42 * ms, Size: 18_432},
		{ID: "2", Method: "POST", Path: "/api/cart/items", Status: 201, Time: 88 * ms, Size: 612},
		{ID: "3", Method: "GET", Path: "/api/shipping/rates", Status: 304, Time: 9 * ms, Size: 0},
		{ID: "4", Method: "PATCH", Path: "/api/cart/items/42", Status: 422, Time: 31 * ms, Size: 204},
		{ID: "5", Method: "POST", Path: "/api/checkout", Status: 502, Time: 1240 * ms, Size: 88},
	}
	root.Children(func() {
		if devtools.Requests(c, reqs, selected, devtools.RequestsOptions{}).Changed() {
			log("request " + *selected)
		}
	})
}

func demoDevRequestDetail(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(500)
	tab := intState(root, "tab", 0)
	req := devtools.RequestInfo{Method: "POST", URL: "/api/cart/items", Status: 201, Time: 88 * time.Millisecond,
		Headers: []devtools.Header{{Name: "content-type", Value: "application/json"}, {Name: "x-request-id", Value: "7c1e2b90"},
			{Name: "cache-control", Value: "no-store"}, {Name: "content-length", Value: "612"}},
		Body: "{\n  \"id\": 42,\n  \"product\": \"Speckled mug\",\n  \"quantity\": 2,\n  \"price\": 24.00\n}"}
	root.Children(func() {
		if devtools.RequestDetail(c, req, tab, devtools.RequestDetailOptions{}).Changed() {
			log("tab " + fmt.Sprint(*tab))
		}
	})
}

func demoDevFlameGraph(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected := textState(root, "frame", "")
	frame := devtools.Frame{Name: "main", Value: 1000, Children: []devtools.Frame{
		{Name: "http.Serve", Value: 720, Children: []devtools.Frame{
			{Name: "checkout", Value: 410, Children: []devtools.Frame{{Name: "json.Encode", Value: 120}, {Name: "stripe.Charge", Value: 240}}},
			{Name: "products", Value: 230, Children: []devtools.Frame{{Name: "sql.Query", Value: 160}, {Name: "json.Encode", Value: 50}}},
		}},
		{Name: "runtime.gc", Value: 180},
	}}
	root.Children(func() {
		if devtools.FlameGraph(c, frame, selected, devtools.FlameGraphOptions{}).Changed() {
			log("frame " + *selected)
		}
	})
}

func demoDevMemory(c *ui.Context, log func(string)) {
	samples := []float64{212, 238, 261, 249, 287, 318, 296, 334, 372, 341, 389, 402}
	ui.Column(c).Width(460).Children(func() { devtools.Memory(c, samples, 512, devtools.MemoryOptions{}) })
}

func demoDevCPU(c *ui.Context, log func(string)) {
	cores := []float64{0.42, 0.67, 0.91, 0.38, 0.55, 0.74, 0.21, 0.48}
	ui.Column(c).Width(440).Children(func() { devtools.CPU(c, cores, devtools.CPUOptions{}) })
}

func demoDevFPS(c *ui.Context, log func(string)) {
	history := []float64{60, 60, 59, 60, 58, 41, 33, 52, 60, 60, 60, 59, 60, 60, 47, 60, 60, 60, 59, 60}
	ui.Column(c).Width(420).Children(func() { devtools.FPS(c, history, devtools.FPSOptions{}) })
}

func demoDevEvents(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	filter := textState(root, "filter", "")
	events := []devtools.Event{
		{ID: "1", Time: "09:14:02.118", Name: "cart:add", Detail: "Speckled mug × 2"},
		{ID: "2", Time: "09:14:05.440", Name: "route:change", Detail: "/cart → /checkout"},
		{ID: "3", Time: "09:14:06.002", Name: "payment:start", Detail: "card ending 4242"},
		{ID: "4", Time: "09:14:07.901", Name: "payment:success", Detail: "order 10482"},
		{ID: "5", Time: "09:14:08.010", Name: "route:change", Detail: "/checkout → /thanks"},
	}
	root.Children(func() {
		if devtools.Events(c, events, filter, devtools.EventsOptions{}).Changed() {
			log("filterRunes " + fmt.Sprint(len([]rune(*filter))))
		}
	})
}

func demoDevStateTree(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	count := intState(root, "count", 2)
	tree := devtools.StateNode{Key: "store", Children: []devtools.StateNode{
		{Key: "user", Children: []devtools.StateNode{{Key: "name", Value: `"Maya Chen"`}, {Key: "signedIn", Value: "true"}}},
		{Key: "cart", Children: []devtools.StateNode{{Key: "items", Value: fmt.Sprint(*count)}, {Key: "coupon", Value: "null"},
			{Key: "total", Value: fmt.Sprintf("%.2f", float64(*count)*24)}}},
		{Key: "theme", Value: `"system"`},
	}}
	root.Children(func() {
		if devtools.StateTree(c, tree, devtools.StateTreeOptions{}).Changed() {
			log("tree toggled")
		}
		if ui.Button(c, "Add a mug").Clicked() {
			*count++
			log("items " + fmt.Sprint(*count))
		}
	})
}

func demoDevFlags(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	flags := ui.Local(root, "flags", func() []devtools.Flag {
		return []devtools.Flag{
			{Key: "new-checkout", Description: "One-page checkout with Apple Pay", On: true, Rollout: 25},
			{Key: "gift-cards", Description: "Sell and redeem gift cards", On: true, Rollout: 100},
			{Key: "ai-descriptions", Description: "Draft product copy from photos", Rollout: 5},
			{Key: "dark-storefront", Description: "Dark theme for the public shop", Rollout: 0},
		}
	})
	root.Children(func() {
		if devtools.Flags(c, flags, devtools.FlagsOptions{}).Changed() {
			on := 0
			for _, f := range *flags {
				if f.On {
					on++
				}
			}
			log("on " + fmt.Sprint(on))
		}
	})
}

func demoDevEnvBadge(c *ui.Context, log func(string)) {
	ui.Row(c).Width(440).Wrap().Gap(8).Children(func() {
		for _, s := range []devtools.Stage{devtools.StageDevelopment, devtools.StagePreview, devtools.StageStaging, devtools.StageProduction} {
			devtools.EnvBadge(c, s, devtools.EnvBadgeOptions{})
		}
	})
}

func demoDevDeploy(c *ui.Context, log func(string)) {
	s := time.Second
	d := devtools.DeployInfo{Commit: "a41f9c2", Message: "Add gift card checkout", Author: "Daniel Okafor", Branch: "main",
		Stage: devtools.StageProduction, Steps: []devtools.DeployStep{
			{Name: "Install dependencies", State: devtools.StepDone, Duration: 18 * s},
			{Name: "Run tests", State: devtools.StepDone, Duration: 74 * s},
			{Name: "Build image", State: devtools.StepDone, Duration: 41 * s},
			{Name: "Migrate database", State: devtools.StepDone, Duration: 6 * s},
			{Name: "Roll out to 3 regions", State: devtools.StepDone, Duration: 52 * s},
		}}
	ui.Column(c).Width(480).Children(func() { devtools.Deploy(c, d, devtools.DeployOptions{}) })
}

func demoDevDeploys(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	rollback := textState(root, "rollback", "")
	deploys := []devtools.DeployRecord{
		{ID: "d5", Version: "v2.14.0", Commit: "a41f9c2", Author: "Daniel Okafor", When: "12 min ago", Current: true},
		{ID: "d4", Version: "v2.13.2", Commit: "9be03d1", Author: "Priya Nair", When: "Yesterday", Failed: true},
		{ID: "d3", Version: "v2.13.1", Commit: "71c5e8a", Author: "Maya Chen", When: "Oct 22"},
		{ID: "d2", Version: "v2.13.0", Commit: "4d2a6f0", Author: "Daniel Okafor", When: "Oct 19"},
	}
	root.Children(func() {
		if devtools.Deploys(c, deploys, rollback, devtools.DeploysOptions{}).Changed() {
			log("rollback " + *rollback)
		}
	})
}

func demoDevPods(c *ui.Context, log func(string)) {
	pods := []devtools.Pod{
		{Name: "kiln-api-7d9f8-x2kq4", Status: devtools.PodRunning, Age: "2d", CPU: 0.34, Memory: 0.58},
		{Name: "kiln-api-7d9f8-m8tz1", Status: devtools.PodRunning, Age: "2d", Restarts: 1, CPU: 0.41, Memory: 0.61},
		{Name: "kiln-worker-5c6b-qp7d", Status: devtools.PodRunning, Age: "6h", CPU: 0.78, Memory: 0.44},
		{Name: "kiln-search-0", Status: devtools.PodFailed, Age: "14m", Restarts: 4, CPU: 0.02, Memory: 0.93},
	}
	ui.Column(c).Width(520).Children(func() { devtools.Pods(c, pods, devtools.PodsOptions{}) })
}
