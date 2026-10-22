package choice

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestIconToggles(t *testing.T) {
	var star, eye, lock, mute bool
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Gap(8).Children(func() {
			IconToggle(c, "Favorite", "star", "star-off", &star)
			VisibilityToggle(c, "Show password", &eye)
			LockToggle(c, "Unlock", &lock)
			MuteToggle(c, "Mute", &mute)
		})
	}, 300, 80)
	before := baseline()
	for _, s := range []string{"Favorite", "Show password", "Unlock", "Mute"} {
		click(t, tt, s)
	}
	if !star || !eye || !lock || !mute {
		t.Fatalf("states %v %v %v %v", star, eye, lock, mute)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle, sound.Toggle)
	settles(t, tt, "icon swaps")
	panics(t, "unknown icon", func() {
		ui.Render(func(c *ui.Context) { IconToggle(c, "x", "nope-nope", "star", &star) }, 50, 50, 1)
	})
}

func TestDayToggles(t *testing.T) {
	var days [7]bool
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if DayToggles(c, "Repeat", &days, DayTogglesOptions{}).Changed() {
			changes++
		}
	}, 400, 80)
	before := baseline()
	click(t, tt, "Monday")
	click(t, tt, "Friday")
	if !days[1] || !days[5] || changes != 2 {
		t.Fatalf("days %v changes %d", days, changes)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
	settles(t, tt, "DayToggles fill")
}

func TestTriCheck(t *testing.T) {
	s := TriMixed
	tt := ui.NewTester(func(c *ui.Context) { TriCheck(c, "All files", &s) }, 200, 80)
	before := baseline()
	click(t, tt, "All files")
	if s != TriOn {
		t.Fatalf("mixed click gave %v, want on", s)
	}
	settles(t, tt, "TriCheck morph")
	click(t, tt, "All files")
	if s != TriOff {
		t.Fatalf("on click gave %v, want off", s)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
}

func TestLabeledSwitch(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) { LabeledSwitch(c, "Autosave", &on, LabeledSwitchOptions{Off: "Nein"}) }, 300, 80)
	find(t, tt, "Nein")
	before := baseline()
	click(t, tt, "Autosave")
	if !on {
		t.Fatal("switch did not turn on")
	}
	settles(t, tt, "LabeledSwitch knob")
	find(t, tt, "On")
	cues(t, before, sound.Toggle)
}
