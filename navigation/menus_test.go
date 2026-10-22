package navigation

import (
	"fmt"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var heads = []string{"Intro", "Setup", "Usage", "Notes"}

// body is a tall section with a unique marker text.
func body(c *ui.Context) func(i int) {
	return func(i int) { ui.Box(c).Height(260).Children(func() { ui.Text(c, fmt.Sprintf("Body %d", i)) }) }
}

func TestReading(t *testing.T) {
	tops := []float32{0, 300, 600, 900}
	for _, c := range []struct {
		y, max float32
		want   int
	}{{0, 1000, 0}, {295, 1000, 1}, {650, 1000, 2}, {1000, 1000, 3}} {
		if got := reading(tops, &ui.ScrollState{Y: c.y, MaxY: c.max}); got != c.want {
			t.Errorf("reading at %v = %d, want %d", c.y, got, c.want)
		}
	}
	if got := reading([]float32{0, 0, 0}, &ui.ScrollState{}); got != 0 {
		t.Errorf("unmeasured reads %d", got)
	}
}

func TestTOC(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { TOC(c, "Guide", heads, body(c), TOCOptions{}).Height(300) }, 700, 320)
	settles(t, tt, "TOC")
	top := find(t, tt, "Body 0").Y
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	settles(t, tt, "jump")
	if y := find(t, tt, "Body 2").Y; y < top-1 || y > top+40 {
		t.Errorf("Usage's body at %v, the top at %v", y, top)
	}
	cues(t, before, sound.Tap)
	title := find(t, tt, "On this page")
	r := find(t, tt, "Body 0")
	tt.Scroll(r.X+20, 150, 0, 5000)
	settles(t, tt, "to the end")
	img, accent := tt.Image(), ui.LightTheme().Accent
	lit := -1
	for y := int(title.Y + title.H); y < 320; y++ {
		if p := img.RGBAAt(int(title.X)+1, y); p.R == accent.R && p.G == accent.G && p.B == accent.B {
			lit = y
			break
		}
	}
	if lit < int(title.Y+title.H)+90 {
		t.Errorf("at the end the marker is at %d, under the title at %v", lit, title.Y+title.H)
	}
}

func TestAnchorNav(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { AnchorNav(c, "Guide", heads, body(c)).Height(300) }, 600, 320)
	settles(t, tt, "AnchorNav")
	top := find(t, tt, "Body 0").Y
	before := baseline()
	click(t, tt, "Notes")
	settles(t, tt, "jump")
	if y := find(t, tt, "Body 3").Y; y > top+120 {
		t.Errorf("Notes' body at %v, the top at %v", y, top)
	}
	cues(t, before, sound.Tap)
}

func TestMenus(t *testing.T) {
	gs := []LinkGroup{{"Product", []string{"Pricing", "Changelog"}}, {"Company", []string{"About", "Careers"}}}
	picked, open, changed := [2]int{-1, -1}, false, 0
	tm := ui.NewTester(func(c *ui.Context) {
		if MegaMenu(c, "Explore", gs, &picked, &open).Changed() {
			changed++
		}
	}, 700, 400)
	before := baseline()
	click(t, tm, "Explore")
	if !open || changed != 1 || !tm.HasText("Company") {
		t.Fatalf("open %v, changed %d, %q", open, changed, tm.Texts())
	}
	tm.Key(0, ui.KeyEscape)
	tm.Advance(time.Second)
	if open || changed != 2 {
		t.Errorf("Escape: open %v, changed %d", open, changed)
	}
	click(t, tm, "Explore")
	click(t, tm, "Careers")
	tm.Advance(time.Second)
	if picked != [2]int{1, 1} || open || tm.HasText("Careers") || changed != 4 {
		t.Errorf("picked %v, open %v, changed %d", picked, open, changed)
	}
	cues(t, before, sound.Open, sound.Open, sound.Tap)

	picked, changed = [2]int{-1, -1}, 0
	opened := -1
	ts := ui.NewTester(func(c *ui.Context) {
		if MenuStrip(c, "Menu", gs, &picked, &opened).Changed() {
			changed++
		}
	}, 700, 400)
	before = baseline()
	click(t, ts, "Company")
	if opened != 1 || changed != 1 {
		t.Errorf("strip open %d, changed %d", opened, changed)
	}
	click(t, ts, "About")
	ts.Advance(time.Second)
	if picked != [2]int{1, 0} || opened != -1 || ts.HasText("About") || changed != 2 {
		t.Errorf("strip picked %v, open %d, changed %d", picked, opened, changed)
	}
	cues(t, before, sound.Open, sound.Tap)

	link, shown, changed := -1, false, 0
	td := ui.NewTester(func(c *ui.Context) {
		if NavDropdown(c, "Docs", []string{"Guides", "API"}, &link, &shown).Changed() {
			changed++
		}
	}, 500, 300)
	before = baseline()
	click(t, td, "Docs")
	click(t, td, "API")
	if link != 1 || shown || changed != 2 {
		t.Errorf("dropdown picked %d, open %v, changed %d", link, shown, changed)
	}
	cues(t, before, sound.Open, sound.Tap)
	panics(t, "bad pick", func() { render(func(c *ui.Context) { p, o := [2]int{0, 5}, false; MegaMenu(c, "x", gs, &p, &o) }) })
	panics(t, "bad open menu", func() { render(func(c *ui.Context) { p, o := [2]int{-1, -1}, 2; MenuStrip(c, "x", gs, &p, &o) }) })
}
func TestCommander(t *testing.T) {
	commands := []string{"Open file", "Save file", "Save all", "Close tab"}
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) { Commander(c, "Commands", commands, &chosen, CommanderOptions{}) }, 500, 400)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("save")
	tt.Frame()
	if tt.HasText("Open file") || !tt.HasText("Save all") {
		t.Fatalf("filter %q", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 2 {
		t.Errorf("chose %d", chosen)
	}
	cues(t, before, sound.Select)
	tt.Command("selectAll")
	tt.Type("zzz")
	tt.Frame()
	if !tt.HasText("No matches") {
		t.Errorf("empty %q", tt.Texts())
	}
}

func TestCommanderParkedPointer(t *testing.T) {
	commands := []string{"Open file", "Save file", "Save all"}
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) { Commander(c, "Commands", commands, &chosen, CommanderOptions{}) }, 500, 400)
	r := find(t, tt, "Open file")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 1 {
		t.Errorf("a parked pointer kept the highlight: chose %d", chosen)
	}
}

func TestCommanderOneHighlight(t *testing.T) {
	commands := []string{"Open file", "Save file", "Save all"}
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) { Commander(c, "Commands", commands, &chosen, CommanderOptions{}) }, 500, 400)
	settles(t, tt, "Commander")
	r := find(t, tt, "Save all")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "hover")
	img, bg := tt.Image(), ui.LightTheme().Background
	var lit []string
	for _, cmd := range commands {
		r := find(t, tt, cmd)
		if p := img.RGBAAt(int(r.X)+3, int(r.Y+r.H/2)); p.R != bg.R || p.G != bg.G || p.B != bg.B {
			lit = append(lit, cmd)
		}
	}
	if len(lit) != 1 || lit[0] != "Save all" {
		t.Errorf("highlighted %q", lit)
	}
}
