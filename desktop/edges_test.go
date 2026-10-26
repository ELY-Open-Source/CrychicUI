package desktop

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// focusHolds clicks name, shows optional copy, then presses Space.
func focusHolds(t *testing.T, what, name string, more *bool, hits *int, view func(c *ui.Context)) {
	t.Helper()
	tt := ui.NewTester(padded(view), 520, 700)
	settles(t, tt)
	click(t, tt, name)
	*more = true
	tt.Frame()
	settles(t, tt)
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if *hits != 2 {
		t.Errorf("%s: %d activations after copy appeared, want 2", what, *hits)
	}
}

func TestActionsKeepFocus(t *testing.T) {
	more, hits := false, 0
	focusHolds(t, "ProgressWindow", "Cancel", &more, &hits, func(c *ui.Context) {
		o := ProgressWindowOptions{}
		if more {
			o.Detail = "Exporting kiln logs to the studio drive"
		}
		if _, cancel := ProgressWindow(c, "Export", 3, 10, o); cancel.Clicked() {
			hits++
		}
	})
	more, hits = false, 0
	focusHolds(t, "UpdatePrompt", "Install update", &more, &hits, func(c *ui.Context) {
		var notes []string
		if more {
			notes = []string{"Faster kiln calendar"}
		}
		if _, install, _ := UpdatePrompt(c, "Glaze", "3.2", notes, UpdatePromptOptions{}); install.Clicked() {
			hits++
		}
	})
}

func TestStatusBarLongFactsWrap(t *testing.T) {
	left := []StatusItem{{Icon: "cloud", Text: "Everything is synced to the studio", Tone: 1}}
	right := []StatusItem{{Icon: "flame", Text: "Three firings scheduled for today", Tone: 2}}
	tt := ui.NewTester(padded(within(320, func(c *ui.Context) { StatusBar(c, left, right) })), 520, 400)
	settles(t, tt)
	a := rectOf(t, tt, left[0].Text)
	b := rectOf(t, tt, right[0].Text)
	inside(t, "left fact", a, 10, 330)
	inside(t, "right fact", b, 10, 330)
	if a.Y < b.Y+b.H && b.Y < a.Y+a.H && a.X < b.X+b.W && b.X < a.X+a.W {
		t.Errorf("facts overlap: %v and %v", a, b)
	}
}

func TestWifiItemEmpty(t *testing.T) {
	current := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { WifiItem(c, nil, &current, WifiItemOptions{}) }), 520, 500)
	settles(t, tt)
	click(t, tt, "Wi-Fi")
	settles(t, tt)
	if !exact(tt, "No networks found") {
		t.Errorf("empty menu shows %q", tt.Texts())
	}
}
