package gallery

import (
	"github.com/ELY-Open-Source/CrychicUI/fx"
	"github.com/egoist/mygo/ui"
)

var fxEntries2 = []Entry{
	{429, "Waves", "fx", "Layered sine lines moving.", demoWaves},
	{430, "Blob", "fx", "A soft shape morphing.", demoBlob},
	{431, "OrbitRings", "fx", "Dots orbiting on concentric rings.", demoOrbitRings},
	{432, "DotGlobe", "fx", "A sphere of dots turning.", demoDotGlobe},
	{433, "Beam", "fx", "Light running between two anchors.", demoBeam},
	{434, "ShimmerSurface", "fx", "A surface with a slow sheen.", demoShimmerSurface},
	{435, "Grain", "fx", "Faint dots dithering over content.", demoGrain},
	{436, "Vignette", "fx", "Edges darkening toward the corners.", demoVignette},
	{437, "Frosted", "fx", "A panel blurring what lies behind it.", demoFrosted},
	{438, "Scanlines", "fx", "Horizontal lines rolling down.", demoScanlines},
	{439, "Halftone", "fx", "Dots sized by a turning gradient.", demoHalftone},
	{440, "Pattern", "fx", "Stripes or checks drifting.", demoPattern},
	{441, "MatrixRain", "fx", "Columns of falling glyphs.", demoMatrixRain},
	{442, "Caret", "fx", "A blinking text caret.", demoCaret},
	{443, "Trail", "fx", "Dots fading behind the pointer.", demoTrail},
	{444, "ClickBurst", "fx", "Particles bursting at each click.", demoClickBurst},
	{445, "LightRays", "fx", "Rays fanning from a corner.", demoLightRays},
	{446, "ConicSpin", "fx", "A conic gradient turning.", demoConicSpin},
	{447, "Plasma", "fx", "Cells colored by moving sine fields.", demoPlasma},
	{448, "WaterRipple", "fx", "A ring expanding from each click.", demoWaterRipple},
	{449, "Flicker", "fx", "Text flickering like a failing sign.", demoFlicker},
	{450, "Mesh", "fx", "A four-corner color mesh shifting.", demoMesh},
}

func demoWaves(c *ui.Context, log func(string)) { fx.Waves(c, stage(c)).Width(360) }

func demoBlob(c *ui.Context, log func(string)) { fx.Blob(c, stage(c)).Width(360) }

func demoOrbitRings(c *ui.Context, log func(string)) { fx.OrbitRings(c, stage(c)).Width(360) }

func demoDotGlobe(c *ui.Context, log func(string)) { fx.DotGlobe(c, stage(c)).Width(360) }

func demoBeam(c *ui.Context, log func(string)) {
	fx.Beam(c, ui.Point{X: 0.1, Y: 0.8}, ui.Point{X: 0.9, Y: 0.8}, stage(c)).Width(360)
}

func demoShimmerSurface(c *ui.Context, log func(string)) { fx.ShimmerSurface(c, stage(c)).Width(360) }

func demoGrain(c *ui.Context, log func(string)) { fx.Grain(c, stage(c)).Width(360) }

func demoVignette(c *ui.Context, log func(string)) { fx.Vignette(c, stage(c)).Width(360) }

func demoFrosted(c *ui.Context, log func(string)) {
	fx.Aurora(c, func() {
		ui.Column(c).Fill().Center().Padding(24).Children(func() {
			fx.Frosted(c, func() { ui.Text(c, "Frosted glass").Bold() }).AlignSelf(ui.Center)
		})
	}).Width(360)
}

func demoScanlines(c *ui.Context, log func(string)) { fx.Scanlines(c, stage(c)).Width(360) }

func demoHalftone(c *ui.Context, log func(string)) { fx.Halftone(c, stage(c)).Width(360) }

func demoPattern(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	checks := ui.Local(root, "checks", func() bool { return false })
	root.Children(func() {
		switched := ui.Button(c, "Switch pattern").AlignSelf(ui.Start).Clicked()
		if switched {
			*checks = !*checks
		}
		kind, name := fx.Stripes, "stripes"
		if *checks {
			kind, name = fx.Checks, "checks"
		}
		if switched {
			log(name)
		}
		fx.Pattern(c, fx.PatternOptions{Kind: kind}, stage(c)).Width(360)
	})
}

func demoMatrixRain(c *ui.Context, log func(string)) { fx.MatrixRain(c, stage(c)).Width(360) }

func demoCaret(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(2).AlignItems(ui.Center).Children(func() {
		ui.Text(c, "Search")
		fx.Caret(c)
	})
}

func demoTrail(c *ui.Context, log func(string)) {
	s := fx.Trail(c, stage(c)).Width(360)
	hover(s, log, "trail on", "trail off")
}

func demoClickBurst(c *ui.Context, log func(string)) {
	if fx.ClickBurst(c, stage(c)).Width(360).Clicked() {
		log("burst")
	}
}

func demoLightRays(c *ui.Context, log func(string)) { fx.LightRays(c, stage(c)).Width(360) }

func demoConicSpin(c *ui.Context, log func(string)) { fx.ConicSpin(c, stage(c)).Width(360) }

func demoPlasma(c *ui.Context, log func(string)) { fx.Plasma(c, stage(c)).Width(360) }

func demoWaterRipple(c *ui.Context, log func(string)) {
	if fx.WaterRipple(c, stage(c)).Width(360).Clicked() {
		log("ripple")
	}
}

func demoFlicker(c *ui.Context, log func(string)) { fx.Flicker(c, "OPEN LATE") }

func demoMesh(c *ui.Context, log func(string)) { fx.Mesh(c, stage(c)).Width(360) }
