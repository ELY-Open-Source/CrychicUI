package button

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestButtonGroup(t *testing.T) {
	selected, changed := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ButtonGroup(c, []string{"A", "B"}, &selected).Changed() {
			changed++
		}
	}, 300, 200)
	taps := sound.Recent()
	if err := tt.Click("B"); err != nil {
		t.Fatal(err)
	}
	if selected != 1 || changed != 1 || !appended(taps, sound.Tap) {
		t.Errorf("selected %d, changed %d, cues %v", selected, changed, sound.Recent())
	}
	panics(t, "ButtonGroup items", func() { ui.Render(func(c *ui.Context) { ButtonGroup(c, nil, new(int)) }, 10, 10, 1) })
}

func TestToggleButton(t *testing.T) {
	on, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ToggleButton(c, "Bold", &on).Changed() {
			changed++
		}
	}, 300, 200)
	taps := sound.Recent()
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	if !on || changed != 1 || !appended(taps, sound.Toggle) {
		t.Errorf("on %v, changed %d", on, changed)
	}
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	if on || changed != 2 || !appended(taps, sound.Toggle, sound.Toggle) {
		t.Errorf("off %v, changed %d", on, changed)
	}
	panics(t, "ToggleButton state", func() { ui.Render(func(c *ui.Context) { ToggleButton(c, "Bold", nil) }, 10, 10, 1) })
}

func TestLoadingButton(t *testing.T) {
	busy, clicks := false, 0
	width := float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		b := LoadingButton(c, "Save", busy)
		width = b.LayoutBounds().W
		if b.Clicked() {
			clicks++
		}
	}, 300, 200)
	taps := sound.Recent()
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 || !appended(taps, sound.Tap) {
		t.Fatal("idle click failed")
	}
	tt.Frame()
	idle := width
	busy = true
	tt.Frame()
	if !tt.FramePending() {
		t.Error("spinner did not start")
	}
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(300 * time.Millisecond)
	if width == idle {
		t.Error("busy content did not change width")
	}
	if clicks != 1 || !appended(taps, sound.Tap) || !tt.FramePending() {
		t.Error("busy button clicked or spinner stopped")
	}
}

func TestConfirmButton(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ConfirmButton(c, "Delete", ConfirmButtonOptions{}).Changed() {
			changed++
		}
	}, 300, 200)
	warnings, deletes := sound.Recent(), sound.Recent()
	if err := tt.Click("Delete"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Sure?") || changed != 0 || !appended(warnings, sound.Warning) {
		t.Error("first click did not arm")
	}
	if err := tt.Click("Sure?"); err != nil {
		t.Fatal(err)
	}
	if changed != 1 || !appended(deletes, sound.Delete) {
		t.Error("second click did not confirm")
	}
	if err := tt.Click("Delete"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(3 * time.Second)
	if !tt.HasText("Delete") || !appended(deletes, sound.Delete) {
		t.Error("armed button did not time out")
	}
}

func TestHoldButton(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if HoldButton(c, "Hold", HoldButtonOptions{Duration: 100 * time.Millisecond}).Changed() {
			changed++
		}
	}, 300, 200)
	r := find(t, tt, "Hold")
	successes := sound.Recent()
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(150 * time.Millisecond)
	if changed != 1 || !appended(successes, sound.Success) {
		t.Error("hold did not complete")
	}
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(150 * time.Millisecond)
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(150 * time.Millisecond)
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	if changed != 2 || !appended(successes, sound.Success, sound.Success) {
		t.Error("second hold did not complete")
	}
}

func TestHoldButtonReducedMotion(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if HoldButton(c, "Hold", HoldButtonOptions{Duration: 100 * time.Millisecond}).Changed() {
			changed++
		}
	}, 300, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	r := find(t, tt, "Hold")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	defer tt.Release(r.X+r.W/2, r.Y+r.H/2)
	if changed != 0 {
		t.Fatal("hold completed at press")
	}
	tt.Advance(50 * time.Millisecond)
	if changed != 0 {
		t.Fatal("hold completed before its duration")
	}
	tt.Advance(60 * time.Millisecond)
	if changed != 1 {
		t.Error("hold did not complete at its duration")
	}
}

func TestCopyButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { CopyButton(c, "Copy", "copy me") }, 300, 200)
	pops := sound.Recent()
	if err := tt.Click("Copy"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != "copy me" || !appended(pops, sound.Pop) {
		t.Errorf("clipboard %q, cues %v", tt.Clipboard(), sound.Recent())
	}
	tt.Advance(2 * time.Second)
}

func TestShareButton(t *testing.T) {
	selected, changed := -1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ShareButton(c, "Share", []string{"Email", "Copy link"}, &selected).Changed() {
			changed++
		}
	}, 300, 200)
	opens := sound.Recent()
	if err := tt.Click("Share"); err != nil {
		t.Fatal(err)
	}
	if len(tt.Menu()) != 2 || !appended(opens, sound.Open) {
		t.Errorf("menu %q, cues %v", tt.Menu(), sound.Recent())
	}
	if err := tt.ChooseMenuItem("Email"); err != nil {
		t.Fatal(err)
	}
	if selected != 0 || changed != 1 || !appended(opens, sound.Open) {
		t.Errorf("selected %d, changed %d", selected, changed)
	}
	panics(t, "ShareButton targets", func() { ui.Render(func(c *ui.Context) { ShareButton(c, "Share", nil, new(int)) }, 10, 10, 1) })
}
