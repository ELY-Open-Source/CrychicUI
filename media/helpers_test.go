package media

import (
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

var (
	red   = photo(220, 40, 40, 64, 48)
	green = photo(40, 180, 60, 64, 48)
	blue  = photo(40, 70, 220, 64, 48)
)

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

// baseline fills Recent with Error, which media never plays.
func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Error)
	}
	return sound.Recent()
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	after := sound.Recent()
	var got []sound.Cue
	if len(before) == 32 && !slices.Equal(before, after) {
		i := 0
		for i < len(after) && after[i] == sound.Error {
			i++
		}
		got = after[i:]
	}
	if !slices.Equal(got, want) {
		t.Errorf("cues %v, want %v", got, want)
	}
}

func near(a, b uint8) bool { return a-b < 24 || b-a < 24 }

func absDiff(a, b uint8) int { return max(int(a)-int(b), int(b)-int(a)) }

// at reports whether the pixel there is near c.
func at(tt *ui.Tester, x, y float32, c color.RGBA) bool {
	p := tt.Image().RGBAAt(int(x), int(y))
	return near(p.R, c.R) && near(p.G, c.G) && near(p.B, c.B)
}

var (
	redInk   = color.RGBA{220, 40, 40, 255}
	greenInk = color.RGBA{40, 180, 60, 255}
	blueInk  = color.RGBA{40, 70, 220, 255}
)
