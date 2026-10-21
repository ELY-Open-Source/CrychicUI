package button

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestChipButton(t *testing.T) {
	on, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ChipButton(c, "Active", &on).Changed() {
			changed++
		}
	}, 300, 200)
	before := baseline()
	if err := tt.Click("Active"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "ChipButton")
	if !on || changed != 1 || !appended(before, sound.Toggle) {
		t.Errorf("on %v, changed %d", on, changed)
	}
	panics(t, "ChipButton state", func() { ui.Render(func(c *ui.Context) { ChipButton(c, "Active", nil) }, 10, 10, 1) })
}

func TestCloseButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { CloseButton(c, "Close") }, 300, 200)
	r := find(t, tt, "Close")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "CloseButton hover")
	before := baseline()
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if !appended(before, sound.Close) {
		t.Error("missing close cue")
	}
	panics(t, "CloseButton label", func() { ui.Render(func(c *ui.Context) { CloseButton(c, "") }, 10, 10, 1) })
}

func TestMoreButton(t *testing.T) {
	selected, changed := -1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if MoreButton(c, "More", []string{"Rename", "Delete"}, &selected).Changed() {
			changed++
		}
	}, 300, 200)
	before := baseline()
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	if len(tt.Menu()) != 2 || !appended(before, sound.Open) {
		t.Errorf("menu %q", tt.Menu())
	}
	if err := tt.ChooseMenuItem("Delete"); err != nil {
		t.Fatal(err)
	}
	if selected != 1 || changed != 1 {
		t.Errorf("selected %d, changed %d", selected, changed)
	}
}

func TestDropdownButton(t *testing.T) {
	open, selected, changed := false, -1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if DropdownButton(c, "Assign", []string{"Email", "Call"}, &open, &selected).Changed() {
			changed++
		}
	}, 400, 300)
	before := baseline()
	if err := tt.Click("Assign"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "DropdownButton open")
	if !open || changed != 1 || !tt.HasText("Email") {
		t.Error("dropdown did not open")
	}
	if err := tt.Click("Email"); err != nil {
		t.Fatal(err)
	}
	if open || selected != 0 || changed != 2 || !appended(before, sound.Open, sound.Tap) {
		t.Errorf("open %v, selected %d, changed %d, cues %v", open, selected, changed, sound.Recent())
	}
	again := baseline()
	if err := tt.Click("Assign"); err != nil {
		t.Fatal(err)
	}
	tt.ClickAt(390, 290)
	tt.Advance(2 * time.Second)
	if open || !appended(again, sound.Open, sound.Close) {
		t.Error("dismissal did not close with a cue")
	}
}

func TestTextButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { TextButton(c, "Learn more") }, 300, 200)
	r := find(t, tt, "Learn more")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "TextButton hover")
	before := baseline()
	if err := tt.Click("Learn more"); err != nil {
		t.Fatal(err)
	}
	if !appended(before, sound.Tap) {
		t.Error("missing tap cue")
	}
}

func TestProviderButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { ProviderButton(c, "Continue with Provider", "globe") }, 400, 200)
	r := find(t, tt, "Continue with Provider")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "ProviderButton hover")
	before := baseline()
	if err := tt.Click("Continue with Provider"); err != nil {
		t.Fatal(err)
	}
	if !appended(before, sound.Tap) {
		t.Error("missing tap cue")
	}
}
