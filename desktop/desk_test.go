package desktop

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestClockItemTurnsOver(t *testing.T) {
	var first string
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if first == "" {
			first = c.Now().Format("15:04")
		}
		ClockItem(c, ClockItemOptions{Format: "15:04"})
	}), 300, 200)
	settles(t, tt)
	if !exact(tt, first) {
		t.Errorf("clock shows %q, want %s", tt.Texts(), first)
	}
	tt.Advance(61 * time.Second)
	if exact(tt, first) {
		t.Errorf("the clock still shows %s a minute later", first)
	}
}

func TestInspectorPaneFolds(t *testing.T) {
	changes := 0
	secs := []InspectorSection{{Title: "Booking", Rows: []InspectorRow{{Label: "Kiln", Value: "Kiln 2"}}}, {Title: "Firing", Rows: []InspectorRow{{Label: "Peak", Value: "2232°F"}}}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if InspectorPane(c, secs).Changed() {
			changes++
		}
	}), 320, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Firing")
	settles(t, tt)
	cues(t, before, sound.Toggle)
	if r, ok := tt.Find("2232°F"); (ok && r.H > 0) || changes != 1 {
		t.Errorf("fold left %v, changes %d", r, changes)
	}
	panics(t, "needs rows", func(c *ui.Context) { InspectorPane(c, []InspectorSection{{Title: "A"}}) })
}

func TestUnifiedToolbarActions(t *testing.T) {
	query, changes, navs := "", 0, 0
	tt := ui.NewTester(within(520, func(c *ui.Context) {
		r, back, forward := UnifiedToolbar(c, "Kiln calendar", &query, UnifiedToolbarOptions{Subtitle: "Week of Oct 26"})
		if r.Changed() {
			changes++
		}
		if back.Clicked() || forward.Clicked() {
			navs++
		}
	}), 600, 200)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Back")
	click(t, tt, "Forward")
	cues(t, before, sound.Tap, sound.Tap)
	click(t, tt, "Search")
	tt.Type("kiln")
	tt.Frame()
	if navs != 2 || changes != 1 || query != "kiln" {
		t.Errorf("navs %d, changes %d, query %q", navs, changes, query)
	}
	panics(t, "a title and a query", func(c *ui.Context) { UnifiedToolbar(c, " ", new(string), UnifiedToolbarOptions{}) })
}

func TestSourceListPicksByName(t *testing.T) {
	groups := []SourceGroup{{Title: "Library", Items: []SourceItem{{Name: "Today", Icon: "sun"}, {Name: "Waitlist", Icon: "clock", Count: 5}, {Name: "Inbox", Icon: "mail"}}}}
	current := "Inbox"
	tt := ui.NewTester(padded(func(c *ui.Context) { SourceList(c, groups, &current) }), 300, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Waitlist", func() { groups[0].Items[0], groups[0].Items[1] = groups[0].Items[1], groups[0].Items[0] })
	if current != "Inbox" {
		t.Errorf("a reorder under the pointer picked %q", current)
	}
	settles(t, tt)
	before := baseline()
	click(t, tt, "Waitlist")
	click(t, tt, "Waitlist")
	cues(t, before, sound.Select)
	if current != "Waitlist" {
		t.Errorf("current %q", current)
	}
	panics(t, "repeated", func(c *ui.Context) {
		s := "A"
		SourceList(c, []SourceGroup{{Title: "X", Items: []SourceItem{{Name: "A", Icon: "sun"}}}, {Title: "Y", Items: []SourceItem{{Name: "A", Icon: "sun"}}}}, &s)
	})
}

func TestPreferencesSwitchesPane(t *testing.T) {
	tabs := []PrefTab{{Name: "General", Icon: "settings-2"}, {Name: "Accounts", Icon: "user"}, {Name: "Notifications", Icon: "bell"}}
	picked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { Preferences(c, tabs, &picked, func() { ui.Text(c, "pane body") }) }), 520, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Notifications")
	click(t, tt, "Notifications")
	cues(t, before, sound.Select)
	settles(t, tt)
	n := 0
	for _, s := range tt.Texts() {
		if s == "Notifications" {
			n++
		}
	}
	// The tab's label and text, and the window title.
	if picked != 2 || n != 3 {
		t.Errorf("picked %d, Notifications shows %d times", picked, n)
	}
	panics(t, "a picked one and a body", func(c *ui.Context) { Preferences(c, tabs, new(int), nil) })
}

func TestAboutAndProgress(t *testing.T) {
	done, closes := 12, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		About(c, "Glaze", AboutOptions{Version: "3.2.0"})
		if _, cancel := ProgressWindow(c, "Exporting firing logs", done, 40, ProgressWindowOptions{}); cancel.Clicked() {
			closes++
		}
	}), 480, 800)
	settles(t, tt)
	if !exact(tt, "Version 3.2.0") || !exact(tt, "12 of 40 items") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Cancel")
	cues(t, before, sound.Tap)
	done = 40
	tt.Frame()
	settles(t, tt)
	if closes != 1 || !exact(tt, "Done") {
		t.Errorf("closes %d, texts %q", closes, tt.Texts())
	}
	panics(t, "done <= total", func(c *ui.Context) { ProgressWindow(c, "Export", 5, 4, ProgressWindowOptions{}) })
	panics(t, "app name", func(c *ui.Context) { About(c, " ", AboutOptions{}) })
}

func TestUpdateAndPermission(t *testing.T) {
	access, changes := Undecided, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		UpdatePrompt(c, "Glaze", "3.3.0", []string{"Faster calendar."}, UpdatePromptOptions{})
		if PermissionPrompt(c, "Glaze", "camera", &access, PermissionPromptOptions{}).Changed() {
			changes++
		}
	}), 480, 900)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Later")
	click(t, tt, "Install update")
	click(t, tt, "Allow")
	click(t, tt, "Allow")
	click(t, tt, "Don't allow")
	cues(t, before, sound.Tap, sound.Success, sound.Tap, sound.Tap)
	settles(t, tt)
	if access != Denied || changes != 2 || !exact(tt, "Access denied. Change it in Settings.") {
		t.Errorf("access %d, changes %d, texts %q", access, changes, tt.Texts())
	}
	panics(t, "known answer", func(c *ui.Context) {
		a := Access(5)
		PermissionPrompt(c, "Glaze", "camera", &a, PermissionPromptOptions{})
	})
}

func TestTrayMenuAndSpaces(t *testing.T) {
	chosen, space, changes := "", 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Spaces(c, []string{"Studio", "Bookkeeping"}, &space).Changed() {
			changes++
		}
		if TrayMenu(c, []TrayItem{{Label: "Open Glaze"}, {Label: "Pause sync", Icon: "pause"}}, &chosen, TrayMenuOptions{Label: "Glaze"}).Changed() {
			changes++
		}
	}), 400, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Bookkeeping")
	click(t, tt, "Bookkeeping")
	click(t, tt, "Glaze")
	settles(t, tt)
	click(t, tt, "Pause sync")
	cues(t, before, sound.Select, sound.Open, sound.Tap)
	if space != 1 || chosen != "Pause sync" || changes != 3 {
		t.Errorf("space %d, chosen %q, changes %d", space, chosen, changes)
	}
	panics(t, "current one", func(c *ui.Context) { n := 2; Spaces(c, []string{"A"}, &n) })
}

func TestSnapLayoutsPicks(t *testing.T) {
	picked := Snap{Layout: -1}
	tt := ui.NewTester(padded(func(c *ui.Context) { SnapLayouts(c, &picked, SnapLayoutsOptions{}) }), 400, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Layout 4, area 3")
	click(t, tt, "Layout 4, area 3")
	cues(t, before, sound.Select)
	if picked != (Snap{3, 2}) {
		t.Errorf("picked %v", picked)
	}
	panics(t, "known layout", func(c *ui.Context) { SnapLayouts(c, &Snap{0, 5}, SnapLayoutsOptions{}) })
}

func TestTaskManagerBusiestFirst(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		TaskManager(c, []Process{{Name: "Mail", Icon: "mail", CPU: 0.8}, {Name: "Browser", Icon: "globe", CPU: 31}}, TaskManagerOptions{})
	}), 480, 400)
	settles(t, tt)
	if b, m := rectOf(t, tt, "Browser"), rectOf(t, tt, "Mail"); b.Y > m.Y || !exact(tt, "31.0%") {
		t.Errorf("order %v, %v; texts %q", b, m, tt.Texts())
	}
	panics(t, "CPU 0 to 100", func(c *ui.Context) {
		TaskManager(c, []Process{{Name: "A", Icon: "cpu", CPU: math.NaN()}}, TaskManagerOptions{})
	})
}

func TestTerminalWindowPrompt(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		TerminalWindow(c, "zsh", []string{"$ glaze login", "Signed in"}, TerminalWindowOptions{})
	}), 480, 400)
	tt.Advance(2 * time.Second)
	if !tt.HasText("glaze login") || !exact(tt, "Signed in") {
		t.Errorf("texts %q", tt.Texts())
	}
	on := bytes.Clone(tt.Image().Pix)
	tt.Advance(550 * time.Millisecond)
	if bytes.Equal(on, tt.Image().Pix) {
		t.Error("the caret does not blink")
	}
}

func TestLockScreenUnlocks(t *testing.T) {
	password, unlocked := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		LockScreen(c, "Maya Chen", &password, func(s string) bool { return s == "glaze" }, &unlocked, LockScreenOptions{})
	}), 480, 700)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Unlock")
	click(t, tt, "Password")
	tt.Type("kiln")
	tt.Frame()
	click(t, tt, "Unlock")
	if password != "" || unlocked {
		t.Errorf("a wrong password: kept %q, unlocked %v", password, unlocked)
	}
	click(t, tt, "Password")
	tt.Type("glaze")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	cues(t, before, sound.Error, sound.Success)
	if !unlocked || !exact(tt, "Welcome back") {
		t.Errorf("unlocked %v, texts %q", unlocked, tt.Texts())
	}
	panics(t, "a check", func(c *ui.Context) { LockScreen(c, "Maya", new(string), nil, new(bool), LockScreenOptions{}) })
}
