package fx

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// lit reports whether the pixel is off the background.
func lit(tt *ui.Tester, x, y int, bg ui.Color) bool {
	p := tt.Image().RGBAAt(x, y)
	return p.R != bg.R || p.G != bg.G || p.B != bg.B
}

func TestSpotlight(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Spotlight(c, func() {}).Height(200) }, 400, 200)
	bg := ui.LightTheme().Background
	if lit(tt, 80, 60, bg) {
		t.Fatal("lit before the pointer came")
	}
	tt.Move(80, 60)
	tt.Frame()
	tt.Advance(time.Second)
	if !lit(tt, 80, 60, bg) || lit(tt, 390, 190, bg) {
		t.Errorf("near %v, far %v", lit(tt, 80, 60, bg), lit(tt, 390, 190, bg))
	}
	tt.Move(300, 150)
	tt.Frame()
	tt.Advance(time.Second)
	if !lit(tt, 300, 150, bg) || lit(tt, 40, 20, bg) {
		t.Error("light did not follow the pointer")
	}
	tt.Move(500, 500)
	tt.Frame()
	tt.Advance(2 * time.Second)
	if tt.FramePending() || lit(tt, 300, 150, bg) {
		t.Error("light stayed after the pointer left")
	}
}

func TestMagicCard(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			MagicCard(c, func() { ui.Box(c).Size(240, 120).Label("face") })
		})
	}, 400, 260)
	th := ui.LightTheme()
	r := find(t, tt, "face")
	edge := tt.Image().RGBAAt(int(r.X)-int(th.Space(4))-1, int(r.Y+60))
	tt.Move(r.X+20, r.Y+20)
	tt.Frame()
	tt.Advance(time.Second)
	if !lit(tt, int(r.X+20), int(r.Y+20), th.Background) || lit(tt, int(r.X+230), int(r.Y+110), th.Background) {
		t.Error("glow not at the pointer")
	}
	if tt.Image().RGBAAt(int(r.X)-int(th.Space(4))-1, int(r.Y+60)) == edge {
		t.Error("border did not take the accent")
	}
}

func TestConfetti(t *testing.T) {
	var run uint64
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { Confetti(c, run, func() {}).Height(200) }, 300, 200)
	rest := shot(tt)
	if tt.FramePending() {
		t.Error("bursts on first build")
	}
	run = 1
	tt.Frame()
	tt.Advance(400 * time.Millisecond)
	if n := ink(tt, ui.Rect{W: 300, H: 200}); n < 200 {
		t.Errorf("%d pixels of pieces mid-burst", n)
	}
	cues(t, before, sound.Pop)
	tt.Advance(2 * time.Second)
	if tt.FramePending() || shot(tt) != rest {
		t.Error("pieces stayed after the burst")
	}
}

func TestFireworks(t *testing.T) {
	var run uint64
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { Fireworks(c, run, func() {}).Height(200) }, 300, 200)
	rest := shot(tt)
	run = 1
	tt.Frame()
	tt.Advance(200 * time.Millisecond)
	rising := ink(tt, ui.Rect{X: 60, W: 30, H: 200})
	tt.Advance(600 * time.Millisecond)
	open := ink(tt, ui.Rect{W: 300, H: 120})
	if rising == 0 || open == 0 {
		t.Errorf("rising shell %d px, open sparks %d px", rising, open)
	}
	cues(t, before, sound.Pop)
	tt.Advance(2 * time.Second)
	if tt.FramePending() || shot(tt) != rest {
		t.Error("sparks stayed after the show")
	}
}

func TestGradientText(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { GradientText(c, "Gradient headline") }, 400, 80)
	tt.Advance(time.Second)
	r := find(t, tt, "Gradient headline")
	left, right := tone(tt, ui.Rect{X: r.X, Y: r.Y, W: r.W / 5, H: r.H}), tone(tt, ui.Rect{X: r.X + r.W*0.8, Y: r.Y, W: r.W / 5, H: r.H})
	if !(left[2] > left[1] && right[1] > right[2]) {
		t.Errorf("left %v should lean blue, right %v green", left, right)
	}
}

// tone is the mean color of r's dark pixels.
func tone(tt *ui.Tester, r ui.Rect) [3]float32 {
	img := tt.Image()
	var sum [3]float32
	n := float32(0)
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); int(p.R)+int(p.G)+int(p.B) < 500 {
				sum[0], sum[1], sum[2], n = sum[0]+float32(p.R), sum[1]+float32(p.G), sum[2]+float32(p.B), n+1
			}
		}
	}
	return [3]float32{sum[0] / max(n, 1), sum[1] / max(n, 1), sum[2] / max(n, 1)}
}

func TestFlowingTextShifts(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { FlowingText(c, "Flowing headline") }, 400, 80)
	r := find(t, tt, "Flowing headline")
	left := ui.Rect{X: r.X, Y: r.Y, W: r.W / 5, H: r.H}
	a := tone(tt, left)
	tt.Advance(time.Second)
	if b := tone(tt, left); a == b {
		t.Errorf("left colors held at %v", a)
	}
}
