package navigation

import (
	"fmt"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// bare reports whether label's center shows the page.
func bare(t *testing.T, tt *ui.Tester, label string) bool {
	t.Helper()
	r, bg := find(t, tt, label), ui.LightTheme().Background
	p := tt.Image().RGBAAt(int(r.X+r.W/2), int(r.Y+r.H/2))
	return p.R == bg.R && p.G == bg.G && p.B == bg.B
}

func TestSkipLink(t *testing.T) {
	var h ui.Handle
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Children(func() {
			SkipLink(c, "Skip to content", &h)
			ui.Column(c).Focusable().Label("Main").Bind(&h).Padding(80).Children(func() { ui.Text(c, "Body") })
		})
	}, 400, 300)
	settles(t, tt, "SkipLink")
	if r, ok := tt.Find("Skip to content"); ok && r.H >= 1 {
		t.Fatalf("the link shows before it has focus: %+v", r)
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	settles(t, tt, "focus")
	if bare(t, tt, "Skip to content") {
		t.Fatal("the focused link stays hidden")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if !tt.Focused("Main") {
		t.Error("focus did not move to the content")
	}
	cues(t, before, sound.Tap)
	panics(t, "no target", func() { render(func(c *ui.Context) { SkipLink(c, "Skip", nil) }) })
}

func TestPagerDots(t *testing.T) {
	page := 0
	tt := ui.NewTester(func(c *ui.Context) { PagerDots(c, "Slides", &page, 4, PagerDotsOptions{}) }, 300, 100)
	before := baseline()
	click(t, tt, "Page 3")
	if page != 2 || !tt.FramePending() {
		t.Fatalf("page %d, widening %v", page, tt.FramePending())
	}
	settles(t, tt, "PagerDots")
	if find(t, tt, "Page 3").W <= find(t, tt, "Page 1").W {
		t.Error("the active dot did not widen")
	}
	tt.Key(0, ui.KeyRight)
	if page != 3 {
		t.Errorf("Right gave page %d", page)
	}
	cues(t, before, sound.Tick, sound.Tick)
	panics(t, "one page", func() { render(func(c *ui.Context) { p := 0; PagerDots(c, "x", &p, 1, PagerDotsOptions{}) }) })
}

func TestPageIndicator(t *testing.T) {
	page, changed := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if PageIndicator(c, "Results", &page, 10, PageIndicatorOptions{}).Changed() {
			changed++
		}
	}, 300, 100)
	if !tt.HasText("/ 10") {
		t.Fatalf("total %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Previous page")
	click(t, tt, "Next page")
	click(t, tt, "Next page")
	click(t, tt, "Previous page")
	if page != 1 || changed != 3 {
		t.Errorf("page %d, changed %d", page, changed)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	panics(t, "page past the end", func() { render(func(c *ui.Context) { p := 10; PageIndicator(c, "x", &p, 10, PageIndicatorOptions{}) }) })
}

// scrolled builds a tall scroll tracked by st beside extra.
func scrolled(st *ui.ScrollState, extra func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) {
		ui.Column(c).Fill().Children(func() {
			extra(c)
			ui.Scroll(c).Height(200).TrackScroll(st).Label("Feed").Children(func() { ui.Box(c).Height(2000).Width(200) })
		})
	}
}

func TestToTop(t *testing.T) {
	var st ui.ScrollState
	ups := 0
	tt := ui.NewTester(scrolled(&st, func(c *ui.Context) {
		if _, up := ToTop(c, "Back to top", &st, ToTopOptions{}); up.Clicked() {
			ups++
		}
	}), 300, 300)
	settles(t, tt, "ToTop")
	if _, ok := tt.Find("Back to top"); ok {
		t.Fatal("shown at the top")
	}
	r := find(t, tt, "Feed")
	tt.Scroll(r.X+20, r.Y+20, 0, 1000)
	settles(t, tt, "scrolled")
	before := baseline()
	click(t, tt, "Back to top")
	settles(t, tt, "back up")
	if ups != 1 || st.Y != 0 {
		t.Errorf("clicks %d, y %v", ups, st.Y)
	}
	if _, ok := tt.Find("Back to top"); ok {
		t.Error("still shown at the top")
	}
	cues(t, before, sound.Tap)
}

func TestSectionSwitch(t *testing.T) {
	active, changed := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SectionSwitch(c, "Views", []string{"Day", "Week"}, &active, func(i int) { ui.Text(c, fmt.Sprint("Section ", i)) }).Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	click(t, tt, "Week")
	if active != 1 || changed != 1 || !tt.FramePending() {
		t.Fatalf("active %d, changed %d, sliding %v", active, changed, tt.FramePending())
	}
	settles(t, tt, "slide")
	if !tt.HasText("Section 1") || tt.HasText("Section 0") {
		t.Errorf("sections %q", tt.Texts())
	}
	cues(t, before, sound.Select)
	panics(t, "no sections", func() { render(func(c *ui.Context) { a := 0; SectionSwitch(c, "x", []string{"A"}, &a, nil) }) })
}

func TestWizardNav(t *testing.T) {
	step, changed, done := 0, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		root, finish := WizardNav(c, "Setup", &step, 3, WizardNavOptions{})
		if root.Changed() {
			changed++
		}
		if finish.Clicked() {
			done++
		}
	}, 500, 100)
	before := baseline()
	click(t, tt, "Back")
	click(t, tt, "Next")
	click(t, tt, "Next")
	if step != 2 || tt.HasText("Next") || !tt.HasText("Finish") {
		t.Fatalf("step %d, %q", step, tt.Texts())
	}
	click(t, tt, "Finish")
	click(t, tt, "Back")
	if step != 1 || changed != 3 || done != 1 {
		t.Errorf("step %d, changed %d, finished %d", step, changed, done)
	}
	cues(t, before, sound.Tap, sound.Tap, sound.Tap, sound.Tap)
}

func TestFloatingNav(t *testing.T) {
	var st ui.ScrollState
	active := 0
	tt := ui.NewTester(scrolled(&st, func(c *ui.Context) {
		FloatingNav(c, "Sections", []string{"Intro", "Docs"}, &active, &st, FloatingNavOptions{})
	}), 400, 300)
	settles(t, tt, "FloatingNav")
	if tt.HasText("Docs") {
		t.Fatal("shown at the top")
	}
	r := find(t, tt, "Feed")
	tt.Scroll(r.X+20, r.Y+20, 0, 1000)
	tt.Frame()
	if !tt.HasText("Docs") || !tt.FramePending() {
		t.Fatalf("did not slide in: %q", tt.Texts())
	}
	settles(t, tt, "slid in")
	before := baseline()
	click(t, tt, "Docs")
	if active != 1 {
		t.Errorf("active %d", active)
	}
	cues(t, before, sound.Tap)
}

func TestStepperNav(t *testing.T) {
	current, changed := 0, 0
	steps := []Step{{"Plan", "Pick a plan"}, {"Pay", "Add a card"}, {"Ship", ""}}
	tt := ui.NewTester(func(c *ui.Context) {
		if StepperNav(c, "Order", steps, &current).Changed() {
			changed++
		}
	}, 400, 400)
	if !tt.HasText("Add a card") {
		t.Fatalf("details %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Ship")
	if current != 2 || changed != 1 || !tt.FramePending() {
		t.Errorf("current %d, changed %d, filling %v", current, changed, tt.FramePending())
	}
	cues(t, before, sound.Select)
	panics(t, "one step", func() { render(func(c *ui.Context) { s := 0; StepperNav(c, "x", steps[:1], &s) }) })
}

func TestTabBarBadges(t *testing.T) {
	active, counts := 0, []int{3, 12}
	tt := ui.NewTester(func(c *ui.Context) { TabBarBadges(c, "Inbox", []string{"All", "Mentions"}, counts, &active) }, 400, 100)
	settles(t, tt, "TabBarBadges")
	before := baseline()
	click(t, tt, "Mentions")
	if active != 1 {
		t.Errorf("active %d", active)
	}
	settles(t, tt, "picked")
	counts[1] = 13
	tt.Frame()
	if !tt.FramePending() {
		t.Error("the count did not roll")
	}
	cues(t, before, sound.Select)
	panics(t, "short counts", func() { render(func(c *ui.Context) { a := 0; TabBarBadges(c, "x", []string{"A", "B"}, []int{1}, &a) }) })
}
