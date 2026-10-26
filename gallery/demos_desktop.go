package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/desktop"
	"github.com/egoist/mygo/ui"
)

var desktopEntries = []Entry{
	{ID: 1456, Name: "WindowFrame", Category: "desktop", Summary: "Window with title bar and body.", Demo: demoDesktopWindowFrame},
	{ID: 1457, Name: "TrafficLights", Category: "desktop", Summary: "Close, minimize, zoom with hover glyphs; cue Tap.", Demo: demoDesktopTrafficLights},
	{ID: 1458, Name: "TitleTabs", Category: "desktop", Summary: "Tabs in a title bar; cue Select.", Demo: demoDesktopTitleTabs},
	{ID: 1459, Name: "MenuBar", Category: "desktop", Summary: "App menu bar; cue Open.", Demo: demoDesktopMenuBar},
	{ID: 1460, Name: "StatusBar", Category: "desktop", Summary: "Window status bar.", Demo: demoDesktopStatusBar},
	{ID: 1461, Name: "Dock", Category: "desktop", Summary: "Dock that magnifies; cue Tap.", Demo: demoDesktopDock},
	{ID: 1462, Name: "Launcher", Category: "desktop", Summary: "App grid with search; cue Select.", Demo: demoDesktopLauncher},
	{ID: 1463, Name: "WindowSwitcher", Category: "desktop", Summary: "Alt-tab strip; cue Tick.", Demo: demoDesktopWindowSwitcher},
	{ID: 1464, Name: "DesktopIcon", Category: "desktop", Summary: "Icon with label; select and open; cue Open.", Demo: demoDesktopIcon},
	{ID: 1465, Name: "NotificationCenter", Category: "desktop", Summary: "Notifications by app.", Demo: demoDesktopNotificationCenter},
	{ID: 1466, Name: "ControlCenter", Category: "desktop", Summary: "Toggles and sliders; cue Toggle.", Demo: demoDesktopControlCenter},
	{ID: 1467, Name: "VolumeHUD", Category: "desktop", Summary: "Volume overlay that fades.", Demo: demoDesktopVolumeHUD},
	{ID: 1468, Name: "BrightnessHUD", Category: "desktop", Summary: "Brightness overlay.", Demo: demoDesktopBrightnessHUD},
	{ID: 1469, Name: "BatteryItem", Category: "desktop", Summary: "Menu bar battery.", Demo: demoDesktopBatteryItem},
	{ID: 1470, Name: "WifiItem", Category: "desktop", Summary: "Menu bar wifi with networks; cue Select.", Demo: demoDesktopWifiItem},
}

func demoDesktopWindowFrame(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		desktop.WindowFrame(c, "Glaze — Kiln calendar", desktop.WindowFrameOptions{}, func() {
			t := c.Theme()
			ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
				for _, row := range [][2]string{{"Kiln 2 · cone 6", "Fri 8:00"}, {"Raku pit", "Sat 10:00"}, {"Kiln 1 · bisque", "Mon 7:30"}} {
					ui.Row(c).FillWidth().Children(func() {
						ui.Text(c, row[0]).FontSize(t.Rem(0.84)).Bold().Grow(1)
						ui.Text(c, row[1]).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
					})
				}
			})
		})
	})
}

func demoDesktopTrafficLights(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Children(func() {
		_, close, minimize, zoom := desktop.TrafficLights(c, desktop.TrafficLightsOptions{})
		if close.Clicked() {
			log("close")
		}
		if minimize.Clicked() {
			log("minimize")
		}
		if zoom.Clicked() {
			log("zoom")
		}
	})
}

func demoDesktopTitleTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	picked := intState(root, "picked", 0)
	root.Children(func() {
		r, add := desktop.TitleTabs(c, []string{"Kiln calendar", "Members", "Invoices"}, picked, desktop.TitleTabsOptions{})
		if r.Changed() {
			log(fmt.Sprintf("tab %d", *picked))
		}
		if add.Clicked() {
			log("new tab")
		}
	})
}

func demoDesktopMenuBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480).Height(260)
	chosen := textState(root, "chosen", "")
	root.Children(func() {
		r := desktop.MenuBar(c, "Glaze", []desktop.Menu{
			{Title: "File", Items: []desktop.MenuItem{{Label: "New booking", Shortcut: "⌘N"}, {Label: "Open…", Shortcut: "⌘O"}, {Label: "Export PDF", Shortcut: "⇧⌘E"}}},
			{Title: "Edit", Items: []desktop.MenuItem{{Label: "Undo", Shortcut: "⌘Z"}, {Label: "Redo", Shortcut: "⇧⌘Z"}}},
			{Title: "View", Items: []desktop.MenuItem{{Label: "Show calendar", Shortcut: "⌘1"}, {Label: "Show members", Shortcut: "⌘2"}}},
		}, chosen)
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

func demoDesktopStatusBar(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		desktop.StatusBar(c,
			[]desktop.StatusItem{{Icon: "check", Text: "Synced", Tone: 1}, {Text: "3 firings today"}},
			[]desktop.StatusItem{{Icon: "wifi", Text: "Online", Tone: 1}, {Text: "Glaze 3.2"}})
	})
}

func demoDesktopDock(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520).Padding(40, 0, 0, 0).AlignItems(ui.Center)
	opened := textState(root, "opened", "")
	root.Children(func() {
		if desktop.Dock(c, []desktop.DockApp{
			{Name: "Glaze", Icon: "flame", Running: true, Badge: 3}, {Name: "Mail", Icon: "mail", Running: true, Badge: 12},
			{Name: "Calendar", Icon: "calendar"}, {Name: "Photos", Icon: "image"}, {Name: "Music", Icon: "music"},
			{Name: "Terminal", Icon: "square-terminal", Running: true}, {Name: "Settings", Icon: "settings-2"},
		}, opened).Changed() {
			log("opened " + *opened)
		}
	})
}

func demoDesktopLauncher(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	query, picked := textState(root, "query", ""), textState(root, "picked", "")
	root.Children(func() {
		before := *picked
		r := desktop.Launcher(c, []desktop.LaunchApp{
			{Name: "Glaze", Icon: "flame", Category: "Studio"}, {Name: "Calendar", Icon: "calendar", Category: "Time"},
			{Name: "Mail", Icon: "mail", Category: "Messages"}, {Name: "Photos", Icon: "image", Category: "Media"},
			{Name: "Music", Icon: "music", Category: "Media"}, {Name: "Notes", Icon: "notebook-pen", Category: "Writing"},
			{Name: "Terminal", Icon: "square-terminal", Category: "Developer"}, {Name: "Settings", Icon: "settings-2", Category: "System"},
		}, query, picked, desktop.LauncherOptions{})
		if r.Changed() {
			if *picked != before {
				log("picked " + *picked)
			} else {
				log(fmt.Sprintf("query %d runes", len([]rune(*query))))
			}
		}
	})
}

func demoDesktopWindowSwitcher(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520).AlignItems(ui.Center)
	current := intState(root, "current", 0)
	root.Children(func() {
		r, _ := desktop.WindowSwitcher(c, []desktop.SwitchWindow{
			{App: "Glaze", Title: "Kiln calendar — week of Oct 26", Icon: "flame"},
			{App: "Mail", Title: "Inbox — 12 unread", Icon: "mail"},
			{App: "Photos", Title: "Spring glaze tests", Icon: "image"},
			{App: "Terminal", Title: "glaze-cli — zsh", Icon: "square-terminal"},
		}, current, desktop.WindowSwitcherOptions{})
		if r.Changed() {
			log(fmt.Sprintf("window %d", *current))
		}
	})
}

func demoDesktopIcon(c *ui.Context, log func(string)) {
	ui.Row(c).Width(480).Gap(8).Children(func() {
		for _, f := range [][2]string{{"Firing logs", "folder"}, {"Shelf map.pdf", "file-text"}, {"Members.csv", "file-spreadsheet"}} {
			c.Key(f[0])
			box := ui.Box(c)
			selected, opened := boolState(box, "sel", false), intState(box, "opened", 0)
			box.Children(func() {
				before := *opened
				if desktop.DesktopIcon(c, f[0], f[1], selected, opened).Changed() {
					if *opened != before {
						log(fmt.Sprintf("opened %s %d", f[0], *opened))
					} else {
						log("selected " + f[0])
					}
				}
			})
		}
	})
}

func demoDesktopNotificationCenter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	cleared := boolState(root, "cleared", false)
	at := func(h, m int) time.Time { return time.Date(2026, 10, 26, h, m, 0, 0, time.UTC) }
	root.Children(func() {
		var notes []desktop.DesktopNote
		if !*cleared {
			notes = []desktop.DesktopNote{
				{ID: "g1", App: "Glaze", Icon: "flame", Title: "Kiln 2 is cooling", Text: "Your cone 6 firing peaked at 2232°F.", Time: at(9, 40)},
				{ID: "g2", App: "Glaze", Icon: "flame", Title: "Waitlist spot open", Text: "Priya Nair can join Tuesday's class.", Time: at(9, 12)},
				{ID: "m1", App: "Mail", Icon: "mail", Title: "Daniel Okafor", Text: "Photos from the shoot are up.", Time: at(8, 55)},
				{ID: "c1", App: "Calendar", Icon: "calendar", Title: "Glaze mixing at 11:00", Time: at(8, 30)},
			}
		}
		_, clear := desktop.NotificationCenter(c, notes, desktop.NotificationCenterOptions{})
		if clear.Clicked() {
			*cleared = true
			log("cleared")
		}
	})
}

func demoDesktopControlCenter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	s := ui.Local(root, "controls", func() desktop.Controls { return desktop.Controls{Wifi: true, Volume: 0.6, Brightness: 0.8} })
	root.Children(func() {
		if desktop.ControlCenter(c, s, desktop.ControlCenterOptions{}).Changed() {
			log(fmt.Sprintf("wifi %v, bluetooth %v, focus %v, dark %v, sound %.0f%%, display %.0f%%",
				s.Wifi, s.Bluetooth, s.Focus, s.Dark, s.Volume*100, s.Brightness*100))
		}
	})
}

func demoDesktopVolumeHUD(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(12).AlignItems(ui.Center)
	level := intState(root, "level", 6)
	root.Children(func() {
		desktop.VolumeHUD(c, float32(*level)/10, desktop.VolumeHUDOptions{})
		if button.Outline(c, "Volume up").Disabled(*level >= 10).Clicked() {
			*level++
			log(fmt.Sprintf("volume %d", *level*10))
		}
	})
}

func demoDesktopBrightnessHUD(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(12).AlignItems(ui.Center)
	level := intState(root, "level", 4)
	root.Children(func() {
		desktop.BrightnessHUD(c, float32(*level)/10, desktop.BrightnessHUDOptions{})
		if button.Outline(c, "Brighter").Disabled(*level >= 10).Clicked() {
			*level++
			log(fmt.Sprintf("brightness %d", *level*10))
		}
	})
}

func demoDesktopBatteryItem(c *ui.Context, log func(string)) {
	ui.Row(c).Width(360).Gap(12).Children(func() {
		desktop.BatteryItem(c, 0.82, desktop.BatteryItemOptions{})
		desktop.BatteryItem(c, 0.14, desktop.BatteryItemOptions{})
		desktop.BatteryItem(c, 0.64, desktop.BatteryItemOptions{Charging: true})
	})
}

func demoDesktopWifiItem(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Height(300)
	current := textState(root, "current", "Kiln Studio")
	root.Children(func() {
		before := *current
		r := desktop.WifiItem(c, []desktop.Network{
			{Name: "Kiln Studio", Bars: 3, Secure: true}, {Name: "Kiln Guest", Bars: 2},
			{Name: "Riverbend Cafe", Bars: 1, Secure: true}, {Name: "Pottery Co-op", Bars: 2, Secure: true},
		}, current, desktop.WifiItemOptions{})
		if r.Changed() {
			if *current == before {
				log("menu open")
			} else {
				log("joined " + *current)
			}
		}
	})
}
