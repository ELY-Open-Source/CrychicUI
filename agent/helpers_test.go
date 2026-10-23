package agent

import (
	"bytes"
	"image/color"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func padded(view func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) { ui.Column(c).Padding(10).Gap(8).Children(func() { view(c) }) }
}

func find(t *testing.T, tt *ui.Tester, s string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(s)
	if !ok {
		t.Fatalf("%q not shown in %q", s, tt.Texts())
	}
	return r
}

func click(t *testing.T, tt *ui.Tester, s string) {
	t.Helper()
	if err := tt.Click(s); err != nil {
		t.Fatalf("click %q: %v", s, err)
	}
	tt.Frame()
}

func panics(t *testing.T, what string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", what)
		}
	}()
	fn()
}

func render(view func(c *ui.Context)) { ui.Render(view, 500, 400, 1) }

func settles(t *testing.T, tt *ui.Tester, what string) {
	t.Helper()
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Errorf("%s still moving", what)
	}
}

// moves reports whether the frame changes over d.
func moves(tt *ui.Tester, d time.Duration) bool {
	a := bytes.Clone(tt.Image().Pix)
	tt.Advance(d)
	return !bytes.Equal(a, tt.Image().Pix)
}

// baseline fills Recent with Warning, which agent never plays.
func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Warning)
	}
	return sound.Recent()
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	after := sound.Recent()
	var got []sound.Cue
	if len(before) == 32 && !slices.Equal(before, after) {
		i := 0
		for i < len(after) && after[i] == sound.Warning {
			i++
		}
		got = after[i:]
	}
	if !slices.Equal(got, want) {
		t.Errorf("cues %v, want %v", got, want)
	}
}

func near(a, b uint8) bool { return a-b < 24 || b-a < 24 }

// at reports whether the pixel there is near c.
func at(tt *ui.Tester, x, y float32, c color.RGBA) bool {
	p := tt.Image().RGBAAt(int(x), int(y))
	return near(p.R, c.R) && near(p.G, c.G) && near(p.B, c.B)
}

func rgba(c ui.Color) color.RGBA { return color.RGBA{c.R, c.G, c.B, 255} }
