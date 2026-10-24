package code

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// goLines is n lines of plausible Go.
func goLines(n int) string {
	var b strings.Builder
	for i := range n {
		fmt.Fprintf(&b, "\tv%d := compute(%d) // step\n", i, i)
	}
	return b.String()
}

func TestMinimap(t *testing.T) {
	at, changed := 0, 0
	src := goLines(200)
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Minimap(c, "Editor map", "go", src, 20, &at, MinimapOptions{}).Changed() {
			changed++
		}
	}), 300, 300)
	settles(t, tt, "minimap")
	r := find(t, tt, "Editor map")
	before := baseline()
	tt.Press(r.X+r.W/2, r.Y+r.H-6)
	tt.Frame()
	if at < 150 || changed == 0 {
		t.Fatalf("at %d changed %d after a press", at, changed)
	}
	tt.Move(r.X+r.W/2, r.Y+r.H-40)
	tt.Frame()
	tt.Release(r.X+r.W/2, r.Y+r.H-40)
	tt.Frame()
	if at >= 190 || at == 0 {
		t.Errorf("at %d after a drag", at)
	}
	cues(t, before, sound.Tick)
	keys := ui.NewTester(padded(func(c *ui.Context) {
		if Minimap(c, "Keys", "go", src, 20, &at, MinimapOptions{}).Changed() {
			changed++
		}
	}), 300, 300)
	keys.Frame()
	keys.Key(0, ui.KeyTab)
	if !keys.Focused("Keys") {
		t.Fatal("minimap does not take the focus")
	}
	at = 0
	keys.Key(0, ui.KeyEnd)
	keys.Frame()
	if at != 180 {
		t.Errorf("at %d after End, want 180", at)
	}
	panics(t, "window over lines", func() {
		n := 0
		render(func(c *ui.Context) { Minimap(c, "x", "go", goLines(10), 20, &n, MinimapOptions{}) })
	})
	panics(t, "start out of range", func() {
		n := 200
		render(func(c *ui.Context) { Minimap(c, "x", "go", src, 20, &n, MinimapOptions{}) })
	})
}

func TestMinimapDrawsTokens(t *testing.T) {
	at := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Minimap(c, "Tokens", "go", "func main() {\n\n\treturn\n}", 1, &at, MinimapOptions{Height: 160})
	}), 300, 300)
	settles(t, tt, "tokens")
	r := find(t, tt, "Tokens")
	ink := func(y float32) bool {
		for x := r.X + 6; x < r.X+r.W/2; x++ {
			if p := tt.Image().RGBAAt(int(x), int(y)); p.R < 200 || p.G < 200 || p.B < 200 {
				return true
			}
		}
		return false
	}
	if !ink(r.Y+r.H*0.625) || ink(r.Y+r.H*0.375) {
		t.Error("lines do not follow the code: blank line inked or code line bare")
	}
}

func TestMinimapSlowDrag(t *testing.T) {
	at := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Minimap(c, "Slow map", "go", goLines(10), 2, &at, MinimapOptions{Height: 160})
	}), 300, 300)
	settles(t, tt, "slow map")
	r := find(t, tt, "Slow map")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	for i := 1; i <= 64; i++ {
		tt.Move(x, y+float32(i))
		tt.Frame()
	}
	tt.Release(x, y+64)
	tt.Frame()
	if at != 8 {
		t.Errorf("at %d after slow drag, want 8", at)
	}
}

func TestMinimapUnscrollable(t *testing.T) {
	at, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Minimap(c, "Whole file", "go", goLines(2), 2, &at, MinimapOptions{}).Changed() {
			changed++
		}
	}), 300, 300)
	settles(t, tt, "whole file")
	before := baseline()
	click(t, tt, "Whole file")
	if at != 0 || changed != 0 {
		t.Errorf("at %d changed %d without scroll range", at, changed)
	}
	cues(t, before)
}

func TestMinimapPressMotion(t *testing.T) {
	at := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Minimap(c, "Motion map", "go", goLines(20), 2, &at, MinimapOptions{})
	}), 300, 300)
	settles(t, tt, "motion map")
	r := find(t, tt, "Motion map")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	if !tt.FramePending() {
		t.Error("press has no motion")
	}
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "motion map")
}
