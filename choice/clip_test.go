package choice

import (
	"testing"

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

func TestSizeChoiceStrikeStaysInside(t *testing.T) {
	tt := ui.NewTester(framed(func(c *ui.Context) { SizeChoice(c, "Size", []string{"L"}, []string{"L"}, new(string)) }), 200, 140)
	r, bg := find(t, tt, "frame"), ui.LightTheme().Background
	if n := ring(tt, r, bg); n != 0 {
		t.Errorf("strike paints %d pixels outside the box", n)
	}
	if p := tt.Image().RGBAAt(int(r.X+r.W)-1, int(r.Y)); p.R != bg.R || p.G != bg.G || p.B != bg.B {
		t.Errorf("strike paints the rounded corner: %v", p)
	}
}
