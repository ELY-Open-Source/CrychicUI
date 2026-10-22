package data

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// band counts pixels near k in rows y0 to y1.
func band(tt *ui.Tester, y0, y1 float32, k ui.Color) int {
	img := tt.Image()
	n := 0
	for y := int(y0); y < int(y1); y++ {
		for x := range img.Bounds().Dx() {
			if p := img.RGBAAt(x, y); near(p.R, k.R) && near(p.G, k.G) && near(p.B, k.B) {
				n++
			}
		}
	}
	return n
}

var who = []Column{{Title: "Name"}, {Title: "Age"}}

func TestEmptyTable(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if _, act := EmptyTable(c, "People", who, "No people yet", "Add person", EmptyTableOptions{}); act.Clicked() {
			clicked++
		}
	}, 500, 300)
	if !tt.HasText("Name") || !tt.HasText("No people yet") || !tt.FramePending() {
		t.Fatalf("shows %q, arriving %v", tt.Texts(), tt.FramePending())
	}
	settles(t, tt, "message")
	before := baseline()
	click(t, tt, "Add person")
	if clicked != 1 {
		t.Errorf("clicked %d", clicked)
	}
	cues(t, before, sound.Tap)
	panics(t, "no action", func() {
		render(func(c *ui.Context) { EmptyTable(c, "x", who, "none", "", EmptyTableOptions{}) })
	})
	panics(t, "bad icon", func() {
		render(func(c *ui.Context) { EmptyTable(c, "x", who, "none", "Add", EmptyTableOptions{Icon: "nope"}) })
	})
}

func TestLoadingTable(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { LoadingTable(c, "People", who, 3) }, 500, 300)
	tt.Advance(time.Second)
	if !tt.HasText("Age") || !tt.FramePending() {
		t.Errorf("shows %q, sheen %v", tt.Texts(), tt.FramePending())
	}
	panics(t, "no rows", func() { render(func(c *ui.Context) { LoadingTable(c, "x", who, 0) }) })
}

func TestStatList(t *testing.T) {
	visits := int64(120)
	tt := ui.NewTester(func(c *ui.Context) {
		StatList(c, "Today", []Stat{{Label: "Visits", Value: visits}, {Label: "Sales", Value: 4, Unit: "orders"}})
	}, 500, 300)
	settles(t, tt, "stats")
	if !tt.HasText("120") || !tt.HasText("orders") {
		t.Fatalf("shows %q", tt.Texts())
	}
	visits = 121
	tt.Frame()
	if !tt.FramePending() {
		t.Error("digits did not roll")
	}
	settles(t, tt, "roll")
	if !tt.HasText("121") {
		t.Errorf("shows %q", tt.Texts())
	}
	panics(t, "twins", func() {
		render(func(c *ui.Context) { StatList(c, "x", []Stat{{Label: "a"}, {Label: "a"}}) })
	})
}

func TestRankList(t *testing.T) {
	th := ui.LightTheme()
	tt := ui.NewTester(func(c *ui.Context) {
		RankList(c, "Top pages", []Ranked{{"Home", 10}, {"Docs", 5}}, RankListOptions{})
	}, 600, 200)
	home, docs := find(t, tt, "Home"), find(t, tt, "Docs")
	// At 150 ms the row shows; the bar still grows.
	tt.Advance(150 * time.Millisecond)
	early := band(tt, home.Y-4, home.Y+home.H+4, th.Accent)
	settles(t, tt, "bars")
	h := band(tt, home.Y-4, home.Y+home.H+4, th.Accent)
	if early == 0 || float64(early) > 0.8*float64(h) {
		t.Errorf("bar at 100 ms has %d of %d pixels; it should grow", early, h)
	}
	d := band(tt, docs.Y-4, docs.Y+docs.H+4, th.Accent)
	if r := float64(h) / float64(d); d == 0 || r < 1.6 || r > 2.4 {
		t.Errorf("bars %d and %d, want about 2:1", h, d)
	}
	for _, v := range []float64{-1, math.NaN(), math.Inf(1)} {
		panics(t, fmt.Sprint(v), func() {
			render(func(c *ui.Context) { RankList(c, "x", []Ranked{{"a", v}}, RankListOptions{}) })
		})
	}
}

func TestTimelineList(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		TimelineList(c, "History", []Moment{{"Mar 2", "Founded", "Two people"}, {"Jun 9", "Launched", ""}})
	}, 500, 300)
	if !tt.HasText("Launched") || !tt.FramePending() {
		t.Fatalf("shows %q, arriving %v", tt.Texts(), tt.FramePending())
	}
	settles(t, tt, "stagger")
	panics(t, "no date", func() {
		render(func(c *ui.Context) { TimelineList(c, "x", []Moment{{Title: "a"}}) })
	})
}

func TestLogList(t *testing.T) {
	var lines []LogLine
	for i := range 200 {
		lines = append(lines, LogLine{Level: Level(i % 4), Time: "12:00", Text: fmt.Sprintf("line %d", i)})
	}
	tt := ui.NewTester(func(c *ui.Context) { LogList(c, "Server", lines, LogListOptions{}) }, 600, 400)
	settles(t, tt, "log")
	if !tt.HasText("line 199") || tt.HasText("line 0") {
		t.Fatalf("not at the end: %q", tt.Texts())
	}
	lines = append(lines, LogLine{Level: LevelError, Time: "12:01", Text: "disk full"})
	tt.Frame()
	th := ui.LightTheme()
	if tinted(tt, find(t, tt, "disk full"), th.Text) {
		t.Error("new line showed at full ink at once")
	}
	settles(t, tt, "new line")
	if !tinted(tt, find(t, tt, "disk full"), th.Text) || !tt.HasText("ERROR") {
		t.Errorf("new line did not fade in: %q", tt.Texts())
	}
	panics(t, "level", func() {
		render(func(c *ui.Context) { LogList(c, "x", []LogLine{{Level: 9}}, LogListOptions{}) })
	})
}

func TestLogListFadesIn(t *testing.T) {
	th := ui.LightTheme()
	lines := []LogLine{{LevelInfo, "12:00", "booted"}}
	tt := ui.NewTester(func(c *ui.Context) { LogList(c, "Server", lines, LogListOptions{}) }, 600, 400)
	settles(t, tt, "log")
	lines = append(lines, LogLine{LevelWarn, "12:01", "slow disk"})
	tt.Frame()
	if tinted(tt, find(t, tt, "slow disk"), th.Text) {
		t.Error("new line showed at full ink at once")
	}
	settles(t, tt, "fade")
	if !tinted(tt, find(t, tt, "slow disk"), th.Text) {
		t.Error("new line never showed")
	}
}

func TestFeed(t *testing.T) {
	posts := []Post{{"1", "Hello", "first"}}
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { Feed(c, "News", posts) }, 500, 400)
	settles(t, tt, "feed")
	cues(t, before)
	posts = append([]Post{{"2", "Update", "second"}}, posts...)
	tt.Frame()
	if !tt.HasText("Update") || !tt.FramePending() {
		t.Fatalf("shows %q, sliding %v", tt.Texts(), tt.FramePending())
	}
	settles(t, tt, "slide")
	tt.Frame()
	cues(t, before, sound.Notify)
	panics(t, "twins", func() {
		render(func(c *ui.Context) { Feed(c, "x", []Post{{"1", "a", ""}, {"1", "b", ""}}) })
	})
}
