package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/egoist/mygo/ui"
)

var layoutEntries2 = []Entry{
	{26, "Collapse", "layout", "Content whose height springs open and closed.", demoCollapse},
	{27, "Reveal", "layout", "Content rising in as it scrolls into view.", demoReveal},
	{28, "Divider", "layout", "Hairline, labeled or not, growing from center.", demoDivider},
	{29, "VDivider", "layout", "Vertical hairline growing from center.", demoVDivider},
	{30, "Gutter", "layout", "Room of n theme units; its size transitions.", demoGutter},
	{31, "Bar", "layout", "Row of start, center and end slots.", demoBar},
	{32, "PageHeader", "layout", "Title, subtitle and actions; the title slides in.", demoPageHeader},
	{33, "PageFooter", "layout", "Muted footer row that fades in.", demoPageFooter},
	{34, "Pane", "layout", "Titled pane that folds; its chevron turns.", demoPane},
	{35, "Dock", "layout", "Panel sliding in from an edge.", demoDock},
	{36, "Rail", "layout", "Icon column with a sliding indicator.", demoRail},
	{37, "Tiles", "layout", "Square tiles that lift on hover.", demoTiles},
	{38, "Showcase", "layout", "Centered stage; content scales in.", demoShowcase},
	{39, "Ratio", "layout", "Row split by weights that spring.", demoRatio},
	{40, "Responsive", "layout", "Row when wide, column when narrow.", demoResponsive},
	{41, "Breakpoint", "layout", "Slot for the window's size class.", demoBreakpoint},
	{42, "Stage", "layout", "Fixed-size content scaled to fit.", demoStage},
	{43, "Corner", "layout", "Content at a corner; it pops in.", demoCorner},
	{44, "Float", "layout", "Content at an offset that springs.", demoFloat},
	{45, "Shelf", "layout", "Titled sideways row with a see-all link.", demoShelf},
	{46, "Island", "layout", "Floating pill that morphs to its content.", demoIsland},
	{47, "Board", "layout", "Sideways columns that slide into place.", demoBoard},
	{48, "Ladder", "layout", "Rows alternating either side of a line.", demoLadder},
	{49, "Grouped", "layout", "Rows in one bordered group.", demoGrouped},
	{50, "Sheet", "layout", "Centered paper page that rises in.", demoSheet},
}

func demoCollapse(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	open := ui.Local(root, "open", func() bool { return true })
	root.Children(func() {
		if ui.Toggle(c, open, "Details").Changed() {
			log(fmt.Sprintf("open: %v", *open))
		}
		layout.Collapse(c, *open, func() { tile(c, "Shipping in two days", 80) })
	})
}

func demoReveal(c *ui.Context, log func(string)) {
	ui.Scroll(c).Height(220).Children(func() {
		ui.Column(c).Gap(12).Padding(4).Children(func() {
			for i := range 12 {
				c.Key(i)
				layout.Reveal(c, func() { tile(c, fmt.Sprintf("Story %d", i+1), 60) })
			}
		})
	})
}

func demoDivider(c *ui.Context, log func(string)) {
	ui.Column(c).Gap(16).Children(func() {
		ui.Text(c, "Sign in with email")
		layout.Divider(c, layout.DividerOptions{Label: "or"})
		ui.Text(c, "Continue with a passkey")
		layout.Divider(c, layout.DividerOptions{})
	})
}

func demoVDivider(c *ui.Context, log func(string)) {
	ui.Row(c).Height(32).Gap(12).Children(func() {
		ui.Text(c, "Edit")
		layout.VDivider(c)
		ui.Text(c, "View")
		layout.VDivider(c)
		ui.Text(c, "Help")
	})
}

func demoGutter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "units", 4, log)
		ui.Row(c).Children(func() {
			tile(c, "Left", 60)
			layout.Gutter(c, float32(*n))
			tile(c, "Right", 60)
		})
	})
}

func demoBar(c *ui.Context, log func(string)) {
	layout.Bar(c, func() { ui.Text(c, "Back") }, func() { ui.Text(c, "Inbox").Bold() }, func() {
		if ui.Button(c, "Edit").Clicked() {
			log("edit")
		}
	})
}

func demoPageHeader(c *ui.Context, log func(string)) {
	layout.PageHeader(c, "Projects", layout.PageHeaderOptions{
		Subtitle: "Everything your team is working on.",
		Actions: func() {
			if ui.Button(c, "New project").Clicked() {
				log("new project")
			}
		},
	})
}

func demoPageFooter(c *ui.Context, log func(string)) {
	layout.PageFooter(c, func() {
		ui.Text(c, "© 2026 Crychic")
		ui.Text(c, "Privacy")
		ui.Text(c, "Terms")
	})
}

func demoPane(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	open := ui.Local(root, "open", func() bool { return true })
	root.Children(func() {
		if layout.Pane(c, "Properties", open, func() {
			ui.Text(c, "Width 320")
			ui.Text(c, "Height 180")
		}).Changed() {
			log(fmt.Sprintf("open: %v", *open))
		}
	})
}

func demoDock(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	open := ui.Local(root, "open", func() bool { return true })
	root.Children(func() {
		if ui.Toggle(c, open, "Inspector").Changed() {
			log(fmt.Sprintf("open: %v", *open))
		}
		ui.Box(c).Height(220).Clip().Children(func() {
			tile(c, "Canvas", 220)
			layout.Dock(c, layout.DockRight, *open, func() { ui.Text(c, "Inspector").Bold() })
		})
	})
}

func demoRail(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	active := ui.Local(root, "active", func() int { return 0 })
	items := []layout.RailItem{
		{Icon: "house", Label: "Home"}, {Icon: "search", Label: "Search"},
		{Icon: "bell", Label: "Alerts"}, {Icon: "settings", Label: "Settings"},
	}
	root.Children(func() {
		if layout.Rail(c, items, active).Changed() {
			log("active: " + items[*active].Label)
		}
	})
}

func demoTiles(c *ui.Context, log func(string)) {
	var items []func()
	for i := range 8 {
		items = append(items, func() { ui.Text(c, fmt.Sprintf("App %d", i+1)) })
	}
	layout.Tiles(c, 96, items)
}

func demoShowcase(c *ui.Context, log func(string)) {
	layout.Showcase(c, func() {
		if ui.Button(c, "Featured").Clicked() {
			log("featured")
		}
	})
}

func demoRatio(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	wide := ui.Local(root, "wide", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, wide, "Wide main").Changed() {
			log(fmt.Sprintf("wide main: %v", *wide))
		}
		weights := []float32{1, 1, 1}
		if *wide {
			weights[1] = 3
		}
		layout.Ratio(c, weights, []func(){
			func() { tile(c, "Nav", 60) }, func() { tile(c, "Main", 60) }, func() { tile(c, "Aside", 60) },
		})
	})
}

func demoResponsive(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	narrow := ui.Local(root, "narrow", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, narrow, "Narrow").Changed() {
			log(fmt.Sprintf("narrow: %v", *narrow))
		}
		w := float32(520)
		if *narrow {
			w = 260
		}
		box := ui.Box(c)
		box.Width(box.AnimateTo("w", w, ui.Snappy)).Children(func() {
			layout.Responsive(c, 400, []func(){
				func() { tile(c, "One", 50) }, func() { tile(c, "Two", 50) }, func() { tile(c, "Three", 50) },
			})
		})
	})
}

func demoBreakpoint(c *ui.Context, log func(string)) {
	layout.Breakpoint(c, func() { tile(c, "Compact window", 60) }, func() { tile(c, "Regular window", 60) },
		func() { tile(c, "Wide window", 60) })
}

func demoStage(c *ui.Context, log func(string)) {
	ui.Box(c).Height(200).Children(func() {
		layout.Stage(c, 320, 240, func() { tile(c, "320 × 240 slide", 240) })
	})
}

func demoCorner(c *ui.Context, log func(string)) {
	ui.Box(c).Height(160).Children(func() {
		tile(c, "Photo", 160)
		layout.Corner(c, layout.CornerTopRight, func() { ui.Text(c, "New").Bold() })
		layout.Corner(c, layout.CornerBottomLeft, func() { ui.Text(c, "4 MB") })
	})
}

func demoFloat(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	spot := ui.Local(root, "spot", func() int { return 0 })
	root.Children(func() {
		if ui.Button(c, "Move").Clicked() {
			*spot = (*spot + 1) % 3
			log(fmt.Sprintf("spot: %d", *spot))
		}
		ui.Box(c).Height(140).Children(func() {
			tile(c, "Map", 140)
			layout.Float(c, float32(40+*spot*120), float32(20+*spot*30), func() { ui.Text(c, "Pin").Bold() })
		})
	})
}

func demoShelf(c *ui.Context, log func(string)) {
	_, seeAll := layout.Shelf(c, "Recently played", layout.ShelfOptions{}, func() {
		for i := range 10 {
			ui.Column(c).Width(120).Children(func() { tile(c, fmt.Sprintf("Album %d", i+1), 120) })
		}
	})
	if seeAll.Clicked() {
		log("see all")
	}
}

func demoIsland(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	root.Children(func() {
		n := counter(c, root, "items", 1, log)
		layout.Island(c, func() {
			for i := range *n {
				ui.Text(c, fmt.Sprintf("Task %d", i+1))
			}
		})
	})
}

func demoBoard(c *ui.Context, log func(string)) {
	var cols []func()
	for _, name := range []string{"Backlog", "Doing", "Review", "Done"} {
		cols = append(cols, func() {
			ui.Text(c, name).Bold()
			tile(c, name+" card", 60)
		})
	}
	layout.Board(c, cols)
}

func demoLadder(c *ui.Context, log func(string)) {
	var rows []func()
	for i, s := range []string{"Founded", "First release", "Ten thousand users", "Version two"} {
		rows = append(rows, func() { ui.Text(c, fmt.Sprintf("%d · %s", 2020+i, s)) })
	}
	layout.Ladder(c, rows)
}

func demoGrouped(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "rows", 3, log)
		var rows []func()
		for i := range *n {
			rows = append(rows, func() { ui.Text(c, fmt.Sprintf("Setting %d", i+1)) })
		}
		layout.Grouped(c, rows)
	})
}

func demoSheet(c *ui.Context, log func(string)) {
	layout.Sheet(c, func() {
		ui.Text(c, "Quarterly letter").Bold()
		ui.Text(c, "Dear shareholders, the quarter went well.")
	})
}
