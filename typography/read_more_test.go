package typography

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestReadMore(t *testing.T) {
	open, changed := false, 0
	height := float32(0)
	text := strings.Repeat("many wrapping words ", 20)
	tt := ui.NewTester(func(c *ui.Context) {
		root := ReadMore(c, text, 2, &open, ReadMoreOptions{})
		height = root.LayoutBounds().H
		if root.Changed() {
			changed++
		}
	}, 240, 500)
	tt.Frame()
	closed := height
	taps := count(sound.Tap)
	if err := tt.Click("Read more"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(50 * time.Millisecond)
	if height <= closed {
		t.Errorf("opening height stayed %v", height)
	}
	settles(t, tt, "ReadMore")
	if !open || changed != 1 || count(sound.Tap) != taps+1 || tt.HasText("Read more") {
		t.Errorf("open %v, changed %d, cues %d, texts %q", open, changed, count(sound.Tap)-taps, tt.Texts())
	}
	tt.Frame()
	if height <= closed {
		t.Errorf("opened height %v, closed %v", height, closed)
	}
}

func TestReadMoreValidation(t *testing.T) {
	panics(t, "ReadMore lines", func() {
		ui.Render(func(c *ui.Context) { ReadMore(c, "Text", 0, new(bool), ReadMoreOptions{}) }, 10, 10, 1)
	})
	panics(t, "ReadMore open", func() {
		ui.Render(func(c *ui.Context) { ReadMore(c, "Text", 2, nil, ReadMoreOptions{}) }, 10, 10, 1)
	})
}
