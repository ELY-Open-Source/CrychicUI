package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/desktop"
	"github.com/egoist/mygo/ui"
)

var desktopEntries2 = []Entry{
	{ID: 1471, Name: "ClockItem", Category: "desktop", Summary: "Menu bar clock.", Demo: demoDesktopClockItem},
	{ID: 1472, Name: "InspectorPane", Category: "desktop", Summary: "Inspector with sections.", Demo: demoDesktopInspectorPane},
	{ID: 1473, Name: "UnifiedToolbar", Category: "desktop", Summary: "Toolbar with title and search.", Demo: demoDesktopUnifiedToolbar},
	{ID: 1474, Name: "SourceList", Category: "desktop", Summary: "Sidebar source list; cue Select.", Demo: demoDesktopSourceList},
	{ID: 1475, Name: "Preferences", Category: "desktop", Summary: "Preferences window with tabs; cue Select.", Demo: demoDesktopPreferences},
	{ID: 1476, Name: "About", Category: "desktop", Summary: "About window.", Demo: demoDesktopAbout},
	{ID: 1477, Name: "ProgressWindow", Category: "desktop", Summary: "Task progress window.", Demo: demoDesktopProgressWindow},
	{ID: 1478, Name: "UpdatePrompt", Category: "desktop", Summary: "Update available with install; cue Success.", Demo: demoDesktopUpdatePrompt},
	{ID: 1479, Name: "PermissionPrompt", Category: "desktop", Summary: "Allow access prompt; cue Tap.", Demo: demoDesktopPermissionPrompt},
	{ID: 1480, Name: "TrayMenu", Category: "desktop", Summary: "Tray menu; cue Open.", Demo: demoDesktopTrayMenu},
	{ID: 1481, Name: "Spaces", Category: "desktop", Summary: "Workspace dots; cue Select.", Demo: demoDesktopSpaces},
	{ID: 1482, Name: "SnapLayouts", Category: "desktop", Summary: "Layout chooser; cue Select.", Demo: demoDesktopSnapLayouts},
	{ID: 1483, Name: "TaskManager", Category: "desktop", Summary: "Processes with CPU.", Demo: demoDesktopTaskManager},
	{ID: 1484, Name: "TerminalWindow", Category: "desktop", Summary: "Terminal in a window frame.", Demo: demoDesktopTerminalWindow},
	{ID: 1485, Name: "LockScreen", Category: "desktop", Summary: "Clock and unlock; cue Success.", Demo: demoDesktopLockScreen},
}

func demoDesktopClockItem(c *ui.Context, log func(string)) {
	t := c.Theme()
	ui.Row(c).Width(480).Gap(4).AlignItems(ui.Center).Padding(2, 8).Background(t.SurfaceHover).Border(1, t.Border).Radius(t.Radius).Children(func() {
		ui.Text(c, "Glaze").FontSize(t.Rem(0.8)).Bold().Grow(1)
		desktop.BatteryItem(c, 0.82, desktop.BatteryItemOptions{})
		desktop.ClockItem(c, desktop.ClockItemOptions{})
	})
}

func demoDesktopInspectorPane(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Children(func() {
		if desktop.InspectorPane(c, []desktop.InspectorSection{
			{Title: "Booking", Rows: []desktop.InspectorRow{{Label: "Kiln", Value: "Kiln 2"}, {Label: "Cone", Value: "6"}, {Label: "Shelves", Value: "2 of 6"}}},
			{Title: "Firing", Rows: []desktop.InspectorRow{{Label: "Start", Value: "Fri 08:00"}, {Label: "Peak", Value: "2232°F"}, {Label: "Hold", Value: "10 min"}}},
			{Title: "Member", Rows: []desktop.InspectorRow{{Label: "Name", Value: "Maya Chen"}, {Label: "Plan", Value: "Studio"}}},
		}).Changed() {
			log("section toggled")
		}
	})
}

func demoDesktopUnifiedToolbar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	query := textState(root, "query", "")
	root.Children(func() {
		r, back, forward := desktop.UnifiedToolbar(c, "Kiln calendar", query, desktop.UnifiedToolbarOptions{Subtitle: "Week of Oct 26"})
		if back.Clicked() {
			log("back")
		}
		if forward.Clicked() {
			log("forward")
		}
		if r.Changed() {
			log(fmt.Sprintf("query %d runes", len([]rune(*query))))
		}
	})
}

func demoDesktopSourceList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	current := textState(root, "current", "Today")
	root.Children(func() {
		if desktop.SourceList(c, []desktop.SourceGroup{
			{Title: "Library", Items: []desktop.SourceItem{{Name: "All bookings", Icon: "calendar", Count: 24}, {Name: "Today", Icon: "sun", Count: 3}, {Name: "Waitlist", Icon: "clock", Count: 5}}},
			{Title: "Kilns", Items: []desktop.SourceItem{{Name: "Kiln 1", Icon: "flame"}, {Name: "Kiln 2", Icon: "flame", Count: 2}, {Name: "Raku pit", Icon: "flame"}}},
		}, current).Changed() {
			log("source " + *current)
		}
	})
}

func demoDesktopPreferences(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	picked := intState(root, "picked", 0)
	panes := [][2]string{{"Language", "English"}, {"Studio", "Riverbend Clay Co."}, {"Email digests", "Daily at 7:00"}, {"Default cone", "6"}}
	root.Children(func() {
		if desktop.Preferences(c, []desktop.PrefTab{{Name: "General", Icon: "settings-2"}, {Name: "Accounts", Icon: "user"},
			{Name: "Notifications", Icon: "bell"}, {Name: "Kilns", Icon: "flame"}}, picked, func() {
			t := c.Theme()
			ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, panes[*picked][0]).FontSize(t.Rem(0.84)).TextColor(t.TextMuted).Grow(1)
				ui.Text(c, panes[*picked][1]).FontSize(t.Rem(0.84)).Bold()
			})
		}).Changed() {
			log(fmt.Sprintf("pane %d", *picked))
		}
	})
}

func demoDesktopAbout(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		desktop.About(c, "Glaze", desktop.AboutOptions{Icon: "flame", Version: "3.2.0 (412)",
			Credits: "Made in Portland by the Glaze team.", Copyright: "© 2026 Glaze Studio Software"})
	})
}

func demoDesktopProgressWindow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	done := intState(root, "done", 12)
	root.Children(func() {
		_, cancel := desktop.ProgressWindow(c, "Exporting firing logs", *done, 40, desktop.ProgressWindowOptions{Detail: "Kiln 2 · October 2026"})
		if cancel.Clicked() {
			log("closed")
		}
		if button.Outline(c, "Export 4 more").Disabled(*done >= 40).Clicked() {
			*done = min(40, *done+4)
			log(fmt.Sprintf("done %d", *done))
		}
	})
}

func demoDesktopUpdatePrompt(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		_, install, later := desktop.UpdatePrompt(c, "Glaze", "3.3.0", []string{
			"Waitlists fill a seat when someone cancels.", "Kiln calendar loads twice as fast.", "Shelf maps print on one page.",
		}, desktop.UpdatePromptOptions{Icon: "flame"})
		if install.Clicked() {
			log("install")
		}
		if later.Clicked() {
			log("later")
		}
	})
}

func demoDesktopPermissionPrompt(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	access := ui.Local(root, "access", func() desktop.Access { return desktop.Undecided })
	root.Children(func() {
		if desktop.PermissionPrompt(c, "Glaze", "camera", access, desktop.PermissionPromptOptions{Icon: "camera",
			Text: "Scan shelf tags to log which piece sits where."}).Changed() {
			log(fmt.Sprintf("access %d", *access))
		}
	})
}

func demoDesktopTrayMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Height(260)
	chosen := textState(root, "chosen", "")
	root.Children(func() {
		r := desktop.TrayMenu(c, []desktop.TrayItem{
			{Label: "Open Glaze", Icon: "flame"}, {Label: "Today's firings", Icon: "calendar", Shortcut: "⌘T"},
			{Label: "Pause sync", Icon: "pause"}, {Label: "Quit Glaze", Icon: "power", Shortcut: "⌘Q"},
		}, chosen, desktop.TrayMenuOptions{Label: "Glaze", Icon: "flame"})
		if r.Changed() {
			if *chosen == "" {
				log("menu open")
			} else {
				log("chose " + *chosen)
				*chosen = ""
			}
		}
	})
}

func demoDesktopSpaces(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).AlignItems(ui.Start)
	current := intState(root, "current", 0)
	root.Children(func() {
		if desktop.Spaces(c, []string{"Studio", "Bookkeeping", "Personal"}, current).Changed() {
			log(fmt.Sprintf("space %d", *current))
		}
	})
}

func demoDesktopSnapLayouts(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).AlignItems(ui.Start)
	picked := ui.Local(root, "snap", func() desktop.Snap { return desktop.Snap{Layout: -1} })
	root.Children(func() {
		if desktop.SnapLayouts(c, picked, desktop.SnapLayoutsOptions{}).Changed() {
			log(fmt.Sprintf("snap %d %d", picked.Layout, picked.Zone))
		}
	})
}

func demoDesktopTaskManager(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		desktop.TaskManager(c, []desktop.Process{
			{Name: "Glaze", Icon: "flame", CPU: 18.4, Memory: 412 << 20},
			{Name: "Kiln sync", Icon: "refresh-cw", CPU: 4.2, Memory: 96 << 20},
			{Name: "Browser", Icon: "globe", CPU: 31.0, Memory: 1310 << 20},
			{Name: "Photos", Icon: "image", CPU: 2.1, Memory: 640 << 20},
			{Name: "Mail", Icon: "mail", CPU: 0.8, Memory: 220 << 20},
		}, desktop.TaskManagerOptions{})
	})
}

func demoDesktopTerminalWindow(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		desktop.TerminalWindow(c, "glaze-cli — zsh", []string{
			"$ glaze login", "Signed in as maya@riverbend.studio",
			"$ glaze book kiln-2 --cone 6", "Booked kiln 2: 2 shelves, Fri Oct 30 at 08:00",
		}, desktop.TerminalWindowOptions{})
	})
}

func demoDesktopLockScreen(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	password, unlocked := textState(root, "password", ""), boolState(root, "unlocked", false)
	tries := intState(root, "tries", 0)
	root.Children(func() {
		before := *tries
		check := func(s string) bool {
			*tries++
			return s == "glaze"
		}
		r := desktop.LockScreen(c, "Maya Chen", password, check, unlocked, desktop.LockScreenOptions{})
		switch {
		case !r.Changed():
		case *unlocked:
			log("unlocked")
		case *tries != before:
			log(fmt.Sprintf("attempt %d refused", *tries))
		default:
			log(fmt.Sprintf("password %d runes", len([]rune(*password))))
		}
	})
}
