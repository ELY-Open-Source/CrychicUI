package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/egoist/mygo/ui"
)

var motionEntries = []Entry{
	{351, "FadeIn", "motion", "Content fades in when built or replayed.", demoFadeIn},
	{352, "SlideIn", "motion", "Content slides in from an edge.", demoSlideIn},
	{353, "ScaleIn", "motion", "Content scales up from 0.9.", demoScaleIn},
	{354, "BlurIn", "motion", "Content sharpens from a blur.", demoBlurIn},
	{355, "FlipIn", "motion", "Content turns in about its top edge.", demoFlipIn},
	{356, "Stagger", "motion", "Items arrive one step after another.", demoStagger},
	{357, "Typewriter", "motion", "Types letter by letter; ticks per word.", demoTypewriter},
	{358, "Scramble", "motion", "Letters scramble, then settle.", demoScramble},
	{359, "WordRotate", "motion", "Words take turns, sliding up.", demoWordRotate},
	{360, "WordFade", "motion", "Words fade in one after another.", demoWordFade},
	{361, "LetterPull", "motion", "Letters drop in from above.", demoLetterPull},
	{362, "WaveText", "motion", "Letters bob in a wave.", demoWaveText},
	{363, "NumberTicker", "motion", "Counts to a value with easing.", demoNumberTicker},
	{364, "RollingDigits", "motion", "Digits roll on change.", demoRollingDigits},
	{365, "Marquee", "motion", "A row scrolling in a loop; hover pauses.", demoMarquee},
	{366, "VerticalMarquee", "motion", "A column scrolling in a loop.", demoVerticalMarquee},
	{367, "FlipCard", "motion", "Click to turn front to back.", demoFlipCard},
	{368, "TiltCard", "motion", "Turns toward the pointer; springs flat.", demoTiltCard},
	{369, "HoverZoom", "motion", "Content zooms inside its clip on hover.", demoHoverZoom},
	{370, "Parallax", "motion", "Layers move against the pointer by depth.", demoParallax},
	{371, "Bob", "motion", "Content floats up and down.", demoBob},
	{372, "Shake", "motion", "Shakes when asked.", demoShake},
	{373, "Wiggle", "motion", "Rocks a few degrees when asked.", demoWiggle},
	{374, "Jump", "motion", "Bounces once when asked.", demoJump},
	{375, "Breathe", "motion", "Opacity pulses in a loop.", demoBreathe},
	{376, "Beat", "motion", "Scales like a heartbeat.", demoBeat},
	{377, "Spinning", "motion", "Turns continuously.", demoSpinning},
	{378, "Orbit", "motion", "Items circle a center.", demoOrbit},
}

// replay builds a Replay button and the arrival it keys.
func replay(c *ui.Context, log func(string), arrive func(c *ui.Context)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	n := intState(root, "n", 0)
	root.Children(func() {
		if ui.Button(c, "Replay").Clicked() {
			*n++
			log(fmt.Sprintf("replay %d", *n))
		}
		arrive(c.Key(*n))
	})
}

func demoFadeIn(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	run := ui.Local(root, "run", func() uint64 { return 0 })
	root.Children(func() {
		if ui.Button(c, "Fade again").Clicked() {
			*run++
			log(fmt.Sprintf("run %d", *run))
		}
		motion.FadeIn(c, *run, func() { tile(c, "Hello again", 44) })
	})
}

func demoSlideIn(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	side := intState(root, "side", 0)
	names := []string{"bottom", "top", "start", "end"}
	root.Children(func() {
		if ui.Button(c, "Next edge").Clicked() {
			*side = (*side + 1) % 4
			log("from " + names[*side])
		}
		motion.SlideIn(c.Key(*side), motion.SlideInOptions{From: motion.Side(*side)}, func() { tile(c, "From the "+names[*side], 44) })
	})
}

func demoScaleIn(c *ui.Context, log func(string)) {
	replay(c, log, func(c *ui.Context) { motion.ScaleIn(c, func() { tile(c, "Scaled in", 44) }) })
}

func demoBlurIn(c *ui.Context, log func(string)) {
	replay(c, log, func(c *ui.Context) { motion.BlurIn(c, func() { tile(c, "In focus", 44) }) })
}

func demoFlipIn(c *ui.Context, log func(string)) {
	replay(c, log, func(c *ui.Context) { motion.FlipIn(c, func() { tile(c, "Turned in", 44) }) })
}

func demoStagger(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	n := intState(root, "n", 4)
	root.Children(func() {
		if ui.Button(c, "Add row").Clicked() {
			*n++
			log(fmt.Sprintf("rows %d", *n))
		}
		motion.Stagger(c, *n, func(i int) { tile(c, fmt.Sprintf("Row %d", i+1), 44) }, motion.StaggerOptions{})
	})
}

func demoTypewriter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	i := intState(root, "i", 0)
	lines := []string{"Small parts, put together well.", "Every part moves.", "Sound on every action."}
	root.Children(func() {
		if ui.Button(c, "Next line").Clicked() {
			*i = (*i + 1) % len(lines)
			log(fmt.Sprintf("line %d", *i+1))
		}
		motion.Typewriter(c, lines[*i], motion.TypewriterOptions{})
	})
}

func demoScramble(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	i := intState(root, "i", 0)
	words := []string{"CRYCHIC", "MOTION", "SIGNAL"}
	root.Children(func() {
		if ui.Button(c, "Next word").Clicked() {
			*i = (*i + 1) % len(words)
			log("word " + words[*i])
		}
		motion.Scramble(c, words[*i], motion.ScrambleOptions{}).FontSize(24).Bold()
	})
}

func demoWordRotate(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(6).Children(func() {
		ui.Text(c, "Build")
		if motion.WordRotate(c, []string{"fast", "calm", "clear"}, motion.WordRotateOptions{}).Changed() {
			log("rotated")
		}
	})
}

func demoWordFade(c *ui.Context, log func(string)) {
	replay(c, log, func(c *ui.Context) { motion.WordFade(c, "Each word arrives on its own beat.").Width(320) })
}

func demoLetterPull(c *ui.Context, log func(string)) {
	replay(c, log, func(c *ui.Context) { motion.LetterPull(c, "Crychic") })
}

func demoWaveText(c *ui.Context, log func(string)) { motion.WaveText(c, "Hello, wave") }

func demoNumberTicker(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	v := numberState(root, "v", 1280)
	root.Children(func() {
		ui.Row(c).Gap(8).Children(func() {
			if ui.Button(c, "Add 250").Clicked() {
				*v += 250
				log(fmt.Sprintf("value %.0f", *v))
			}
			if ui.Button(c, "Reset").Clicked() {
				*v = 0
				log("value 0")
			}
		})
		motion.NumberTicker(c, *v, motion.NumberTickerOptions{}).FontSize(32).Bold()
	})
}

func demoRollingDigits(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(12).AlignItems(ui.Center)
	v := intState(root, "v", 128)
	root.Children(func() {
		if ui.Button(c, "−1").Clicked() {
			*v--
			log(fmt.Sprintf("value %d", *v))
		}
		motion.RollingDigits(c, int64(*v))
		if ui.Button(c, "+1").Clicked() {
			*v++
			log(fmt.Sprintf("value %d", *v))
		}
	})
}

// hover logs the pointer entering and leaving e.
func hover(e ui.Element, log func(string), in, out string) {
	was := ui.Local(e, "demo.hover", func() bool { return false })
	if h := e.Hovered(); h != *was {
		*was = h
		if h {
			log(in)
		} else {
			log(out)
		}
	}
}

func demoMarquee(c *ui.Context, log func(string)) {
	m := motion.Marquee(c, motion.MarqueeOptions{}, func() {
		ui.Row(c).Gap(24).Children(func() {
			for _, s := range []string{"Free shipping", "New colors", "Members save ten percent"} {
				ui.Text(c, s)
			}
		})
	}).Width(360)
	hover(m, log, "paused", "resumed")
}

func demoVerticalMarquee(c *ui.Context, log func(string)) {
	m := motion.VerticalMarquee(c, motion.VerticalMarqueeOptions{Height: 120}, func() {
		ui.Column(c).Gap(10).Children(func() {
			for _, s := range []string{"Ada joined", "Build 42 passed", "Lin left a review", "Deploy finished"} {
				ui.Text(c, s)
			}
		})
	}).Width(240)
	hover(m, log, "paused", "resumed")
}

func demoFlipCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(240)
	flipped := boolState(root, "flipped", false)
	root.Children(func() {
		if motion.FlipCard(c, "Flip card", flipped, func() { ui.Text(c, "Question: what moves?") }, func() {
			ui.Text(c, "Answer: every part.")
		}).Changed() {
			log(fmt.Sprintf("flipped: %v", *flipped))
		}
	})
}

func demoTiltCard(c *ui.Context, log func(string)) {
	card := motion.TiltCard(c, func() {
		ui.Text(c, "Tilt me").Bold()
		ui.Text(c, "The card turns toward the pointer.")
	}).Width(260)
	hover(card, log, "tilting", "springs flat")
}

func demoHoverZoom(c *ui.Context, log func(string)) {
	t := c.Theme()
	z := motion.HoverZoom(c, func() {
		ui.Box(c).Size(260, 140).Background(t.Accent.Mix(t.Background, 0.6)).Center().Children(func() {
			ui.Icon(c, icon.Must("mountain")).Size(48, 48).TextColor(t.Accent)
		})
	})
	hover(z, log, "zoom in", "zoom out")
}

func demoParallax(c *ui.Context, log func(string)) {
	t := c.Theme()
	shades := []ui.Color{t.Accent.Mix(t.Background, 0.85), t.Accent.Mix(t.Background, 0.55), t.Accent}
	sizes := []float32{220, 120, 48}
	p := motion.Parallax(c, []float32{0.2, 0.6, 1.4}, func(i int) {
		ui.Box(c).Fill().Center().Children(func() { ui.Box(c).Size(sizes[i], sizes[i]).Radius(999).Background(shades[i]) })
	}, motion.ParallaxOptions{Height: 180}).Width(320)
	hover(p, log, "pointer in", "pointer out")
}

func demoBob(c *ui.Context, log func(string)) {
	motion.Bob(c, func() { ui.Icon(c, icon.Must("rocket")).Size(32, 32) })
}

// trigger builds a button that bumps a run for move.
func trigger(c *ui.Context, log func(string), label string, move func(run uint64)) {
	root := ui.Row(c).Gap(16).AlignItems(ui.Center)
	run := ui.Local(root, "run", func() uint64 { return 0 })
	root.Children(func() {
		if ui.Button(c, label).Clicked() {
			*run++
			log(fmt.Sprintf("run %d", *run))
		}
		move(*run)
	})
}

func demoShake(c *ui.Context, log func(string)) {
	trigger(c, log, "Wrong code", func(run uint64) { motion.Shake(c, run, func() { tile(c, "PIN", 44) }) })
}

func demoWiggle(c *ui.Context, log func(string)) {
	trigger(c, log, "Ring", func(run uint64) { motion.Wiggle(c, run, func() { ui.Icon(c, icon.Must("bell")).Size(24, 24) }) })
}

func demoJump(c *ui.Context, log func(string)) {
	trigger(c, log, "New message", func(run uint64) { motion.Jump(c, run, func() { ui.Icon(c, icon.Must("mail")).Size(24, 24) }) })
}

func demoBreathe(c *ui.Context, log func(string)) {
	motion.Breathe(c, func() { tile(c, "Listening…", 44) })
}

func demoBeat(c *ui.Context, log func(string)) {
	motion.Beat(c, func() { ui.Icon(c, icon.Must("heart")).Size(32, 32).TextColor(c.Theme().Danger) })
}

func demoSpinning(c *ui.Context, log func(string)) {
	motion.Spinning(c, func() { ui.Icon(c, icon.Must("settings")).Size(32, 32) })
}

func demoOrbit(c *ui.Context, log func(string)) {
	glyphs := []string{"star", "moon", "sun", "cloud"}
	motion.Orbit(c, len(glyphs), func(i int) { ui.Icon(c, icon.Must(glyphs[i])).Size(16, 16) }, motion.OrbitOptions{})
}
