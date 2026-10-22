package gallery

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/navigation"
	"github.com/egoist/mygo/ui"
)

var navigationEntries = []Entry{
	{556, "Navbar", "navigation", "Brand, links and actions.", demoNavbar},
	{557, "TopBar", "navigation", "Title with back and actions.", demoTopBar},
	{558, "SideNav", "navigation", "A column of destinations.", demoSideNav},
	{559, "CollapsibleNav", "navigation", "A side nav folding to icons.", demoCollapsibleNav},
	{560, "NavRail", "navigation", "Icons with labels under them.", demoNavRail},
	{561, "BottomTabs", "navigation", "Mobile tabs; the active bounces.", demoBottomTabs},
	{562, "UnderlineTabs", "navigation", "Tabs whose underline slides.", demoUnderlineTabs},
	{563, "PillTabs", "navigation", "Tabs whose pill slides.", demoPillTabs},
	{564, "VerticalTabs", "navigation", "Tabs in a column.", demoVerticalTabs},
	{565, "ClosableTabs", "navigation", "Tabs you close and add.", demoClosableTabs},
	{566, "ScrollTabs", "navigation", "Overflowing tabs with arrows.", demoScrollTabs},
	{567, "Crumbs", "navigation", "A path whose middle folds.", demoCrumbs},
	{568, "Pagination", "navigation", "Pages with gaps.", demoPagination},
	{569, "Pager", "navigation", "Page x of y.", demoPager},
	{570, "LoadMore", "navigation", "Items and a button for more.", demoLoadMore},
	{571, "InfiniteList", "navigation", "A list asking for more at its end.", demoInfiniteList},
	{572, "Steps", "navigation", "Numbered steps joined by lines.", demoSteps},
	{573, "TOC", "navigation", "Headings that follow the scroll.", demoTOC},
	{574, "AnchorNav", "navigation", "Links scrolling to sections.", demoAnchorNav},
	{575, "MegaMenu", "navigation", "A wide panel of grouped links.", demoMegaMenu},
	{576, "MenuStrip", "navigation", "A bar of menus.", demoMenuStrip},
	{577, "NavDropdown", "navigation", "A link with a submenu.", demoNavDropdown},
	{578, "Commander", "navigation", "A searchable command list.", demoCommander},
	{579, "NavTree", "navigation", "Nested links that expand.", demoNavTree},
	{580, "FileNav", "navigation", "A folder tree with file icons.", demoFileNav},
	{581, "AccordionNav", "navigation", "Groups opening one at a time.", demoAccordionNav},
	{582, "Drilldown", "navigation", "Menus sliding into submenus.", demoDrilldown},
	{583, "BackLink", "navigation", "An arrow and a label.", demoBackLink},
}

var (
	navDests  = []navigation.NavItem{{Label: "Home", Icon: "house"}, {Label: "Search", Icon: "search"}, {Label: "Inbox", Icon: "inbox"}, {Label: "Settings", Icon: "settings"}}
	navGroups = []navigation.LinkGroup{{Title: "Product", Links: []string{"Features", "Pricing", "Changelog"}},
		{Title: "Company", Links: []string{"About", "Careers", "Press"}}, {Title: "Help", Links: []string{"Docs", "Support"}}}
	navSite = []navigation.Node{{Label: "Getting started", Children: []navigation.Node{{Label: "Install"}, {Label: "First app"}}},
		{Label: "Guides", Children: []navigation.Node{{Label: "Layout"}, {Label: "Motion", Children: []navigation.Node{{Label: "Springs"}}}}},
		{Label: "Changelog"}}
	navFiles = []navigation.Node{{Label: "src", Children: []navigation.Node{{Label: "main.go"}, {Label: "app.go"}}},
		{Label: "docs", Children: []navigation.Node{{Label: "README.md"}, {Label: "logo.svg"}}}, {Label: "LICENSE"}}
)

// index is a demo choice kept across builds.
func index(root ui.Element, start int) *int { return ui.Local(root, "i", func() int { return start }) }

func demoNavbar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(560)
	i := index(root, 0)
	links := []string{"Home", "Docs", "Blog"}
	root.Children(func() {
		if navigation.Navbar(c, "Crychic", links, i, func() {
			if button.Solid(c, "Sign up").Clicked() {
				log("sign up")
			}
		}).Changed() {
			log(links[*i])
		}
	})
}

func demoTopBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	deep := ui.Local(root, "deep", func() bool { return true })
	root.Children(func() {
		title := "Inbox"
		if *deep {
			title = "Message"
		}
		if _, back := navigation.TopBar(c, title, navigation.TopBarOptions{}, nil); back.Clicked() {
			*deep = !*deep
			log("back")
		}
	})
}

func demoSideNav(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, 0)
	root.Children(func() {
		if navigation.SideNav(c, "Main", navDests, i).Changed() {
			log(navDests[*i].Label)
		}
	})
}

func demoCollapsibleNav(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, 0)
	open := ui.Local(root, "open", func() bool { return true })
	root.Children(func() {
		if navigation.CollapsibleNav(c, "Main", navDests, i, open, navigation.CollapsibleNavOptions{}).Changed() {
			log(fmt.Sprintf("%s, open %v", navDests[*i].Label, *open))
		}
	})
}

func demoNavRail(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, 0)
	root.Children(func() {
		if navigation.NavRail(c, "Main", navDests, i).Changed() {
			log(navDests[*i].Label)
		}
	})
}

func demoBottomTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	i := index(root, 0)
	root.Children(func() {
		if navigation.BottomTabs(c, "Main", navDests, i).Changed() {
			log(navDests[*i].Label)
		}
	})
}

var navTabs = []string{"Overview", "Activity", "Settings"}

func demoUnderlineTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, 0)
	root.Children(func() {
		if navigation.UnderlineTabs(c, "Sections", navTabs, i).Changed() {
			log(navTabs[*i])
		}
	})
}

func demoPillTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, 0)
	root.Children(func() {
		if navigation.PillTabs(c, "Sections", navTabs, i).Changed() {
			log(navTabs[*i])
		}
	})
}

func demoVerticalTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, 0)
	root.Children(func() {
		if navigation.VerticalTabs(c, "Sections", navTabs, i).Changed() {
			log(navTabs[*i])
		}
	})
}

func demoClosableTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, 0)
	names := ui.Local(root, "names", func() []string { return []string{"main.go", "app.go", "README.md"} })
	root.Children(func() {
		if navigation.ClosableTabs(c, "Files", names, i, navigation.ClosableTabsOptions{}).Changed() {
			log(fmt.Sprintf("%d tabs, active %d", len(*names), *i))
		}
	})
}

func demoScrollTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	i := index(root, 0)
	var names []string
	for m := range 12 {
		names = append(names, fmt.Sprintf("Month %d", m+1))
	}
	root.Children(func() {
		tabs, back, forth := navigation.ScrollTabs(c, "Months", names, i, navigation.ScrollTabsOptions{})
		if tabs.Changed() {
			log(names[*i])
		}
		if back.Clicked() {
			log("scrolled back")
		}
		if forth.Clicked() {
			log("scrolled forward")
		}
	})
}

func demoCrumbs(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, -1)
	path := []string{"Home", "Work", "Projects", "Crychic", "Docs", "Readme"}
	root.Children(func() {
		if navigation.Crumbs(c, "Path", path, i, navigation.CrumbsOptions{}).Changed() {
			log(path[*i])
		}
	})
}

func demoPagination(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	p := index(root, 6)
	root.Children(func() {
		if navigation.Pagination(c, "Results", p, 20, navigation.PaginationOptions{}).Changed() {
			log(fmt.Sprintf("page %d", *p))
		}
	})
}

func demoPager(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	p := index(root, 1)
	root.Children(func() {
		if navigation.Pager(c, "Results", p, 8, navigation.PagerOptions{}).Changed() {
			log(fmt.Sprintf("page %d", *p))
		}
	})
}

func demoLoadMore(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	shown := index(root, 4)
	root.Children(func() {
		if navigation.LoadMore(c, "Posts", shown, 14, func(i int) { ui.Text(c, fmt.Sprintf("Post %d", i+1)) },
			navigation.LoadMoreOptions{Step: 4}).Changed() {
			log(fmt.Sprintf("%d shown", *shown))
		}
	})
}

func demoInfiniteList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	count := index(root, 12)
	root.Children(func() {
		navigation.InfiniteList(c, "Feed", *count, false, func(i int) { ui.Text(c, fmt.Sprintf("Story %d", i+1)) }, func() {
			*count += 8
			log(fmt.Sprintf("%d stories", *count))
		}, navigation.InfiniteListOptions{}).Height(200)
	})
}

func demoSteps(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	i := index(root, 1)
	steps := []string{"Cart", "Address", "Payment", "Review"}
	root.Children(func() {
		if navigation.Steps(c, "Checkout", steps, i).Changed() {
			log(steps[*i])
		}
	})
}

var navHeads = []string{"Introduction", "Install", "Usage", "Notes"}

// section is demo body text for one heading.
func section(c *ui.Context) func(i int) {
	return func(i int) {
		ui.Text(c, "Each section holds a few lines of text, long enough to scroll past as you read on.").MaxWidth(320)
		ui.Box(c).Height(120)
	}
}

func demoTOC(c *ui.Context, log func(string)) {
	if navigation.TOC(c, "Guide", navHeads, section(c), navigation.TOCOptions{}).Width(560).Height(220).Changed() {
		log("jumped")
	}
}

func demoAnchorNav(c *ui.Context, log func(string)) {
	if navigation.AnchorNav(c, "Guide", navHeads, section(c)).Width(480).Height(220).Changed() {
		log("jumped")
	}
}

func demoMegaMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	pick := ui.Local(root, "pick", func() [2]int { return [2]int{-1, -1} })
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		if navigation.MegaMenu(c, "Explore", navGroups, pick, open).Changed() {
			log(fmt.Sprintf("open %v, picked %v", *open, *pick))
		}
	})
}

func demoMenuStrip(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	pick := ui.Local(root, "pick", func() [2]int { return [2]int{-1, -1} })
	open := index(root, -1)
	root.Children(func() {
		if navigation.MenuStrip(c, "Menu", navGroups, pick, open).Changed() {
			log(fmt.Sprintf("open menu %d, picked %v", *open, *pick))
		}
	})
}

func demoNavDropdown(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, -1)
	open := ui.Local(root, "open", func() bool { return false })
	links := []string{"Guides", "API", "Examples"}
	root.Children(func() {
		if navigation.NavDropdown(c, "Docs", links, i, open).Changed() {
			log(fmt.Sprintf("open %v, picked %d", *open, *i))
		}
	})
}

func demoCommander(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	i := index(root, -1)
	commands := []string{"Open file", "Save file", "Save all", "Close tab", "Toggle sidebar", "Go to line"}
	root.Children(func() {
		if navigation.Commander(c, "Commands", commands, i, navigation.CommanderOptions{}).Changed() {
			log(commands[*i])
		}
	})
}

func demoNavTree(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := ui.Local(root, "chosen", func() string { return "" })
	open := ui.Local(root, "open", func() map[string]bool { return map[string]bool{} })
	root.Children(func() {
		if navigation.NavTree(c, "Docs", navSite, chosen, *open).Changed() {
			log(fmt.Sprintf("open %v, chosen %q", *open, *chosen))
		}
	})
}

func demoFileNav(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := ui.Local(root, "chosen", func() string { return "" })
	open := ui.Local(root, "open", func() map[string]bool { return map[string]bool{} })
	root.Children(func() {
		if navigation.FileNav(c, "Files", navFiles, chosen, *open).Changed() {
			log(fmt.Sprintf("open %v, chosen %q", *open, *chosen))
		}
	})
}

func demoAccordionNav(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	pick := ui.Local(root, "pick", func() [2]int { return [2]int{-1, -1} })
	open := index(root, 0)
	root.Children(func() {
		if navigation.AccordionNav(c, "Help", navGroups, pick, open).Changed() {
			log(fmt.Sprintf("open group %d, picked %v", *open, *pick))
		}
	})
}

func demoDrilldown(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := ui.Local(root, "chosen", func() string { return "" })
	trail := ui.Local(root, "trail", func() []string { return nil })
	root.Children(func() {
		if navigation.Drilldown(c, "Menu", navSite, chosen, trail, navigation.DrilldownOptions{}).Changed() {
			log(fmt.Sprintf("at %q, chosen %q", strings.Join(*trail, "/"), *chosen))
		}
	})
}

func demoBackLink(c *ui.Context, log func(string)) {
	if navigation.BackLink(c, "All posts").Clicked() {
		log("back")
	}
}
