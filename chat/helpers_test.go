package chat

import (
	"bytes"
	"image"
	"image/color"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// photo is a solid w×h bitmap of one color.
func photo(r, g, b uint8, w, h int) *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return ui.NewBitmap(img)
}

var red = photo(220, 40, 40, 64, 48)

var redInk = color.RGBA{220, 40, 40, 255}

var now = time.Date(2026, 10, 22, 9, 30, 0, 0, time.UTC)

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

// baseline fills Recent with Warning, which chat never plays.
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

// count is how many shown texts equal s.
func count(tt *ui.Tester, s string) int {
	n := 0
	for _, x := range tt.Texts() {
		if x == s {
			n++
		}
	}
	return n
}

// hasInk scans r for a pixel near c.
func hasInk(tt *ui.Tester, r ui.Rect, c color.RGBA) bool {
	for y := int(r.Y); y < int(r.Y+r.H); y += 2 {
		for x := int(r.X); x < int(r.X+r.W); x += 2 {
			if at(tt, float32(x), float32(y), c) {
				return true
			}
		}
	}
	return false
}
