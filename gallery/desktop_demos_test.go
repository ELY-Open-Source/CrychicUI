package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDesktopDemosLog(t *testing.T) {
	press := func(names ...string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, name := range names {
				if err := tt.Click(name); err != nil {
					return err
				}
				tt.Frame()
			}
			return nil
		}
	}
	unlock := func(tt *ui.Tester) error {
		if err := tt.Click("Password"); err != nil {
			return err
		}
		tt.Type("glaze")
		tt.Frame()
		return press("Unlock")(tt)
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"TrafficLights":      {press("Close"), "close"},
		"TitleTabs":          {press("Members"), "tab 1"},
		"MenuBar":            {press("File", "New booking"), "chose New booking"},
		"Dock":               {press("Mail"), "opened Mail"},
		"Launcher":           {press("Calendar"), "picked Calendar"},
		"WindowSwitcher":     {press("Next window"), "window 1"},
		"DesktopIcon":        {press("Firing logs"), "selected Firing logs"},
		"NotificationCenter": {press("Clear all"), "cleared"},
		"ControlCenter":      {press("Bluetooth"), "bluetooth true"},
		"VolumeHUD":          {press("Volume up"), "volume 70"},
		"BrightnessHUD":      {press("Brighter"), "brightness 50"},
		"WifiItem":           {press("Wi-Fi", "Kiln Guest"), "joined Kiln Guest"},
		"InspectorPane":      {press("Firing"), "section toggled"},
		"UnifiedToolbar":     {press("Back"), "back"},
		"SourceList":         {press("Waitlist"), "source Waitlist"},
		"Preferences":        {press("Notifications"), "pane 2"},
		"ProgressWindow":     {press("Cancel"), "closed"},
		"UpdatePrompt":       {press("Install update"), "install"},
		"PermissionPrompt":   {press("Allow"), "access 1"},
		"TrayMenu":           {press("Glaze", "Pause sync"), "chose Pause sync"},
		"Spaces":             {press("Bookkeeping"), "space 1"},
		"SnapLayouts":        {press("Layout 1, area 2"), "snap 0 1"},
		"LockScreen":         {unlock, "unlocked"},
	}
	quiet := map[string]bool{"WindowFrame": true, "StatusBar": true, "BatteryItem": true, "ClockItem": true, "About": true,
		"TaskManager": true, "TerminalWindow": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "desktop" {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 900)
		tt.Advance(2 * time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 30 {
		t.Errorf("checked %d desktop demos", seen)
	}
}

func TestDesktopSecondaryLogs(t *testing.T) {
	for _, c := range []struct {
		name  string
		press []string
		want  string
	}{
		{"MenuBar", []string{"File"}, "menu open"},
		{"WifiItem", []string{"Wi-Fi"}, "menu open"},
		{"TrayMenu", []string{"Glaze"}, "menu open"},
		{"TitleTabs", []string{"New tab"}, "new tab"},
		{"UnifiedToolbar", []string{"Forward"}, "forward"},
		{"ProgressWindow", []string{"Export 4 more"}, "done 16"},
		{"UpdatePrompt", []string{"Later"}, "later"},
	} {
		for _, e := range Entries() {
			if e.Category != "desktop" || e.Name != c.name {
				continue
			}
			var logs []string
			tt := ui.NewTester(func(ctx *ui.Context) { e.Demo(ctx, func(s string) { logs = append(logs, s) }) }, 900, 900)
			tt.Advance(2 * time.Second)
			for _, p := range c.press {
				if err := tt.Click(p); err != nil {
					t.Fatalf("%s: %v", c.name, err)
				}
				tt.Frame()
			}
			if !strings.Contains(strings.Join(logs, "\n"), c.want) {
				t.Errorf("%s logs %q, want %q", c.name, logs, c.want)
			}
		}
	}
}

func TestDesktopEventLogs(t *testing.T) {
	run := func(name string, do func(tt *ui.Tester)) string {
		for _, e := range Entries() {
			if e.Category == "desktop" && e.Name == name {
				var logs []string
				tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 900)
				tt.Advance(2 * time.Second)
				do(tt)
				tt.Frame()
				return strings.Join(logs, "\n")
			}
		}
		t.Fatalf("no %s demo", name)
		return ""
	}
	click := func(tt *ui.Tester, names ...string) {
		for _, n := range names {
			if err := tt.Click(n); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
		}
	}
	for _, m := range []struct{ name, menu, item string }{{"MenuBar", "File", "New booking"}, {"TrayMenu", "Glaze", "Pause sync"}} {
		got := run(m.name, func(tt *ui.Tester) { click(tt, m.menu, m.item, m.menu, m.item) })
		if n := strings.Count(got, "chose "+m.item); n != 2 {
			t.Errorf("%s repeated command logged %d times: %q", m.name, n, got)
		}
		if n := strings.Count(got, "menu open"); n != 2 {
			t.Errorf("%s logged %d openings: %q", m.name, n, got)
		}
	}
	got := run("LockScreen", func(tt *ui.Tester) {
		click(tt, "Password")
		tt.Type("clay")
		tt.Frame()
		click(tt, "Unlock")
	})
	if !strings.Contains(got, "password 4 runes") || !strings.Contains(got, "attempt 1 refused") || strings.Contains(got, "clay") {
		t.Errorf("LockScreen logs %q", got)
	}
	got = run("ControlCenter", func(tt *ui.Tester) {
		r, ok := tt.Find("Sound")
		if !ok {
			t.Fatal("no Sound slider")
		}
		// The slider row sits just below its label.
		x, y := r.X+r.W/2, r.Y+r.H+14
		tt.Press(x, y)
		tt.Release(x, y)
	})
	if !strings.Contains(got, "sound ") || strings.Contains(got, "sound 60%") {
		t.Errorf("ControlCenter slider logs %q", got)
	}
}
