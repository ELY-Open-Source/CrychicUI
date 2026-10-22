package navigation

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestQuickLinks(t *testing.T) {
	picked := -1
	tt := ui.NewTester(func(c *ui.Context) {
		QuickLinks(c, "Shortcuts", []NavItem{{"Home", "house"}, {"Docs", "book"}}, &picked)
	}, 500, 200)
	r := find(t, tt, "Docs")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("the tile did not lift")
	}
	before := baseline()
	click(t, tt, "Docs")
	if picked != 1 {
		t.Errorf("picked %d", picked)
	}
	cues(t, before, sound.Tap)
}

func TestFooterNav(t *testing.T) {
	picked := [2]int{-1, -1}
	gs := []LinkGroup{{"Product", []string{"Pricing"}}, {"Company", []string{"About", "Careers"}}}
	tt := ui.NewTester(func(c *ui.Context) { FooterNav(c, "Footer", gs, &picked) }, 600, 300)
	if !tt.FramePending() {
		t.Error("the columns did not arrive")
	}
	settles(t, tt, "FooterNav")
	before := baseline()
	click(t, tt, "Careers")
	if picked != [2]int{1, 1} {
		t.Errorf("picked %v", picked)
	}
	cues(t, before, sound.Tap)
}

func TestSitemap(t *testing.T) {
	picked := ""
	tt := ui.NewTester(func(c *ui.Context) { Sitemap(c, "Pages", site, &picked) }, 400, 400)
	if !tt.FramePending() || !tt.HasText("Deploy") {
		t.Fatalf("rows did not fade in: %q", tt.Texts())
	}
	settles(t, tt, "Sitemap")
	before := baseline()
	click(t, tt, "Deploy")
	click(t, tt, "Deploy")
	if picked != "Guides/Deploy" {
		t.Errorf("picked %q", picked)
	}
	cues(t, before, sound.Tap, sound.Tap)
}

func TestKeyHints(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		KeyHints(c, "Shortcuts", []KeyHint{{[]string{"⌘", "K"}, "Search"}, {[]string{"Esc"}, "Close"}})
	}, 500, 100)
	if !tt.HasText("Search") || !tt.HasText("Esc") || !tt.FramePending() {
		t.Errorf("hints %q, arriving %v", tt.Texts(), tt.FramePending())
	}
	panics(t, "no keys", func() { render(func(c *ui.Context) { KeyHints(c, "x", []KeyHint{{nil, "Search"}}) }) })
}

func TestRelatedNav(t *testing.T) {
	picked := -1
	tt := ui.NewTester(func(c *ui.Context) { RelatedNav(c, "Related", []string{"Install", "Deploy"}, &picked) }, 400, 200)
	r := find(t, tt, "Deploy")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("the arrow did not nudge")
	}
	before := baseline()
	click(t, tt, "Deploy")
	if picked != 1 {
		t.Errorf("picked %d", picked)
	}
	cues(t, before, sound.Tap)
}

// pages routes /inbox and /sent to titled text.
func pages(rt *ui.Route, c *ui.Context) {
	switch {
	case rt.Match("/inbox"):
		rt.Title("Inbox")
		ui.Text(c, "Inbox page")
	case rt.Match("/sent"):
		rt.Title("Sent")
		ui.Text(c, "Sent page")
	}
}

var mail = []Page{{"Inbox", "/inbox"}, {"Sent", "/sent"}}

func TestRecentNav(t *testing.T) {
	r, open, changed := ui.NewRouter("/inbox"), false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if RecentNav(c, "History", r, &open, RecentNavOptions{}).Changed() {
			changed++
		}
		r.View(c, func(rt *ui.Route) { pages(rt, c) })
	}, 500, 300)
	r.Push("/sent")
	settles(t, tt, "pushed")
	before := baseline()
	click(t, tt, "Back")
	click(t, tt, "Forward")
	if r.Path() != "/sent" || changed != 2 {
		t.Fatalf("path %q, changed %d", r.Path(), changed)
	}
	click(t, tt, "Recent pages")
	if !open || changed != 3 {
		t.Fatalf("open %v, changed %d", open, changed)
	}
	click(t, tt, "Inbox")
	settles(t, tt, "picked")
	if r.Path() != "/inbox" || open || changed != 4 {
		t.Errorf("recent pick: path %q, open %v, changed %d", r.Path(), open, changed)
	}
	cues(t, before, sound.Tap, sound.Tap, sound.Open, sound.Tap)
	shut := countText(tt, "Inbox")
	click(t, tt, "Recent pages")
	if n := countText(tt, "Inbox") - shut; n != 2 {
		t.Errorf("a page seen twice should be one row, its label and text: %q", tt.Texts())
	}
	panics(t, "no open", func() { render(func(c *ui.Context) { RecentNav(c, "x", ui.NewRouter("/"), nil, RecentNavOptions{}) }) })
}

// countText counts entries equal to s.
func countText(tt *ui.Tester, s string) int {
	n := 0
	for _, x := range tt.Texts() {
		if x == s {
			n++
		}
	}
	return n
}

func TestRoutedPicksTheDeepestPage(t *testing.T) {
	ps := []Page{{"Home", "/"}, {"Inbox", "/inbox"}, {"Sent", "/sent"}}
	for path, want := range map[string]int{"/inbox/42": 1, "/inboxes": 0, "/": 0, "/sent": 2} {
		if got := routed("T", "x", ui.NewRouter(path), ps); got != want {
			t.Errorf("%s: page %d, want %d", path, got, want)
		}
	}
}

func TestSegmentNav(t *testing.T) {
	r := ui.NewRouter("/inbox")
	tt := ui.NewTester(func(c *ui.Context) { SegmentNav(c, "Mail", r, mail) }, 400, 100)
	before := baseline()
	click(t, tt, "Sent")
	if r.Path() != "/sent" {
		t.Errorf("path %q", r.Path())
	}
	cues(t, before, sound.Select)
	panics(t, "path under no page", func() { render(func(c *ui.Context) { SegmentNav(c, "x", ui.NewRouter("/settings"), mail) }) })
	panics(t, "twin paths", func() {
		render(func(c *ui.Context) { SegmentNav(c, "x", ui.NewRouter("/a"), []Page{{"A", "/a"}, {"B", "/a"}}) })
	})
}

func TestRouterTabs(t *testing.T) {
	r, changed := ui.NewRouter("/inbox"), 0
	tt := ui.NewTester(func(c *ui.Context) {
		if RouterTabs(c, "Mail", r, mail, func(rt *ui.Route) { pages(rt, c) }).Changed() {
			changed++
		}
	}, 500, 300)
	settles(t, tt, "RouterTabs")
	rest := find(t, tt, "Inbox page")
	before := baseline()
	click(t, tt, "Sent")
	tt.Advance(50 * time.Millisecond)
	if r.Path() != "/sent" || changed != 1 || find(t, tt, "Sent page").X <= rest.X {
		t.Fatalf("path %q, changed %d, did not slide in from the end", r.Path(), changed)
	}
	settles(t, tt, "slid")
	r.Back()
	tt.Advance(50 * time.Millisecond)
	if p := find(t, tt, "Inbox page"); r.Path() != "/inbox" || p.W >= rest.W {
		t.Errorf("back to %q did not slide in from the start: %+v", r.Path(), p)
	}
	cues(t, before, sound.Select)
	rr := ui.NewRouter("/inbox")
	tr := ui.NewTester(func(c *ui.Context) { RouterTabs(c, "Mail", rr, mail, func(rt *ui.Route) { pages(rt, c) }) }, 500, 300)
	tr.SetDirection(ui.RightToLeft)
	settles(t, tr, "right to left")
	rest = find(t, tr, "Inbox page")
	click(t, tr, "Sent")
	tr.Advance(50 * time.Millisecond)
	if p := find(t, tr, "Sent page"); p.X+p.W >= rest.X+rest.W-1 {
		t.Errorf("right to left: the next page did not come from the left: %+v", p)
	}
}
