package icon

import (
	"image"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestEveryIconParses(t *testing.T) {
	n := Names()
	if len(n) != 2159 {
		t.Fatalf("icons = %d, want 2159", len(n))
	}
	if !slices.IsSorted(n) {
		t.Error("Names not sorted")
	}
	for _, name := range n {
		if Must(name) == nil {
			t.Errorf("%s parsed to nil", name)
		}
	}
}

func TestMustCaches(t *testing.T) {
	if Must("check") != Must("check") {
		t.Error("Must parsed check twice")
	}
}

func TestMustPanicsOnUnknown(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Must(nope) did not panic")
		}
	}()
	Must("nope")
}

func TestNamesIsACopy(t *testing.T) {
	n := Names()
	n[0] = "changed"
	if Names()[0] == "changed" {
		t.Error("Names exposed its slice")
	}
}

func TestIconDraws(t *testing.T) {
	img := ui.Render(func(c *ui.Context) {
		ui.Icon(c, Must("check")).Size(24, 24)
	}, 32, 32, 1)
	bg := img.RGBAAt(0, 0)
	inked := 0
	for y := range 32 {
		for x := range 32 {
			if img.RGBAAt(x, y) != bg {
				inked++
			}
		}
	}
	if inked < 20 {
		t.Errorf("check drew %d pixels", inked)
	}
}

func TestFilledHeart(t *testing.T) {
	outline := ui.Render(func(c *ui.Context) { ui.Icon(c, Must("heart")).Size(24, 24) }, 32, 32, 1)
	filled := ui.Render(func(c *ui.Context) { ui.Icon(c, Filled("heart")).Size(24, 24) }, 32, 32, 1)
	count := func(img *image.RGBA) int {
		bg, n := img.RGBAAt(0, 0), 0
		for y := range 32 {
			for x := range 32 {
				if img.RGBAAt(x, y) != bg {
					n++
				}
			}
		}
		return n
	}
	if count(filled) <= count(outline) {
		t.Errorf("filled %d pixels, outline %d", count(filled), count(outline))
	}
	if Filled("heart") != Filled("heart") {
		t.Error("Filled parsed heart twice")
	}
}

func TestFilledDoesNotAliasNames(t *testing.T) {
	Filled("heart")
	defer func() {
		if recover() == nil {
			t.Error("Must accepted a filled cache key")
		}
	}()
	Must("filled:heart")
}
