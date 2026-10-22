package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/navigation"
	"github.com/egoist/mygo/ui"
)

var navigationEntries2 = []Entry{
	{584, "SkipLink", "navigation", "Shows on focus; jumps to content.", demoSkipLink},
	{585, "PagerDots", "navigation", "Dots for pages; the active widens.", demoPagerDots},
	{586, "PageIndicator", "navigation", "3 / 10 with arrows.", demoPageIndicator},
	{587, "ToTop", "navigation", "Appears past a scroll; scrolls up.", demoToTop},
	{588, "SectionSwitch", "navigation", "Segments over sliding sections.", demoSectionSwitch},
	{589, "WizardNav", "navigation", "Back, step dots, next.", demoWizardNav},
	{590, "IconToolbar", "navigation", "Grouped icon buttons.", demoIconToolbar},
	{591, "DockBar", "navigation", "A dock that magnifies by pointer.", demoDockBar},
	{592, "AppSwitcher", "navigation", "A grid of apps in a popover.", demoAppSwitcher},
	{593, "WorkspaceSwitcher", "navigation", "Workspaces with avatars.", demoWorkspaceSwitcher},
	{594, "AccountNav", "navigation", "An avatar opening a menu.", demoAccountNav},
	{595, "LanguageMenu", "navigation", "Languages with a check.", demoLanguageMenu},
	{596, "QuickLinks", "navigation", "Icon tiles that lift.", demoQuickLinks},
	{597, "FooterNav", "navigation", "Columns of links under titles.", demoFooterNav},
	{598, "Sitemap", "navigation", "A nested page list.", demoSitemap},
	{599, "FilterBar", "navigation", "Chips and a clear-all.", demoFilterBar},
	{600, "SortMenu", "navigation", "Sort key and direction.", demoSortMenu},
	{601, "ViewSwitch", "navigation", "List, grid, board.", demoViewSwitch},
	{602, "ZoomControl", "navigation", "Minus, percent, plus, fit.", demoZoomControl},
	{603, "RecentNav", "navigation", "Back and forward with recent pages.", demoRecentNav},
	{604, "KeyHints", "navigation", "Shortcut hints for the view.", demoKeyHints},
	{605, "RelatedNav", "navigation", "Related links after content.", demoRelatedNav},
	{606, "FloatingNav", "navigation", "A pill nav appearing on scroll.", demoFloatingNav},
	{607, "StepperNav", "navigation", "Vertical steps with details.", demoStepperNav},
	{608, "TabBarBadges", "navigation", "Tabs with rolling counts.", demoTabBarBadges},
	{609, "SegmentNav", "navigation", "Segments switching routes.", demoSegmentNav},
	{610, "RouterTabs", "navigation", "Tabs bound to a router.", demoRouterTabs},
}

var (
	navApps = []navigation.NavItem{{Label: "Mail", Icon: "mail"}, {Label: "Calendar", Icon: "calendar"}, {Label: "Files", Icon: "folder"}, {Label: "Notes", Icon: "notebook-pen"}, {Label: "Music", Icon: "music"}, {Label: "Photos", Icon: "image"}}
	navMail = []navigation.Page{{Name: "Inbox", Path: "/inbox"}, {Name: "Sent", Path: "/sent"}, {Name: "Drafts", Path: "/drafts"}}
)

// navLines fills a demo scroll.
func navLines(c *ui.Context) {
	for i := range 30 {
		ui.Text(c, fmt.Sprintf("Line %d of the story", i+1))
	}
}

// mailPages builds the demo routers' pages.
func mailPages(c *ui.Context) func(rt *ui.Route) {
	return func(rt *ui.Route) {
		for _, p := range navMail {
			if rt.Match(p.Path) {
				rt.Title(p.Name)
				ui.Text(c, p.Name+" page").Padding(16)
			}
		}
	}
}

// opens logs a menu's open state and its pick.
func opens(log func(string), open bool, pick any) { log(fmt.Sprintf("open %v, picked %v", open, pick)) }

func demoSkipLink(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(8)
	h := ui.Local(root, "main", func() *ui.Handle { return &ui.Handle{} })
	root.Children(func() {
		if navigation.SkipLink(c, "Skip to content", *h).Clicked() {
			log("jumped to content")
		}
		kit.Muted(c, "Press Tab to show the link.")
		ui.Column(c).Focusable().Label("Content").Bind(*h).Padding(12).Children(func() { ui.Text(c, "Main content") })
	})
}

func demoPagerDots(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	p := index(root, 0)
	root.Children(func() {
		if navigation.PagerDots(c, "Slides", p, 5, navigation.PagerDotsOptions{}).Changed() {
			log(fmt.Sprintf("page %d", *p+1))
		}
	})
}

func demoPageIndicator(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	p := index(root, 2)
	root.Children(func() {
		if navigation.PageIndicator(c, "Results", p, 10, navigation.PageIndicatorOptions{}).Changed() {
			log(fmt.Sprintf("page %d", *p+1))
		}
	})
}

func demoToTop(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(8)
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	root.Children(func() {
		ui.Scroll(c).Height(200).TrackScroll(st).Label("Story").Children(func() { navLines(c) })
		ui.Row(c).FillWidth().Height(48).Justify(ui.End).Children(func() {
			if _, up := navigation.ToTop(c, "Back to top", st, navigation.ToTopOptions{After: 120}); up.Clicked() {
				log("scrolled up")
			}
		})
	})
}

func demoSectionSwitch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	i := index(root, 0)
	names := []string{"Day", "Week", "Month"}
	root.Children(func() {
		if navigation.SectionSwitch(c, "Range", names, i, func(n int) { ui.Text(c, names[n]+" view").Padding(12) }).Changed() {
			log(names[*i])
		}
	})
}

func demoWizardNav(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	step := index(root, 0)
	root.Children(func() {
		nav, finish := navigation.WizardNav(c, "Setup", step, 4, navigation.WizardNavOptions{})
		if nav.Changed() {
			log(fmt.Sprintf("step %d", *step+1))
		}
		if finish.Clicked() {
			log("finished")
		}
	})
}

func demoIconToolbar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	used := ui.Local(root, "used", func() [2]int { return [2]int{-1, -1} })
	groups := [][]navigation.NavItem{{{Label: "Bold", Icon: "bold"}, {Label: "Italic", Icon: "italic"}, {Label: "Underline", Icon: "underline"}},
		{{Label: "Link", Icon: "link"}, {Label: "Image", Icon: "image"}}}
	root.Children(func() {
		if navigation.IconToolbar(c, "Format", groups, used).Changed() {
			log(groups[used[0]][used[1]].Label)
		}
	})
}

func demoDockBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	opened := index(root, -1)
	root.Children(func() {
		if navigation.DockBar(c, "Dock", navApps, opened).Changed() {
			log(navApps[*opened].Label)
		}
	})
}

func demoAppSwitcher(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	app := index(root, -1)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		if navigation.AppSwitcher(c, "Apps", navApps, app, open).Changed() {
			opens(log, *open, *app)
		}
	})
}

func demoWorkspaceSwitcher(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	space := index(root, 0)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		if navigation.WorkspaceSwitcher(c, "Workspace", []string{"Acme Labs", "Globex", "Initech"}, space, open).Changed() {
			opens(log, *open, *space)
		}
	})
}

func demoAccountNav(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	item := index(root, -1)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		if navigation.AccountNav(c, "Ada Lovelace", nil, []string{"Profile", "Settings", "Sign out"}, item, open).Changed() {
			opens(log, *open, *item)
		}
	})
}

func demoLanguageMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	lang := index(root, 0)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		if navigation.LanguageMenu(c, "Language", []string{"English", "Deutsch", "Français", "Español"}, lang, open).Changed() {
			opens(log, *open, *lang)
		}
	})
}

func demoQuickLinks(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	p := index(root, -1)
	root.Children(func() {
		if navigation.QuickLinks(c, "Shortcuts", navApps, p).Changed() {
			log(navApps[*p].Label)
		}
	})
}

func demoFooterNav(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(560)
	pick := ui.Local(root, "pick", func() [2]int { return [2]int{-1, -1} })
	root.Children(func() {
		if navigation.FooterNav(c, "Footer", navGroups, pick).Changed() {
			log(navGroups[pick[0]].Links[pick[1]])
		}
	})
}

func demoSitemap(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	picked := ui.Local(root, "picked", func() string { return "" })
	root.Children(func() {
		if navigation.Sitemap(c, "Pages", navSite, picked).Changed() {
			log(*picked)
		}
	})
}

func demoFilterBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	sel := ui.Local(root, "sel", func() []string { return nil })
	root.Children(func() {
		if navigation.FilterBar(c, "Status", []string{"Open", "In review", "Closed"}, sel, navigation.FilterBarOptions{}).Changed() {
			log(fmt.Sprintf("filters %v", *sel))
		}
	})
}

func demoSortMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	key := index(root, 0)
	desc := ui.Local(root, "desc", func() bool { return false })
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		if navigation.SortMenu(c, "Sort", []string{"Name", "Date", "Size"}, key, desc, open, navigation.SortMenuOptions{}).Changed() {
			log(fmt.Sprintf("open %v, key %d, descending %v", *open, *key, *desc))
		}
	})
}

func demoViewSwitch(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	view := index(root, 0)
	root.Children(func() {
		if navigation.ViewSwitch(c, "View", view, navigation.ViewSwitchOptions{}).Changed() {
			log(fmt.Sprintf("view %d", *view))
		}
	})
}

func demoZoomControl(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	zoom := ui.Local(root, "zoom", func() float64 { return 1 })
	root.Children(func() {
		if navigation.ZoomControl(c, "Zoom", zoom, navigation.ZoomControlOptions{}).Changed() {
			log(fmt.Sprintf("zoom %.0f%%", *zoom*100))
		}
	})
}

func demoRecentNav(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(8)
	r := ui.Local(root, "router", func() *ui.Router { return ui.NewRouter("/inbox") })
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		if navigation.RecentNav(c, "History", *r, open, navigation.RecentNavOptions{}).Changed() {
			log(fmt.Sprintf("open %v, at %s", *open, (*r).Path()))
		}
		ui.Row(c).Gap(8).Children(func() {
			for _, p := range navMail {
				c.Key(p.Path)
				if button.Ghost(c, "Open "+p.Name).Clicked() {
					(*r).Push(p.Path)
					log("opened " + p.Path)
				}
			}
		})
		(*r).View(c, mailPages(c)).Height(60)
	})
}

func demoKeyHints(c *ui.Context, log func(string)) {
	navigation.KeyHints(c, "Shortcuts", []navigation.KeyHint{{Keys: []string{"Ctrl", "K"}, Label: "Search"},
		{Keys: []string{"J"}, Label: "Next"}, {Keys: []string{"K"}, Label: "Previous"}, {Keys: []string{"Esc"}, Label: "Close"}}).Width(440)
}

func demoRelatedNav(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	p := index(root, -1)
	links := []string{"Install the CLI", "Your first window", "Deploy to the web"}
	root.Children(func() {
		if navigation.RelatedNav(c, "Related", links, p).Changed() {
			log(links[*p])
		}
	})
}

func demoFloatingNav(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(8)
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	i := index(root, 0)
	links := []string{"Intro", "Usage", "API"}
	root.Children(func() {
		ui.Box(c).Height(48).FillWidth().Children(func() {
			if navigation.FloatingNav(c, "Sections", links, i, st, navigation.FloatingNavOptions{After: 80}).Changed() {
				log(links[*i])
			}
		})
		ui.Scroll(c).Height(200).TrackScroll(st).Label("Story").Children(func() { navLines(c) })
	})
}

func demoStepperNav(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	cur := index(root, 1)
	steps := []navigation.Step{{Title: "Account", Detail: "Name and email"}, {Title: "Plan", Detail: "Pick a plan"},
		{Title: "Payment", Detail: "Add a card"}, {Title: "Done", Detail: "Start using it"}}
	root.Children(func() {
		if navigation.StepperNav(c, "Signup", steps, cur).Changed() {
			log(steps[*cur].Title)
		}
	})
}

func demoTabBarBadges(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	i := index(root, 0)
	counts := ui.Local(root, "counts", func() []int { return []int{4, 12, 0} })
	names := []string{"All", "Mentions", "Archived"}
	root.Children(func() {
		if navigation.TabBarBadges(c, "Inbox", names, *counts, i).Changed() {
			log(names[*i])
		}
		if button.Ghost(c, "New mention").Clicked() {
			(*counts)[1]++
			log(fmt.Sprintf("mentions %d", (*counts)[1]))
		}
	})
}

func demoSegmentNav(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(8)
	r := ui.Local(root, "router", func() *ui.Router { return ui.NewRouter("/inbox") })
	root.Children(func() {
		if navigation.SegmentNav(c, "Mail", *r, navMail).Changed() {
			log((*r).Path())
		}
		(*r).View(c, mailPages(c)).Height(60)
	})
}

func demoRouterTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	r := ui.Local(root, "router", func() *ui.Router { return ui.NewRouter("/inbox") })
	root.Children(func() {
		if navigation.RouterTabs(c, "Mail", *r, navMail, mailPages(c)).Changed() {
			log((*r).Path())
		}
	})
}
