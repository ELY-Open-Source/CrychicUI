package devtools

import (
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func logLines() []LogLine {
	return []LogLine{{ID: "1", Time: "09:12:01", Level: LevelInfo, Message: "server listening"},
		{ID: "2", Time: "09:12:21", Level: LevelError, Message: "inventory sync failed"}}
}

func TestLogViewSearchAndFollow(t *testing.T) {
	query, follow := "", true
	tt := ui.NewTester(padded(func(c *ui.Context) { LogView(c, logLines(), &query, &follow, LogViewOptions{Height: 200}) }), 600, 400)
	settles(t, tt)
	if !exact(tt, "ERROR") || !exact(tt, "2") {
		t.Fatalf("lines: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Follow")
	if follow {
		t.Error("follow did not turn off")
	}
	cues(t, before, sound.Toggle)
	click(t, tt, "Search logs")
	tt.Type("zzz")
	tt.Frame()
	if !exact(tt, "No matching logs") || !exact(tt, "0") {
		t.Errorf("empty search: %q", tt.Texts())
	}
}

func TestLevelFilterAndMetrics(t *testing.T) {
	on := [4]bool{false, true, true, true}
	tt := ui.NewTester(padded(func(c *ui.Context) { LevelFilter(c, &on, [4]int{2, 3, 4, 5}, LevelFilterOptions{}) }), 600, 200)
	settles(t, tt)
	click(t, tt, "Debug")
	if !on[0] {
		t.Error("Debug did not turn on")
	}
	metrics := []Metric{{Name: "Requests", Unit: "req/s", Points: []float64{100, 150}},
		{Name: "Queue", Unit: "jobs", Points: []float64{0, 5}}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Metrics(c, metrics, MetricsOptions{}) }), 600, 400)
	settles(t, tt)
	if !exact(tt, "+50%") || !exact(tt, "150") {
		t.Errorf("metrics: %q", tt.Texts())
	}
	for _, s := range tt.Texts() {
		if s == "+Inf%" || s == "NaN%" {
			t.Errorf("a zero start showed %q", s)
		}
	}
	panics(t, "two points", func(c *ui.Context) { Metrics(c, []Metric{{Name: "A", Points: []float64{1}}}, MetricsOptions{}) })
}

func TestTracesPickAndGuards(t *testing.T) {
	ms := time.Millisecond
	spans := []Span{{ID: "a", Name: "POST /checkout", Service: "gateway", Duration: 100 * ms},
		{ID: "b", Parent: "a", Name: "auth.verify", Service: "auth", Start: 4 * ms, Duration: 20 * ms},
		{ID: "c", Parent: "a", Name: "cart.load", Service: "api", Start: 30 * ms, Duration: 40 * ms}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Traces(c, spans, &selected, TracesOptions{}) }), 600, 400)
	settles(t, tt)
	before := baseline()
	swapUnderPress(t, tt, "cart.load", func() { spans[1], spans[2] = spans[2], spans[1] })
	if selected == "b" {
		t.Error("Traces picked the span that slid in")
	}
	click(t, tt, "auth.verify")
	if selected != "b" {
		t.Errorf("selected %q", selected)
	}
	cues(t, before, sound.Select)
	panics(t, "before its parent", func(c *ui.Context) {
		Traces(c, []Span{{ID: "b", Parent: "a", Name: "x", Service: "s", Duration: time.Millisecond}}, new(string), TracesOptions{})
	})
	panics(t, "past the largest duration", func(c *ui.Context) {
		Traces(c, []Span{{ID: "a", Name: "x", Service: "s", Start: math.MaxInt64, Duration: 1}}, new(string), TracesOptions{})
	})
}

func TestRequestsAndDetail(t *testing.T) {
	reqs := []Request{{ID: "1", Method: "GET", Path: "/api/products", Status: 200, Time: 40 * time.Millisecond, Size: 1024},
		{ID: "2", Method: "POST", Path: "/api/checkout", Status: 502, Time: 1240 * time.Millisecond, Size: 88}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Requests(c, reqs, &selected, RequestsOptions{}) }), 700, 400)
	settles(t, tt)
	if !exact(tt, "1.24 s") || !exact(tt, "40 ms") {
		t.Errorf("times: %q", tt.Texts())
	}
	click(t, tt, "POST /api/checkout")
	if selected != "2" {
		t.Errorf("selected %q", selected)
	}
	panics(t, "status 100-599", func(c *ui.Context) {
		Requests(c, []Request{{ID: "1", Method: "GET", Path: "/", Status: 700}}, new(string), RequestsOptions{})
	})
	tab := 0
	info := RequestInfo{Method: "POST", URL: "/api/cart", Status: 201, Headers: []Header{{Name: "content-type", Value: "json"}}}
	tt = ui.NewTester(padded(func(c *ui.Context) { RequestDetail(c, info, &tab, RequestDetailOptions{}) }), 600, 400)
	settles(t, tt)
	click(t, tt, "Body")
	settles(t, tt)
	if tab != 1 || !exact(tt, "No body") || !exact(tt, "Headers 1") {
		t.Errorf("tab %d texts %q", tab, tt.Texts())
	}
}

func TestFlameGraphPicksByName(t *testing.T) {
	frame := Frame{Name: "main", Value: 1000, Children: []Frame{
		{Name: "checkout", Value: 400, Children: []Frame{{Name: "json.Encode", Value: 120}}},
		{Name: "products", Value: 300, Children: []Frame{{Name: "json.Encode", Value: 50}}}}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { FlameGraph(c, frame, &selected, FlameGraphOptions{}) }), 700, 400)
	settles(t, tt)
	click(t, tt, "checkout")
	settles(t, tt)
	if selected != "checkout" || !exact(tt, "checkout · 40.0% of samples") {
		t.Errorf("selected %q texts %q", selected, tt.Texts())
	}
	selected = "json.Encode"
	tt.Frame()
	if !exact(tt, "json.Encode · 17.0% of samples") {
		t.Errorf("same name sums: %q", tt.Texts())
	}
	nested := Frame{Name: "main", Value: 1000, Children: []Frame{{Name: "render", Value: 500, Children: []Frame{{Name: "render", Value: 200}}}}}
	pick := "render"
	tt = ui.NewTester(padded(func(c *ui.Context) { FlameGraph(c, nested, &pick, FlameGraphOptions{}) }), 700, 400)
	settles(t, tt)
	if !exact(tt, "render · 50.0% of samples") {
		t.Errorf("recursion counted twice: %q", tt.Texts())
	}
	panics(t, "children over its value", func(c *ui.Context) {
		FlameGraph(c, Frame{Name: "a", Value: 1, Children: []Frame{{Name: "b", Value: 2}}}, new(string), FlameGraphOptions{})
	})
}

func TestGauges(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Memory(c, []float64{200, 402, 380}, 512, MemoryOptions{}) }), 600, 400)
	settles(t, tt)
	if !exact(tt, "Peak 402 MB") || !exact(tt, "380 MB") {
		t.Errorf("memory: %q", tt.Texts())
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { CPU(c, []float64{0.5, 0.6}, CPUOptions{}) }), 600, 400)
	settles(t, tt)
	if !exact(tt, "55% average") {
		t.Errorf("cpu: %q", tt.Texts())
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { FPS(c, []float64{60, 30, 59, 1e300}, FPSOptions{}) }), 600, 300)
	settles(t, tt)
	if !exact(tt, "1 dropped") || !exact(tt, "1e+300") {
		t.Errorf("fps: %q", tt.Texts())
	}
	panics(t, "outside 0..1", func(c *ui.Context) { CPU(c, []float64{1.5}, CPUOptions{}) })
	panics(t, "positive limit", func(c *ui.Context) { Memory(c, []float64{1, 2}, 0, MemoryOptions{}) })
	panics(t, "is negative", func(c *ui.Context) { FPS(c, []float64{-1, 2}, FPSOptions{}) })
}

func TestStateTreeFlashAndFold(t *testing.T) {
	items := "2"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		StateTree(c, StateNode{Key: "store", Children: []StateNode{{Key: "cart", Children: []StateNode{{Key: "items", Value: items}}}}}, StateTreeOptions{})
	}), 600, 400)
	settles(t, tt)
	if !exact(tt, "cart") || exact(tt, "items") {
		t.Fatalf("root open, cart shut: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "cart")
	settles(t, tt)
	if !exact(tt, "items") {
		t.Fatalf("cart did not open: %q", tt.Texts())
	}
	cues(t, before, sound.Toggle)
	items = "3"
	tt.Frame()
	if !tt.FramePending() {
		t.Error("a changed value did not flash")
	}
}

func TestFlagsDeploysPods(t *testing.T) {
	flags := []Flag{{Key: "new-checkout", On: true, Rollout: 25}, {Key: "gift-cards", Rollout: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Flags(c, &flags, FlagsOptions{}) }), 600, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "gift-cards", func() { flags[0], flags[1] = flags[1], flags[0] })
	for _, f := range flags {
		if f.Key == "new-checkout" && !f.On {
			t.Error("Flags toggled the row that slid in")
		}
	}
	panics(t, "outside 0..100", func(c *ui.Context) { Flags(c, &[]Flag{{Key: "a", Rollout: 101}}, FlagsOptions{}) })
	panics(t, "Development to Production", func(c *ui.Context) { EnvBadge(c, Stage(9), EnvBadgeOptions{}) })

	steps := []DeployStep{{Name: "Build", State: StepDone}, {Name: "Ship", State: StepRunning}}
	tt = ui.NewTester(padded(func(c *ui.Context) {
		Deploy(c, DeployInfo{Commit: "a41f9c2", Message: "Gift cards", Author: "Daniel Okafor", Branch: "main", Steps: steps}, DeployOptions{})
	}), 600, 500)
	tt.Advance(2 * time.Second)
	if !exact(tt, "Building") || !tt.FramePending() {
		t.Errorf("running deploy: %q", tt.Texts())
	}

	rollback := ""
	deploys := []DeployRecord{{ID: "d2", Version: "v2", Commit: "b", Author: "Maya", When: "now", Current: true},
		{ID: "d1", Version: "v1", Commit: "a", Author: "Maya", When: "then"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Deploys(c, deploys, &rollback, DeploysOptions{}) }), 600, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Roll back to v1")
	if rollback != "d1" {
		t.Errorf("rollback %q", rollback)
	}
	cues(t, before, sound.Warning)
	panics(t, "one current deploy", func(c *ui.Context) {
		Deploys(c, []DeployRecord{{ID: "a", Version: "v", Commit: "c", Author: "M", When: "w"}}, new(string), DeploysOptions{})
	})

	pods := []Pod{{Name: "api-1", Age: "2d", CPU: 0.3, Memory: 0.5}, {Name: "search-0", Status: PodFailed, Age: "1m", Restarts: 4},
		{Name: "worker-2", Status: PodPending, Age: "5s"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Pods(c, pods, PodsOptions{}) }), 700, 400)
	settles(t, tt)
	tt.Advance(2 * time.Second)
	if !exact(tt, "1 of 3 ready") || !exact(tt, "4 restarts") {
		t.Errorf("pods: %q", tt.Texts())
	}
}

func TestDevtoolsRepeatsPanic(t *testing.T) {
	cases := map[string]func(c *ui.Context){
		"LogView": func(c *ui.Context) {
			l := LogLine{ID: "a", Time: "t", Message: "m"}
			LogView(c, []LogLine{l, l}, new(string), new(bool), LogViewOptions{})
		},
		"Requests": func(c *ui.Context) {
			r := Request{ID: "a", Method: "GET", Path: "/", Status: 200}
			Requests(c, []Request{r, r}, new(string), RequestsOptions{})
		},
		"Events": func(c *ui.Context) {
			e := Event{ID: "a", Time: "t", Name: "n"}
			Events(c, []Event{e, e}, new(string), EventsOptions{})
		},
		"Pods": func(c *ui.Context) {
			p := Pod{Name: "a", Age: "1d"}
			Pods(c, []Pod{p, p}, PodsOptions{})
		},
		"StateTree": func(c *ui.Context) {
			StateTree(c, StateNode{Key: "r", Children: []StateNode{{Key: "a"}, {Key: "a"}}}, StateTreeOptions{})
		},
	}
	for name, view := range cases {
		t.Run(name, func(t *testing.T) { panics(t, "devtools: "+name, view) })
	}
}
