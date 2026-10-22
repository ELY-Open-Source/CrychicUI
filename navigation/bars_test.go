package navigation

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var dests = []NavItem{{"Inbox", "inbox"}, {"Drafts", "file-text"}, {"Sent", "send"}}

func TestNavbar(t *testing.T) {
	active, changed := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Navbar(c, "Crychic", []string{"Home", "Docs", "Blog"}, &active, func() { button.Ghost(c, "Sign in") }).Changed() {
			changed++
		}
	}, 700, 120)
	find(t, tt, "Sign in")
	before := baseline()
	click(t, tt, "Docs")
	if active != 1 || changed != 1 || !tt.FramePending() {
		t.Errorf("active %d, changed %d, sliding %v", active, changed, tt.FramePending())
	}
	cues(t, before, sound.Tap)
	settles(t, tt, "Navbar")
	panics(t, "no links", func() { render(func(c *ui.Context) { a := 0; Navbar(c, "x", nil, &a, nil) }) })
}

func TestTopBar(t *testing.T) {
	title, backs := "Inbox", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if _, back := TopBar(c, title, TopBarOptions{}, nil); back.Clicked() {
			backs++
		}
	}, 500, 100)
	before := baseline()
	click(t, tt, "Back")
	if backs != 1 {
		t.Errorf("back clicks %d", backs)
	}
	cues(t, before, sound.Tap)
	title = "Message"
	tt.Frame()
	if !tt.HasText("Message") || !tt.FramePending() {
		t.Errorf("title did not cross-fade: %q", tt.Texts())
	}
}

func TestSideNavAndRail(t *testing.T) {
	for name, build := range map[string]func(*ui.Context, *int) ui.Element{
		"SideNav":    func(c *ui.Context, a *int) ui.Element { return SideNav(c, "Mail", dests, a) },
		"NavRail":    func(c *ui.Context, a *int) ui.Element { return NavRail(c, "Mail", dests, a) },
		"BottomTabs": func(c *ui.Context, a *int) ui.Element { return BottomTabs(c, "Mail", dests, a) },
	} {
		active, changed := 0, 0
		tt := ui.NewTester(func(c *ui.Context) {
			if build(c, &active).Changed() {
				changed++
			}
		}, 500, 400)
		before := baseline()
		click(t, tt, "Sent")
		if active != 2 || changed != 1 || !tt.FramePending() {
			t.Errorf("%s: active %d, changed %d, moving %v", name, active, changed, tt.FramePending())
		}
		cues(t, before, sound.Select)
		settles(t, tt, name)
	}
	panics(t, "bad icon", func() { render(func(c *ui.Context) { a := 0; SideNav(c, "x", []NavItem{{"A", "no-such-icon"}}, &a) }) })
}

func TestCollapsibleNav(t *testing.T) {
	active, open := 0, true
	tt := ui.NewTester(func(c *ui.Context) { CollapsibleNav(c, "Mail", dests, &active, &open, CollapsibleNavOptions{}) }, 500, 400)
	wide := find(t, tt, "Mail").W
	before := baseline()
	click(t, tt, "Collapse")
	settles(t, tt, "folding")
	if open || find(t, tt, "Mail").W >= wide/2 {
		t.Errorf("open %v, width %v of %v", open, find(t, tt, "Mail").W, wide)
	}
	click(t, tt, "Expand")
	settles(t, tt, "unfolding")
	if !open || find(t, tt, "Mail").W != wide {
		t.Errorf("open %v, width %v", open, find(t, tt, "Mail").W)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
}

func TestBackLink(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if BackLink(c, "All posts").Clicked() {
			clicks++
		}
	}, 300, 100)
	r := find(t, tt, "All posts")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("the arrow did not nudge")
	}
	before := baseline()
	click(t, tt, "All posts")
	if clicks != 1 {
		t.Errorf("clicks %d", clicks)
	}
	cues(t, before, sound.Tap)
}
