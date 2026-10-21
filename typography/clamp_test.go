package typography

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestClamp(t *testing.T) {
	open, changed := false, 0
	height := float32(0)
	text := strings.Repeat("many wrapping words ", 20)
	tt := ui.NewTester(func(c *ui.Context) {
		root := Clamp(c, text, 2, &open, ClampOptions{})
		height = root.LayoutBounds().H
		if root.Changed() {
			changed++
		}
	}, 240, 500)
	tt.Frame()
	closed := height
	if closed <= 0 {
		t.Fatal("clamp has no closed height")
	}
	taps := count(sound.Toggle)
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(50 * time.Millisecond)
	if height <= closed {
		t.Errorf("opening height stayed %v", height)
	}
	settles(t, tt, "Clamp open")
	if !open || changed != 1 || count(sound.Toggle) != taps+1 || !tt.HasText("Less") {
		t.Errorf("open %v, changed %d, cues %d", open, changed, count(sound.Toggle)-taps)
	}
	tt.Frame()
	opened := height
	if opened <= closed {
		t.Errorf("opened height %v, closed %v", opened, closed)
	}
	if err := tt.Click("Less"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(50 * time.Millisecond)
	if height <= closed || height >= opened {
		t.Errorf("closing height %v, want %v to %v", height, closed, opened)
	}
	settles(t, tt, "Clamp close")
	if open || changed != 2 || count(sound.Toggle) != taps+2 || !tt.HasText("More") {
		t.Errorf("closed %v, changed %d", open, changed)
	}
}

func TestClampValidation(t *testing.T) {
	panics(t, "Clamp lines", func() { ui.Render(func(c *ui.Context) { Clamp(c, "Text", 0, new(bool), ClampOptions{}) }, 10, 10, 1) })
	panics(t, "Clamp open", func() { ui.Render(func(c *ui.Context) { Clamp(c, "Text", 2, nil, ClampOptions{}) }, 10, 10, 1) })
}
