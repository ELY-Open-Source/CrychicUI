package gallery

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestNavigationDemosLog(t *testing.T) {
	click := func(labels ...string) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, l := range labels {
				if err := tt.Click(l); err != nil {
					return err
				}
				tt.Advance(300 * time.Millisecond)
			}
			return nil
		}
	}
	keys := func(ks ...ui.Key) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, k := range ks {
				tt.Key(0, k)
			}
			return nil
		}
	}
	scrollThen := func(area, label string) func(*ui.Tester) error {
		return func(tt *ui.Tester) error {
			r, _ := tt.Find(area)
			tt.Scroll(r.X+20, r.Y+20, 0, 2000)
			tt.Advance(time.Second)
			return tt.Click(label)
		}
	}
	acts := map[string]func(*ui.Tester) error{
		"Navbar":            click("Docs"),
		"TopBar":            click("Back"),
		"SideNav":           click("Inbox"),
		"CollapsibleNav":    click("Collapse"),
		"NavRail":           click("Inbox"),
		"BottomTabs":        click("Inbox"),
		"UnderlineTabs":     click("Activity"),
		"PillTabs":          click("Activity"),
		"VerticalTabs":      click("Activity"),
		"ClosableTabs":      click("New tab"),
		"ScrollTabs":        click("Month 2"),
		"Crumbs":            click("Home"),
		"Pagination":        click("Next"),
		"Pager":             click("Next"),
		"LoadMore":          click("Show more"),
		"InfiniteList":      func(tt *ui.Tester) error { r, _ := tt.Find("Feed"); tt.Scroll(r.X+20, r.Y+20, 0, 5000); return nil },
		"Steps":             click("Payment"),
		"TOC":               keys(ui.KeyTab, ui.KeyDown, ui.KeyEnter),
		"AnchorNav":         keys(ui.KeyTab, ui.KeyRight, ui.KeyEnter),
		"MegaMenu":          click("Explore", "Pricing"),
		"MenuStrip":         click("Company", "Careers"),
		"NavDropdown":       click("Docs", "API"),
		"Commander":         keys(ui.KeyTab, ui.KeyDown, ui.KeyEnter),
		"NavTree":           click("Getting started", "Install"),
		"FileNav":           click("src", "main.go"),
		"AccordionNav":      click("Company", "Careers"),
		"Drilldown":         click("Guides", "Layout"),
		"BackLink":          click("All posts"),
		"SkipLink":          keys(ui.KeyTab, ui.KeyEnter),
		"PagerDots":         click("Page 3"),
		"PageIndicator":     click("Next page"),
		"ToTop":             scrollThen("Story", "Back to top"),
		"SectionSwitch":     click("Week"),
		"WizardNav":         click("Next"),
		"IconToolbar":       click("Link"),
		"DockBar":           click("Mail"),
		"AppSwitcher":       click("Apps"),
		"WorkspaceSwitcher": click("Workspace"),
		"AccountNav":        click("Ada Lovelace"),
		"LanguageMenu":      click("Language"),
		"QuickLinks":        click("Calendar"),
		"FooterNav":         click("Careers"),
		"Sitemap":           click("Install"),
		"FilterBar":         click("Open"),
		"SortMenu":          click("Sort"),
		"ViewSwitch":        click("Grid"),
		"ZoomControl":       click("Zoom in"),
		"RecentNav":         click("Open Sent", "Back"),
		"RelatedNav":        click("Deploy to the web"),
		"FloatingNav":       scrollThen("Story", "Usage"),
		"StepperNav":        click("Payment"),
		"TabBarBadges":      click("Mentions", "New mention"),
		"SegmentNav":        click("Sent"),
		"RouterTabs":        click("Sent"),
	}
	for _, e := range Entries() {
		act, ok := acts[e.Name]
		if !ok || e.Category != "navigation" {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 600)
		tt.Advance(time.Second)
		if err := act(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Advance(time.Second)
		if len(logs) == 0 {
			t.Errorf("%s: the action logged nothing", e.Name)
		}
	}
}

func TestNavigationDemosLogMoves(t *testing.T) {
	moves := map[string][2]string{
		"Navbar":       {"Sign up", "sign up"},
		"ScrollTabs":   {"Scroll forward", "scrolled forward"},
		"MegaMenu":     {"Explore", "open true"},
		"MenuStrip":    {"Company", "open menu 1"},
		"NavDropdown":  {"Docs", "open true"},
		"AccordionNav": {"Company", "open group 1"},
		"Drilldown":    {"Guides", `at "Guides"`},
		"NavTree":      {"Getting started", "open map[Getting started:true]"},
		"FileNav":      {"src", "open map[src:true]"},
	}
	for _, e := range Entries() {
		m, ok := moves[e.Name]
		if !ok || e.Category != "navigation" {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 600)
		tt.Advance(time.Second)
		if err := tt.Click(m[0]); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Advance(time.Second)
		if !slices.ContainsFunc(logs, func(s string) bool { return strings.HasPrefix(s, m[1]) }) {
			t.Errorf("%s: %q logged %q, want %q", e.Name, m[0], logs, m[1])
		}
	}
}

func TestRecentNavDemoLogsBack(t *testing.T) {
	for _, e := range Entries() {
		if e.Name != "RecentNav" {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 600)
		for _, l := range []string{"Open Sent", "Back"} {
			if err := tt.Click(l); err != nil {
				t.Fatal(err)
			}
			tt.Advance(time.Second)
		}
		if len(logs) == 0 || logs[len(logs)-1] != "open false, at /inbox" {
			t.Errorf("Back logged %q", logs)
		}
	}
}
