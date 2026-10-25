package devtools

import (
	"fmt"
	"image"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// plain paints the page the card's color.
func plain(view func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) {
		ui.Column(c).FillWidth().FillHeight().Padding(10).Background(c.Theme().Background).Children(func() { view(c) })
	}
}

// inked counts pixels in r unlike the page.
func inked(img *image.RGBA, r image.Rectangle) int {
	bg := img.RGBAAt(2, 2)
	n := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			p := img.RGBAAt(x, y)
			d := max(absDiff(p.R, bg.R), absDiff(p.G, bg.G), absDiff(p.B, bg.B))
			if d > 8 {
				n++
			}
		}
	}
	return n
}

func absDiff(a, b uint8) int { return int(max(a, b) - min(a, b)) }

func TestFollowResumes(t *testing.T) {
	lines := make([]LogLine, 40)
	for i := range lines {
		lines[i] = LogLine{ID: fmt.Sprint(i), Time: "10:00", Message: fmt.Sprintf("job %02d done", i), Level: LevelInfo}
	}
	query, follow := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) { LogView(c, lines, &query, &follow, LogViewOptions{}) }), 600, 420)
	settles(t, tt)
	if r, ok := tt.Find("job 39 done"); ok && r.H > 0 && r.Y+r.H <= 420 {
		t.Fatalf("the last line shows before Follow: %v", r)
	}
	click(t, tt, "Follow")
	settles(t, tt)
	r, ok := tt.Find("job 39 done")
	if !ok || r.H <= 0 || r.Y+r.H > 420 {
		t.Errorf("Follow did not bring the end into view: %v %v", r, ok)
	}
}

func TestFlameReorderKeepsPick(t *testing.T) {
	root := Frame{Name: "main", Value: 100, Children: []Frame{{Name: "checkout", Value: 50}, {Name: "products", Value: 50}}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { FlameGraph(c, root, &selected, FlameGraphOptions{}) }), 600, 300)
	settles(t, tt)
	r, ok := tt.Find("checkout")
	if !ok {
		t.Fatal("no checkout frame")
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	root.Children[0], root.Children[1] = root.Children[1], root.Children[0]
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
	if selected == "products" {
		t.Error("a press on checkout picked products after the swap")
	}
}

func TestTracesAxisMatchesBars(t *testing.T) {
	spans := []Span{
		{ID: "a", Name: "checkout", Service: "api", Duration: 800 * time.Millisecond},
		{ID: "b", Parent: "a", Name: "charge", Service: "pay", Duration: 400 * time.Millisecond},
	}
	selected := ""
	tt := ui.NewTester(plain(func(c *ui.Context) { Traces(c, spans, &selected, TracesOptions{}) }), 600, 300)
	settles(t, tt)
	tick, ok := tt.Find("400 ms")
	if !ok {
		t.Fatal("no 400 ms tick")
	}
	start, _ := tt.Find("0 ms")
	row, _ := tt.Find("charge")
	img := tt.Image()
	y := int(row.Y + row.H/2)
	bar := img.RGBAAt(int(tick.X)-10, y)
	lo, hi := int(tick.X)-10, int(tick.X)-10
	for img.RGBAAt(lo-1, y) == bar {
		lo--
	}
	for img.RGBAAt(hi+1, y) == bar {
		hi++
	}
	// Labels sit the same distance past their ticks.
	if at := tick.X - (start.X - float32(lo)); math.Abs(float64(float32(hi+1)-at)) > 2 {
		t.Errorf("the 400 ms bar runs %d..%d, its tick sits at %v", lo, hi, at)
	}
}

func TestTraceTicksDoNotOverflow(t *testing.T) {
	spans := []Span{{ID: "a", Name: "batch", Service: "jobs", Duration: math.MaxInt64 / 2}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Traces(c, spans, &selected, TracesOptions{}) }), 600, 300)
	settles(t, tt)
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, "-") {
			t.Errorf("tick %q is negative", s)
		}
	}
}

func TestMetricsChangeBounds(t *testing.T) {
	panics(t, `"Latency" change overflows`, func(c *ui.Context) {
		Metrics(c, []Metric{{Name: "Latency", Points: []float64{1e-306, 10}}}, MetricsOptions{})
	})
	metrics := []Metric{{Name: "Signups", Points: []float64{1, 1.234567e100}}, {Name: "Queue", Points: []float64{3, 4}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Metrics(c, metrics, MetricsOptions{}) }), 320, 300)
	settles(t, tt)
	// Names fill their tiles, so they mark each tile's edges.
	left, _ := tt.Find("Signups")
	right, _ := tt.Find("Queue")
	for _, s := range tt.Texts() {
		r, _ := tt.Find(s)
		tile := left
		if r.X >= right.X {
			tile = right
		}
		if r.X+r.W > tile.X+tile.W+1 {
			t.Errorf("%q ends at %v past its tile at %v", s, r.X+r.W, tile.X+tile.W)
		}
	}
}

func TestChartsScaleHugeValues(t *testing.T) {
	huge := 1.7e308
	tt := ui.NewTester(plain(func(c *ui.Context) { Memory(c, []float64{0.3 * huge, 0.6 * huge}, huge, MemoryOptions{}) }), 400, 320)
	settles(t, tt)
	heap, _ := tt.Find("Heap")
	if n := inked(tt.Image(), image.Rect(20, int(heap.Y)+80, 380, 320)); n < 4000 {
		t.Errorf("Memory chart drew %d pixels at huge values", n)
	}
	tt = ui.NewTester(plain(func(c *ui.Context) { FPS(c, []float64{0.5 * huge, huge}, FPSOptions{}) }), 320, 160)
	settles(t, tt)
	if n := inked(tt.Image(), image.Rect(160, 0, 305, 160)); n < 2000 {
		t.Errorf("FPS chart drew %d pixels at huge values", n)
	}
}

func TestCPUManyCores(t *testing.T) {
	cores := make([]float64, 32)
	for i := range cores {
		cores[i] = float64(i+50) / 100
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { CPU(c, cores, CPUOptions{}) }), 320, 300)
	settles(t, tt)
	a, _ := tt.Find("50%")
	b, _ := tt.Find("51%")
	if b.X < a.X+a.W {
		t.Errorf("core labels overlap: %v %v", a, b)
	}
	tt.Scroll(a.X+a.W/2, a.Y+a.H/2+40, 2000, 0)
	settles(t, tt)
	last, ok := tt.Find("81%")
	if !ok || last.W <= 0 || last.X+last.W > 320 {
		t.Errorf("the last core is out of reach: %v %v", last, ok)
	}
}

func TestLongTextKeepsNeighbors(t *testing.T) {
	long := strings.Repeat("checkout_", 10)
	state := StateNode{Key: "cart", Children: []StateNode{{Key: long, Value: "42"}}}
	filter := ""
	events := []Event{{ID: "1", Time: "10:00", Name: long, Detail: "paid"}}
	query, follow := "", false
	lines := []LogLine{{ID: "1", Time: "10:00", Message: "charged card", Source: "internal/" + long + ".go:88", Level: LevelInfo}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		StateTree(c, state, StateTreeOptions{})
		Events(c, events, &filter, EventsOptions{})
		LogView(c, lines, &query, &follow, LogViewOptions{Height: 60})
	}), 320, 700)
	settles(t, tt)
	for _, s := range []string{"42", "paid", "charged card"} {
		r, ok := tt.Find(s)
		if !ok || r.W <= 0 || r.X+r.W > 320 {
			t.Errorf("%q lost its room: %v %v", s, r, ok)
		}
	}
}
