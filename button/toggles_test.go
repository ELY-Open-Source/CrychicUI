package button

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestBinaryIconButtons(t *testing.T) {
	cases := []struct {
		name  string
		label string
		loops bool
		build func(*ui.Context, *bool) ui.Element
	}{
		{"ThemeToggle", "Toggle theme", false, func(c *ui.Context, on *bool) ui.Element { return ThemeToggle(c, on, ThemeToggleOptions{}) }},
		{"VoiceButton", "Voice", true, func(c *ui.Context, on *bool) ui.Element { return VoiceButton(c, on, VoiceButtonOptions{}) }},
		{"PlayButton", "Play", false, func(c *ui.Context, on *bool) ui.Element { return PlayButton(c, on, PlayButtonOptions{}) }},
		{"ExpandButton", "Expand", false, func(c *ui.Context, on *bool) ui.Element { return ExpandButton(c, on, ExpandButtonOptions{}) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			on, changed := false, 0
			tt := ui.NewTester(func(c *ui.Context) {
				if tc.build(c, &on).Changed() {
					changed++
				}
			}, 300, 200)
			before := baseline()
			if err := tt.Click(tc.label); err != nil {
				t.Fatal(err)
			}
			if !on || changed != 1 || !appended(before, sound.Toggle) {
				t.Errorf("on %v, changed %d", on, changed)
			}
			if tc.loops {
				if !tt.FramePending() {
					t.Error("active rings did not loop")
				}
			} else {
				settles(t, tt, tc.name)
			}
			if err := tt.Click(tc.label); err != nil {
				t.Fatal(err)
			}
			settles(t, tt, tc.name+" off")
			if on || changed != 2 {
				t.Errorf("off %v, changed %d", on, changed)
			}
		})
	}
}

func TestSortButton(t *testing.T) {
	direction, changed := SortNone, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SortButton(c, "Sort", &direction).Changed() {
			changed++
		}
	}, 300, 200)
	before := baseline()
	for _, want := range []SortDirection{SortAscending, SortDescending, SortNone} {
		if err := tt.Click("Sort"); err != nil {
			t.Fatal(err)
		}
		if direction != want {
			t.Fatalf("direction %v, want %v", direction, want)
		}
	}
	settles(t, tt, "SortButton")
	if changed != 3 || !appended(before, sound.Tick, sound.Tick, sound.Tick) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
	panics(t, "SortButton direction", func() {
		ui.Render(func(c *ui.Context) { d := SortDirection(9); SortButton(c, "Sort", &d) }, 10, 10, 1)
	})
}

func TestRefreshButton(t *testing.T) {
	refreshing, clicks := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if RefreshButton(c, "Refresh", refreshing).Clicked() {
			clicks++
			refreshing = true
		}
	}, 300, 200)
	before := baseline()
	if err := tt.Click("Refresh"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 || !refreshing || !appended(before, sound.Tap) || !tt.FramePending() {
		t.Error("refresh did not start its loop")
	}
	refreshing = false
	tt.Frame()
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Error("refresh loop did not stop")
	}
}

func TestShortcutButton(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ShortcutButton(c, "Command", ui.Cmd, ui.KeyK, []string{"Ctrl", "K"}, ShortcutButtonOptions{}).Changed() {
			changed++
		}
	}, 300, 200)
	before := baseline()
	if err := tt.Click("Command"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyK)
	if changed != 2 || !appended(before, sound.Tap, sound.Tap) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
	panics(t, "ShortcutButton keys", func() {
		ui.Render(func(c *ui.Context) { ShortcutButton(c, "Command", 0, ui.KeyK, nil, ShortcutButtonOptions{}) }, 10, 10, 1)
	})
}

func TestShortcutButtonDisabledAncestor(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Disabled(true).Children(func() {
			if ShortcutButton(c, "Command", ui.Cmd, ui.KeyK, []string{"Ctrl", "K"}, ShortcutButtonOptions{}).Changed() {
				changed++
			}
		})
	}, 300, 200)
	before := baseline()
	tt.Key(ui.Cmd, ui.KeyK)
	if changed != 0 || !appended(before) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
}

func TestSlideConfirmCatchSpring(t *testing.T) {
	changed := 0
	var accent ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		accent = c.Theme().Accent
		if SlideConfirm(c, "Slide to confirm", SlideConfirmOptions{}).Changed() {
			changed++
		}
	}, 300, 200)
	tt.Press(20, 20)
	tt.Move(60, 20)
	tt.Release(60, 20)
	tt.Advance(50 * time.Millisecond)
	position := func() int {
		for x := 0; x < 300; x++ {
			if near(tt.Image(), x, 20, accent) {
				return x
			}
		}
		return -1
	}
	before := position()
	if before < 0 {
		t.Fatal("thumb was not visible")
	}
	tt.Press(float32(before), 20)
	tt.Move(float32(before+5), 20)
	if after := position(); after < before-1 {
		t.Fatalf("thumb jumped from %d to %d", before, after)
	}
	tt.Move(270, 20)
	tt.Release(270, 20)
	tt.Advance(2 * time.Second)
	if changed != 1 {
		t.Error("caught drag did not confirm")
	}
}

func TestShortcutButtonNested(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			if ShortcutButton(c, "Command", ui.Cmd, ui.KeyK, []string{"Ctrl", "K"}, ShortcutButtonOptions{}).Changed() {
				changed++
			}
			ui.Button(c, "Focus").AutoFocus()
		})
	}, 300, 200)
	before := baseline()
	tt.Key(ui.Cmd, ui.KeyK)
	if changed != 1 || !appended(before, sound.Tap) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
}

func TestShortcutButtonDisabled(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		b := ShortcutButton(c, "Command", ui.Cmd, ui.KeyK, []string{"Ctrl", "K"}, ShortcutButtonOptions{Disabled: true})
		if b.Changed() {
			changed++
		}
	}, 300, 200)
	before := baseline()
	tt.Key(ui.Cmd, ui.KeyK)
	if changed != 0 || !appended(before) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
}

func TestSlideConfirm(t *testing.T) {
	t.Run("short drag", func(t *testing.T) {
		changed := 0
		tt := ui.NewTester(func(c *ui.Context) {
			if SlideConfirm(c, "Slide to confirm", SlideConfirmOptions{}).Changed() {
				changed++
			}
		}, 300, 200)
		before := baseline()
		tt.Press(20, 20)
		tt.Move(60, 20)
		tt.Release(60, 20)
		tt.Advance(2 * time.Second)
		if changed != 0 || !appended(before) {
			t.Errorf("changed %d, cues %v", changed, sound.Recent())
		}
	})
	t.Run("short fling", func(t *testing.T) {
		changed := 0
		tt := ui.NewTester(func(c *ui.Context) {
			if SlideConfirm(c, "Slide to confirm", SlideConfirmOptions{}).Changed() {
				changed++
			}
		}, 300, 200)
		before := baseline()
		tt.Press(20, 20)
		tt.Move(140, 20)
		tt.Advance(time.Millisecond)
		tt.Release(140, 20)
		tt.Advance(2 * time.Second)
		if changed != 0 || !appended(before) {
			t.Errorf("changed %d, cues %v", changed, sound.Recent())
		}
	})
	t.Run("full drag", func(t *testing.T) {
		changed := 0
		tt := ui.NewTester(func(c *ui.Context) {
			if SlideConfirm(c, "Slide to confirm", SlideConfirmOptions{}).Changed() {
				changed++
			}
		}, 300, 200)
		before := baseline()
		tt.Press(20, 20)
		tt.Move(270, 20)
		tt.Release(270, 20)
		tt.Advance(2 * time.Second)
		if changed != 1 || !appended(before, sound.Success) {
			t.Errorf("changed %d, cues %v", changed, sound.Recent())
		}
	})
}
