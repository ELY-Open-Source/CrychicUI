package typography

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestKbd(t *testing.T) {
	pressed := false
	tt := ui.NewTester(func(c *ui.Context) { Kbd(c, "S", pressed) }, 300, 200)
	before := find(t, tt, "S").Y
	pressed = true
	tt.Frame()
	pressed = false
	if !tt.FramePending() {
		t.Fatal("match did not start the press")
	}
	tt.Advance(60 * time.Millisecond)
	if got := find(t, tt, "S").Y; got <= before {
		t.Errorf("cap moved from %v to %v", before, got)
	}
	settles(t, tt, "Kbd release")
	if got := find(t, tt, "S").Y; got != before {
		t.Errorf("cap rests at %v, want %v", got, before)
	}
}

func TestKeyboardHintsAllowShortcuts(t *testing.T) {
	s, combo := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Kbd(c, "S", false)
		KeyCombo(c, []string{"Ctrl", "K"}, false)
		if c.Shortcut(0, ui.KeyS) {
			s++
		}
		if c.Shortcut(ui.Cmd, ui.KeyK) {
			combo++
		}
	}, 300, 200)
	tt.Key(0, ui.KeyS)
	tt.Key(ui.Cmd, ui.KeyK)
	if s != 1 || combo != 1 {
		t.Errorf("application actions: S %d, combo %d", s, combo)
	}
}

func TestKeyCombo(t *testing.T) {
	matched := false
	tt := ui.NewTester(func(c *ui.Context) {
		KeyCombo(c, []string{"Ctrl", "K"}, matched)
	}, 300, 200)
	ctrl, k := find(t, tt, "Ctrl").Y, find(t, tt, "K").Y
	matched = true
	tt.Frame()
	matched = false
	tt.Advance(60 * time.Millisecond)
	if find(t, tt, "Ctrl").Y <= ctrl {
		t.Error("first cap did not press")
	}
	if find(t, tt, "K").Y != k {
		t.Error("second cap pressed before its turn")
	}
	tt.Advance(60 * time.Millisecond)
	if find(t, tt, "K").Y <= k {
		t.Error("second cap did not press")
	}
	settles(t, tt, "KeyCombo")
	if find(t, tt, "Ctrl").Y != ctrl || find(t, tt, "K").Y != k {
		t.Error("caps did not return to rest")
	}
	panics(t, "KeyCombo labels", func() {
		ui.Render(func(c *ui.Context) { KeyCombo(c, nil, false) }, 10, 10, 1)
	})
}
