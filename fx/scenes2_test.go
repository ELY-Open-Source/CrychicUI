package fx

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestScenesMove2(t *testing.T) {
	hi := func(c *ui.Context) func() { return func() { ui.Text(c, "Hi") } }
	ends := ui.Point{X: 0.1, Y: 0.7}
	scenes := []struct {
		name string
		view func(c *ui.Context) ui.Element
	}{
		{"Waves", func(c *ui.Context) ui.Element { return Waves(c, hi(c)) }},
		{"Blob", func(c *ui.Context) ui.Element { return Blob(c, hi(c)) }},
		{"OrbitRings", func(c *ui.Context) ui.Element { return OrbitRings(c, hi(c)) }},
		{"DotGlobe", func(c *ui.Context) ui.Element { return DotGlobe(c, hi(c)) }},
		{"Beam", func(c *ui.Context) ui.Element { return Beam(c, ends, ui.Point{X: 0.9, Y: 0.7}, hi(c)) }},
		{"ShimmerSurface", func(c *ui.Context) ui.Element { return ShimmerSurface(c, hi(c)) }},
		{"Grain", func(c *ui.Context) ui.Element { return Grain(c, hi(c)) }},
		{"Vignette", func(c *ui.Context) ui.Element { return Vignette(c, hi(c)) }},
		{"Scanlines", func(c *ui.Context) ui.Element { return Scanlines(c, hi(c)) }},
		{"Halftone", func(c *ui.Context) ui.Element { return Halftone(c, hi(c)) }},
		{"Stripes", func(c *ui.Context) ui.Element { return Pattern(c, PatternOptions{}, hi(c)) }},
		{"Checks", func(c *ui.Context) ui.Element { return Pattern(c, PatternOptions{Kind: Checks}, hi(c)) }},
		{"MatrixRain", func(c *ui.Context) ui.Element { return MatrixRain(c, hi(c)) }},
		{"LightRays", func(c *ui.Context) ui.Element { return LightRays(c, hi(c)) }},
		{"ConicSpin", func(c *ui.Context) ui.Element { return ConicSpin(c, hi(c)) }},
		{"Plasma", func(c *ui.Context) ui.Element { return Plasma(c, hi(c)) }},
		{"Mesh", func(c *ui.Context) ui.Element { return Mesh(c, hi(c)) }},
	}
	for _, s := range scenes {
		t.Run(s.name, func(t *testing.T) { checkScene(t, s.name, "scene", s.view) })
	}
}

// still renders an empty scene 300 by 160 at rest.
func still(view func(c *ui.Context, content func()) ui.Element) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) { view(c, func() {}).Height(160) }, 300, 160)
}

// weight sums how far r's pixels sit from bg.
func weight(tt *ui.Tester, r ui.Rect, bg ui.Color) int {
	img, n := tt.Image(), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			p := img.RGBAAt(x, y)
			n += abs(int(p.R)-int(bg.R)) + abs(int(p.G)-int(bg.G)) + abs(int(p.B)-int(bg.B))
		}
	}
	return n
}

func abs(v int) int { return max(v, -v) }

func TestWavesSpanTheWidth(t *testing.T) {
	tt := still(Waves)
	if ink(tt, ui.Rect{W: 6, H: 160}) == 0 || ink(tt, ui.Rect{X: 294, W: 6, H: 160}) == 0 {
		t.Error("lines do not reach both edges")
	}
}

func TestBlobSitsInTheMiddle(t *testing.T) {
	tt := still(Blob)
	if ink(tt, ui.Rect{X: 145, Y: 75, W: 10, H: 10}) != 100 || ink(tt, ui.Rect{W: 60, H: 30}) != 0 {
		t.Error("blob not centered")
	}
}

func TestOrbitRingsDrawRings(t *testing.T) {
	tt := still(OrbitRings)
	for i := range 3 {
		d := 80 * float32(i+1) / 3.5
		if ink(tt, ui.Rect{X: 150 - d - 1, Y: 79, W: 3, H: 3}) == 0 || ink(tt, ui.Rect{X: 149, Y: 80 + d - 1, W: 3, H: 3}) == 0 {
			t.Errorf("ring %d missing at radius %.1f", i, d)
		}
	}
	if ink(tt, ui.Rect{X: 140, Y: 70, W: 4, H: 4}) != 0 {
		t.Error("ink between the center and the first ring")
	}
}

func TestDotGlobeIsRound(t *testing.T) {
	tt := still(DotGlobe)
	if ink(tt, ui.Rect{W: 80, H: 160}) != 0 || ink(tt, ui.Rect{X: 100, W: 100, H: 160}) == 0 {
		t.Error("dots fall outside the sphere")
	}
}

func TestBeam(t *testing.T) {
	th := ui.LightTheme()
	tt := ui.NewTester(func(c *ui.Context) {
		Beam(c, ui.Point{X: 0.1, Y: 0.8}, ui.Point{X: 0.9, Y: 0.8}, func() {}).Height(160)
	}, 300, 160)
	if !lit(tt, 30, 128, th.Background) || !lit(tt, 270, 128, th.Background) {
		t.Error("anchors not drawn")
	}
	lightest := 0
	for range 10 {
		tt.Advance(250 * time.Millisecond)
		lightest = max(lightest, count(tt, ui.Rect{W: 300, H: 160}, th.Accent))
	}
	if lightest == 0 {
		t.Error("no light ran along the beam")
	}
	panics(t, "anchor past the box", func() { render(func(c *ui.Context) { Beam(c, ui.Point{X: 1.2}, ui.Point{}, func() {}) }) })
	nan := float32(math.NaN())
	panics(t, "NaN anchor", func() { render(func(c *ui.Context) { Beam(c, ui.Point{}, ui.Point{X: nan}, func() {}) }) })
}

func TestShimmerSurfaceSheen(t *testing.T) {
	tt := still(ShimmerSurface)
	s := ui.LightTheme().Surface
	if !near(tt.Image().RGBAAt(150, 80), s) {
		t.Fatalf("surface %v at rest", tt.Image().RGBAAt(150, 80))
	}
	tt.Advance(2 * time.Second)
	if p := tt.Image().RGBAAt(150, 80); int(p.R)+int(p.G)+int(p.B) <= int(s.R)+int(s.G)+int(s.B) {
		t.Errorf("sheen %v did not lift the surface", p)
	}
}

// near reports whether p matches k within 2 per channel.
func near(p interface{ RGBA() (r, g, b, a uint32) }, k ui.Color) bool {
	r, g, b, _ := p.RGBA()
	return abs(int(r>>8)-int(k.R)) <= 2 && abs(int(g>>8)-int(k.G)) <= 2 && abs(int(b>>8)-int(k.B)) <= 2
}

func TestGrainIsFaintAndSparse(t *testing.T) {
	n := ink(still(Grain), ui.Rect{W: 300, H: 160})
	if n == 0 || n > 300*160/4 {
		t.Errorf("%d grain pixels", n)
	}
}

func TestVignetteDarkensCorners(t *testing.T) {
	tt := still(Vignette)
	bg := ui.LightTheme().Background
	if !near(tt.Image().RGBAAt(150, 80), bg) || !(weight(tt, ui.Rect{X: 12, Y: 12, W: 4, H: 4}, bg) > 16*30) {
		t.Error("edges not darker than the middle")
	}
}

func TestScanlinesAreLevel(t *testing.T) {
	tt := still(Scanlines)
	tt.Advance(time.Second)
	img := tt.Image()
	for y := 10; y < 150; y++ {
		if img.RGBAAt(20, y) != img.RGBAAt(280, y) {
			t.Fatalf("row %d differs across", y)
		}
	}
	if ink(tt, ui.Rect{X: 20, W: 1, H: 160}) == 0 {
		t.Error("no lines")
	}
}

func TestHalftoneFollowsItsGradient(t *testing.T) {
	tt := still(Halftone)
	if l, r := ink(tt, ui.Rect{W: 75, H: 160}), ink(tt, ui.Rect{X: 225, W: 75, H: 160}); l >= r {
		t.Errorf("left %d px, right %d px: dots should grow rightward", l, r)
	}
}

func TestPatternKinds(t *testing.T) {
	stripes := still(func(c *ui.Context, f func()) ui.Element { return Pattern(c, PatternOptions{}, f) })
	checks := still(func(c *ui.Context, f func()) ui.Element { return Pattern(c, PatternOptions{Kind: Checks}, f) })
	if shot(stripes) == shot(checks) || ink(checks, ui.Rect{W: 300, H: 160}) == 0 {
		t.Error("kinds look alike")
	}
	panics(t, "unknown kind", func() {
		render(func(c *ui.Context) { Pattern(c, PatternOptions{Kind: 7}, func() {}) })
	})
}

func TestMatrixRainIsGreen(t *testing.T) {
	tt := still(MatrixRain)
	img, n := tt.Image(), 0
	for y := range 160 {
		for x := range 300 {
			if p := img.RGBAAt(x, y); int(p.G) > int(p.R)+20 && int(p.G) > int(p.B)+20 {
				n++
			}
		}
	}
	if n < 20 {
		t.Errorf("%d green glyph pixels", n)
	}
}

func TestLightRaysFanFromTheCorner(t *testing.T) {
	tt := still(LightRays)
	bg := ui.LightTheme().Background
	if a, b := weight(tt, ui.Rect{X: 4, Y: 4, W: 30, H: 30}, bg), weight(tt, ui.Rect{X: 266, Y: 126, W: 30, H: 30}, bg); a <= b {
		t.Errorf("corner %d, far %d", a, b)
	}
}

func TestConicSpinTurns(t *testing.T) {
	tt := still(ConicSpin)
	a := tt.Image().RGBAAt(150, 20)
	tt.Advance(2500 * time.Millisecond)
	if b := tt.Image().RGBAAt(150, 20); a == b {
		t.Errorf("top stayed %v after a quarter turn", a)
	}
}

func TestPlasmaFillsCells(t *testing.T) {
	tt := still(Plasma)
	img := tt.Image()
	if img.RGBAAt(10, 18) != img.RGBAAt(14, 22) {
		t.Error("one cell holds two colors")
	}
	if img.RGBAAt(100, 40) == img.RGBAAt(220, 120) && img.RGBAAt(100, 40) == img.RGBAAt(40, 120) {
		t.Error("cells share one color")
	}
}

func TestMeshCornersDiffer(t *testing.T) {
	img := still(Mesh).Image()
	tl, br := img.RGBAAt(10, 10), img.RGBAAt(290, 150)
	if tl == br {
		t.Errorf("corners both %v", tl)
	}
}
