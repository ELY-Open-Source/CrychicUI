package media

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MetronomeOptions tune a Metronome.
type MetronomeOptions struct {
	// Min and Max bound tempo; defaults 40, 220.
	Min, Max float64
	// Tempo names the stepper, Play the button.
	Tempo, Play string
	// Format writes the tempo; "%.0f BPM" when empty.
	Format string
}

// Metronome swings an arm at bpm and ticks each beat.
func Metronome(c *ui.Context, label string, bpm *float64, playing *bool, o MetronomeOptions) ui.Element {
	if o.Min == 0 {
		o.Min = 40
	}
	if o.Max == 0 {
		o.Max = 220
	}
	if label == "" || playing == nil || bpm == nil || !(o.Min > 0 && o.Min < o.Max) || math.IsInf(o.Max, 0) ||
		!(*bpm >= o.Min && *bpm <= o.Max) {
		panic(fmt.Sprintf("media: Metronome %q needs a label, state and a tempo within 0 < Min < Max", label))
	}
	def(&o.Tempo, "Tempo")
	def(&o.Play, "Play")
	def(&o.Format, "%.0f BPM")
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(4)).AlignItems(ui.Center).Label(label)
	now := c.Now()
	start := ui.Local(root, "start", func() time.Time { return time.Time{} })
	beat := ui.Local(root, "beat", func() int { return 0 })
	tempo := ui.Local(root, "tempo", func() float64 { return *bpm })
	period := func(b float64) float64 { return float64(time.Minute) / b }
	count := func() {
		switch {
		case !*playing:
			*start = time.Time{}
		case start.IsZero():
			*start, *beat, *tempo = now, 0, *bpm
		default:
			if *tempo != *bpm {
				ph := float64(now.Sub(*start)) / period(*tempo)
				*start, *tempo = now.Add(-time.Duration(ph*period(*bpm))), *bpm
			}
			for n := int(float64(now.Sub(*start)) / period(*bpm)); *beat < n; *beat++ {
				sound.Play(sound.Tick)
			}
		}
	}
	count()
	return root.Children(func() {
		body := ui.Box(c).Size(t.Space(24), t.Space(32)).Shrink(0)
		ui.Column(c).Gap(t.Space(2)).Children(func() {
			if button.PlayButton(c, playing, button.PlayButtonOptions{Label: o.Play}).Changed() {
				root.SetChanged()
				kit.Log("Metronome", "toggled", "playing", *playing)
				count()
			}
			if choice.NumberStepper(c, o.Tempo, bpm, o.Min, o.Max, 1, choice.NumberStepperOptions{Format: o.Format}).Changed() {
				root.SetChanged()
				count()
			}
		})
		const swing = 25
		var turn float64
		if *playing {
			p := period(*bpm)
			ph := float64(now.Sub(*start)) / p
			if c.Preferences().ReduceMotion {
				turn = swing
				if int(ph)%2 == 1 {
					turn = -swing
				}
			} else {
				turn = swing * math.Cos(math.Pi*ph)
				c.AnimationFrame()
			}
			c.After(time.Duration(p) - time.Duration(math.Mod(float64(now.Sub(*start)), p)))
		}
		body.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			var shell ui.Path
			shell.MoveTo(r.X+r.W*0.12, r.Y+r.H).LineTo(r.X+r.W*0.88, r.Y+r.H).LineTo(r.X+r.W*0.66, r.Y+r.H*0.04).
				LineTo(r.X+r.W*0.34, r.Y+r.H*0.04).Close()
			p.FillPath(&shell, t.Surface)
			p.StrokePath(&shell, 1, t.Border)
			px, py := r.X+r.W/2, r.Y+r.H*0.86
			a := turn * math.Pi / 180
			dx, dy := float32(math.Sin(a)), -float32(math.Cos(a))
			l := r.H * 0.78
			var arm ui.Path
			arm.MoveTo(px, py).LineTo(px+dx*l, py+dy*l)
			p.StrokePath(&arm, t.Space(0.5), t.Text)
			wx, wy, ws := px+dx*l*0.62, py+dy*l*0.62, t.Space(3)
			p.Fill(ui.Rect{X: wx - ws/2, Y: wy - ws/2, W: ws, H: ws}, t.Accent, ws*0.3)
			p.Fill(ui.Rect{X: px - t.Space(1), Y: py - t.Space(1), W: t.Space(2), H: t.Space(2)}, t.Text, t.Space(1))
		}))
	})
}
