package data

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// drag moves the pointer from a to b, pressed.
func drag(tt *ui.Tester, a, b ui.Rect, dy float32) {
	x0, y0 := a.X+a.W/2, a.Y+a.H/2
	x1, y1 := b.X+b.W/2, b.Y+b.H/2+dy
	tt.Press(x0, y0)
	for i := 1; i <= 8; i++ {
		tt.Move(x0+(x1-x0)*float32(i)/8, y0+(y1-y0)*float32(i)/8)
		tt.Advance(16 * time.Millisecond)
	}
	tt.Release(x1, y1)
	tt.Frame()
}

func TestKanban(t *testing.T) {
	lanes := []Lane{{"To do", []string{"A", "B"}}, {"Done", []string{"C"}}}
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Kanban(c, "Board", &lanes).Changed() {
			changed++
		}
	}, 900, 500)
	settles(t, tt, "board")
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Alt, ui.KeyRight)
	tt.Frame()
	if !slices.Equal(lanes[1].Cards, []string{"A", "C"}) {
		t.Fatalf("Alt+Right: %v", lanes)
	}
	// The moved card keeps focus in its new lane.
	tt.Key(ui.Alt, ui.KeyLeft)
	tt.Frame()
	if !slices.Equal(lanes[0].Cards, []string{"A", "B"}) || !slices.Equal(lanes[1].Cards, []string{"C"}) {
		t.Fatalf("Alt+Left: %v", lanes)
	}
	settles(t, tt, "moves")
	drag(tt, find(t, tt, "B"), find(t, tt, "C"), 0)
	if !slices.Equal(lanes[1].Cards, []string{"B", "C"}) {
		t.Fatalf("drop on C: %v", lanes)
	}
	settles(t, tt, "drop")
	drag(tt, find(t, tt, "A"), find(t, tt, "C"), find(t, tt, "C").H)
	if len(lanes[0].Cards) != 0 || !slices.Equal(lanes[1].Cards, []string{"B", "C", "A"}) {
		t.Fatalf("drop on lane: %v", lanes)
	}
	if changed != 4 {
		t.Errorf("changed %d", changed)
	}
	cues(t, before, sound.Drop, sound.Drop, sound.Drop, sound.Drop)
	panics(t, "twin cards", func() {
		render(func(c *ui.Context) { ls := []Lane{{"a", []string{"x"}}, {"b", []string{"x"}}}; Kanban(c, "x", &ls) })
	})
}

func TestKanbanEdgesAndFlight(t *testing.T) {
	lanes := []Lane{{"To do", []string{"A"}}, {"Done", nil}}
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Kanban(c, "Board", &lanes).Changed() {
			changed++
		}
	}, 900, 500)
	settles(t, tt, "board")
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Key(ui.Alt, ui.KeyLeft)
	tt.Frame()
	if changed != 0 {
		t.Fatalf("edge moves reported %d changes: %v", changed, lanes)
	}
	cues(t, before)
	from := find(t, tt, "A").X
	tt.Key(ui.Alt, ui.KeyRight)
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	mid := find(t, tt, "A").X
	settles(t, tt, "flight")
	to := find(t, tt, "A").X
	if len(lanes[1].Cards) != 1 || mid <= from+1 || mid >= to-1 {
		t.Errorf("card at %v mid-flight, from %v to %v", mid, from, to)
	}
}

func TestMasterDetail(t *testing.T) {
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) {
		MasterDetail(c, "Mail", []string{"Alpha", "Beta"}, &chosen, func(i int) { ui.Text(c, "Body "+[]string{"A", "B"}[i]) }, MasterDetailOptions{})
	}, 700, 400)
	if !tt.HasText("Choose an item") {
		t.Fatalf("shows %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Beta")
	th := ui.LightTheme()
	if chosen != 1 || tinted(tt, find(t, tt, "Body B"), th.Text) {
		t.Fatalf("chosen %d, detail did not fade in: %q", chosen, tt.Texts())
	}
	settles(t, tt, "fade")
	if !tinted(tt, find(t, tt, "Body B"), th.Text) {
		t.Fatal("detail never showed")
	}
	click(t, tt, "Beta")
	cues(t, before, sound.Select)
	panics(t, "range", func() {
		render(func(c *ui.Context) {
			n := 2
			MasterDetail(c, "x", []string{"a"}, &n, func(int) {}, MasterDetailOptions{})
		})
	})
}

func TestAccordionList(t *testing.T) {
	open := -1
	tt := ui.NewTester(func(c *ui.Context) {
		AccordionList(c, "FAQ", []string{"Shipping", "Returns"}, &open, func(i int) { ui.Text(c, []string{"Two days", "Thirty days"}[i]) })
	}, 500, 400)
	before := baseline()
	click(t, tt, "Returns")
	if open != 1 || !tt.HasText("Thirty days") || !tt.FramePending() {
		t.Fatalf("open %d, %q", open, tt.Texts())
	}
	settles(t, tt, "open")
	click(t, tt, "Shipping")
	if open != 0 || !tt.HasText("Two days") {
		t.Fatalf("open %d", open)
	}
	settles(t, tt, "switch")
	if tt.HasText("Thirty days") {
		t.Error("two rows open")
	}
	click(t, tt, "Shipping")
	if open != -1 {
		t.Errorf("open %d", open)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle)
}

func TestTagCloud(t *testing.T) {
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) {
		TagCloud(c, "Topics", []Tag{{"go", 10}, {"rust", 1}}, &chosen)
	}, 500, 300)
	settles(t, tt, "cloud")
	if g, r := find(t, tt, "go"), find(t, tt, "rust"); g.H <= r.H*1.5 {
		t.Errorf("go %v not larger than rust %v", g.H, r.H)
	}
	before := baseline()
	click(t, tt, "rust")
	if chosen != 1 {
		t.Errorf("chosen %d", chosen)
	}
	cues(t, before, sound.Tap)
	panics(t, "NaN", func() {
		render(func(c *ui.Context) { n := -1; TagCloud(c, "x", []Tag{{"a", math.NaN()}}, &n) })
	})
}

func TestDataCards(t *testing.T) {
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) {
		DataCards(c, "Clients", []Entry{{"Acme", []Pair{{"City", "Oslo"}}}, {"Globex", nil}}, &chosen)
	}, 700, 400)
	settles(t, tt, "cards")
	if !tt.HasText("Oslo") {
		t.Fatalf("shows %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Globex")
	if chosen != 1 {
		t.Errorf("chosen %d", chosen)
	}
	cues(t, before, sound.Tap)
	panics(t, "twins", func() {
		render(func(c *ui.Context) { n := -1; DataCards(c, "x", []Entry{{Title: "a"}, {Title: "a"}}, &n) })
	})
}

func TestRecordViewAndAuditTrail(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		RecordView(c, "Order", []Section{{"Buyer", []Pair{{"Name", "Ada"}}}, {"Payment", []Pair{{"Card", "•• 42"}}}})
		AuditTrail(c, "Changes", []Audit{{"Grace", "Renamed the file", "2 min ago"}})
	}, 600, 500)
	if !tt.FramePending() {
		t.Error("nothing arrives")
	}
	settles(t, tt, "record")
	for _, s := range []string{"Buyer", "Ada", "•• 42", "Renamed the file", "2 min ago"} {
		if !tt.HasText(s) {
			t.Errorf("no %q in %q", s, tt.Texts())
		}
	}
	panics(t, "no who", func() { render(func(c *ui.Context) { AuditTrail(c, "x", []Audit{{What: "a"}}) }) })
	panics(t, "no sections", func() { render(func(c *ui.Context) { RecordView(c, "x", nil) }) })
}

func TestVersionList(t *testing.T) {
	current := 1
	versions := []Version{{"v1", "May 1", "First draft"}, {"v2", "May 3", ""}}
	tt := ui.NewTester(func(c *ui.Context) { VersionList(c, "History", versions, &current, VersionListOptions{}) }, 500, 300)
	settles(t, tt, "versions")
	before := baseline()
	if _, ok := tt.Find("Restore v2"); ok {
		t.Fatal("the current version offers restore")
	}
	click(t, tt, "Restore v1")
	if current != 0 {
		t.Fatalf("current %d", current)
	}
	tt.Frame()
	if _, ok := tt.Find("Restore v2"); !ok {
		t.Error("v2 does not offer restore")
	}
	cues(t, before, sound.Success)
	panics(t, "range", func() {
		render(func(c *ui.Context) { n := 2; VersionList(c, "x", versions, &n, VersionListOptions{}) })
	})
}
