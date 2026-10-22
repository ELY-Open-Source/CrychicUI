package button

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// ring counts pixels just outside r unlike bg.
func ring(tt *ui.Tester, r ui.Rect, bg ui.Color) int {
	img, n := tt.Image(), 0
	for y := int(r.Y) - 5; y < int(r.Y+r.H)+5; y++ {
		for x := int(r.X) - 5; x < int(r.X+r.W)+5; x++ {
			in := x >= int(r.X)-1 && x <= int(r.X+r.W)+1 && y >= int(r.Y)-1 && y <= int(r.Y+r.H)+1
			if p := img.RGBAAt(x, y); !in && (p.R != bg.R || p.G != bg.G || p.B != bg.B) {
				n++
			}
		}
	}
	return n
}

// framed builds v in a labeled box that hugs it.
func framed(v func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) {
		ui.Column(c).Padding(40).AlignItems(ui.Start).Children(func() { ui.Box(c).Label("frame").Children(func() { v(c) }) })
	}
}

func TestRippleStaysInside(t *testing.T) {
	ripple := ui.NewTester(framed(func(c *ui.Context) { RippleButton(c, "Ripple") }), 300, 160)
	r := find(t, ripple, "frame")
	ripple.Press(r.X+4, r.Y+r.H/2)
	ripple.Release(r.X+4, r.Y+r.H/2)
	for range 4 {
		ripple.Advance(60 * time.Millisecond)
		if n := ring(ripple, r, ui.LightTheme().Background); n != 0 {
			t.Fatalf("ripple paints %d pixels outside the button", n)
		}
	}
}

func TestShimmerStaysInside(t *testing.T) {
	shimmer := ui.NewTester(framed(func(c *ui.Context) { ShimmerButton(c, "Shimmer") }), 300, 160)
	shimmer.SetDark(true)
	shimmer.Frame()
	r := find(t, shimmer, "frame")
	for range 30 {
		shimmer.Advance(100 * time.Millisecond)
		if n := ring(shimmer, r, ui.DarkTheme().Background); n != 0 {
			t.Fatalf("shimmer paints %d pixels outside the button", n)
		}
	}
}

func TestDownloadFillKeepsCorners(t *testing.T) {
	progress := 0.0
	dl := ui.NewTester(framed(func(c *ui.Context) { DownloadButton(c, "Download", progress) }), 300, 160)
	r := find(t, dl, "frame")
	corner := dl.Image().RGBAAt(int(r.X), int(r.Y))
	progress = 1
	dl.Frame()
	dl.Advance(time.Second)
	if p := dl.Image().RGBAAt(int(r.X), int(r.Y)); p != corner {
		t.Errorf("full fill paints the rounded corner: %v, was %v", p, corner)
	}
}
