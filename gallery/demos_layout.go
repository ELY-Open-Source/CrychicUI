package gallery

import (
	"fmt"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/layout"
	"github.com/egoist/mygo/ui"
)

var layoutEntries = []Entry{
	{1, "Stack", "layout", "Column gapped in theme units; its size animates.", demoStack},
	{2, "Inline", "layout", "Wrapping row; its size animates.", demoInline},
	{3, "Layers", "layout", "Layers on one box, each fading in.", demoLayers},
	{4, "Centered", "layout", "Centers content, which scales in.", demoCentered},
	{5, "Container", "layout", "Centered page column with a max width.", demoContainer},
	{6, "Section", "layout", "Muted heading, growing divider, body.", demoSection},
	{7, "Panel", "layout", "Bordered surface; its size animates.", demoPanel},
	{8, "Card", "layout", "Card surface that lifts on hover.", demoCard},
	{9, "Inset", "layout", "Recessed surface following the theme.", demoInset},
	{10, "Well", "layout", "Sunken box; content fades in.", demoWell},
	{11, "SplitPane", "layout", "Two resizable panes side by side.", demoSplitPane},
	{12, "SplitStack", "layout", "Two resizable panes stacked.", demoSplitStack},
	{13, "SidebarLayout", "layout", "Side column that springs open and closed.", demoSidebarLayout},
	{14, "HolyGrail", "layout", "Header, nav, main, aside and footer.", demoHolyGrail},
	{15, "AppFrame", "layout", "Top bar, content and status bar.", demoAppFrame},
	{16, "Masonry", "layout", "Items in the shortest of n columns.", demoMasonry},
	{17, "AutoGrid", "layout", "As many columns as a minimum width allows.", demoAutoGrid},
	{18, "Bento", "layout", "Grid of spanning tiles that lift.", demoBento},
	{19, "Cluster", "layout", "Wrapping items aligned start, center or end.", demoCluster},
	{20, "Columns", "layout", "Equal columns whose widths animate.", demoColumns},
	{21, "AspectBox", "layout", "Box held at a ratio.", demoAspectBox},
	{22, "Frame", "layout", "Clips content to the theme radius.", demoFrame},
	{23, "ScrollArea", "layout", "Scroll view with fading edge shadows.", demoScrollArea},
	{24, "Scroller", "layout", "Sideways scroll with paging arrows.", demoScroller},
	{25, "StickyHeader", "layout", "Header over a scroll; shadow once scrolled.", demoStickyHeader},
}

// counter is a demo's count with add and remove buttons.
func counter(c *ui.Context, root ui.Element, name string, start int, log func(string)) *int {
	n := ui.Local(root, "n", func() int { return start })
	ui.Row(c).Gap(8).Children(func() {
		if ui.Button(c, "Add").Clicked() {
			*n++
			log(fmt.Sprintf("%s: %d", name, *n))
		}
		if ui.Button(c, "Remove").Disabled(*n == 0).Clicked() {
			*n--
			log(fmt.Sprintf("%s: %d", name, *n))
		}
	})
	return n
}

func tile(c *ui.Context, label string, h float32) {
	t := c.Theme()
	ui.Column(c).Height(h).Padding(t.Space(2)).Radius(t.Radius).Background(t.Surface).Children(func() {
		ui.Text(c, label)
	})
}

func demoStack(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "rows", 3, log)
		layout.Stack(c, 2, func() {
			for i := range *n {
				ui.Text(c, fmt.Sprintf("Row %d", i+1))
			}
		})
	})
}

func demoInline(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "tags", 6, log)
		layout.Inline(c, 2, func() {
			for i := range *n {
				ui.Badge(c, fmt.Sprintf("Tag %d", i+1))
			}
		})
	})
}

func demoLayers(c *ui.Context, log func(string)) {
	t := c.Theme()
	layout.Layers(c, []func(){
		func() { ui.Box(c).Size(240, 120).Radius(t.Radius).Background(t.Surface) },
		func() { ui.Text(c, "Caption over the backdrop").Padding(16) },
	})
}

func demoCentered(c *ui.Context, log func(string)) {
	ui.Box(c).Height(160).Children(func() {
		layout.Centered(c, func() { ui.Text(c, "In the middle").Bold() })
	})
}

func demoContainer(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	wide := ui.Local(root, "wide", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, wide, "Wide").Changed() {
			log(fmt.Sprintf("wide: %v", *wide))
		}
		w := c.Theme().Rem(20)
		if *wide {
			w = c.Theme().Rem(40)
		}
		layout.Container(c, layout.ContainerOptions{MaxWidth: w}, func() {
			ui.Text(c, "Content keeps to the container's width, centered on the page.")
		}).Background(c.Theme().Surface)
	})
}

func demoSection(c *ui.Context, log func(string)) {
	layout.Section(c, "Notifications", func() {
		ui.Text(c, "Choose what reaches you.")
		kit.Muted(c, "The divider grows in as the section appears.")
	})
}

func demoPanel(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "lines", 1, log)
		layout.Panel(c, func() {
			for i := range *n {
				ui.Text(c, fmt.Sprintf("Line %d", i+1))
			}
		})
	})
}

func demoCard(c *ui.Context, log func(string)) {
	layout.Card(c, func() {
		ui.Text(c, "Hover me").Bold()
		kit.Muted(c, "The card lifts 2 DIP.")
	}).Width(220)
}

func demoInset(c *ui.Context, log func(string)) {
	layout.Inset(c, func() {
		ui.Text(c, "Grouped settings")
		kit.Muted(c, "Switch the theme to watch the surface ease over.")
	})
}

func demoWell(c *ui.Context, log func(string)) {
	layout.Well(c, func() { ui.Text(c, "Sunken content") })
}

func demoSplitPane(c *ui.Context, log func(string)) {
	root := ui.Box(c).Height(200)
	size := ui.Local(root, "size", func() float32 { return 180 })
	root.Children(func() {
		if layout.SplitPane(c, size, func() { tile(c, "Files", 200) }, func() { tile(c, "Editor", 200) }).Changed() {
			log(fmt.Sprintf("width: %.0f", *size))
		}
	})
}

func demoSplitStack(c *ui.Context, log func(string)) {
	root := ui.Box(c).Height(240)
	size := ui.Local(root, "size", func() float32 { return 120 })
	root.Children(func() {
		if layout.SplitStack(c, size, func() { tile(c, "Editor", 240) }, func() { tile(c, "Console", 240) }).Changed() {
			log(fmt.Sprintf("height: %.0f", *size))
		}
	})
}

func demoSidebarLayout(c *ui.Context, log func(string)) {
	root := ui.Box(c).Height(220)
	open := ui.Local(root, "open", func() bool { return true })
	root.Children(func() {
		if layout.SidebarLayout(c, open, func() {
			ui.Text(c, "Inbox")
			ui.Text(c, "Drafts")
		}, func() { ui.Text(c, "Message list").Padding(12) }).Changed() {
			log(fmt.Sprintf("open: %v", *open))
		}
	})
}

func demoHolyGrail(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	aside := ui.Local(root, "aside", func() bool { return true })
	root.Children(func() {
		if ui.Toggle(c, aside, "Aside").Changed() {
			log(fmt.Sprintf("aside: %v", *aside))
		}
		s := layout.HolyGrailSlots{
			Header: func() { ui.Text(c, "Header") }, Nav: func() { ui.Text(c, "Nav") },
			Main: func() { ui.Text(c, "Main") }, Footer: func() { ui.Text(c, "Footer") },
		}
		if *aside {
			s.Aside = func() { ui.Text(c, "Aside") }
		}
		ui.Box(c).Height(220).Children(func() { layout.HolyGrail(c, s) })
	})
}

func demoAppFrame(c *ui.Context, log func(string)) {
	ui.Box(c).Height(200).Border(1, c.Theme().Border).Children(func() {
		layout.AppFrame(c, func() { ui.Text(c, "Project").Bold() },
			func() { ui.Text(c, "Content").Padding(12) },
			func() { kit.Muted(c, "Ready") })
	})
}

func demoMasonry(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "items", 7, log)
		var items []func()
		for i := range *n {
			items = append(items, func() { tile(c, fmt.Sprintf("Item %d", i+1), float32(40+i%3*30)) })
		}
		layout.Masonry(c, 3, items)
	})
}

func demoAutoGrid(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "cells", 8, log)
		var items []func()
		for i := range *n {
			items = append(items, func() { tile(c, fmt.Sprintf("Cell %d", i+1), 48) })
		}
		layout.AutoGrid(c, 120, items)
	})
}

func demoBento(c *ui.Context, log func(string)) {
	cell := func(s string) func() { return func() { ui.Text(c, s).Bold() } }
	layout.Bento(c, 3, []layout.BentoTile{
		{Cols: 2, Content: cell("Analytics")}, {Rows: 2, Content: cell("Team")},
		{Content: cell("Inbox")}, {Content: cell("Billing")},
	})
}

func demoCluster(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	align := ui.Local(root, "align", func() int { return 0 })
	root.Children(func() {
		if ui.Segmented(c, align, "Start", "Center", "End").Changed() {
			log("align: " + []string{"start", "center", "end"}[*align])
		}
		var items []func()
		for _, s := range []string{"Design", "Code", "Review", "Ship", "Measure"} {
			items = append(items, func() { ui.Badge(c, s) })
		}
		layout.Cluster(c, []ui.Align{ui.Start, ui.Center, ui.End}[*align], items)
	})
}

func demoColumns(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "columns", 3, log)
		var items []func()
		for i := range max(*n, 1) {
			items = append(items, func() { tile(c, fmt.Sprintf("Column %d", i+1), 60) })
		}
		layout.Columns(c, items)
	})
}

func demoAspectBox(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(320)
	wide := ui.Local(root, "wide", func() bool { return true })
	root.Children(func() {
		if ui.Toggle(c, wide, "16:9").Changed() {
			log(fmt.Sprintf("16:9: %v", *wide))
		}
		ratio := float32(1)
		if *wide {
			ratio = 16.0 / 9
		}
		layout.AspectBox(c, ratio, func() { tile(c, "Video", 400) })
	})
}

func demoFrame(c *ui.Context, log func(string)) {
	layout.Frame(c, func() { ui.Box(c).Size(240, 120).Background(c.Theme().Accent) })
}

func demoScrollArea(c *ui.Context, log func(string)) {
	layout.ScrollArea(c, func() {
		for i := range 30 {
			ui.Text(c, fmt.Sprintf("Entry %d", i+1)).Padding(4, 8)
		}
	}).Height(180)
}

func demoScroller(c *ui.Context, log func(string)) {
	_, back, forward := layout.Scroller(c, func() {
		for i := range 14 {
			ui.Column(c).Size(110, 70).Children(func() { tile(c, fmt.Sprintf("Card %d", i+1), 70) })
		}
	})
	if back.Clicked() {
		log("back")
	}
	if forward.Clicked() {
		log("forward")
	}
}

func demoStickyHeader(c *ui.Context, log func(string)) {
	layout.StickyHeader(c, func() { ui.Text(c, "Contacts").Bold() }, func() {
		for i := range 30 {
			ui.Text(c, fmt.Sprintf("Person %d", i+1)).Padding(4, 12)
		}
	}).Height(200)
}
