package display

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

func TestAchievementShineStaysInside(t *testing.T) {
	on := false
	tt := ui.NewTester(framed(func(c *ui.Context) { Achievement(c, "First week", "trophy", on) }), 400, 160)
	tt.SetDark(true)
	tt.Frame()
	tt.Advance(time.Second)
	r := find(t, tt, "frame")
	on = true
	tt.Frame()
	for range 10 {
		tt.Advance(50 * time.Millisecond)
		if n := ring(tt, r, ui.DarkTheme().Background); n != 0 {
			t.Fatalf("shine paints %d pixels outside the badge", n)
		}
	}
}
