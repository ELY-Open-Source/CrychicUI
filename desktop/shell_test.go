package desktop

import (
	"bytes"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestWindowFrameAndLights(t *testing.T) {
	presses := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		WindowFrame(c, "Kiln calendar", WindowFrameOptions{}, func() { ui.Text(c, "Kiln 2 · cone 6") })
		_, close, _, zoom := TrafficLights(c, TrafficLightsOptions{})
		if close.Clicked() || zoom.Clicked() {
			presses++
		}
	}), 480, 400)
	settles(t, tt)
	if !exact(tt, "Kiln 2 · cone 6") {
		t.Errorf("body missing: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Close")
	click(t, tt, "Zoom")
	cues(t, before, sound.Tap, sound.Tap)
	if presses != 2 {
		t.Errorf("presses %d", presses)
	}
	panics(t, "a title and a body", func(c *ui.Context) { WindowFrame(c, "Kiln", WindowFrameOptions{}, nil) })
}

func TestTitleTabsPickAndReorder(t *testing.T) {
	tabs, picked, adds := []string{"Kiln calendar", "Members", "Invoices"}, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, add := TitleTabs(c, tabs, &picked, TitleTabsOptions{}); add.Clicked() {
			adds++
		}
	}), 520, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Members")
	click(t, tt, "Members")
	click(t, tt, "New tab")
	cues(t, before, sound.Select, sound.Tap)
	if picked != 1 || adds != 1 {
		t.Errorf("picked %d, adds %d", picked, adds)
	}
	swapUnderPress(t, tt, "Invoices", func() { tabs[0], tabs[2] = tabs[2], tabs[0] })
	if picked != 1 {
		t.Errorf("a reorder under the pointer picked %q", tabs[picked])
	}
	panics(t, "picked one", func(c *ui.Context) { n := 3; TitleTabs(c, []string{"A"}, &n, TitleTabsOptions{}) })
	panics(t, "repeated", func(c *ui.Context) { TitleTabs(c, []string{"A", "A"}, new(int), TitleTabsOptions{}) })
}

func menus() []Menu {
	return []Menu{
		{Title: "File", Items: []MenuItem{{Label: "New booking", Shortcut: "⌘N"}, {Label: "Export PDF"}}},
		{Title: "Edit", Items: []MenuItem{{Label: "Undo", Shortcut: "⌘Z"}}},
	}
}

func TestMenuBarOpensOneAndChooses(t *testing.T) {
	chosen, changes := "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if MenuBar(c, "Glaze", menus(), &chosen).Changed() {
			changes++
		}
	}), 480, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "File")
	settles(t, tt)
	click(t, tt, "Edit")
	settles(t, tt)
	if r, ok := tt.Find("New booking"); ok && r.H > 0 {
		t.Errorf("opening Edit left File open: %v", r)
	}
	click(t, tt, "Undo")
	cues(t, before, sound.Open, sound.Open, sound.Tap)
	if chosen != "Undo" || changes != 3 {
		t.Errorf("chosen %q, changes %d", chosen, changes)
	}
	panics(t, "repeated", func(c *ui.Context) {
		MenuBar(c, "Glaze", []Menu{{Title: "A", Items: []MenuItem{{Label: "Close"}}}, {Title: "B", Items: []MenuItem{{Label: "Close"}}}}, new(string))
	})
}

func TestStatusBarTones(t *testing.T) {
	tt := ui.NewTester(within(480, func(c *ui.Context) {
		StatusBar(c, []StatusItem{{Icon: "check", Text: "Synced", Tone: 1}}, []StatusItem{{Text: "Glaze 3.2"}})
	}), 600, 200)
	settles(t, tt)
	if !exact(tt, "Synced") || !exact(tt, "Glaze 3.2") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "tone 0 to 3", func(c *ui.Context) { StatusBar(c, []StatusItem{{Text: "x", Tone: 4}}, nil) })
}

func TestDockReportsEveryClick(t *testing.T) {
	apps := []DockApp{{Name: "Glaze", Icon: "flame", Badge: 3}, {Name: "Mail", Icon: "mail", Badge: 120}, {Name: "Calendar", Icon: "calendar"}}
	opened, changes := "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Dock(c, apps, &opened).Changed() {
			changes++
		}
	}), 480, 300)
	settles(t, tt)
	if !exact(tt, "99") {
		t.Errorf("a badge over 99 is not capped: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Mail")
	settles(t, tt)
	click(t, tt, "Mail")
	cues(t, before, sound.Tap, sound.Tap)
	if opened != "Mail" || changes != 2 {
		t.Errorf("opened %q, changes %d", opened, changes)
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { Dock(c, apps, &opened) }), 480, 300)
	settles(t, tt)
	still := bytes.Clone(tt.Image().Pix)
	r := rectOf(t, tt, "Glaze")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	tt.Advance(time.Second)
	if bytes.Equal(still, tt.Image().Pix) {
		t.Error("hovering an app did not magnify it")
	}
	panics(t, "badge of 0 or more", func(c *ui.Context) { Dock(c, []DockApp{{Name: "A", Icon: "flame", Badge: -1}}, new(string)) })
}

func TestLauncherFiltersAndPicks(t *testing.T) {
	apps := []LaunchApp{{Name: "Calendar", Icon: "calendar", Category: "Time"}, {Name: "Mail", Icon: "mail", Category: "Messages"}}
	query, picked := "", ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Launcher(c, apps, &query, &picked, LauncherOptions{}) }), 480, 500)
	settles(t, tt)
	click(t, tt, "Search apps")
	tt.Type("tim")
	settles(t, tt)
	if !tt.HasText("Calendar") || tt.HasText("Mail") {
		t.Errorf("category filter: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Calendar")
	click(t, tt, "Calendar")
	cues(t, before, sound.Select)
	click(t, tt, "Search apps")
	tt.Type("zzz")
	settles(t, tt)
	if picked != "Calendar" || !exact(tt, "No apps match") {
		t.Errorf("picked %q, texts %q", picked, tt.Texts())
	}
	panics(t, "repeated", func(c *ui.Context) {
		Launcher(c, []LaunchApp{apps[0], apps[0]}, new(string), new(string), LauncherOptions{})
	})
}

func TestWindowSwitcherCycles(t *testing.T) {
	windows := []SwitchWindow{{App: "Glaze", Title: "Kiln calendar", Icon: "flame"}, {App: "Mail", Title: "Inbox", Icon: "mail"}, {App: "Photos", Title: "Glaze tests", Icon: "image"}}
	current := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { WindowSwitcher(c, windows, &current, WindowSwitcherOptions{}) }), 600, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Next window")
	click(t, tt, "Glaze tests")
	click(t, tt, "Next window")
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	if current != 0 {
		t.Errorf("next from the last window gave %d", current)
	}
	one := 0
	tt = ui.NewTester(padded(func(c *ui.Context) { WindowSwitcher(c, windows[:1], &one, WindowSwitcherOptions{}) }), 600, 400)
	settles(t, tt)
	before = baseline()
	click(t, tt, "Next window")
	cues(t, before)
}

func TestDesktopIconSelectsThenOpens(t *testing.T) {
	selected, opened, changes := false, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if DesktopIcon(c, "Firing logs", "folder", &selected, &opened).Changed() {
			changes++
		}
	}), 300, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Firing logs")
	tt.Advance(time.Second)
	r := rectOf(t, tt, "Firing logs")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Release(x, y)
	tt.Press(x, y)
	tt.Release(x, y)
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	cues(t, before, sound.Select, sound.Open, sound.Open)
	if !selected || opened != 2 || changes != 3 {
		t.Errorf("selected %v, opened %d, changes %d", selected, opened, changes)
	}
	panics(t, "0 or more", func(c *ui.Context) { n := -1; DesktopIcon(c, "A", "folder", new(bool), &n) })
}

func TestNotificationCenterGroupsAndClears(t *testing.T) {
	notes := []DesktopNote{{ID: "a", App: "Glaze", Icon: "flame", Title: "Kiln 2 is cooling"}, {ID: "b", App: "Glaze", Icon: "flame", Title: "Waitlist"},
		{ID: "c", App: "Mail", Icon: "mail", Title: "Daniel Okafor"}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, clear := NotificationCenter(c, notes, NotificationCenterOptions{}); clear.Clicked() {
			notes = nil
		}
	}), 480, 600)
	settles(t, tt)
	if !exact(tt, "2") || !exact(tt, "1") {
		t.Errorf("group counts: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Clear all")
	settles(t, tt)
	cues(t, before, sound.Delete)
	if !exact(tt, "No notifications") {
		t.Errorf("no empty state: %q", tt.Texts())
	}
	panics(t, "repeated", func(c *ui.Context) {
		NotificationCenter(c, []DesktopNote{notes0(), notes0()}, NotificationCenterOptions{})
	})
}

func notes0() DesktopNote { return DesktopNote{ID: "a", App: "Glaze", Icon: "flame", Title: "x"} }

func TestControlCenterToggles(t *testing.T) {
	s := Controls{Volume: 0.5, Brightness: 0.5}
	tt := ui.NewTester(padded(func(c *ui.Context) { ControlCenter(c, &s, ControlCenterOptions{}) }), 400, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Bluetooth")
	cues(t, before, sound.Toggle)
	settles(t, tt)
	if !s.Bluetooth || !exact(tt, "On") {
		t.Errorf("bluetooth %v, texts %q", s.Bluetooth, tt.Texts())
	}
	panics(t, "level 0 to 1", func(c *ui.Context) { ControlCenter(c, &Controls{Volume: 1.5}, ControlCenterOptions{}) })
}

func TestHUDFadesAfterChange(t *testing.T) {
	level := float32(0.5)
	tt := ui.NewTester(padded(func(c *ui.Context) { VolumeHUD(c, level, VolumeHUDOptions{}) }), 300, 300)
	tt.Advance(3 * time.Second)
	gone := bytes.Clone(tt.Image().Pix)
	level = 0.6
	tt.Frame()
	tt.Advance(300 * time.Millisecond)
	if bytes.Equal(gone, tt.Image().Pix) {
		t.Error("a level change did not show the overlay")
	}
	tt.Advance(3 * time.Second)
	if !bytes.Equal(gone, tt.Image().Pix) {
		t.Error("the overlay did not fade")
	}
	panics(t, "level 0 to 1", func(c *ui.Context) { BrightnessHUD(c, 2, BrightnessHUDOptions{}) })
}

func TestBatteryPercent(t *testing.T) {
	for _, c := range []struct {
		level float32
		want  string
	}{{0.996, "99%"}, {0.001, "1%"}, {1, "100%"}, {0, "0%"}, {0.5, "50%"}} {
		tt := ui.NewTester(padded(func(ctx *ui.Context) { BatteryItem(ctx, c.level, BatteryItemOptions{}) }), 300, 200)
		settles(t, tt)
		if !exact(tt, c.want) {
			t.Errorf("level %v shows %q, want %s", c.level, tt.Texts(), c.want)
		}
	}
	panics(t, "level 0 to 1", func(c *ui.Context) { BatteryItem(c, 1.2, BatteryItemOptions{}) })
}

func TestWifiItemJoins(t *testing.T) {
	nets := []Network{{Name: "Kiln Studio", Bars: 3, Secure: true}, {Name: "Kiln Guest", Bars: 1}}
	current, changes := "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if WifiItem(c, nets, &current, WifiItemOptions{}).Changed() {
			changes++
		}
	}), 400, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Wi-Fi")
	settles(t, tt)
	click(t, tt, "Kiln Guest")
	cues(t, before, sound.Open, sound.Select)
	if current != "Kiln Guest" || changes != 2 {
		t.Errorf("current %q, changes %d", current, changes)
	}
	panics(t, "or none", func(c *ui.Context) { s := "Cafe"; WifiItem(c, nets, &s, WifiItemOptions{}) })
	panics(t, "0 to 3 bars", func(c *ui.Context) { WifiItem(c, []Network{{Name: "A", Bars: 4}}, new(string), WifiItemOptions{}) })
}
