package gallery

import (
	"github.com/ELY-Open-Source/CrychicUI/fx"
	"github.com/egoist/mygo/ui"
)

var fxEntries = []Entry{
	{406, "GridBackdrop", "fx", "Faint grid lines drifting behind content.", demoGridBackdrop},
	{407, "DotBackdrop", "fx", "A dot grid drifting behind content.", demoDotBackdrop},
	{408, "GradientBackdrop", "fx", "Accent to surface; the angle turns.", demoGradientBackdrop},
	{409, "Aurora", "fx", "Soft bands of theme colors moving.", demoAurora},
	{410, "Spotlight", "fx", "A soft light that follows the pointer.", demoSpotlight},
	{411, "Starfield", "fx", "Points flying out from the center.", demoStarfield},
	{412, "Particles", "fx", "Drifting dots that wrap at the edges.", demoParticles},
	{413, "Confetti", "fx", "A burst of falling, spinning pieces.", demoConfetti},
	{414, "Fireworks", "fx", "Shells rise and open into sparks.", demoFireworks},
	{415, "Snowfall", "fx", "Flakes falling with sway.", demoSnowfall},
	{416, "Rainfall", "fx", "Streaks falling at a slant.", demoRainfall},
	{417, "Bubbles", "fx", "Circles rising and fading.", demoBubbles},
	{418, "Meteors", "fx", "Streaks crossing diagonally.", demoMeteors},
	{419, "Sparkles", "fx", "Stars blinking around content.", demoSparkles},
	{420, "BorderBeam", "fx", "A light segment running round a border.", demoBorderBeam},
	{421, "ConicBorder", "fx", "A turning conic gradient border.", demoConicBorder},
	{422, "ShineBorder", "fx", "A sheen passing along the border.", demoShineBorder},
	{423, "MagicCard", "fx", "A card whose glow follows the pointer.", demoMagicCard},
	{424, "GradientText", "fx", "Text filled with a gradient.", demoGradientText},
	{425, "FlowingText", "fx", "Gradient text whose stops move.", demoFlowingText},
	{426, "ShinyText", "fx", "Text with a passing sheen.", demoShinyText},
	{427, "RetroGrid", "fx", "A perspective grid scrolling closer.", demoRetroGrid},
	{428, "RippleBackdrop", "fx", "Rings expanding behind content.", demoRippleBackdrop},
}

// stage is the sample content the backdrops sit behind.
func stage(c *ui.Context) func() {
	return func() {
		ui.Column(c).Fill().Center().Padding(24).Children(func() { ui.Text(c, "Every part moves").Bold().FontSize(20) })
	}
}

func demoGridBackdrop(c *ui.Context, log func(string)) { fx.GridBackdrop(c, stage(c)).Width(360) }

func demoDotBackdrop(c *ui.Context, log func(string)) { fx.DotBackdrop(c, stage(c)).Width(360) }

func demoGradientBackdrop(c *ui.Context, log func(string)) {
	fx.GradientBackdrop(c, stage(c)).Width(360)
}

func demoAurora(c *ui.Context, log func(string)) { fx.Aurora(c, stage(c)).Width(360) }

func demoSpotlight(c *ui.Context, log func(string)) {
	s := fx.Spotlight(c, stage(c)).Width(360)
	hover(s, log, "light on", "light off")
}

func demoStarfield(c *ui.Context, log func(string)) { fx.Starfield(c, stage(c)).Width(360) }

func demoParticles(c *ui.Context, log func(string)) { fx.Particles(c, stage(c)).Width(360) }

func demoConfetti(c *ui.Context, log func(string)) {
	trigger(c, log, "Celebrate", func(run uint64) { fx.Confetti(c, run, stage(c)).Width(360) })
}

func demoFireworks(c *ui.Context, log func(string)) {
	trigger(c, log, "Launch", func(run uint64) { fx.Fireworks(c, run, stage(c)).Width(360) })
}

func demoSnowfall(c *ui.Context, log func(string)) { fx.Snowfall(c, stage(c)).Width(360) }

func demoRainfall(c *ui.Context, log func(string)) { fx.Rainfall(c, stage(c)).Width(360) }

func demoBubbles(c *ui.Context, log func(string)) { fx.Bubbles(c, stage(c)).Width(360) }

func demoMeteors(c *ui.Context, log func(string)) { fx.Meteors(c, stage(c)).Width(360) }

func demoSparkles(c *ui.Context, log func(string)) {
	fx.Sparkles(c, func() { ui.Text(c, "New").Bold().FontSize(24) })
}

func demoBorderBeam(c *ui.Context, log func(string)) {
	fx.BorderBeam(c, func() { ui.Text(c, "A beam runs round this card.") })
}

func demoConicBorder(c *ui.Context, log func(string)) {
	fx.ConicBorder(c, func() { ui.Text(c, "Colors turn round this card.") })
}

func demoShineBorder(c *ui.Context, log func(string)) {
	fx.ShineBorder(c, func() { ui.Text(c, "A sheen crosses this border.") })
}

func demoMagicCard(c *ui.Context, log func(string)) {
	m := fx.MagicCard(c, func() {
		ui.Text(c, "Magic card").Bold()
		ui.Text(c, "Move the pointer across it.")
	}).Width(300)
	hover(m, log, "glow on", "glow off")
}

func demoGradientText(c *ui.Context, log func(string)) {
	fx.GradientText(c, "Small parts, put together")
}

func demoFlowingText(c *ui.Context, log func(string)) { fx.FlowingText(c, "Colors that keep moving") }

func demoShinyText(c *ui.Context, log func(string)) { fx.ShinyText(c, "Loading your feed") }

func demoRetroGrid(c *ui.Context, log func(string)) { fx.RetroGrid(c, stage(c)).Width(360) }

func demoRippleBackdrop(c *ui.Context, log func(string)) { fx.RippleBackdrop(c, stage(c)).Width(360) }
